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

// SignUp is a method for creating an admin user to SmartHome
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

// DeleteUser is a method that allows to remove an admin user from SmartHome
func (cl *Client) DeleteUser(c echo.Context) error {
	authParams := new(Auth)
	if err := json.NewDecoder(c.Request().Body).Decode(&authParams); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("Invalid payload: %s", err.Error()),
		)
	}
	if err := cl.SmartHomeInterface.DeleteUser(authParams.Username); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Error when deleting user from DynamoDB: %s", err.Error(),
		)
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Successfully deleted user",
		"user":    authParams.Username,
	})
}
