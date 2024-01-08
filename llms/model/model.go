package model

type Category int
type Type int

type InteractionMetadata struct {
	Instruction string   `json:"instruction"`
	Category    Category `json:"category"`
	Type        Type     `json:"type"`
}
