package routes

import (
	"pelita/controller"
	"pelita/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteRoom(api *gin.RouterGroup, roomController *controller.RoomController, redisClient *redis.Client, db *gorm.DB) {
	// All Role
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(redisClient, "admin", "technician", "guest"))
	{
		room := protected.Group("/rooms")
		{
			room.GET("/", roomController.GetAllRoom)
		}
	}
	// Admin Only
	protected_admin := api.Group("/")
	protected_admin.Use(middleware.AuthMiddleware(redisClient, "admin"))
	{
		room := protected_admin.Group("/rooms")
		{
			room.GET("/most-context/:targetCol", roomController.GetMostContext)
			room.POST("/", roomController.Create, middleware.AuditTrailMiddleware(db, "create_room"))
			room.DELETE("/:id", roomController.DeleteById, middleware.AuditTrailMiddleware(db, "delete_room_by_id"))
			room.PUT("/:id", roomController.UpdateById, middleware.AuditTrailMiddleware(db, "update_room_by_id"))
		}
	}
	// Admin & Technician Only
	protected_admin_technician := api.Group("/")
	protected_admin_technician.Use(middleware.AuthMiddleware(redisClient, "admin", "technician"))
	{
		room := protected_admin_technician.Group("/rooms/assets")
		{
			room.GET("/detail/:floor/:roomName", roomController.GetRoomAssetByFloorAndRoomName)
			room.GET("/short/:floor/:roomName", roomController.GetRoomAssetShortByFloorAndRoomName)
		}
	}
}
