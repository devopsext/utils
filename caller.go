package utils

import (
	"path"
	"runtime"
)

func callerLastPath(s string, limit int) string {

	index := 0
	dir := s
	var arr []string

	for !IsEmpty(dir) {
		if index >= limit {
			break
		}
		index++
		arr = append([]string{path.Base(dir)}, arr...)
		dir = path.Dir(dir)
	}
	return path.Join(arr...)
}

func CallerGetInfo(offset int) (string, string, int) {

	pc := make([]uintptr, 15)
	n := runtime.Callers(offset, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	function := callerLastPath(frame.Function, 1)
	file := callerLastPath(frame.File, 3)
	line := frame.Line

	return function, file, line
}
