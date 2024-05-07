package query

import (
	"fmt"
	"strings"
	"sync"
)

type Query struct {
	persona      Persona
	queryContext string
}

func New() Query {
	return Query{}
}

func (q Query) String() string {
	return strings.TrimSpace(fmt.Sprintf(`
	As a %s you will answer briefly and truthfully as possible.
	`,
		q.persona,
	))
}

func Enrich(query Query, features ...Feature) Query {
	query = enrichmentDefaults(query)
	var wg sync.WaitGroup
	for _, feature := range features {
		wg.Add(1)
		go func(f Feature) {
			defer wg.Done()
			f(&query)
		}(feature)
	}
	wg.Wait()
	return query
}

func enrichmentDefaults(query Query) Query {
	query.persona = Assistant
	return query
}
