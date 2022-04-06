package utils

import (
	"os"
	"strconv"
)

func EnvGet(key string, defaultValue interface{}) interface{} {

	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}

	switch defaultValue.(type) {
	case int:
		r, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return r
	case bool:
		r, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return r
	default:
		return value
	}
}
