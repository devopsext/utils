package utils

import (
	"os"
)

func FileExists(path string) bool {

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if info == nil {
		return false
	}
	return !info.IsDir()
}

func DirExists(path string) bool {

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if info == nil {
		return false
	}
	return info.IsDir()
}
