package handlers

import (
	"github.com/qwark97/assistant/llms/openai"
)

func defaultLLM() openai.LLM {
	return openai.New()
}
