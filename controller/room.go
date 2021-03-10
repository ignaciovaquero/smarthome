package controller

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// RoomOptions is a struct that represents the options available for a room
type RoomOptions struct {
	Enabled      bool    `json:"enabled"`
	ThresholdOn  float32 `json:"threshold_on"`
	ThresholdOff float32 `json:"threshold_off"`
}

// SetRoomOptions can enable or disable automating temperature
// adjust for a particular room or the whole home.
func (s *SmartHome) SetRoomOptions(room string, options *RoomOptions) error {
	s.Debugw("saving item in DynamoDB", "room", room, "options", options)

	_, err := s.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &s.Config.ControlPlaneTable,
		Item: map[string]types.AttributeValue{
			"room":          &types.AttributeValueMemberS{Value: string(room)},
			"enabled":       &types.AttributeValueMemberBOOL{Value: options.Enabled},
			"threshold_on":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%.1f", options.ThresholdOn)},
			"threshold_off": &types.AttributeValueMemberN{Value: fmt.Sprintf("%.1f", options.ThresholdOff)},
		},
	})

	if err != nil {
		return fmt.Errorf("error setting room %s with values %v in DynamoDB: %w", room, options, err)
	}

	s.Debugw("successfully saved item in DynamoDB", "room", room, "options", options)

	return nil
}

// GetRoomOptions Gets the current temperature options for a given room
func (s *SmartHome) GetRoomOptions(room string) (map[string]types.AttributeValue, error) {
	s.Debugw("getting item from DynamoDB", "room", room)
	output, err := s.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"room": &types.AttributeValueMemberS{Value: string(room)}},
		TableName: &s.Config.ControlPlaneTable,
	})

	if err != nil {
		return map[string]types.AttributeValue{}, fmt.Errorf("error getting room %s: %w", room, err)
	}

	s.Debugw("successfully retrieved item from DynamoDB", "room", room, "item", output.Item)

	return output.Item, nil
}
