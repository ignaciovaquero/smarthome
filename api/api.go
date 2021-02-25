package api

import (
	"github.com/igvaquero18/smarthome/controller"
)

// Client is the API client for SmartHome
type Client struct {
	controller.SmartHomeInterface
}

// NewClient returns a new SmartHome API Client
func NewClient(smartHome controller.SmartHomeInterface) *Client {
	return &Client{smartHome}
}

type itemNotFound interface {
	NotFound() bool
}

func isNotFound(err error) bool {
	notFound, ok := err.(itemNotFound)
	return ok && notFound.NotFound()
}

type badRequest interface {
	badRequest() bool
}

func isBadRequest(err error) bool {
	br, ok := err.(badRequest)
	return ok && br.badRequest()
}
