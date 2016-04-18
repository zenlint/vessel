package utils

func TransMapToStr(sourceMap map[string]string) string {
	result := "{"

	for k, v := range sourceMap {
		result += `"` + k + `":"` + v + `"`
	}

	result += "}"
	return result
}
