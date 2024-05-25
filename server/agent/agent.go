package agent

import (
	"context"
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

type Embedding interface {
	Load(ctx context.Context, question string) (string, error)
}

type Agent struct {
	store     Store
	embedding Embedding
	llms      LLMsGroup
	log       wlog.Logger
}

type LLMsGroup struct {
	openai openai.LLM
}

func NewLLMsGroup(env map[string]string, log wlog.Logger) LLMsGroup {
	return LLMsGroup{
		openai: openai.New(env["OPENAI_KEY"], log),
	}
}

func New(store Store, embedding Embedding, llms LLMsGroup, log wlog.Logger) Agent {
	return Agent{
		store:     store,
		embedding: embedding,
		llms:      llms,
		log:       log,
	}
}

func (a Agent) Interact(ctx context.Context, request model.InteractionRequest) <-chan string {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	responsesPipe := make(ResponsesPipe)
	go func() {
		defer cancel()
		defer close(responsesPipe)

		if request.ConversationID == uuid.Nil {
			request.ConversationID = uuid.New()
		}

		eg, ctx := errgroup.WithContext(ctx)

		conversationHistory, err := a.store.LoadHistoryRecords(ctx, request.ConversationID)
		if err != nil {
			responsesPipe.SendError(ctx, err)
			return
		}

		requestContext, err := a.embedding.Load(ctx, request.Instruction)
		if err != nil {
			responsesPipe.SendError(ctx, err)
			return
		}

		eg.Go(func() error {
			err := a.store.SaveHistoryRecord(ctx, model.NewHistoryMessage(request.ConversationID, model.User, request.Instruction))
			if err != nil {
				a.log.Error(err)
			}
			return err
		})

		interactionType, err := defineTypeOfInteraction(ctx, request.Instruction, a.llms.openai)
		if err != nil {
			responsesPipe.SendError(ctx, err)
			return
		}
		a.log.Debugf("interaction type: %s", interactionType)

		switch interactionType {
		case Question:
			generatedMessages, err := a.answerQuestion(ctx, request.Instruction, requestContext, responsesPipe, conversationHistory, a.llms.openai)
			if err != nil {
				responsesPipe.SendError(ctx, err)
				return
			}
			eg.Go(func() error {
				for _, msg := range generatedMessages {
					err := a.store.SaveHistoryRecord(ctx, model.NewHistoryMessage(request.ConversationID, model.Assistant, msg))
					if err != nil {
						a.log.Error(err)
					}
				}
				return nil
			})
		default:
			a.reactToUnrecognisedInteraction(ctx, request.Instruction, responsesPipe)
		}
		eg.Wait()
	}()
	return responsesPipe
}

func (a Agent) answerQuestion(ctx context.Context, instruction, instructionContext string, responsesPipe ResponsesPipe, conversationHistory []model.HistoryMessage, llm LLM) ([]string, error) {
	generatedMessages, err := answerQuestion(ctx, instruction, instructionContext, responsesPipe, conversationHistory, llm)
	if err != nil {
		a.log.Error(err)
		return generatedMessages, err
	}
	a.log.Debug("answered question successfully")
	return generatedMessages, nil
}

func (a Agent) reactToUnrecognisedInteraction(ctx context.Context, _ string, responsesPipe ResponsesPipe) {
	responsesPipe.Send(ctx, "I'm not sure what you mean... Could you repeat your question?")
}
