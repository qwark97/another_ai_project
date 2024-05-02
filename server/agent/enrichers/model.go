package enrichers

import (
	"time"

	"github.com/google/uuid"
)

type Type int

const (
	Action Type = iota - 1
	Unknown
	Query
)

var types = map[int]Type{-1: Action, 0: Unknown, 1: Query}

func toType(t int) Type {
	return types[t]
}

type Category int

const (
	Other Category = iota
	Todos
	Notes
)

var categories = map[string]Category{"other": Other, "todos": Todos, "notes": Notes}

func toCategory(c string) Category {
	return categories[c]
}

type Message struct {
	ID         uuid.UUID
	GroupID    uuid.UUID
	Producer   string
	InsertTime time.Time
}
