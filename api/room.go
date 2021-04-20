package api

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/igvaquero18/smarthome/controller"
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

	if r.ThresholdOn > r.ThresholdOff {
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

	rooms := []string{room}

	if room == "all" {
		rooms = utils.AllButOne(validRooms, "all")
	}

	for _, roomName := range rooms {
		if err := cl.SmartHomeInterface.SetRoomOptions(roomName, r); err != nil {
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

	if room == "all" {
		rooms := utils.AllButOne(validRooms, "all")
		roomOpts := []map[string]types.AttributeValue{}
		for _, roomName := range rooms {
			item, err := cl.SmartHomeInterface.GetRoomOptions(roomName)
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
				continue
			}
			roomOpts = append(roomOpts, item)
		}
		if len(roomOpts) == 0 {
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
		return c.JSON(http.StatusOK, roomOpts)
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
