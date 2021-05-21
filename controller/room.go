package controller

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// SetRoomOptions can enable or disable automating temperature
// adjust for a particular room or the whole home.
func (s *SmartHome) SetRoomOptions(room string, enabled bool, thresholdOn, thresholdOff float32) error {
	s.Debugw("saving item in DynamoDB",
		"room", room,
		"enabled", enabled,
		"threshold_on", thresholdOn,
		"threshold_off", thresholdOff,
	)

	_, err := s.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &s.Config.ControlPlaneTable,
		Item: map[string]types.AttributeValue{
			"Room":         &types.AttributeValueMemberS{Value: room},
			"Enabled":      &types.AttributeValueMemberBOOL{Value: enabled},
			"ThresholdOn":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%.1f", thresholdOn)},
			"ThresholdOff": &types.AttributeValueMemberN{Value: fmt.Sprintf("%.1f", thresholdOff)},
		},
	})

	if err != nil {
		return fmt.Errorf(
			"error setting room %s with values Enabled=%t, ThresholdOn=%.1f, ThresholdOff=%.1f in DynamoDB: %w",
			room,
			enabled,
			thresholdOn,
			thresholdOff,
			err,
		)
	}

	s.Debugw("successfully saved item in DynamoDB",
		"room", room,
		"enabled", enabled,
		"threshold_on", thresholdOn,
		"threshold_off", thresholdOff,
	)

	return nil
}

// GetRoomOptions Gets the current temperature options for a given room
func (s *SmartHome) GetRoomOptions(room string) (map[string]types.AttributeValue, error) {
	s.Debugw("getting item from DynamoDB", "room", room)
	roomOptions, err := s.get("Room", room, s.Config.ControlPlaneTable)
	if err != nil {
		return map[string]types.AttributeValue{}, fmt.Errorf("error getting room %s: %w", room, err)
	}

	s.Debugw("successfully retrieved item from DynamoDB", "room", room, "item", roomOptions)

	return roomOptions, nil
}
