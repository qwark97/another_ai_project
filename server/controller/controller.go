package controller

import (
	"context"

	"github.com/vargspjut/wlog"

	"github.com/qwark97/assistant/llms/model"
	"github.com/qwark97/assistant/server/controller/enrichers"
	"github.com/qwark97/assistant/server/controller/query"
)

type LLM interface {
	DetermineInteraction(ctx context.Context, instruction string) (model.InteractionMetadata, error)
	query.LLM
}

type Controller struct {
	log wlog.Logger
}

func New(log wlog.Logger) Controller {
	return Controller{
		log: log,
	}
}

func (c Controller) Interact(ctx context.Context, instruction string, llm LLM) string {
	interactionMetadata, err := llm.DetermineInteraction(ctx, instruction)
	if err != nil {
		c.log.Error(err.Error())
		return ""
	}
	instructionEnricher := enrichers.New(instruction)
	instructionEnricher.Type(interactionMetadata.Type)
	instructionEnricher.Category(interactionMetadata.Category)
	instructionEnricher.Tags(interactionMetadata.Tags)

	enrichedInstruction := instructionEnricher.Instruction()

	var response string
	switch enrichedInstruction.Type() {
	case enrichers.Action:
		response = "Nice action"
	case enrichers.Query:
		response = query.Answer(ctx, enrichedInstruction, llm)
	default:
		response = "Sorry, something went wrong"
	}

	return response
}
