package agents

import (
	"context"
	"fmt"

	"github.com/qwark97/another_ai_project/server/model"
)

type Alice model.Agent

func (a Alice) String() string {
	return "Alice"
}

func (a Alice) Do(ctx context.Context, instruction string, history []model.HistoryMessage) string {
	return fmt.Sprintf("I'm %s and I will manage your tasks to do", a)
}
