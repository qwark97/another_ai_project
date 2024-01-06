package controller

type Type int

const (
	Query Type = iota + 1
	Action
)

type Category int

const (
	Unrecognized Category = iota
	TODO
)

var ActionsDescriptions = map[Category]string{
	Unrecognized: `Category which should be used as a fallback when no other category fits to the instructions`,
	TODO:         `Category should be used when the instructions say about some task that needs to be done. It doesn't include instructions that say something about meetings or some scheduled appointments`,
}

func (a Category) String() string {
	switch a {
	case TODO:
		return "TODO"
	default:
		return "Unrecognized"
	}
}
