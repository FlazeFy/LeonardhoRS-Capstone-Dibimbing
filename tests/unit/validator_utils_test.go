package unit

import (
	"pelita/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	// Test 1: Should return true if item found in list
	targetColTest := "B"
	validTargetTest := []string{"A", "B", "C"}

	validTest := utils.Contains(validTargetTest, targetColTest)
	assert.Equal(t, true, validTest, "should return true")

	// Test 2: Should return false if item not found in list
	targetColTest2 := "D"
	validTargetTest2 := []string{"A", "B", "C"}

	validTest2 := utils.Contains(validTargetTest2, targetColTest2)
	assert.Equal(t, false, validTest2, "should return true")
}
