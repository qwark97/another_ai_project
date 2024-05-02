package enrichers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToType(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected Type
	}{
		{
			name:     "test action",
			value:    -1,
			expected: Action,
		},
		{
			name:     "test unknown",
			value:    0,
			expected: Unknown,
		},
		{
			name:     "test query",
			value:    1,
			expected: Query,
		},
		{
			name:     "test unexpected positive value",
			value:    5,
			expected: Unknown,
		},
		{
			name:     "test unexpected negative value",
			value:    5,
			expected: Unknown,
		},
	}

	assertion := assert.New(t)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertion.Equal(tc.expected, toType(tc.value))
		})
	}

}

func TestToCategory(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected Category
	}{
		{
			name:     "test other",
			value:    "other",
			expected: Other,
		},
		{
			name:     "test todos",
			value:    "todos",
			expected: Todos,
		},
		{
			name:     "test notes",
			value:    "notes",
			expected: Notes,
		},
		{
			name:     "test unexpected empty value",
			value:    "",
			expected: Other,
		},
		{
			name:     "test unexpected value",
			value:    "something",
			expected: Other,
		},
	}

	assertion := assert.New(t)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertion.Equal(tc.expected, toCategory(tc.value))
		})
	}

}
