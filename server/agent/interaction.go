package agent

import (
	"context"
	"fmt"
	"strconv"

	llmTypes "github.com/qwark97/assistant/llms/model"
)

type InteractionType uint

const (
	Question InteractionType = iota
	Command
	Information
	Other
)

func (it InteractionType) String() string {
	return []string{"Question", "Command", "Information", "Other"}[it]
}

func defineTypeOfInteraction(ctx context.Context, text string, llm LLM) (InteractionType, error) {

	question := llmTypes.Question{
		SystemPrompt: `
		You need to categorize user's input as one of the following category:
		0 = question - input which question or seems like it (it asks is about giving, showing, saying, explaining etc. something)
		1 = command - input which wants you to do something or perform some action
		2 = information - input which is some statement, new information or is about remembering something
		3 = other - input that does not fit into any other category

		As an answer you must return only one numer, nothing else. The number must be related with the category.
		`,
		UserQuestion: text,
	}

	response, err := llm.Ask(ctx, question)
	if err != nil {
		return Other, err
	}
	interactionType, err := strconv.Atoi(response)
	if err != nil {
		return Other, fmt.Errorf("failed to convert %s to number: %w", response, err)
	}
	return InteractionType(interactionType), nil
}
