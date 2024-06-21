package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	llms "github.com/qwark97/another_ai_project/llms/model"
	"github.com/qwark97/another_ai_project/server/ai/integrations/todoist"
	"github.com/qwark97/another_ai_project/server/model"
	"github.com/vargspjut/wlog"
)

const failureResponse = "Sorry, I can't now help you with that, please try again"

type LLM interface {
	Ask(ctx context.Context, request llms.Request) (llms.Response, error)
}

type Todoist interface {
	GetActiveTasks(ctx context.Context, params todoist.GetActiveTasksParams) ([]todoist.Task, error)
	GetProjects(ctx context.Context) ([]todoist.Project, error)
}

type Alice struct {
	model.Agent
	llm     LLM
	log     wlog.Logger
	todoist Todoist
}

func NewAlice(llm LLM, todoist Todoist, log wlog.Logger) *Alice {
	return &Alice{
		llm:     llm,
		log:     log,
		todoist: todoist,
	}
}

func (a *Alice) String() string {
	return "Alice"
}

type functions map[string]skill
type skill struct {
	functionCall llms.Function
	action       func(ctx context.Context, argumentsJSON string) (string, error)
}

func (a *Alice) definedFunctions() functions {
	return functions{
		"create_task": skill{
			functionCall: llms.Function{
				Name:        "create_task",
				Description: "Creates new task to do",
				Parameters: llms.Parameters{
					Type: "object",
					Properties: map[string]llms.Parameter{
						"content": {
							Type:        llms.String,
							Description: "Information which describes what needs to be done",
						},
						"date": {
							Type:        llms.String,
							Description: fmt.Sprintf("Information about the time when task should be done in format YYYY-MM-DD. By default pick '%s'.", time.Now().Format(time.DateOnly)),
						},
					},
					Required: []string{"content", "date"},
				},
			},
			action: func(ctx context.Context, argumentsJSON string) (string, error) {
				a.log.Info("create_task")
				return "create_task", nil
			},
		},
		"list_tasks": skill{
			functionCall: llms.Function{
				Name:        "list_tasks",
				Description: "Lists tasks to do for particular date",
				Parameters: llms.Parameters{
					Type: "object",
					Properties: map[string]llms.Parameter{
						"date": {
							Type:        llms.String,
							Description: fmt.Sprintf("Information about the time in format YYYY-MM-DD. By default pick '%s'.", time.Now().Format(time.DateOnly)),
						},
					},
					Required: []string{"date"},
				},
			},
			action: func(ctx context.Context, argumentsJSON string) (string, error) {
				var params todoist.GetActiveTasksParams
				err := json.Unmarshal([]byte(argumentsJSON), &params)
				if err != nil {
					return "", err
				}
				tasks, err := a.todoist.GetActiveTasks(ctx, params)
				if err != nil {
					return "", err
				}

				formattedTasks := func() string {
					var res string
					for _, task := range tasks {
						res += ", \"" + task.Content + "\""
					}
					return res
				}()

				return "Tasks to do: " + formattedTasks, nil
			},
		},
		"close_task": skill{
			functionCall: llms.Function{
				Name:        "close_task",
				Description: "Closes task using it's name and taskID. Task is closed when it has been done",
				Parameters: llms.Parameters{
					Type: "object",
					Properties: map[string]llms.Parameter{
						"name": {
							Type:        llms.String,
							Description: "Information about task's name",
						},
						"taskID": {
							Type:        llms.String,
							Description: "Information about task's ID",
						},
					},
					Required: []string{"name", "taskID"},
				},
			},
			action: func(ctx context.Context, argumentsJSON string) (string, error) {
				a.log.Info("close_task")
				return "close_task", nil
			},
		},
	}
}

// todo: in next step such conversation should be possible
// U: What tasks I have today
// A: X, Y, Z
// U: close X
// A: closed X
func (a *Alice) Do(ctx context.Context, instruction string, history []model.HistoryMessage) string {
	request := llms.Request{
		Model: llms.GPT_3_5_Turbo,
		Messages: func() []llms.Message {
			messages := []llms.Message{}
			for _, m := range history {
				messages = append(messages, llms.Message{
					Role:    string(m.Producer),
					Content: m.Content,
				})
			}
			messages = append(messages, llms.Message{
				Role:    "user",
				Content: instruction,
			})
			return messages
		}(),
		Functions: func() []llms.Function {
			var res []llms.Function
			for _, f := range a.definedFunctions() {
				res = append(res, f.functionCall)
			}
			return res
		}(),
	}
	response, err := a.llm.Ask(ctx, request)
	if err != nil {
		a.log.Error("alice:", err)
		return failureResponse
	}
	if len(response.Choices) == 0 {
		return failureResponse
	}
	functionCall := response.Choices[0].Message.FunctionCall

	if functionCall != nil {
		action := a.definedFunctions()[functionCall.Name].action
		if action == nil {
			return failureResponse
		}
		answer, err := action(ctx, functionCall.Arguments)
		if err != nil {
			return failureResponse
		}
		return answer
	} else {
		return response.Choices[0].Message.Content
	}
}
