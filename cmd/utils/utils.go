package utils

import (
	"fmt"
	"os"
)

func GetEnv(key string, fallback any) any {
	if value, ok := os.LookupEnv(key); ok {
		fmt.Println(key, " found from env.")
		return value
	}
	fmt.Println(key, " not found from env, settting fallback value.")
	return fallback
}
