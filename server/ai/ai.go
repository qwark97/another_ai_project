package ai

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/qwark97/another_ai_project/alog"
	"github.com/qwark97/another_ai_project/server/model"
)

type HistoryKeeper interface {
	Load(ctx context.Context, conversationID uuid.UUID) ([]model.HistoryMessage, error)
	Save(ctx context.Context, message model.HistoryMessage) error
}

type AgentSelector interface {
	Select(ctx context.Context, instruction string, history []model.HistoryMessage) Doer
}

type AI struct {
	history HistoryKeeper
	agents  AgentSelector
	log     alog.Logger
}

func New(h HistoryKeeper, a AgentSelector, log alog.Logger) AI {
	return AI{
		history: h,
		agents:  a,
		log:     log,
	}
}

func (ai AI) Act(ctx context.Context, request model.Request) model.Response {
	var response model.Response
	response.ConversationID = assureConversationID(request.ConversationID)

	historyMessages, err := ai.history.Load(ctx, response.ConversationID)
	if err != nil {
		ai.log.Error("failed to load user message from history:", err)
		response.Answer = "sorry, something went wrong"
		return response
	}
	go ai.remember(context.WithoutCancel(ctx), response.ConversationID, model.User, request.Instruction)

	agent := ai.agents.Select(ctx, request.Instruction, historyMessages)
	answer := agent.Do(ctx, request.Instruction, historyMessages)
	response.Answer = answer

	go ai.remember(context.WithoutCancel(ctx), response.ConversationID, model.Assistant, response.Answer)

	return response
}

func (ai AI) remember(ctx context.Context, conversationID uuid.UUID, owner model.Owner, information string) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	msg := model.NewHistoryMessage(conversationID, owner, information)
	err := ai.history.Save(ctx, msg)
	if err != nil {
		ai.log.Error("failed to save %s's message in history:", owner, err)
	}
}

func assureConversationID(id uuid.UUID) uuid.UUID {
	if id == uuid.Nil {
		return uuid.New()
	}
	return id
}
