package query

import (
	"context"

	apiModel "github.com/qwark97/assistant/llms/openai/model"
)

type Instruction interface {
}

type LLM interface {
	Ask(ctx context.Context, query string, history ...apiModel.Message) (string, error)
}

func Answer(ctx context.Context, enrichedInstruction Instruction, llm LLM) string {
	// prepare query context

	// ask model

	return ""
}
