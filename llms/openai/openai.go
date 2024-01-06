package openai

import (
	"github.com/qwark97/assistant/llms/model"
)

type LLM struct {
}

func New() LLM {
	return LLM{}
}

func (llm LLM) DetermineInteraction(instruction string) (model.InteractionMetadata, error) {
	return model.InteractionMetadata{}, nil
}
