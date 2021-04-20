package api

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login returns a valid JWT token
func (cl *Client) Login(c echo.Context) error {
	auth := new(auth)
	if err := json.NewDecoder(c.Request().Body).Decode(&auth); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	return nil
}
