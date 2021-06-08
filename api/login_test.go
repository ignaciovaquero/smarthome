package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	testCases := []struct {
		name          string
		ctx           mockContext
		cl            *Client
		errorExpected bool
	}{
		{
			name: "Request with correct payload, no authentication errors",
			ctx: &baseMockContext{
				Body: `{"username": "admin", "password": "admin"}`,
			},
			cl: NewClient(
				JWTConfig{
					JWTSecret:     "secret",
					JWTExpiration: time.Minute,
				},
				&mockSmartHome{},
			),
			errorExpected: false,
		},
		{
			name: "Request with invalid payload",
			ctx: &baseMockContext{
				Body: "Invalid payload",
			},
			cl: NewClient(
				JWTConfig{
					JWTSecret:     "secret",
					JWTExpiration: time.Minute,
				},
				&mockSmartHome{},
			),
			errorExpected: true,
		},
		{
			name: "Request with correct payload, authentication error",
			ctx: &baseMockContext{
				Body: `{"username": "admin", "password": "admin"}`,
			},
			cl: NewClient(
				JWTConfig{
					JWTSecret:     "secret",
					JWTExpiration: time.Minute,
				},
				&mockSmartHome{
					Err: fmt.Errorf("Error"),
				},
			),
			errorExpected: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.cl.Login(tc.ctx)
			if tc.errorExpected {
				assert.Error(tt, err)
				assert.IsType(tt, &echo.HTTPError{}, err)
				return
			}
			assert.NoError(tt, err)
			tok := tc.ctx.GetToken(tc.cl.Config.JWTSecret)
			assert.True(tt, tok.Valid)
		})
	}
}

func TestSignUp(t *testing.T) {
	testCases := []struct {
		name          string
		ctx           mockContext
		cl            *Client
		errorExpected bool
	}{
		{
			name: "Request with correct payload, no errors",
			ctx: &baseMockContext{
				Body: `{"username": "admin", "password": "admin"}`,
			},
			cl: NewClient(
				JWTConfig{},
				&mockSmartHome{},
			),
			errorExpected: false,
		},
		{
			name: "Request with invalid payload",
			ctx: &baseMockContext{
				Body: "Invalid payload",
			},
			cl: NewClient(
				JWTConfig{},
				&mockSmartHome{},
			),
			errorExpected: true,
		},
		{
			name: "Request with correct payload, controller error",
			ctx: &baseMockContext{
				Body: `{"username": "admin", "password": "admin"}`,
			},
			cl: NewClient(
				JWTConfig{},
				&mockSmartHome{
					Err: fmt.Errorf("Error"),
				},
			),
			errorExpected: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.cl.SignUp(tc.ctx)
			if tc.errorExpected {
				assert.Error(tt, err)
				assert.IsType(tt, &echo.HTTPError{}, err)
				return
			}
			assert.NoError(tt, err)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	testCases := []struct {
		name          string
		ctx           mockContext
		cl            *Client
		errorExpected bool
	}{
		{
			name: "Request with correct payload, no errors",
			ctx: &baseMockContext{
				Body: `{"username": "admin", "password": "admin"}`,
			},
			cl: NewClient(
				JWTConfig{},
				&mockSmartHome{},
			),
			errorExpected: false,
		},
		{
			name: "Request with invalid payload",
			ctx: &baseMockContext{
				Body: "Invalid payload",
			},
			cl: NewClient(
				JWTConfig{},
				&mockSmartHome{},
			),
			errorExpected: true,
		},
		{
			name: "Request with correct payload, controller error",
			ctx: &baseMockContext{
				Body: `{"username": "admin", "password": "admin"}`,
			},
			cl: NewClient(
				JWTConfig{},
				&mockSmartHome{
					Err: fmt.Errorf("Error"),
				},
			),
			errorExpected: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.cl.DeleteUser(tc.ctx)
			if tc.errorExpected {
				assert.Error(tt, err)
				assert.IsType(tt, &echo.HTTPError{}, err)
				return
			}
			assert.NoError(tt, err)
		})
	}
}
