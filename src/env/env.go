package env

import "os"

// LookupEnvWithDefault calls os.LookupEnv, and returns its response, but replaces the
// value returned by `default_value`, if the value is empty or not present in the os ENV
func LookupEnvWithDefault(key string, default_value string) (string, bool) {
	value, exists := os.LookupEnv(key)

	if value == "" {
		return default_value, exists
	}

	return value, exists
}
