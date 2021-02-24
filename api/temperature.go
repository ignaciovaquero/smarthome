package api

import "github.com/labstack/echo/v4"

// SetOutsideTemperature allows to set the current temperature
func (s *SmartHome) SetOutsideTemperature(c echo.Context) error {
	return nil
}

// GetOutsideTemperature gets the outside temperature for a given date
func (s *SmartHome) GetOutsideTemperature(c echo.Context) error {
	return nil
}

// SetInsideTemperature allows to set the current temperature inside
func (s *SmartHome) SetInsideTemperature(c echo.Context) error {
	return nil
}

// GetInsideTemperature gets the inside temperature for a given date and room
func (s *SmartHome) GetInsideTemperature(c echo.Context) error {
	return nil
}
