package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/labstack/echo/v4"
)

const roomParam = "room"

var validRooms = []Room{"all", "bedroom", "livingroom"}

type roomItem struct {
	Enabled      bool    `json:"enabled,omitempty"`
	ThresholdOn  float32 `json:"threshold_on,omitempty"`
	ThresholdOff float32 `json:"threshold_off,omitempty"`
}

// Room is a wrapper around the string type to define what
// set of rooms are available to the temperature actions
type Room string

// IsValid returns whether a particular room is valid or not
func (r Room) IsValid() bool {
	for _, room := range validRooms {
		if r == room {
			return true
		}
	}
	return false
}

// AutoAdjustTemperature can enable or disable automating temperature
// adjust for a particular room or the whole home.
func (a *API) AutoAdjustTemperature(c echo.Context) error {
	var room Room
	room = Room(c.Param(roomParam))
	if !room.IsValid() {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid room name provided")
	}
	r := new(roomItem)
	if err := json.NewDecoder(c.Request().Body).Decode(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := a.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &a.TableName,
		Item: map[string]types.AttributeValue{
			"Room":          &types.AttributeValueMemberS{Value: string(room)},
			"enabled":       &types.AttributeValueMemberBOOL{Value: r.Enabled},
			"threshold_on":  &types.AttributeValueMemberS{Value: fmt.Sprintf("%f", r.ThresholdOn)},
			"threshold_off": &types.AttributeValueMemberS{Value: fmt.Sprintf("%f", r.ThresholdOff)},
		},
	})

	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("error putting value in DynamoDB: %s", err.Error()),
		)
	}

	a.Infow("successfully added item to database", "item", r)

	return nil
}
