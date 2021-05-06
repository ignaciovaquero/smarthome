package api

import (
	"time"

	"github.com/igvaquero18/smarthome/controller"
)

// Client is the API client for SmartHome
type Client struct {
	Config JWTConfig
	controller.SmartHomeInterface
}

// JWTConfig is the configuration of the JWT parameters.
// This includes the JWT secret and the duration for the JWT
// before it expires.
type JWTConfig struct {
	JWTSecret     string
	JWTExpiration time.Duration
}

// NewClient returns a new SmartHome API Client
func NewClient(config JWTConfig, smartHome controller.SmartHomeInterface) *Client {
	return &Client{
		Config:             config,
		SmartHomeInterface: smartHome,
	}
}
