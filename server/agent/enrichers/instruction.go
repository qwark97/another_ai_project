package enrichers

type Instruction struct {
	instruction string
	typ         Type
	category    Category
	tags        []string
	history     []Message
}

func (i Instruction) String() string {
	return i.instruction
}

func (i Instruction) Type() Type {
	return i.typ
}

func (i Instruction) Category() Category {
	return i.category
}

func (i Instruction) Tags() []string {
	return i.tags
}

func (i Instruction) History() []Message {
	return i.history
}
