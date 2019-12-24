package utils

import "strings"

// IsEmpty ...
func IsEmpty(v interface{}) bool {

	switch v.(type) {
	case string:
		return len(strings.TrimSpace(v.(string))) == 0
	case int:
		return v.(int) == 0
	case bool:
		return v.(bool)
	default:
		return true
	}

}

// Contains ...
func Contains(items []interface{}, item interface{}) bool {

	for i := range items {
		if items[i] == item {
			return true
		}
	}
	return false
}
