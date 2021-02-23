package api

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func TestSetLogger(t *testing.T) {
	testCases := []struct {
		name     string
		logger   Logger
		expected *API
	}{
		{
			name: "Testing non setting anything",
			expected: &API{
				Logger:    &DefaultLogger{},
				TableName: DefaultTableName,
			},
			logger: nil,
		},
		{
			name: "Testing setting a default Logger",
			expected: &API{
				Logger:    &DefaultLogger{},
				TableName: DefaultTableName,
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

func TestSetTableName(t *testing.T) {
	testCases := []struct {
		name      string
		tableName func() *string
		expected  *API
	}{
		{
			name: "Testing setting a custom name",
			tableName: func() *string {
				s := "custom"
				return &s
			},
			expected: &API{
				Logger:    &DefaultLogger{},
				TableName: "custom",
			},
		},
		{
			name: "Testing setting an empty name",
			tableName: func() *string {
				s := ""
				return &s
			},
			expected: &API{
				Logger:    &DefaultLogger{},
				TableName: DefaultTableName,
			},
		},
		{
			name:      "Testing non setting anything",
			tableName: func() *string { return nil },
			expected: &API{
				Logger:    &DefaultLogger{},
				TableName: DefaultTableName,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			var actual *API
			if tc.tableName() != nil {
				actual = NewAPI(SetTableName(*tc.tableName()))
			} else {
				actual = NewAPI()
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
		expected *API
	}{
		{
			name:   "Set local DynamoDB Client",
			client: localClient,
			expected: &API{
				Logger:    &DefaultLogger{},
				TableName: DefaultTableName,
				Client:    localClient,
			},
		},
		{
			name: "Don't set anything DynamoDB Client",
			expected: &API{
				Logger:    &DefaultLogger{},
				TableName: DefaultTableName,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			var actual *API
			if tc.client != nil {
				actual = NewAPI(SetDynamoDBClient(tc.client))
			} else {
				actual = NewAPI()
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
