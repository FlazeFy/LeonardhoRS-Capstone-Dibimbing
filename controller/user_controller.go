package controller

import (
	"net/http"
	"pelita/service"
	"pelita/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{UserService: userService}
}

// @Summary      Get My Profile
// @Description  Returns a paginated list of handle
// @Tags         Use
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetMyProfile
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/profile [get]
func (ac *UserController) GetMyProfile(c *gin.Context) {
	// Get User Id
	userID, err := utils.GetCurrentUserID(c)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Get Role
	role, err := utils.GetCurrentRole(c)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Service: Get Profile by User ID
	user, err := ac.UserService.GetMyProfile(userID, role)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "user", "get", http.StatusOK, user, nil)
}
