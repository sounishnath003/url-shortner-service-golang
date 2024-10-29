package utils

import (
	"fmt"
	"os"
)

// GetEnv Returns the value of the environment variable with the given key.
// If the environment variable is not set, it returns the fallback value.
func GetEnv(key string, fallback any) any {
	if value, ok := os.LookupEnv(key); ok {
		fmt.Println(key, " found from env.")
		return value
	}
	fmt.Println(key, " not found from env, settting fallback value.")
	return fallback
}
