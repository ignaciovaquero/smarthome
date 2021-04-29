package controller

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang.org/x/crypto/bcrypt"
)

// Authenticate returns an error if the combination of the username and
// password is incorrect.
func (s *SmartHome) Authenticate(username, password string) error {
	s.Debugw("Getting credentials for user", "user", username)
	credentials, err := s.get("username", username, s.Config.AuthTable)
	if err != nil {
		return fmt.Errorf("error getting user %s: %w", username, err)
	}

	hashedPassword, ok := credentials["password"].(*types.AttributeValueMemberS)

	if !ok {
		return fmt.Errorf("user %s not found", username)
	}

	s.Debugw("successfully retrieved credentials for user", "user", username)
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword.Value), []byte(password))
}

// SetCredentials stores the username and password for the user in the DynamoDB table.
// It takes care of hashing the password using the bcrypt package before storing it.
func (s *SmartHome) SetCredentials(username, password string) error {
	s.Debugw("Storing credentials for user", "user", username)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing the password: %w", err)
	}
	_, err = s.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &s.Config.AuthTable,
		Item: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{Value: username},
			"password": &types.AttributeValueMemberS{Value: string(hashedPassword)},
		},
	})

	if err != nil {
		return fmt.Errorf("error storing the user and password in the database: %w", err)
	}

	s.Debugw("successfully stored credentials for user", "user", username)
	return nil
}
