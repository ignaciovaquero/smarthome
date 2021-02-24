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
	Enabled      bool    `json:"enabled"`
	ThresholdOn  float32 `json:"threshold_on"`
	ThresholdOff float32 `json:"threshold_off"`
}

type errorResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"status_code"`
	Params  interface{} `json:"params,omitempty"`
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

// SetRoomOptions can enable or disable automating temperature
// adjust for a particular room or the whole home.
func (s *SmartHome) SetRoomOptions(c echo.Context) error {
	room := Room(c.Param(roomParam))
	if !room.IsValid() {
		return c.JSON(http.StatusBadRequest, errorResponse{
			Message: "invalid room name provided",
			Code:    http.StatusBadRequest,
			Params: struct {
				Name string `json:"room_name"`
			}{
				Name: string(room),
			},
		})
	}
	r := new(roomItem)
	if err := json.NewDecoder(c.Request().Body).Decode(&r); err != nil {
		return c.JSON(http.StatusBadRequest, errorResponse{
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	if r.ThresholdOn >= r.ThresholdOff {
		return c.JSON(http.StatusBadRequest, errorResponse{
			Message: "threshold_on should be lower or equal to threshold_off",
			Code:    http.StatusBadRequest,
			Params: struct {
				ThresholdOn  float32 `json:"threshold_on"`
				ThresholdOff float32 `json:"threshold_off"`
			}{
				ThresholdOn:  r.ThresholdOn,
				ThresholdOff: r.ThresholdOff,
			},
		})
	}

	_, err := s.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &s.Config.ControlPlaneTable,
		Item: map[string]types.AttributeValue{
			"room":          &types.AttributeValueMemberS{Value: string(room)},
			"enabled":       &types.AttributeValueMemberBOOL{Value: r.Enabled},
			"threshold_on":  &types.AttributeValueMemberS{Value: fmt.Sprintf("%f", r.ThresholdOn)},
			"threshold_off": &types.AttributeValueMemberS{Value: fmt.Sprintf("%f", r.ThresholdOff)},
		},
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{
			Message: fmt.Sprintf("error putting value in DynamoDB: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, struct {
		Message string   `json:"message"`
		Code    int      `json:"status_code"`
		Item    roomItem `json:"item"`
	}{
		Message: "successfully added item to database",
		Code:    http.StatusOK,
		Item:    *r,
	})
}

// GetRoomOptions Gets the current temperature options for a given valid room
func (s *SmartHome) GetRoomOptions(c echo.Context) error {
	room := Room(c.Param(roomParam))
	if !room.IsValid() {
		return c.JSON(http.StatusBadRequest, errorResponse{
			Message: "invalid room name provided",
			Code:    http.StatusBadRequest,
			Params: struct {
				Name string `json:"room_name"`
			}{
				Name: string(room),
			},
		})
	}

	output, err := s.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{"room": &types.AttributeValueMemberS{Value: string(room)}},
		TableName: &s.Config.ControlPlaneTable,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{
			Message: fmt.Sprintf("error getting value from DynamoDB: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		})
	}

	if output.Item == nil {
		return c.JSON(http.StatusNotFound, errorResponse{
			Message: "Not found",
			Code:    http.StatusNotFound,
			Params: struct {
				Room string `json:"room"`
			}{
				Room: string(room),
			},
		})
	}

	return c.JSON(http.StatusOK, output.Item)
}
