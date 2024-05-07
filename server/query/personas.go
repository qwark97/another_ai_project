package query

type Persona uint

func (p Persona) String() string {
	switch p {
	case Assistant:
		return "helpful assistant"
	default:
		return "AI model"
	}
}

const (
	Assistant = iota
)
