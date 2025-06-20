package controller

import (
	"net/http"
	"pelita/entity"
	"pelita/service"
	"pelita/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

// @Summary      Post Register
// @Description  Register User / Guest
// @Tags         Auth
// @Accept       application/json
// @Produce      json
// @Param        request  body  entity.RequestPostRegister true  "Post Register Request Body"
// @Success      201  {object}  entity.ResponsePostLogin
// @Failure      400  {object}  entity.ResponseBadRequest
// @Router       /api/v1/auths/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	// Model
	var req entity.User

	// Validator
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Register Token
	token, err := ac.AuthService.Register(&req)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "user", "register", http.StatusCreated, gin.H{
		"data": gin.H{
			"access_token": token,
		},
	}, nil)
}

// @Summary      Post Login
// @Description  Login to the Apps
// @Tags         Auth
// @Accept       application/json
// @Produce      json
// @Param        request  body  entity.UserAuth true  "Post Login Request Body"
// @Success      200  {object}  entity.ResponsePostLogin
// @Failure      400  {object}  entity.ResponseBadRequest
// @Router       /api/v1/auths/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	// Model
	var req entity.UserAuth

	// Validator
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Token Generate
	token, role, err := ac.AuthService.Login(req.Email, req.Password)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "user", "login", http.StatusOK, gin.H{
		"role":         role,
		"access_token": token,
	}, nil)
}

// @Summary      Post Sign Out
// @Description  Sign Out to the Apps
// @Tags         Auth
// @Accept       application/json
// @Produce      json
// @Success      200  {object}  entity.ResponsePostSignOut
// @Failure      400  {object}  entity.ResponseBadRequest
// @Router       /api/v1/auths/signout [post]
func (ac *AuthController) SignOut(c *gin.Context) {
	// Header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.BuildErrorMessage(c, http.StatusUnauthorized, "missing authorization header")
		return
	}

	// Clean Bearer
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)
	if token == "" {
		utils.BuildErrorMessage(c, http.StatusUnauthorized, "invalid authorization header")
		return
	}

	// Reset Token By Adding Blacklist Redis
	err := ac.AuthService.SignOut(token)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "user", "sign out", http.StatusOK, nil, nil)
}
