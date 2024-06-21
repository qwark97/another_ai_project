package openai

import (
	"context"
	"fmt"

	"github.com/vargspjut/wlog"

	"github.com/qwark97/another_ai_project/llms/model"
)

const (
	url          = "https://api.openai.com/v1/"
	embeddingURL = "https://api.openai.com/v1/embeddings"
)

type LLM struct {
	api api
	log wlog.Logger
}

func New(key string, log wlog.Logger) LLM {
	return LLM{
		api: api{
			log:          log,
			url:          url,
			embeddingURL: embeddingURL,
			key:          key,
		},
		log: log,
	}
}

func (llm LLM) Ask(ctx context.Context, request model.Request) (model.Response, error) {
	response, err := llm.api.askModel(ctx, request)
	if err != nil || response.Error != nil {
		return response, err
	}

	if len(response.Choices) < 1 {
		return response, fmt.Errorf("no choice")
	}

	return response, nil
}

func (llm LLM) GetEmbeddings(ctx context.Context, instruction string) ([]float32, error) {
	embeddingRequest := model.EmbeddingRequest{
		Input: instruction,
		Model: model.Text_Embedding_Ada_002,
	}
	response, err := llm.api.getEmbeddings(ctx, embeddingRequest)
	if err != nil {
		return nil, err
	}

	if len(response.Data) < 1 {
		return nil, fmt.Errorf("no embeddings")
	}

	return response.Data[0].Vector, nil
}
