package utils

import (
	"encoding/json"
	"unicode"
)

func OptionalString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func NullSafeString(val *string) string {
	if val != nil {
		return *val
	}
	return "-"
}

func CleanResponse(data interface{}, keysToRemove ...string) map[string]interface{} {
	jsonData, _ := json.Marshal(data)

	var result map[string]interface{}
	json.Unmarshal(jsonData, &result)

	for _, key := range keysToRemove {
		delete(result, key)
	}

	return result
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
