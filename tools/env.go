package tools

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// extracting to a variable that may be changed during a unit test to be able to capture the "exit" output
var logFatal = log.Fatal

// ErrEnvironmentIsMissing returns when table name is missing
type ErrEnvironmentIsMissing struct {
	name string
}

// Error return the error in string format.
func (e ErrEnvironmentIsMissing) Error() string {
	return fmt.Sprintf("the environment [%s] is missing.", e.name)
}

// GetEnv returns environment string
func GetEnv(key string) string {
	r, ok := os.LookupEnv(key)
	if !ok {
		logFatal(ErrEnvironmentIsMissing{name: key})
	}
	return r
}

// GetEnvIntOrDefault returns the environment int if found, or the default value if not found, as well as a flag that specifies if the env key was found or not.
func EnvIntOrDefault(key string, defaultValue int) int {
	value, _ := GetEnvIntOrDefault(key, defaultValue)
	return value
}

// GetEnvIntOrDefault returns the environment int if found, or the default value if not found, as well as a flag that specifies if the env key was found or not.
func GetEnvIntOrDefault(key string, defaultValue int) (int, bool) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue, ok
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue, false
	}

	return intValue, ok
}
