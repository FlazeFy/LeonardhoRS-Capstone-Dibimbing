package routes

import (
	"pelita/controller"
	"pelita/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetUpRouteHistory(api *gin.RouterGroup, historyController *controller.HistoryController, redisClient *redis.Client) {
	// All Role
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(redisClient, "admin", "technician", "guest"))
	{
		history := protected.Group("/histories")
		{
			history.GET("/my", historyController.GetMyHistory)
		}
	}
	// Admin Only
	protected_admin := api.Group("/")
	protected_admin.Use(middleware.AuthMiddleware(redisClient, "admin"))
	{
		history := protected_admin.Group("/histories")
		{
			history.GET("/all", historyController.GetAllHistory)
			history.GET("/:targetCol", historyController.GetMostContext)
		}
	}
}
