package ai

import (
	"context"

	"github.com/qwark97/another_ai_project/llms/openai"
	"github.com/qwark97/another_ai_project/server/ai/agents"
	"github.com/qwark97/another_ai_project/server/ai/integrations/todoist"
	"github.com/qwark97/another_ai_project/server/model"
	"github.com/vargspjut/wlog"
)

type Doer interface {
	Do(ctx context.Context, instruction string, history []model.HistoryMessage) string
}

type Agents struct {
	llm     openai.LLM
	log     wlog.Logger
	todoist todoist.Todoist
}

func NewAgents(llm openai.LLM, todoist todoist.Todoist, log wlog.Logger) Agents {
	return Agents{
		llm:     llm,
		log:     log,
		todoist: todoist,
	}
}

func (a Agents) Select(ctx context.Context, instruction string, history []model.HistoryMessage) Doer {
	var selectedAgent Doer
	// todo: selection should be done by AI
	switch {
	default:
		selectedAgent = agents.NewAlice(a.llm, a.todoist, a.log)
	}
	return selectedAgent
}
