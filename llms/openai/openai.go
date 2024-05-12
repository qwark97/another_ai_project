package openai

import (
	"context"
	"fmt"

	"github.com/vargspjut/wlog"

	"github.com/qwark97/assistant/llms/model"
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

func (llm LLM) Ask(ctx context.Context, question model.Question, history ...model.Message) (string, error) {
	messages := prepareMessages(question, history)
	questionPrompt := model.Request{
		Model:    model.GPT_3_5_Turbo,
		Messages: messages,
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

func prepareMessages(question model.Question, history []model.Message) []model.Message {
	messages := []model.Message{
		{
			Role:    "system",
			Content: question.SystemPrompt,
		},
	}

	messages = append(messages, history...)

	messages = append(messages, model.Message{
		Role:    "user",
		Content: question.UserQuestion,
	})
	return messages
}
