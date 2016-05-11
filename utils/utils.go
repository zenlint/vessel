package utils

import (
	"runtime"
	"strconv"
)

func TransMapToStr(sourceMap map[string]string) string {
	result := "{"

	for k, v := range sourceMap {
		result += `"` + k + `":"` + v + `"`
	}

	result += "}"
	return result
}

func CurrentLocation() string {
	location := "["
	if _, file, line, ok := runtime.Caller(1); ok {
		location += file + ":" + strconv.Itoa(line)
	}else{
		location += "???:???"
	}
	location += "] "
	return location
}
