package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldEnrichQuery(t *testing.T) {
	assertion := assert.New(t)

	// given
	query := New()

	// when
	query = Enrich(query)

	// then
	assertion.Equal("", query.queryContext)
}
