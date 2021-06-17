package controller

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type mockDynamoClient struct {
	getItemOutput    *dynamodb.GetItemOutput
	putItemOutput    *dynamodb.PutItemOutput
	deleteItemOutput *dynamodb.DeleteItemOutput
	err              error
}

func (m *mockDynamoClient) GetItem(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	return m.getItemOutput, m.err
}

func (m *mockDynamoClient) PutItem(ctx context.Context, input *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m.putItemOutput, m.err
}

func (m *mockDynamoClient) DeleteItem(ctx context.Context, input *dynamodb.DeleteItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return m.deleteItemOutput, m.err
}

type mockLogger struct{}

func (m mockLogger) Debug(...interface{}) {
	return
}

func (m mockLogger) Debugf(string, ...interface{}) {
	return
}

func (m mockLogger) Debugw(string, ...interface{}) {
	return
}

func (m mockLogger) Error(...interface{}) {
	return
}

func (m mockLogger) Errorf(string, ...interface{}) {
	return
}

func (m mockLogger) Errorw(string, ...interface{}) {
	return
}

func (m mockLogger) Info(...interface{}) {
	return
}

func (m mockLogger) Infof(string, ...interface{}) {
	return
}

func (m mockLogger) Infow(string, ...interface{}) {
	return
}
