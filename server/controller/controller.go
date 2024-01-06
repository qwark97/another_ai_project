package controller

import (
	"fmt"
	"log/slog"

	"github.com/qwark97/assistant/llms/model"
)

type LLM interface {
	DetermineInteraction(instruction string) (model.InteractionMetadata, error)
}

type Controller struct {
	log *slog.Logger
}

func New(log *slog.Logger) Controller {
	return Controller{
		log: log,
	}
}

func (c Controller) RecogniseInteraction(instruction string, llm LLM) (Category, Type) {
	interactionMetadata, err := llm.DetermineInteraction(instruction)
	if err != nil {
		c.log.Error(err.Error())
		return Unrecognized, Query
	}
	fmt.Printf("%+v\n", interactionMetadata)

	return Unrecognized, Query
}
