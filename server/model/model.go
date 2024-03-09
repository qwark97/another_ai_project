package model

import "github.com/google/uuid"

type InteractionRequest struct {
	Instruction    string    `json:"instruction"`
	ConversationID uuid.UUID `json:"conversation_id"`
}
