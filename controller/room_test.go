package controller

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

func TestSetRoomOptions(t *testing.T) {
	testCases := []struct {
		name,
		room string
		enabled bool
		thresholdOn,
		thresholdOff float32
		client      DynamoDBInterface
		expectedErr bool
	}{
		{
			name:         "Save a room",
			room:         "room",
			enabled:      true,
			thresholdOn:  19.3,
			thresholdOff: 19.5,
			client: &mockDynamoClient{
				putItemOutput: &dynamodb.PutItemOutput{},
			},
			expectedErr: false,
		},
		{
			name:         "Get an error from DynamoDB",
			room:         "room",
			enabled:      true,
			thresholdOn:  19.3,
			thresholdOff: 19.5,
			client: &mockDynamoClient{
				err: fmt.Errorf("Error"),
			},
			expectedErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			sh := NewSmartHome(SetDynamoDBClient(tc.client), SetLogger(mockLogger{}))
			err := sh.SetRoomOptions(tc.room, tc.enabled, tc.thresholdOn, tc.thresholdOff)
			if tc.expectedErr {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
		})
	}
}

func TestGetRoomOptions(t *testing.T) {
	testCases := []struct {
		name,
		room string
		client      DynamoDBInterface
		expected    map[string]types.AttributeValue
		expectedErr bool
	}{
		{
			name: "Get bedroom options",
			room: "bedroom",
			client: &mockDynamoClient{
				getItemOutput: &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"Room": &types.AttributeValueMemberS{Value: "bedroom"},
					},
				},
			},
			expected: map[string]types.AttributeValue{
				"Room": &types.AttributeValueMemberS{Value: "bedroom"},
			},
			expectedErr: false,
		},
		{
			name: "Room not found",
			room: "some_room",
			client: &mockDynamoClient{
				getItemOutput: &dynamodb.GetItemOutput{},
			},
			expectedErr: false,
		},
		{
			name: "Error from DynamoDB client",
			room: "bedroom",
			client: &mockDynamoClient{
				err: fmt.Errorf("Error"),
			},
			expectedErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			sh := NewSmartHome(SetDynamoDBClient(tc.client), SetLogger(mockLogger{}))
			actual, err := sh.GetRoomOptions(tc.room)
			if tc.expectedErr {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
			assert.EqualValues(tt, tc.expected, actual)
		})
	}
}

func TestDeleteRoomOptions(t *testing.T) {
	testCases := []struct {
		name,
		room string
		client      DynamoDBInterface
		expectedErr bool
	}{
		{
			name: "Delete room",
			room: "bedroom",
			client: &mockDynamoClient{
				deleteItemOutput: &dynamodb.DeleteItemOutput{},
			},
			expectedErr: false,
		},
		{
			name: "Error from DynamoDB",
			room: "bedroom",
			client: &mockDynamoClient{
				err: fmt.Errorf("Error"),
			},
			expectedErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			sh := NewSmartHome(SetDynamoDBClient(tc.client), SetLogger(mockLogger{}))
			err := sh.DeleteRoomOptions(tc.room)
			if tc.expectedErr {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
		})
	}
}
