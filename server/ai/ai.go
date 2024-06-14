package ai

import (
	"context"

	"github.com/qwark97/another_ai_project/server/model"
)

type AI struct{}

func New() AI {
	return AI{}
}

func (ai AI) Ask(ctx context.Context, request model.Request) model.Response {
	return model.Response{
		Answer:         "good job",
		ConversationID: request.ConversationID,
	}
}
