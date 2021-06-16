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

func (m *mockDynamoClient) DeleteItem(ctx context.Context, input *dynamodb.DeleteItemOutput, opts ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	return m.deleteItemOutput, m.err
}
