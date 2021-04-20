package controller

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// GetCredentials returns the credentials for a given username from the
// DynamoDB table.
func (s *SmartHome) GetCredentials(username string) (map[string]types.AttributeValue, error) {
	s.Debugw("Getting credentials for user", "user", username)
	credentials, err := s.get(username, "username", s.Config.AuthTable)
	if err != nil {
		return map[string]types.AttributeValue{}, fmt.Errorf("error getting user %s: %w", username, err)
	}
	s.Debugw("successfully retrieved credentials for user", "user", username)

	return credentials, nil
}
