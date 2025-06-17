package unit

import (
	"pelita/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentRole(t *testing.T) {
	t.Run("should return role successfully when role exists and is valid", func(t *testing.T) {
		c := &gin.Context{}
		c.Set("role", "admin")

		role, err := utils.GetCurrentRole(c)
		assert.NoError(t, err)
		assert.Equal(t, "admin", role)
	})

	t.Run("should return error when role not set in context", func(t *testing.T) {
		c := &gin.Context{}

		role, err := utils.GetCurrentRole(c)
		assert.Error(t, err)
		assert.Equal(t, "", role)
		assert.Equal(t, "role not found in context", err.Error())
	})

	t.Run("should return error when role is not a string", func(t *testing.T) {
		c := &gin.Context{}
		c.Set("role", 123)

		role, err := utils.GetCurrentRole(c)
		assert.Error(t, err)
		assert.Equal(t, "", role)
		assert.Equal(t, "invalid role format in context", err.Error())
	})
}

func TestGetCurrentUserID(t *testing.T) {
	t.Run("should return UUID when userID is valid", func(t *testing.T) {
		c := &gin.Context{}
		expectedUUID := uuid.New()
		c.Set("userID", expectedUUID.String())

		userID, err := utils.GetCurrentUserID(c)
		assert.NoError(t, err)
		assert.Equal(t, expectedUUID, userID)
	})

	t.Run("should return error when userID not in context", func(t *testing.T) {
		c := &gin.Context{}

		userID, err := utils.GetCurrentUserID(c)
		assert.Error(t, err)
		assert.Equal(t, uuid.UUID{}, userID)
		assert.Equal(t, "user id not found in context", err.Error())
	})

	t.Run("should return error when userID is not a valid UUID", func(t *testing.T) {
		c := &gin.Context{}
		c.Set("userID", "123-123")

		userID, err := utils.GetCurrentUserID(c)
		assert.Error(t, err)
		assert.Equal(t, uuid.UUID{}, userID)
		assert.Equal(t, "user id is not a valid UUID", err.Error())
	})
}
