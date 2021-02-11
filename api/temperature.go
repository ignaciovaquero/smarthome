package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const roomParam = "room"

var validRooms = []Room{"all", "bedroom", "livingroom"}

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
	if room == "" {
		room = "all"
	}
	if !room.IsValid() {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid room name provided")
	}
	return nil
}
