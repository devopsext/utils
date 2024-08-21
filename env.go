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
	case int, int32:
		r, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return r
	case int64:
		r, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return int64(r)
	case bool:
		r, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return r
	case float32:
		r, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return defaultValue
		}
		return r
	case float64:
		r, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return defaultValue
		}
		return r
	default:
		return value
	}
}
