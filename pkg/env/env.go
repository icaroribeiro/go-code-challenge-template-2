package env

import (
	"os"
)

// GetEnvWithDefaultValue is the function that tries to get the contents of an environment variable.
// In case of it does not exist, it is set a default value.
func GetEnvWithDefaultValue(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return value
}
