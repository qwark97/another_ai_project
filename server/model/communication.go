package model

import "github.com/google/uuid"

type Request struct {
	Instruction    string    `json:"instruction"`
	ConversationID uuid.UUID `json:"conversation_id"`
}

type Response struct {
	Answer         string    `json:"answer"`
	ConversationID uuid.UUID `json:"conversation_id"`
}
