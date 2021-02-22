package api

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

// DefaultTableName is the default Table name value for the
// SmartHome DynamoDB table.
const DefaultTableName = "SmartHome"

// API is a struct that defines the API actions for
// Smarthome app
type API struct {
	Logger
	*dynamodb.Client
	TableName string
}

// Option is a function to apply settings to Scraper structure
type Option func(a *API) Option

// NewAPI returns a new instance of an API
func NewAPI(opts ...Option) *API {
	a := &API{
		Logger:    &DefaultLogger{},
		TableName: DefaultTableName,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// SetLogger sets the Logger for the API
func SetLogger(logger Logger) Option {
	return func(a *API) Option {
		prev := a.Logger
		if logger != nil {
			a.Logger = logger
		}
		return SetLogger(prev)
	}
}

// SetDynamoDBClient sets the DynamoDB client for the API
func SetDynamoDBClient(client *dynamodb.Client) Option {
	return func(a *API) Option {
		prev := a.Client
		a.Client = client
		return SetDynamoDBClient(prev)
	}
}

// SetTableName sets the DynamoDB table name to be used
func SetTableName(name string) Option {
	return func(a *API) Option {
		prev := a.TableName
		if name != "" {
			a.TableName = name
		}
		return SetTableName(prev)
	}
}
