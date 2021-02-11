package api

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	testCases := []struct {
		name     string
		room     Room
		expected bool
	}{
		{
			name:     "Test home room",
			room:     "home",
			expected: true,
		},
		{
			name:     "Test bedroom room",
			room:     "bedroom",
			expected: true,
		},
		{
			name:     "Test hall room",
			room:     "hall",
			expected: true,
		},
		{
			name:     "Test invalid room",
			room:     "invalid",
			expected: false,
		},
		{
			name:     "Test empty room",
			room:     "empty",
			expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			actual := tc.room.IsValid()
			assert.Equal(tt, tc.expected, actual)
		})
	}
}

func TestAutoAdjustTemperature(t *testing.T) {
	testCases := []struct {
		name     string
		context  echo.Context
		expected error
	}{
		{},
	}
	_ = testCases
}
