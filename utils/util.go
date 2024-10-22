package utils

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
