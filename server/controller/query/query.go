package query

import (
	"context"
	"time"

	"github.com/qwark97/assistant/llms/openai/model"
	"github.com/qwark97/assistant/server/controller/enrichers"
	"github.com/qwark97/assistant/server/controller/query/prompts"
	"github.com/vargspjut/wlog"
)

//go:generate mockery --name LLM --with-expecter
type LLM interface {
	Ask(ctx context.Context, system model.SystemPrompt, question model.UserPrompt, history ...model.Message) (string, error)
}

func Answer(ctx context.Context, enrichedInstruction enrichers.Instruction, llm LLM, log wlog.Logger) string {
	data := prompts.AskQuestionVariables{
		Facts: prompts.Facts{
			TodayRFC850: time.Now().Format(time.RFC850),
		},
		ContextData: "Current Golang version is 1.21",
	}
	system, err := prompts.GenerateAskQuestionPrompt(data)
	if err != nil {
		return err.Error()
	}

	response, err := llm.Ask(ctx, model.SystemPrompt(system), model.UserPrompt(enrichedInstruction.String()))
	if err != nil {
		log.Error(err)
		return "No idea"
	}

	return response
}
