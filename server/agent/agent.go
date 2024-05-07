package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vargspjut/wlog"
	"golang.org/x/sync/errgroup"

	llmTypes "github.com/qwark97/assistant/llms/model"

	"github.com/qwark97/assistant/llms/openai"

	"github.com/qwark97/assistant/server/model"
)

type LLM interface {
	Ask(ctx context.Context, question llmTypes.Question, history ...llmTypes.Message) (string, error)
}

type Store interface {
	SaveHistoryRecord(ctx context.Context, message model.HistoryMessage) error
	LoadHistoryRecords(ctx context.Context, groupID uuid.UUID) ([]model.HistoryMessage, error)
}

type Agent struct {
	store Store
	llms  LLMsGroup
	log   wlog.Logger
}

type LLMsGroup struct {
	openai openai.LLM
}

func NewLLMsGroup(env map[string]string, log wlog.Logger) LLMsGroup {
	return LLMsGroup{
		openai: openai.New(env["OPENAI_KEY"], log),
	}
}

func New(store Store, llms LLMsGroup, log wlog.Logger) Agent {
	return Agent{
		store: store,
		llms:  llms,
		log:   log,
	}
}

func (a Agent) Interact(ctx context.Context, request model.InteractionRequest) <-chan string {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	responsesCh := make(chan string)
	go func() {
		defer cancel()
		defer close(responsesCh)

		if request.ConversationID == uuid.Nil {
			request.ConversationID = uuid.New()
		}

		eg, ctx := errgroup.WithContext(ctx)

		_, err := a.store.LoadHistoryRecords(ctx, request.ConversationID)
		if err != nil {
			responsesCh <- err.Error()
			return
		}

		eg.Go(func() error {
			return a.store.SaveHistoryRecord(ctx, model.NewHistoryMessage(request.ConversationID, model.User, request.Instruction))
		})

		interactionType, err := defineTypeOfInteraction(ctx, request.Instruction, a.llms.openai)
		if err != nil {
			responsesCh <- err.Error()
			return
		}
		a.log.Debugf("interaction type: %s", interactionType)

		responsesCh <- fmt.Sprintf("%d", interactionType)

	}()
	return responsesCh
}
