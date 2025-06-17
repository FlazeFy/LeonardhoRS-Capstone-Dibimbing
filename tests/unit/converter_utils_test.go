package unit

import (
	"pelita/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionalString(t *testing.T) {
	// Test 1: Should return pointer when string is not empty
	str := "hello"
	ptr := utils.OptionalString(str)
	assert.Equal(t, "hello", *ptr, "pointer value should match input")

	// Test 2: Should return nil when input is an empty string
	emptyPtr := utils.OptionalString("")
	assert.Nil(t, emptyPtr, "should return nil when string is empty")
}

func TestNullSafeString(t *testing.T) {
	// Test 1: Should return actual value when pointer is not nil
	value := "test"
	ptr := &value
	result := utils.NullSafeString(ptr)
	assert.Equal(t, "test", result, "should return value from pointer")

	// Test 2: Should return hypen when pointer is nil
	nilPtr := (*string)(nil)
	result = utils.NullSafeString(nilPtr)
	assert.Equal(t, "-", result, "should return fallback dash when pointer is nil")
}

func TestCleanResponse(t *testing.T) {
	type Dummy struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	d := Dummy{
		ID:    1,
		Name:  "tester123",
		Email: "tester@gmail.com",
	}

	// Test 1 : Email should not exist in object after removed
	result := utils.CleanResponse(d, "email")
	assert.Equal(t, float64(1), result["id"], "id should be present")
	assert.Equal(t, "tester123", result["name"], "name should be present")
	_, exists := result["email"]
	assert.False(t, exists, "email should be removed from result")
}

func TestCapitalize(t *testing.T) {
	// Test 1 : Only Uppercase first letter in first word
	input := "hello world"
	output := utils.Capitalize(input)
	assert.Equal(t, "Hello world", output, "first letter should be capitalized")

	// Test 2 : Should returning empty string when capitalize empty string
	empty := utils.Capitalize("")
	assert.Equal(t, "", empty, "empty string should return empty string")
}
