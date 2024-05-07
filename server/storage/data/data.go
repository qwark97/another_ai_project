package data

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/qwark97/assistant/server/model"
)

type Data struct {
	memory map[string]any
}

func New() *Data {
	return &Data{memory: make(map[string]any)}
}

func (s *Data) SaveHistoryRecord(ctx context.Context, message model.HistoryMessage) error {
	s.memory[fmt.Sprintf("history_%s", message.ID.String())] = message
	return nil
}

func (s *Data) LoadHistoryRecords(ctx context.Context, conversationID uuid.UUID) ([]model.HistoryMessage, error) {
	var res []model.HistoryMessage
	for key, val := range s.memory {
		if strings.HasPrefix(key, "history") {
			msg := val.(model.HistoryMessage)
			if msg.ConversationID == conversationID {
				res = append(res, msg)
			}
		}
	}
	return res, nil
}
