package utils

import (
	"os"
	"strconv"
)

type Environment struct{}

var env = Environment{}

func GetEnvironment() *Environment {
	return &env
}

func (e *Environment) Get(key string, defaultValue interface{}) interface{} {

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
