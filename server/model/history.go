package model

import (
	"time"

	"github.com/google/uuid"
)

type owner string

const (
	System    owner = "system"
	User      owner = "user"
	Assistant owner = "assistant"
)

type HistoryMessage struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	Producer       owner
	InsertTime     time.Time
	Content        string
}

func NewHistoryMessage(conversationID uuid.UUID, producer owner, content string) HistoryMessage {
	return HistoryMessage{
		ID:             uuid.New(),
		ConversationID: conversationID,
		Producer:       producer,
		InsertTime:     time.Now(),
		Content:        content,
	}
}
