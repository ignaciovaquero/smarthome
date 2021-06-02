package api

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
