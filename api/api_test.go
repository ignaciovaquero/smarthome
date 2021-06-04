package api

import (
	"testing"
	"time"

	"github.com/igvaquero18/smarthome/controller"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	m := &mockSmartHome{}
	testCases := []struct {
		name      string
		config    JWTConfig
		smarthome controller.SmartHomeInterface
		expected  *Client
	}{
		{
			name: "Creating a new client with config and interface",
			config: JWTConfig{
				JWTSecret:     "secret",
				JWTExpiration: time.Hour,
			},
			smarthome: m,
			expected: &Client{
				Config: JWTConfig{
					JWTSecret:     "secret",
					JWTExpiration: time.Hour,
				},
				SmartHomeInterface: m,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			actual := NewClient(tc.config, tc.smarthome)
			assert.Equal(tt, tc.expected, actual)
		})
	}
}
