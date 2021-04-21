package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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

	if err := cl.Authenticate(auth.Username, auth.Password); err != nil {
		return c.JSON(http.StatusForbidden, errorResponse{
			Message: fmt.Sprintf("Wrong username or password: %s", err.Error()),
			Code: http.StatusForbidden,
		})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// t, err := token.SignedString([]byte())

	return c.JSON(http.StatusOK, i interface{})
}
