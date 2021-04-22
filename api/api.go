package api

import (
	"github.com/igvaquero18/smarthome/controller"
)

// Client is the API client for SmartHome
type Client struct {
	JWTSecret string
	controller.SmartHomeInterface
}

// NewClient returns a new SmartHome API Client
func NewClient(jwtSecret string, smartHome controller.SmartHomeInterface) *Client {
	return &Client{
		JWTSecret:          jwtSecret,
		SmartHomeInterface: smartHome,
	}
}
