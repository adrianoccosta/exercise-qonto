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

// GetEnvInt64 returns environment int64
func GetEnvInt64(key string) int64 {
	r, ok := GetEnvInt64OrDefault(key, 0)
	if !ok {
		logFatal(ErrEnvironmentIsMissing{name: key})
	}

	return r
}

// GetEnvInt returns environment int
func GetEnvInt(key string) int {
	r, ok := GetEnvIntOrDefault(key, 0)
	if !ok {
		logFatal(ErrEnvironmentIsMissing{name: key})
	}

	return r
}

// GetEnvBool returns environment bool
func GetEnvBool(key string) bool {
	r, ok := GetEnvBoolOrDefault(key, false)
	if !ok {
		logFatal(ErrEnvironmentIsMissing{name: key})
	}

	return r
}

// EnvOrDefault returns the environment string if found, or the default value if not found.
func EnvOrDefault(key, defaultValue string) string {
	value, _ := GetEnvOrDefault(key, defaultValue)
	return value
}

// GetEnvOrDefault returns the environment string if found, or the default value if not found, as well as a flag that specifies if the env key was found or not.
func GetEnvOrDefault(key, defaultValue string) (string, bool) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue, ok
	}
	return value, ok
}

// EnvInt64OrDefault returns the environment int64 if found, or the default value if not found
func EnvInt64OrDefault(key string, defaultValue int64) int64 {
	intValue, _ := GetEnvIntOrDefault(key, int(defaultValue))
	return int64(intValue)
}

// GetEnvInt64OrDefault returns the environment int64 if found, or the default value if not found, as well as a flag that specifies if the env key was found or not.
func GetEnvInt64OrDefault(key string, defaultValue int64) (int64, bool) {
	intValue, err := GetEnvIntOrDefault(key, int(defaultValue))
	return int64(intValue), err
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

// EnvBoolOrDefault returns the environment bool if found, or the default value if not found
func EnvBoolOrDefault(key string, defaultValue bool) bool {
	value, _ := GetEnvBoolOrDefault(key, defaultValue)
	return value
}

// GetEnvBoolOrDefault returns the environment bool if found, or the default value if not found, as well as a flag that specifies if the env key was found or not.
func GetEnvBoolOrDefault(key string, defaultValue bool) (bool, bool) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue, ok
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue, false
	}

	return boolValue, ok
}
