package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetLogger(t *testing.T) {
	testCases := []struct {
		name     string
		expected *API
		logger   Logger
	}{
		{
			name: "Testing non setting anything",
			expected: &API{
				Logger: &DefaultLogger{},
			},
			logger: nil,
		},
		{
			name: "Testing setting a default Logger",
			expected: &API{
				Logger: &DefaultLogger{},
			},
			logger: &DefaultLogger{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			actual := NewAPI(SetLogger(tc.logger))
			assert.Equal(tt, tc.expected, actual)
		})
	}
}
