package model

type Category string
type Type int

type InteractionMetadata struct {
	Instruction string   `json:"instruction"`
	Category    Category `json:"category"`
	Type        Type     `json:"type"`
	Tags        []string `json:"tags"`
}
