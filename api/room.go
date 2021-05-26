package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/igvaquero18/smarthome/utils"
	"github.com/labstack/echo/v4"
)

const roomParam = "room"

var validRooms = []string{"all", "bedroom", "livingroom"}

type validRoom string

func (r validRoom) isValid() bool {
	for _, room := range validRooms {
		if string(r) == room {
			return true
		}
	}
	return false
}

// RoomOptions is a struct that represents the options available for a room
type RoomOptions struct {
	Enabled      bool    `json:"enabled"`
	ThresholdOn  float32 `json:"threshold_on"`
	ThresholdOff float32 `json:"threshold_off"`
}

// SetRoomOptions can enable or disable automating temperature
// adjust for a particular room or the whole home.
func (cl *Client) SetRoomOptions(c echo.Context) error {
	room := c.Param(roomParam)

	if !validRoom(room).isValid() {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("Invalid room name %s. Valid rooms: %v", room, validRooms),
		)
	}

	r := new(RoomOptions)
	if err := json.NewDecoder(c.Request().Body).Decode(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if r.ThresholdOn > r.ThresholdOff {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf(
				"threshold_on should be lower or equal to threshold_off. However we have: threshold_on = %f; threshold_off = %f",
				r.ThresholdOn,
				r.ThresholdOff,
			),
		)
	}

	rooms := []string{room}

	if room == "all" {
		rooms = utils.AllButOne(validRooms, "all")
	}

	for _, roomName := range rooms {
		if err := cl.SmartHomeInterface.SetRoomOptions(roomName, r.Enabled, r.ThresholdOn, r.ThresholdOff); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, struct {
		Message string      `json:"message"`
		Code    int         `json:"status_code"`
		Item    RoomOptions `json:"item"`
	}{
		Message: "successfully set room options",
		Code:    http.StatusOK,
		Item:    *r,
	})
}

// GetRoomOptions Gets the current temperature options for a given valid room
func (cl *Client) GetRoomOptions(c echo.Context) error {
	room := c.Param(roomParam)

	if !validRoom(room).isValid() {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("Invalid room name %s. Valid rooms: %v", room, validRooms),
		)
	}

	if room == "all" {
		rooms := utils.AllButOne(validRooms, "all")
		roomOpts := []map[string]RoomOptions{}
		for _, roomName := range rooms {
			item, err := cl.SmartHomeInterface.GetRoomOptions(roomName)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			if item == nil {
				continue
			}
			roomOpt := RoomOptions{}
			if err = attributevalue.UnmarshalMap(item, &roomOpt); err != nil {
				return echo.NewHTTPError(
					http.StatusInternalServerError,
					fmt.Sprintf("Error unmarshalling DynamoDB item: %s", err.Error()),
				)
			}
			roomOpts = append(roomOpts, map[string]RoomOptions{
				roomName: roomOpt,
			})
		}
		if len(roomOpts) == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "No rooms were found")
		}
		return c.JSON(http.StatusOK, roomOpts)
	}

	item, err := cl.SmartHomeInterface.GetRoomOptions(room)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if item == nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Room %s not found", room))
	}

	roomOpt := RoomOptions{}
	if err = attributevalue.UnmarshalMap(item, &roomOpt); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Error unmarshalling DynamoDB item: %s", err.Error()),
		)
	}

	return c.JSON(http.StatusOK, roomOpt)
}
