package utils

import (
	"encoding/json"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

type token struct {
	Token string `json:"token"`
}

// ValidateTokenFromBody validates a JWT token from a request body, with
// the secret used to sign the token.
func ValidateTokenFromBody(body, secret string) error {
	if secret == "" {
		return nil
	}

	t := new(token)

	if err := json.Unmarshal([]byte(body), &t); err != nil {
		return fmt.Errorf("error unmarshalling the token: %w", err)
	}

	tok, err := jwt.Parse(t.Token, func(to *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return fmt.Errorf("error parsing token: %w", err)
	}

	if !tok.Valid {
		return fmt.Errorf("token is invalid")
	}

	return nil
}
