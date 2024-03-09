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

type Message struct {
	ID         uuid.UUID
	GroupID    uuid.UUID
	Producer   owner
	InsertTime time.Time
}
