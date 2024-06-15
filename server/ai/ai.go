package ai

import (
	"context"

	"github.com/google/uuid"
	"github.com/qwark97/another_ai_project/server/model"
)

type HistoryLoader interface {
	Load(ctx context.Context, conversationID uuid.UUID) ([]model.HistoryMessage, error)
}

type AgentSelector interface {
	Select(ctx context.Context, instruction string, history []model.HistoryMessage) Doer
}

type AI struct {
	history HistoryLoader
	agents  AgentSelector
}

func New(h HistoryLoader, a AgentSelector) AI {
	return AI{
		history: h,
		agents:  a,
	}
}

func (ai AI) Act(ctx context.Context, request model.Request) model.Response {
	var response model.Response
	response.ConversationID = assureConversationID(request.ConversationID)

	historyMessages, err := ai.history.Load(ctx, response.ConversationID)
	if err != nil {
		response.Answer = "sorry, something went wrong"
		return response
	}

	agent := ai.agents.Select(ctx, request.Instruction, historyMessages)
	answer := agent.Do(ctx, request.Instruction, historyMessages)
	response.Answer = answer

	return response
}

func assureConversationID(id uuid.UUID) uuid.UUID {
	if id == uuid.Nil {
		return uuid.New()
	}
	return id
}
