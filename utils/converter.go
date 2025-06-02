package utils

import (
	"encoding/json"
)

func OptionalString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
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
