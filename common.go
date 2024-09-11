package utils

import (
	"os"
	"reflect"
	"strings"
)

// IsEmpty ...
func IsEmpty(v interface{}) bool {

	switch v.(type) {
	case string:
		return len(strings.TrimSpace(v.(string))) == 0
	case int:
		return v.(int) == 0
	case float32:
		return v.(float32) == 0.0
	case float64:
		return v.(float64) == 0.0
	case bool:
		return v.(bool)
	case []string:
		arr := v.([]string)
		if len(arr) == 1 && len(strings.TrimSpace(arr[0])) == 0 {
			return true
		}
		return len(arr) == 0
	default:
		if v == nil {
			return true
		}
		return reflect.ValueOf(v).IsNil()
	}
}

// Contains ...
func Contains(items interface{}, item interface{}) bool {

	arrV := reflect.ValueOf(items)
	kind := arrV.Kind()

	if kind == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			// XXX - panics if slice element points to an unexported struct field
			// see https://golang.org/pkg/reflect/#Value.Interface
			if arrV.Index(i).Interface() == item {
				return true
			}
		}
	}
	return false
}

func Index(items interface{}, item interface{}) int {

	arrV := reflect.ValueOf(items)

	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			if arrV.Index(i).Interface() == item {
				return i
			}
		}
	}
	return -1
}

func Content(contentOrPath string) ([]byte, error) {

	if FileExists(contentOrPath) {
		b, err := os.ReadFile(contentOrPath)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	return []byte(contentOrPath), nil
}
