package controller

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

var defaultConfig *SmartHomeConfig = &SmartHomeConfig{
	AuthTable:         DefaultAuthTable,
	ControlPlaneTable: DefaultControlPlaneTable,
	TempOutsideTable:  DefaultTempOutsideTable,
	TempInsideTable:   DefaultTempInsideTable,
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
				Logger: &DefaultLogger{},
				Config: defaultConfig,
				Client: localClient,
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
