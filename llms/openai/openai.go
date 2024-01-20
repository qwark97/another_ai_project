package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/vargspjut/wlog"

	"github.com/qwark97/assistant/llms/model"
	apiModel "github.com/qwark97/assistant/llms/openai/model"
)

const url = "https://api.openai.com/v1/"

type LLM struct {
	api api
	log wlog.Logger
}

func New(key string, log wlog.Logger) LLM {
	return LLM{
		api: api{
			log: log,
			url: url,
			key: key,
		},
		log: log,
	}
}

func (llm LLM) Ask(ctx context.Context, system apiModel.SystemPrompt, question apiModel.UserPrompt, history ...apiModel.Message) (string, error) {
	questionPrompt := apiModel.Request{
		Model: apiModel.GPT_3_5_Turbo_0613,
		Messages: []apiModel.Message{
			{
				Role:    apiModel.System,
				Content: string(system),
			},
			{
				Role:    apiModel.User,
				Content: string(question),
			},
		},
	}
	response, err := llm.api.askModel(ctx, questionPrompt)
	if err != nil {
		return err.Error(), err
	}

	if len(response.Choices) < 1 {
		return "", fmt.Errorf("no choice")
	}

	return response.Choices[0].Message.Content, nil
}

func (llm LLM) DetermineInteraction(ctx context.Context, instruction string) (model.InteractionMetadata, error) {
	var result = model.InteractionMetadata{Instruction: instruction}

	tResp, err := llm.determineType(ctx, instruction)
	if err != nil {
		return result, err
	}
	var c1 map[string]string
	if err := parseTo(tResp, &c1); err != nil {
		return result, err
	}
	t, err := strconv.Atoi(c1["type"])
	if err != nil {
		return result, err
	}
	result.Type = t

	cResp, err := llm.determineCategory(ctx, instruction)
	if err != nil {
		return result, err
	}
	var c2 map[string]string
	if err := parseTo(cResp, &c2); err != nil {
		return result, err
	}
	result.Category = c2["category"]

	eResp, err := llm.enrichWithTags(ctx, instruction)
	if err != nil {
		return result, err
	}
	var c3 map[string][]string
	if err := parseTo(eResp, &c3); err != nil {
		return result, err
	}
	result.Tags = c3["tags"]

	return result, nil
}

func parseTo[T any](source apiModel.Response, container *map[string]T) error {
	if len(source.Choices) < 1 {
		return fmt.Errorf("not enough choices in source")
	}
	return json.Unmarshal([]byte(source.Choices[0].Message.FunctionCall.Arguments), container)
}

func (llm LLM) determineType(ctx context.Context, instruction string) (apiModel.Response, error) {
	typePrompt := apiModel.Request{
		Model: apiModel.GPT_3_5_Turbo_0613,
		Messages: []apiModel.Message{
			{
				Role:    apiModel.User,
				Content: instruction,
			},
		},
		FunctionCall: &apiModel.FunctionCall{
			Name: "determineInteractionType",
		},
		Functions: []apiModel.Function{
			{
				Name:        "determineInteractionType",
				Description: `Decide what is the type of the user's query`,
				Parameters: apiModel.Parameters{
					Type: "object",
					Properties: map[string]apiModel.Parameter{
						"type": {
							Type:        apiModel.Integer,
							Description: `Value equals to -1 when query is an ACTION to take (action other than saying/explaining/translating) Value is 1 when query is a QUESTION or asks to say/explain/translate something. Value is number`,
							Enum:        []any{-1, 1},
						},
					},
					Required: []string{"type"},
				},
			},
		},
	}
	return llm.api.askModel(ctx, typePrompt)
}

func (llm LLM) determineCategory(ctx context.Context, instruction string) (apiModel.Response, error) {
	categoryPrompt := apiModel.Request{
		Model: apiModel.GPT_3_5_Turbo_0613,
		Messages: []apiModel.Message{
			{
				Role:    apiModel.User,
				Content: instruction,
			},
		},
		FunctionCall: &apiModel.FunctionCall{
			Name: "determineInteractionCategory",
		},
		Functions: []apiModel.Function{
			{
				Name:        "determineInteractionCategory",
				Description: `Decide what is the category of the user's query`,
				Parameters: apiModel.Parameters{
					Type: "object",
					Properties: map[string]apiModel.Parameter{
						"category": {
							Type: apiModel.String,
							Description: `Category of the user's query. Category must describe about what is the query. Available categories:
							notes - when query is about creating new note or memorising something
							todos - when query is about creating a reminder
							other - when query does not fit to any other category`,
							Enum: []any{"notes", "todos", "other"},
						},
					},
					Required: []string{"category"},
				},
			},
		},
	}
	return llm.api.askModel(ctx, categoryPrompt)
}

func (llm LLM) enrichWithTags(ctx context.Context, instruction string) (apiModel.Response, error) {
	tagsPrompt := apiModel.Request{
		Model: apiModel.GPT_3_5_Turbo_0613,
		Messages: []apiModel.Message{
			{
				Role:    apiModel.User,
				Content: instruction,
			},
		},
		FunctionCall: &apiModel.FunctionCall{
			Name: "enrichWithTags",
		},
		Functions: []apiModel.Function{
			{
				Name:        "enrichWithTags",
				Description: `Enrich user's query with the list of tags related with the meaning of the query`,
				Parameters: apiModel.Parameters{
					Type: "object",
					Properties: map[string]apiModel.Parameter{
						"tags": {
							Type:        apiModel.Array,
							Description: `list of tags that describe meaning of the user's query; there must be at least one enrichment tag for the query`,
							Items: &apiModel.Item{
								Type: apiModel.String,
							},
						},
					},
					Required: []string{"tags"},
				},
			},
		},
	}
	return llm.api.askModel(ctx, tagsPrompt)
}
