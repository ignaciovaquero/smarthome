package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// Auth is a struct that holds the credentials (username and password)
// of a particular user
type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login returns a valid JWT token
func (cl *Client) Login(c echo.Context) error {
	authParams := new(Auth)
	if err := json.NewDecoder(c.Request().Body).Decode(&authParams); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("No valid username or password provided: %s", err.Error()),
		)
	}

	if err := cl.Authenticate(authParams.Username, authParams.Password); err != nil {
		return echo.NewHTTPError(
			http.StatusForbidden,
			fmt.Sprintf("Wrong username or password: %s", err.Error()),
		)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = authParams.Username
	claims["exp"] = time.Now().Add(cl.Config.JWTExpiration).Unix()

	t, err := token.SignedString([]byte(cl.Config.JWTSecret))
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Error signing token: %s", err.Error()),
		)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

// SignUp is a method for creating an admin user to the SmartHome Interface
func (cl *Client) SignUp(c echo.Context) error {
	authParams := new(Auth)
	if err := json.NewDecoder(c.Request().Body).Decode(&authParams); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("No valid username or password provided: %s", err.Error()),
		)
	}

	if err := cl.SetCredentials(authParams.Username, authParams.Password); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Error saving the credentials in the database: %s", err.Error()),
		)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Successfully signed up",
	})
}
