package utils

import (
	"os"
)

func FileExists(path string) bool {

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func DirExists(path string) bool {

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
