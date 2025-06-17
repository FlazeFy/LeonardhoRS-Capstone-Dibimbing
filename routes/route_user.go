package routes

import (
	"pelita/controller"
	"pelita/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetUpRouteUser(api *gin.RouterGroup, userController *controller.UserController, redisClient *redis.Client) {
	// All Role
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(redisClient, "admin", "technician", "guest"))
	{
		protected.GET("/profile", userController.GetMyProfile)
	}
}
