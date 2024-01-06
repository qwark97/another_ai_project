package openai

import (
	"github.com/qwark97/assistant/llms/model"
)

type LLM struct {
	key string
}

func New(key string) LLM {
	return LLM{
		key: key,
	}
}

func (llm LLM) DetermineInteraction(instruction string) (model.InteractionMetadata, error) {
	return model.InteractionMetadata{}, nil
}
