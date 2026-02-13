package env

import (
	"os"
	"strconv"
)

// GetString returns the string value from an environment variable
func GetString(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}

	return value
}

// GetInt returns the int value from an environment variable
func GetInt(name string, defaultValue int) int {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}

	// Convert string to int
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

// GetBool returns the bool value from an environment variable
func GetBool(name string, defaultValue bool) bool {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}

	// Convert string to bool
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}
