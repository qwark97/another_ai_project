package controller

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vargspjut/wlog"
	"golang.org/x/sync/errgroup"

	"github.com/qwark97/assistant/llms/model"
	"github.com/qwark97/assistant/server/controller/enrichers"
	"github.com/qwark97/assistant/server/controller/query"
	serverModel "github.com/qwark97/assistant/server/model"
	storeModel "github.com/qwark97/assistant/server/storage/data/model"
)

type LLM interface {
	DetermineInteraction(ctx context.Context, instruction string) (model.InteractionMetadata, error)
	query.LLM
}

type Store interface {
	SaveHistoryRecord(ctx context.Context, message storeModel.Message) error
	LoadHistoryRecords(ctx context.Context, groupID uuid.UUID) ([]storeModel.Message, error)
}

type Controller struct {
	store Store
	log   wlog.Logger
}

func New(store Store, log wlog.Logger) Controller {
	return Controller{
		store: store,
		log:   log,
	}
}

func (c Controller) Interact(ctx context.Context, request serverModel.InteractionRequest, llm LLM) string {
	ctx, cacnel := context.WithTimeout(ctx, 5*time.Minute)
	defer cacnel()

	if request.ConversationID == uuid.Nil {
		request.ConversationID = uuid.New()
	}

	eg, ctx := errgroup.WithContext(ctx)

	instruction := request.Instruction
	conversationID := request.ConversationID

	history, err := c.store.LoadHistoryRecords(ctx, conversationID)
	if err != nil {
		c.log.Error(err)
		return "History failed"
	}
	eg.Go(
		func() error {
			return c.store.SaveHistoryRecord(ctx, storeModel.NewMessage(conversationID, storeModel.User, instruction))
		},
	)

	interactionMetadata, err := llm.DetermineInteraction(ctx, instruction)
	if err != nil {
		c.log.Error(err.Error())
		return ""
	}
	instructionEnricher := enrichers.New(instruction)
	instructionEnricher.Type(interactionMetadata.Type)
	instructionEnricher.Category(interactionMetadata.Category)
	instructionEnricher.Tags(interactionMetadata.Tags)

	instructionEnricher.History(transformHistory(history))

	enrichedInstruction := instructionEnricher.Instruction()

	var response string
	switch enrichedInstruction.Type() {
	case enrichers.Action:
		response = "Nice action"
		fallthrough
	case enrichers.Query:
		response = query.Answer(ctx, enrichedInstruction, llm, c.log)
	default:
		response = "Sorry, something went wrong"
	}
	if err = eg.Wait(); err != nil {
		c.log.Error(err)
	}
	return response
}

func transformHistory(history []storeModel.Message) []enrichers.Message {
	var res []enrichers.Message
	for _, msg := range history {
		res = append(res, enrichers.Message{
			ID:         msg.ID,
			GroupID:    msg.GroupID,
			Producer:   string(msg.Producer),
			InsertTime: msg.InsertTime,
		})
	}
	return res
}
