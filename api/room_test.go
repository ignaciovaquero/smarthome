package api

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/labstack/echo/v4"
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

func TestSetRoomOptions(t *testing.T) {
	testCases := []struct {
		name          string
		ctx           mockContext
		cl            *Client
		errorExpected bool
	}{
		{
			name: "Valid Payload, all rooms, no controller errors",
			ctx: &baseMockContext{
				Body:      `{"enabled": true, "threshold_on": 19.5, "threshold_off": 19.7}`,
				Parameter: "all",
			},
			cl:            NewClient(JWTConfig{}, &mockSmartHome{}),
			errorExpected: false,
		},
		{
			name: "Valid Payload, bedroom, no controller errors",
			ctx: &baseMockContext{
				Body:      `{"enabled": true, "threshold_on": 19.5, "threshold_off": 19.7}`,
				Parameter: "bedroom",
			},
			cl:            NewClient(JWTConfig{}, &mockSmartHome{}),
			errorExpected: false,
		},
		{
			name: "Valid Payload, livingroom, no controller errors",
			ctx: &baseMockContext{
				Body:      `{"enabled": true, "threshold_on": 19.5, "threshold_off": 19.7}`,
				Parameter: "livingroom",
			},
			cl:            NewClient(JWTConfig{}, &mockSmartHome{}),
			errorExpected: false,
		},
		{
			name: "Invalid Payload",
			ctx: &baseMockContext{
				Body:      "Invalid Payload",
				Parameter: "all",
			},
			cl:            NewClient(JWTConfig{}, &mockSmartHome{}),
			errorExpected: true,
		},
		{
			name: "Invalid room parameter",
			ctx: &baseMockContext{
				Body:      `{"enabled": true, "threshold_on": 19.5, "threshold_off": 19.7}`,
				Parameter: "fakeroom",
			},
			cl:            NewClient(JWTConfig{}, &mockSmartHome{}),
			errorExpected: true,
		},
		{
			name: "Threshold on is greater than threshold off",
			ctx: &baseMockContext{
				Body:      `{"enabled": true, "threshold_on": 19.6, "threshold_off": 19.5}`,
				Parameter: "all",
			},
			cl:            NewClient(JWTConfig{}, &mockSmartHome{}),
			errorExpected: true,
		},
		{
			name: "Controller errors",
			ctx: &baseMockContext{
				Body:      `{"enabled": true, "threshold_on": 19.5, "threshold_off": 19.6}`,
				Parameter: "all",
			},
			cl: NewClient(JWTConfig{}, &mockSmartHome{
				Err: fmt.Errorf("Error"),
			}),
			errorExpected: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.cl.SetRoomOptions(tc.ctx)
			if tc.errorExpected {
				assert.Error(tt, err)
				assert.IsType(tt, &echo.HTTPError{}, err)
				return
			}
			assert.NoError(tt, err)
		})
	}
}

func TestGetRoomOptions(t *testing.T) {
	testCases := []struct {
		name          string
		ctx           mockContext
		cl            *Client
		errorExpected bool
		expected      interface{}
	}{
		{
			name: "All rooms, no controller errors",
			ctx: &baseMockContext{
				Parameter: "all",
			},
			cl: NewClient(JWTConfig{}, &mockSmartHome{
				BedroomOpts: map[string]types.AttributeValue{
					"Name":         &types.AttributeValueMemberS{Value: "bedroom"},
					"ThresholdOn":  &types.AttributeValueMemberN{Value: "19.3"},
					"ThresholdOff": &types.AttributeValueMemberN{Value: "19.5"},
					"Enabled":      &types.AttributeValueMemberBOOL{Value: true},
				},
				LivingRoomOpts: map[string]types.AttributeValue{
					"Name":         &types.AttributeValueMemberS{Value: "livingroom"},
					"ThresholdOn":  &types.AttributeValueMemberN{Value: "19.3"},
					"ThresholdOff": &types.AttributeValueMemberN{Value: "19.5"},
					"Enabled":      &types.AttributeValueMemberBOOL{Value: true},
				},
			}),
			errorExpected: false,
			expected: []RoomOptions{
				{
					Name:         "bedroom",
					ThresholdOn:  19.3,
					ThresholdOff: 19.5,
					Enabled:      true,
				},
				{
					Name:         "livingroom",
					ThresholdOn:  19.3,
					ThresholdOff: 19.5,
					Enabled:      true,
				},
			},
		},
		{
			name: "All rooms, no results found",
			ctx: &baseMockContext{
				Parameter: "all",
			},
			cl:            NewClient(JWTConfig{}, &mockSmartHome{}),
			errorExpected: true,
		},
		{
			name: "Bedroom, results are found",
			ctx: &baseMockContext{
				Parameter: "bedroom",
			},
			cl: NewClient(JWTConfig{}, &mockSmartHome{
				BedroomOpts: map[string]types.AttributeValue{
					"Name":         &types.AttributeValueMemberS{Value: "bedroom"},
					"ThresholdOn":  &types.AttributeValueMemberN{Value: "19.3"},
					"ThresholdOff": &types.AttributeValueMemberN{Value: "19.5"},
					"Enabled":      &types.AttributeValueMemberBOOL{Value: true},
				},
			}),
			errorExpected: false,
			expected: RoomOptions{
				Name:         "bedroom",
				ThresholdOn:  19.3,
				ThresholdOff: 19.5,
				Enabled:      true,
			},
		},
		{
			name: "Bedroom, no results found",
			ctx: &baseMockContext{
				Parameter: "bedroom",
			},
			cl:            NewClient(JWTConfig{}, &mockSmartHome{}),
			errorExpected: true,
		},
		{
			name: "Livingroom, results are found",
			ctx: &baseMockContext{
				Parameter: "livingroom",
			},
			cl: NewClient(JWTConfig{}, &mockSmartHome{
				LivingRoomOpts: map[string]types.AttributeValue{
					"Name":         &types.AttributeValueMemberS{Value: "livingroom"},
					"ThresholdOn":  &types.AttributeValueMemberN{Value: "19.3"},
					"ThresholdOff": &types.AttributeValueMemberN{Value: "19.5"},
					"Enabled":      &types.AttributeValueMemberBOOL{Value: true},
				},
			}),
			errorExpected: false,
			expected: RoomOptions{
				Name:         "livingroom",
				ThresholdOn:  19.3,
				ThresholdOff: 19.5,
				Enabled:      true,
			},
		},
		{
			name: "Invalid room parameter",
			ctx: &baseMockContext{
				Parameter: "fakeroom",
			},
			cl:            NewClient(JWTConfig{}, &mockSmartHome{}),
			errorExpected: true,
		},
		{
			name: "Controller errors",
			ctx: &baseMockContext{
				Parameter: "all",
			},
			cl: NewClient(JWTConfig{}, &mockSmartHome{
				Err: fmt.Errorf("Error"),
			}),
			errorExpected: true,
		}}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.cl.GetRoomOptions(tc.ctx)
			if tc.errorExpected {
				assert.Error(tt, err)
				assert.IsType(tt, &echo.HTTPError{}, err)
				return
			}
			if tc.ctx.GetParameter() == "all" {
				actual, ok := tc.ctx.GetJSONPayload().([]RoomOptions)
				if !ok {
					assert.Fail(tt, "Actual should be of type []RoomOptions")
				}
				expected, ok := tc.expected.([]RoomOptions)
				if !ok {
					assert.Fail(tt, "Expected should be of type []RoomOptions")
				}
				assert.Equal(tt, expected, actual)
			} else {
				actual, ok := tc.ctx.GetJSONPayload().(RoomOptions)
				if !ok {
					assert.Fail(tt, "Actual should be of type RoomOptions")
				}
				expected, ok := tc.expected.(RoomOptions)
				if !ok {
					assert.Fail(tt, "Expected should be of type RoomOptions")
				}
				assert.Equal(tt, expected, actual)
			}
			assert.NoError(tt, err)
		})
	}
}

func TestDeleteRoomOptions(t *testing.T) {
	testCases := []struct {
		name          string
		ctx           mockContext
		cl            *Client
		errorExpected bool
	}{
		{
			name: "All rooms, no controller errors",
			ctx: &baseMockContext{
				Parameter: "all",
			},
			cl:            NewClient(JWTConfig{}, &mockSmartHome{}),
			errorExpected: false,
		},
		{
			name: "Bedroom",
			ctx: &baseMockContext{
				Parameter: "bedroom",
			},
			cl:            NewClient(JWTConfig{}, &mockSmartHome{}),
			errorExpected: false,
		},
		{
			name: "Livingroom",
			ctx: &baseMockContext{
				Parameter: "livingroom",
			},
			cl:            NewClient(JWTConfig{}, &mockSmartHome{}),
			errorExpected: false,
		},
		{
			name: "Invalid room parameter",
			ctx: &baseMockContext{
				Parameter: "fakeroom",
			},
			cl:            NewClient(JWTConfig{}, &mockSmartHome{}),
			errorExpected: true,
		},
		{
			name: "All rooms, controller errors",
			ctx: &baseMockContext{
				Parameter: "all",
			},
			cl: NewClient(JWTConfig{}, &mockSmartHome{
				Err: fmt.Errorf("Error"),
			}),
			errorExpected: true,
		},
		{
			name: "Bedroom, controller errors",
			ctx: &baseMockContext{
				Parameter: "bedroom",
			},
			cl: NewClient(JWTConfig{}, &mockSmartHome{
				Err: fmt.Errorf("Error"),
			}),
			errorExpected: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.cl.DeleteRoomOptions(tc.ctx)
			if tc.errorExpected {
				assert.Error(tt, err)
				assert.IsType(tt, &echo.HTTPError{}, err)
				return
			}
			assert.NoError(tt, err)
		})
	}
}
