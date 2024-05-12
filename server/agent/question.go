package agent

import (
	"context"

	llmTypes "github.com/qwark97/assistant/llms/model"
	"github.com/qwark97/assistant/server/model"
)

func answerQuestion(ctx context.Context, instruction string, responsesPipe ResponsesPipe, conversationHistory []model.HistoryMessage, llm LLM) ([]string, error) {
	question := llmTypes.Question{
		SystemPrompt: `
		You are helpfull assistant, answer the question using provided >>CONTEXT<< and history of our conversation. 
		Your answers should be sincere and truthful. If you don't know an answer, respond with "Sorry, I don't know"

		>>CONTEXT<<
		
		>>CONTEXT<<
		`,
		UserQuestion: instruction,
	}

	history := prepareConversationHistory(conversationHistory)
	response, err := llm.Ask(ctx, question, history...)
	if err != nil {
		return []string{}, err
	}
	responsesPipe.Send(ctx, response)
	return []string{response}, nil
}

func prepareConversationHistory(conversationHistory []model.HistoryMessage) []llmTypes.Message {
	var res []llmTypes.Message
	for _, msg := range conversationHistory {
		historyMsg := llmTypes.Message{
			Role:    string(msg.Producer),
			Content: msg.Content,
		}
		res = append(res, historyMsg)
	}
	return res
}
