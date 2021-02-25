package api

import (
	"encoding/json"
	"net/http"

	"github.com/igvaquero18/smarthome/controller"
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

type errorResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"status_code"`
	Params  interface{} `json:"params"`
}

// SetRoomOptions can enable or disable automating temperature
// adjust for a particular room or the whole home.
func (cl *Client) SetRoomOptions(c echo.Context) error {
	room := c.Param(roomParam)

	if !validRoom(room).isValid() {
		return c.JSON(http.StatusBadRequest, errorResponse{
			Message: "Invalid room name",
			Code:    http.StatusBadRequest,
			Params: struct {
				ValidRooms []string `json:"valid_rooms"`
				Room       string   `json:"room"`
			}{
				ValidRooms: validRooms,
				Room:       room,
			},
		})
	}

	r := new(controller.RoomOptions)
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

	err := cl.SmartHomeInterface.SetRoomOptions(room, r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{
			Message: "Internal Server Error",
			Code:    http.StatusInternalServerError,
			Params: struct {
				Error string `json:"error"`
			}{
				Error: err.Error(),
			},
		})
	}

	return c.JSON(http.StatusOK, struct {
		Message string                 `json:"message"`
		Code    int                    `json:"status_code"`
		Item    controller.RoomOptions `json:"item"`
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
		return c.JSON(http.StatusBadRequest, errorResponse{
			Message: "Invalid room name",
			Code:    http.StatusBadRequest,
			Params: struct {
				ValidRooms []string `json:"valid_rooms"`
				Room       string   `json:"room"`
			}{
				ValidRooms: validRooms,
				Room:       room,
			},
		})
	}

	item, err := cl.SmartHomeInterface.GetRoomOptions(room)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResponse{
			Message: "Internal Server Error",
			Code:    http.StatusInternalServerError,
			Params: struct {
				Error string `json:"error"`
			}{
				Error: err.Error(),
			},
		})
	}

	if item == nil {
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

	return c.JSON(http.StatusOK, item)
}
