package prompts

import (
	"bytes"
	"embed"
	"html/template"
)

type AskQuestionVariables struct {
	Facts       Facts
	ContextData string
}

type Facts struct {
	TodayRFC850 string
}

//go:embed templateComponents
//go:embed generalQuestionPrompt
var templates embed.FS

func GenerateAskQuestionPrompt(variables AskQuestionVariables) (string, error) {
	templ, err := template.ParseFS(templates, "generalQuestionPrompt", "templateComponents")
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	err = templ.Execute(buffer, variables)
	if err != nil {
		return "", err
	}

	return buffer.String(), err
}
