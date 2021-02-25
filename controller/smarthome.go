package controller

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	// DefaultControlPlaneTable is the default table name
	// for the Control Plane DynamoDB table.
	DefaultControlPlaneTable = "ControlPlane"

	// DefaultTempOutsideTable is the default table name
	// for the Temperature Outside DynamoDB table.
	DefaultTempOutsideTable = "TemperatureOutside"

	// DefaultTempInsideTable is the default table name
	// for the Temperature Inside DynamoDB table.
	DefaultTempInsideTable = "TemperatureInside"
)

// SmartHomeInterface is the interface implemented by the SmartHome API
type SmartHomeInterface interface {
	SetRoomOptions(room string, options *RoomOptions) error
	GetRoomOptions(room string) (map[string]types.AttributeValue, error)
}

// SmartHome is a struct that defines the API actions for
// Smarthome app
type SmartHome struct {
	Logger
	*dynamodb.Client
	Config *SmartHomeConfig
}

// SmartHomeConfig is a struct that allows to set all the configuration
// options for the SmartHome API
type SmartHomeConfig struct {
	// ControlPlaneTable is the name of the ControlPlane table in DynamoDB
	ControlPlaneTable string

	// TempOutsideTable is the name of the TemperatureOutside table in DynamoDB
	TempOutsideTable string

	// TempInsideTable is the name of the TemperatureInside table in DynamoDB
	TempInsideTable string
}

// Option is a function to apply settings to Scraper structure
type Option func(s *SmartHome) Option

// NewSmartHome returns a new instance of SmartHome
func NewSmartHome(opts ...Option) *SmartHome {
	a := &SmartHome{
		Logger: &DefaultLogger{},
		Config: &SmartHomeConfig{
			ControlPlaneTable: DefaultControlPlaneTable,
			TempOutsideTable:  DefaultTempOutsideTable,
			TempInsideTable:   DefaultTempInsideTable,
		},
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// SetLogger sets the Logger for the API
func SetLogger(logger Logger) Option {
	return func(s *SmartHome) Option {
		prev := s.Logger
		if logger != nil {
			s.Logger = logger
		}
		return SetLogger(prev)
	}
}

// SetDynamoDBClient sets the DynamoDB client for the API
func SetDynamoDBClient(client *dynamodb.Client) Option {
	return func(s *SmartHome) Option {
		prev := s.Client
		s.Client = client
		return SetDynamoDBClient(prev)
	}
}

// SetConfig sets the DynamoDB config
func SetConfig(c *SmartHomeConfig) Option {
	return func(s *SmartHome) Option {
		prev := s.Config

		if c.ControlPlaneTable == "" {
			c.ControlPlaneTable = DefaultControlPlaneTable
		}

		if c.TempOutsideTable == "" {
			c.TempOutsideTable = DefaultTempOutsideTable
		}

		if c.TempInsideTable == "" {
			c.TempInsideTable = DefaultTempInsideTable
		}

		s.Config = c
		return SetConfig(prev)
	}
}
