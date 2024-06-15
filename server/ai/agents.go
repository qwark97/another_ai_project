package ai

import (
	"context"

	"github.com/qwark97/another_ai_project/server/ai/agents"
	"github.com/qwark97/another_ai_project/server/model"
)

type Doer interface {
	Do(ctx context.Context, instruction string, history []model.HistoryMessage) string
}

type Agents struct {
}

func NewAgents() Agents {
	return Agents{}
}

func (a Agents) Select(ctx context.Context, instruction string, history []model.HistoryMessage) Doer {
	var selectedAgent Doer
	// todo: selection should be done by AI
	switch {
	default:
		selectedAgent = new(agents.Alice)
	}
	return selectedAgent
}
