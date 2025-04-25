package utils

import (
	"os"
)

// Environment is a utility class to manage environment variables.
type Environment struct{}

// Get retrieves the value of the environment variable named by the key.
// If the variable is not set, it returns the provided default value.
func (e *Environment) Get(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
