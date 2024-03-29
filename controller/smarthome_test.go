package controller

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

var defaultConfig *SmartHomeConfig = &SmartHomeConfig{
	AuthTable:         DefaultAuthTable,
	ControlPlaneTable: DefaultControlPlaneTable,
	TempOutsideTable:  DefaultTempOutsideTable,
	TempInsideTable:   DefaultTempInsideTable,
}

func getLocalClient() *dynamodb.Client {
	customResolver := aws.EndpointResolverFunc(func(service, awsRegion string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           "http://127.0.0.1:8000",
			SigningRegion: awsRegion,
		}, nil
	})
	cfg, _ := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolver(customResolver),
	)
	return dynamodb.NewFromConfig(cfg)
}

func TestSetLogger(t *testing.T) {
	testCases := []struct {
		name     string
		logger   Logger
		expected *SmartHome
	}{
		{
			name: "Testing non setting anything",
			expected: &SmartHome{
				Logger: &DefaultLogger{},
				Config: defaultConfig,
			},
			logger: nil,
		},
		{
			name: "Testing setting a default Logger",
			expected: &SmartHome{
				Logger: &DefaultLogger{},
				Config: defaultConfig,
			},
			logger: &DefaultLogger{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			actual := NewSmartHome(SetLogger(tc.logger))
			assert.Equal(tt, tc.expected, actual)
		})
	}
}

func TestSetConfig(t *testing.T) {
	testCases := []struct {
		name     string
		config   *SmartHomeConfig
		expected *SmartHome
	}{
		{
			name: "Testing setting custom names for all tables",
			config: &SmartHomeConfig{
				AuthTable:         "Auth",
				ControlPlaneTable: "Control",
				TempOutsideTable:  "Outside",
				TempInsideTable:   "Inside",
			},
			expected: &SmartHome{
				Logger: &DefaultLogger{},
				Config: &SmartHomeConfig{
					AuthTable:         "Auth",
					ControlPlaneTable: "Control",
					TempOutsideTable:  "Outside",
					TempInsideTable:   "Inside",
				},
			},
		},
		{
			name: "Testing setting an empty auth table",
			config: &SmartHomeConfig{
				AuthTable:         "",
				ControlPlaneTable: DefaultControlPlaneTable,
				TempOutsideTable:  DefaultTempOutsideTable,
				TempInsideTable:   DefaultTempInsideTable,
			},
			expected: &SmartHome{
				Logger: &DefaultLogger{},
				Config: &SmartHomeConfig{
					AuthTable:         DefaultAuthTable,
					ControlPlaneTable: DefaultControlPlaneTable,
					TempOutsideTable:  DefaultTempOutsideTable,
					TempInsideTable:   DefaultTempInsideTable,
				},
			},
		},
		{
			name: "Testing setting an empty control table",
			config: &SmartHomeConfig{
				AuthTable:         DefaultAuthTable,
				ControlPlaneTable: "",
				TempOutsideTable:  DefaultTempOutsideTable,
				TempInsideTable:   DefaultTempInsideTable,
			},
			expected: &SmartHome{
				Logger: &DefaultLogger{},
				Config: &SmartHomeConfig{
					AuthTable:         DefaultAuthTable,
					ControlPlaneTable: DefaultControlPlaneTable,
					TempOutsideTable:  DefaultTempOutsideTable,
					TempInsideTable:   DefaultTempInsideTable,
				},
			},
		},
		{
			name: "Testing setting an empty temperature outside table",
			config: &SmartHomeConfig{
				AuthTable:         DefaultAuthTable,
				ControlPlaneTable: DefaultControlPlaneTable,
				TempOutsideTable:  "",
				TempInsideTable:   DefaultTempInsideTable,
			},
			expected: &SmartHome{
				Logger: &DefaultLogger{},
				Config: &SmartHomeConfig{
					AuthTable:         DefaultAuthTable,
					ControlPlaneTable: DefaultControlPlaneTable,
					TempOutsideTable:  DefaultTempOutsideTable,
					TempInsideTable:   DefaultTempInsideTable,
				},
			},
		},
		{
			name: "Testing setting an empty temperature inside table",
			config: &SmartHomeConfig{
				AuthTable:         DefaultAuthTable,
				ControlPlaneTable: DefaultControlPlaneTable,
				TempOutsideTable:  DefaultTempOutsideTable,
				TempInsideTable:   "",
			},
			expected: &SmartHome{
				Logger: &DefaultLogger{},
				Config: &SmartHomeConfig{
					AuthTable:         DefaultAuthTable,
					ControlPlaneTable: DefaultControlPlaneTable,
					TempOutsideTable:  DefaultTempOutsideTable,
					TempInsideTable:   DefaultTempInsideTable,
				},
			},
		},
		{
			name: "Testing non setting anything",
			expected: &SmartHome{
				Logger: &DefaultLogger{},
				Config: defaultConfig,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			var actual *SmartHome
			if tc.config != nil {
				actual = NewSmartHome(SetConfig(tc.config))
			} else {
				actual = NewSmartHome()
			}
			assert.Equal(tt, tc.expected, actual)
		})
	}
}

func TestSetDynamoDBClient(t *testing.T) {
	localClient := getLocalClient()
	testCases := []struct {
		name     string
		client   *dynamodb.Client
		expected *SmartHome
	}{
		{
			name:   "Set local DynamoDB Client",
			client: localClient,
			expected: &SmartHome{
				Logger:            &DefaultLogger{},
				Config:            defaultConfig,
				DynamoDBInterface: localClient,
			},
		},
		{
			name: "Don't set anything DynamoDB Client",
			expected: &SmartHome{
				Logger: &DefaultLogger{},
				Config: defaultConfig,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			var actual *SmartHome
			if tc.client != nil {
				actual = NewSmartHome(SetDynamoDBClient(tc.client))
			} else {
				actual = NewSmartHome()
			}
			assert.Equal(tt, tc.expected, actual)
		})
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		name,
		hashkey,
		object,
		table string
		client      DynamoDBInterface
		expected    map[string]types.AttributeValue
		expectedErr bool
	}{
		{
			name:    "Get existing item",
			hashkey: "Item",
			object:  "1",
			table:   "table",
			client: &mockDynamoClient{
				getItemOutput: &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"Item": &types.AttributeValueMemberS{Value: "1"},
					},
				},
			},
			expected: map[string]types.AttributeValue{
				"Item": &types.AttributeValueMemberS{Value: "1"},
			},
			expectedErr: false,
		},
		{
			name:    "Error from DynamoDB client",
			hashkey: "Item",
			object:  "1",
			table:   "table",
			client: &mockDynamoClient{
				getItemOutput: &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"Item": &types.AttributeValueMemberS{Value: "1"},
					},
				},
				err: fmt.Errorf("Error"),
			},
			expected: map[string]types.AttributeValue{
				"Item": &types.AttributeValueMemberS{Value: "1"},
			},
			expectedErr: true,
		},
		{
			name:    "Non existing item",
			hashkey: "Item",
			object:  "2",
			table:   "table",
			client: &mockDynamoClient{
				getItemOutput: &dynamodb.GetItemOutput{
					Item: nil,
				},
			},
			expected:    nil,
			expectedErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			sh := NewSmartHome(SetDynamoDBClient(tc.client))
			actual, err := sh.get(tc.hashkey, tc.object, tc.table)
			if tc.expectedErr {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
			assert.EqualValues(tt, tc.expected, actual)
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		name,
		hashkey,
		object,
		table string
		client      DynamoDBInterface
		expectedErr bool
	}{
		{
			name:    "Delete existing item",
			hashkey: "Item",
			object:  "1",
			table:   "table",
			client: &mockDynamoClient{
				deleteItemOutput: &dynamodb.DeleteItemOutput{},
			},
			expectedErr: false,
		},
		{
			name:    "Error deleting item",
			hashkey: "Item",
			object:  "1",
			table:   "table",
			client: &mockDynamoClient{
				deleteItemOutput: &dynamodb.DeleteItemOutput{},
				err:              fmt.Errorf("Error"),
			},
			expectedErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			sh := NewSmartHome(SetDynamoDBClient(tc.client))
			err := sh.delete(tc.hashkey, tc.object, tc.table)
			if tc.expectedErr {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
		})
	}
}
