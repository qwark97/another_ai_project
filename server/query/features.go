package query

type Feature func(query *Query)

func WithPersona(persona Persona) Feature {
	return func(query *Query) {
		query.persona = persona
	}
}

func WithQueryContext() Feature {
	return func(query *Query) {
		query.queryContext = prepareContext()
	}
}

func prepareContext() string {
	return ""
}
