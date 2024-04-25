package http_header

import (
	"slices"
	"strings"
)

func Match(headers map[string][]string, key string, value string) (m bool) {
	kk := strings.ToLower(key)
	vv := strings.ToLower(value)
	for k, list := range headers {
		if strings.ToLower(k) != kk {
			continue
		}
		m = slices.ContainsFunc(list, func(s string) bool {
			return strings.ToLower(s) == vv
		})
	}
	return
}

func FirstValue(headers map[string][]string, key string) string {
	kk := strings.ToLower(key)
	for k := range headers {
		if strings.ToLower(k) == kk {
			return strings.ToLower(headers[k][0])
		}
	}
	return ""
}
