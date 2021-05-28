package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	testCases := []struct {
		name     string
		r        string
		expected bool
	}{
		{
			name:     "Test home room",
			r:        "all",
			expected: true,
		},
		{
			name:     "Test bedroom room",
			r:        "bedroom",
			expected: true,
		},
		{
			name:     "Test living room",
			r:        "livingroom",
			expected: true,
		},
		{
			name:     "Test invalid room",
			r:        "invalid",
			expected: false,
		},
		{
			name:     "Test empty room",
			r:        "",
			expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			room := ValidRoom(tc.r)
			actual := room.IsValid()
			assert.Equal(tt, tc.expected, actual)
		})
	}
}
