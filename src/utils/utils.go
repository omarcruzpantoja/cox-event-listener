package utils

import "os"

func GetEnv(key string, raise bool) string {
	// Dummy implementation for illustration purposes
	value := os.Getenv(key)
	if value == "" && raise {
		panic("Environment variable not set: " + key)
	}
	return value
}

