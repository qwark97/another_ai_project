package model

import (
	"time"

	"github.com/google/uuid"
)

type Owner string

const (
	System    Owner = "system"
	User      Owner = "user"
	Assistant Owner = "assistant"
)

type HistoryMessage struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	Producer       Owner
	InsertTime     time.Time
	Content        string
}

func NewHistoryMessage(conversationID uuid.UUID, producer Owner, content string) HistoryMessage {
	return HistoryMessage{
		ID:             uuid.New(),
		ConversationID: conversationID,
		Producer:       producer,
		InsertTime:     time.Now(),
		Content:        content,
	}
}
