package controller

import (
	"net/http"
	"pelita/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{UserService: userService}
}

func (ac *UserController) GetMyProfile(c *gin.Context) {
	// Get User Id in Middleware
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "something went wrong",
			"status":  "failed",
		})
		return
	}
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid user ID",
			"status":  "failed",
		})
		return
	}

	// Service: Get Profile by User ID
	user, err := ac.UserService.GetMyProfile(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "user fetched",
		"status":  "success",
		"data":    user,
	})
}
