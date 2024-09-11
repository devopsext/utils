package utils

import (
	"fmt"
	"strings"
)

func MapGetKeyValuesEx(s string, delim string) map[string]string {

	pairs := strings.Split(s, ",")

	m := make(map[string]string)

	for _, p := range pairs {

		if IsEmpty(p) {
			continue
		}
		kv := strings.SplitN(p, delim, 2)
		k := strings.TrimSpace(kv[0])
		if len(kv) > 1 {
			v := strings.TrimSpace(kv[1])
			if strings.HasPrefix(v, "${") && strings.HasSuffix(v, "}") {
				ed := strings.SplitN(v[2:len(v)-1], ":", 2)
				e, d := ed[0], ed[1]
				v = EnvGet(e, "").(string)
				if v == "" && d != "" {
					v = d
				}
			}
			m[k] = v
		} else {
			m[k] = ""
		}
	}
	return m
}

func MapGetKeyValues(s string) map[string]string {

	return MapGetKeyValuesEx(s, "=")
}

func MapToArrayWithSeparator(m map[string]string, s string) []string {

	var arr []string
	if m == nil {
		return arr
	}
	for k, v := range m {
		if IsEmpty(v) {
			arr = append(arr, k)
		} else {
			arr = append(arr, fmt.Sprintf("%s%s%v", k, s, v))
		}
	}
	return arr
}

func MapToArray(m map[string]string) []string {
	return MapToArrayWithSeparator(m, "=")
}
