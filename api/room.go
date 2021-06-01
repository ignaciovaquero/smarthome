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

var ValidRooms = []string{"all", "bedroom", "livingroom"}

// ValidRoom is an alias to string that allow us to check whether a particular room
// name is valid
type ValidRoom string

// IsValid checks whether the name of the room is valid
func (r ValidRoom) IsValid() bool {
	for _, room := range ValidRooms {
		if string(r) == room {
			return true
		}
	}
	return false
}

// RoomOptions is a struct that represents the options available for a room
type RoomOptions struct {
	Name         string  `json:"name,omitempty"`
	Enabled      bool    `json:"enabled"`
	ThresholdOn  float32 `json:"threshold_on"`
	ThresholdOff float32 `json:"threshold_off"`
}

// SetRoomOptions can enable or disable automating temperature
// adjust for a particular room or the whole home.
func (cl *Client) SetRoomOptions(c echo.Context) error {
	room := c.Param(roomParam)

	if !ValidRoom(room).IsValid() {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("Invalid room name %s. Valid rooms: %v", room, ValidRooms),
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
		rooms = utils.AllButOne(ValidRooms, "all")
	}

	for _, roomName := range rooms {
		if err := cl.SmartHomeInterface.SetRoomOptions(roomName, r.Enabled, r.ThresholdOn, r.ThresholdOff); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, struct {
		Message string      `json:"message"`
		Code    int         `json:"status_code"`
		Options RoomOptions `json:"options"`
	}{
		Message: "successfully set room options",
		Code:    http.StatusOK,
		Options: *r,
	})
}

// GetRoomOptions Gets the current temperature options for a given valid room
func (cl *Client) GetRoomOptions(c echo.Context) error {
	room := c.Param(roomParam)

	if !ValidRoom(room).IsValid() {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("Invalid room name %s. Valid rooms: %v", room, ValidRooms),
		)
	}

	if room == "all" {
		rooms := utils.AllButOne(ValidRooms, "all")
		roomOpts := []RoomOptions{}
		for _, roomName := range rooms {
			item, err := cl.SmartHomeInterface.GetRoomOptions(roomName)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			if item == nil {
				continue
			}
			roomOpt := RoomOptions{Name: roomName}
			if err = attributevalue.UnmarshalMap(item, &roomOpt); err != nil {
				return echo.NewHTTPError(
					http.StatusInternalServerError,
					fmt.Sprintf("Error unmarshalling DynamoDB item: %s", err.Error()),
				)
			}
			roomOpts = append(roomOpts, roomOpt)
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

	roomOpt := RoomOptions{Name: room}
	if err = attributevalue.UnmarshalMap(item, &roomOpt); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			fmt.Sprintf("Error unmarshalling DynamoDB item: %s", err.Error()),
		)
	}

	return c.JSON(http.StatusOK, roomOpt)
}

func (cl *Client) DeleteRoomOptions(c echo.Context) error {
	room := c.Param(roomParam)
	if !ValidRoom(room).IsValid() {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			fmt.Sprintf("Invalid room name %s. Valid rooms: %v", room, ValidRooms),
		)
	}
	if room == "all" {
		for _, r := range ValidRooms {
			if err := cl.SmartHomeInterface.DeleteRoomOptions(r); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
		}
	} else {
		if err := cl.SmartHomeInterface.DeleteRoomOptions(room); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "successfully deleted room options",
		"status_code": http.StatusOK,
		"room":        room,
	})
}
