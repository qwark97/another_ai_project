package enrichers

type Enricher struct {
	instruction string
	typ         Type
	category    Category
	tags        []string
}

func New(instruction string) *Enricher {
	return &Enricher{
		instruction: instruction,
	}
}

func (e *Enricher) Type(t int) {
	e.typ = toType(t)
}

func (e *Enricher) Category(c string) {
	e.category = toCategory(c)
}

func (e *Enricher) Tags(t []string) {
	e.tags = t
}

func (e *Enricher) Instruction() Instruction {
	return Instruction{
		instruction: e.instruction,
		typ:         e.typ,
		category:    e.category,
		tags:        e.tags,
	}
}
