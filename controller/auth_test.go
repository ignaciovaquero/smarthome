package controller

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	testCases := []struct {
		name,
		username,
		password string
		client      DynamoDBInterface
		expectedErr bool
	}{
		{
			name:     "Authenticate valid username and password",
			username: "admin",
			password: "admin",
			client: &mockDynamoClient{
				getItemOutput: &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"Username": &types.AttributeValueMemberS{Value: "admin"},
						"Password": &types.AttributeValueMemberS{Value: "$2a$10$Xj/KDbv0lD/k0.WV7UxFq.tfHTEcnTCoowkKyMiIWCoj2cIobPF1C"},
					},
				},
			},
			expectedErr: false,
		},
		{
			name:     "Authenticate invalid password",
			username: "admin",
			password: "admin",
			client: &mockDynamoClient{
				getItemOutput: &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"Username": &types.AttributeValueMemberS{Value: "admin"},
						"Password": &types.AttributeValueMemberS{Value: "$2a$10$Xj/KDbv0lD/k0.WV7UxFq.tfHTEcnTCoowkKyMiIWCoj2cIabPF1C"},
					},
				},
			},
			expectedErr: true,
		},
		{
			name:     "Username not found",
			username: "admin2",
			password: "admin",
			client: &mockDynamoClient{
				getItemOutput: &dynamodb.GetItemOutput{
					Item: nil,
				},
			},
			expectedErr: true,
		},
		{
			name:     "Empty username",
			username: "",
			password: "admin",
			client: &mockDynamoClient{
				err: fmt.Errorf("Wrong username or password: error getting user : operation error DynamoDB: GetItem, https response error StatusCode: 400, RequestID: 990a39b9-92d2-42f1-bf4c-2650a213dbc9, api error ValidationException: One or more parameter values are not valid. The AttributeValue for a key attribute cannot contain an empty string value. Key: Username"),
			},
			expectedErr: true,
		},
		{
			name:     "Client error",
			username: "admin",
			password: "admin",
			client: &mockDynamoClient{
				err: fmt.Errorf("Error"),
			},
			expectedErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			sh := NewSmartHome(SetDynamoDBClient(tc.client), SetLogger(mockLogger{}))
			err := sh.Authenticate(tc.username, tc.password)
			if tc.expectedErr {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
		})
	}
}

func TestSetCredentials(t *testing.T) {
	testCases := []struct {
		name,
		username,
		password string
		client        DynamoDBInterface
		expectedError bool
	}{
		{
			name:          "Set valid username and password",
			username:      "admin",
			password:      "password",
			client:        &mockDynamoClient{},
			expectedError: false,
		},
		{
			name:          "Set empty password",
			username:      "admin",
			password:      "",
			client:        &mockDynamoClient{},
			expectedError: false,
		},
		{
			name:     "Set empty username",
			username: "admin",
			password: "",
			client: &mockDynamoClient{
				err: fmt.Errorf("Error saving the credentials in the database: error storing the user and password in the database: operation error DynamoDB: PutItem, https response error StatusCode: 400, RequestID: c309b965-d0da-4cdc-9dc2-7244cdad2312, api error ValidationException: One or more parameter values are not valid. The AttributeValue for a key attribute cannot contain an empty string value. Key: Username"),
			},
			expectedError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			sh := NewSmartHome(SetDynamoDBClient(tc.client), SetLogger(mockLogger{}))
			err := sh.SetCredentials(tc.username, tc.password)
			if tc.expectedError {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	testCases := []struct {
		name,
		username string
		client      DynamoDBInterface
		expectedErr bool
	}{
		{
			name:        "Delete existing user",
			username:    "admin",
			client:      &mockDynamoClient{},
			expectedErr: false,
		},
		{
			name:        "Delete non-existing user",
			username:    "admin",
			client:      &mockDynamoClient{},
			expectedErr: false,
		},
		{
			name:     "Delete empty user",
			username: "",
			client: &mockDynamoClient{
				err: fmt.Errorf("Error"),
			},
			expectedErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			sh := NewSmartHome(SetDynamoDBClient(tc.client), SetLogger(mockLogger{}))
			err := sh.DeleteUser(tc.username)
			if tc.expectedErr {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
		})
	}
}
