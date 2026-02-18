package handlers

import (
	"strconv"

	"github.com/labstack/echo/v5"
)

// Helper method to get int param or default value
func GetIntParam(c *echo.Context, name string, defaultValue int) int {
	value, err := strconv.Atoi(c.ParamOr(name, strconv.Itoa(defaultValue)))
	if err != nil {
		return defaultValue
	}

	return value
}

// Helper method to get bool param or default value
func GetBoolParam(c *echo.Context, name string, defaultValue bool) bool {
	value, err := strconv.ParseBool(c.ParamOr(name, strconv.FormatBool(defaultValue)))
	if err != nil {
		return defaultValue
	}

	return value
}
