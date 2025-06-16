package routes

import (
	"pelita/controller"

	"github.com/gin-gonic/gin"
)

func SetUpRouteAuth(api *gin.RouterGroup, authController *controller.AuthController) {
	// Public Routes
	auth := api.Group("/auths")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/signout", authController.SignOut)
	}
}
