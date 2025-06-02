package utils

func OptionalString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
