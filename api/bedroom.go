package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Something is a test function
func Something(c echo.Context) error {
	return c.String(http.StatusOK, "tutto bene")
}
