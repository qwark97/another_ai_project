package controller

import (
	"context"
	"fmt"

	"github.com/vargspjut/wlog"

	"github.com/qwark97/assistant/llms/model"
)

type LLM interface {
	DetermineInteraction(ctx context.Context, instruction string) (model.InteractionMetadata, error)
}

type Controller struct {
	log wlog.Logger
}

func New(log wlog.Logger) Controller {
	return Controller{
		log: log,
	}
}

func (c Controller) RecogniseInteraction(ctx context.Context, instruction string, llm LLM) (Category, Type) {
	interactionMetadata, err := llm.DetermineInteraction(ctx, instruction)
	if err != nil {
		c.log.Error(err.Error())
		return Unrecognized, Query
	}
	fmt.Printf("interactionMetadata: %+v\n", interactionMetadata)

	return Unrecognized, Query
}
