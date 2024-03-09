package data

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/qwark97/assistant/server/storage/data/model"
)

type Store struct {
	memory map[string]any
}

func New() *Store {
	return &Store{memory: make(map[string]any)}
}

func (s *Store) SaveHistoryRecord(ctx context.Context, message model.Message) error {
	s.memory[fmt.Sprintf("history_%s", message.ID.String())] = message
	return nil
}

func (s *Store) LoadHistoryRecords(ctx context.Context, groupID uuid.UUID) ([]model.Message, error) {
	var res []model.Message
	for key, val := range s.memory {
		if strings.HasPrefix(key, "history") {
			msg := val.(model.Message)
			res = append(res, msg)
		}
	}
	return res, nil
}
