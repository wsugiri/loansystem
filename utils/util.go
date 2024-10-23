package utils

import (
	"fmt"
	"strings"
)

func GetNestedValue(data map[string]interface{}, keys ...string) interface{} {
	var current interface{} = data

	for _, key := range keys {
		// Check if current is a map and contains the key
		if m, ok := current.(map[string]interface{}); ok {
			if val, exists := m[key]; exists {
				current = val
			} else {
				return nil // Key not found
			}
		} else {
			return nil // Not a map
		}
	}

	return current
}

// insertCommas formats an integer with thousand separators.
func InsertCommas(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}

	var result []string
	for i := len(s); i > 0; i -= 3 {
		start := max(0, i-3)
		result = append([]string{s[start:i]}, result...)
	}
	return strings.Join(result, ",")
}
