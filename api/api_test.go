package api

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/igvaquero18/smarthome/controller"
	"github.com/stretchr/testify/assert"
)

type mockSmartHome struct{}

func (m *mockSmartHome) Authenticate(username, password string) error {
	return nil
}
func (m *mockSmartHome) SetCredentials(username, password string) error {
	return nil
}
func (m *mockSmartHome) SetRoomOptions(room string, enabled bool, thresholdOn, thresholdOff float32) error {
	return nil
}
func (m *mockSmartHome) GetRoomOptions(room string) (map[string]types.AttributeValue, error) {
	return map[string]types.AttributeValue{}, nil
}
func (m *mockSmartHome) DeleteRoomOptions(room string) error {
	return nil
}
func (m *mockSmartHome) DeleteUser(username string) error {
	return nil
}

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
