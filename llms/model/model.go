package model

type InteractionMetadata struct {
	Instruction string   `json:"instruction"`
	Category    string   `json:"category"`
	Type        int      `json:"type"`
	Tags        []string `json:"tags"`
}
