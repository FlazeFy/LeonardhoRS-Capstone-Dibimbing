package routes

import (
	"pelita/controller"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRoutes(r *gin.Engine, db *gorm.DB, redisClient *redis.Client,
	authController *controller.AuthController,
	technicianController *controller.TechnicianController,
	userController *controller.UserController,
	roomController *controller.RoomController,
	assetController *controller.AssetController,
	assetPlacementController *controller.AssetPlacementController,
	assetMaintenanceController *controller.AssetMaintenanceController,
	assetFindingController *controller.AssetFindingController,
	historyController *controller.HistoryController) {

	// V1 Endpoint
	api := r.Group("/api/v1")

	// Routes Endpoint
	SetUpRouteAuth(api, authController)
	SetUpRouteUser(api, userController, redisClient)
	SetUpRouteTechnician(api, technicianController, redisClient, db)
	SetUpRouteRoom(api, roomController, redisClient, db)
	SetUpRouteAsset(api, assetController, assetFindingController, assetMaintenanceController, assetPlacementController, redisClient, db)
	SetUpRouteHistory(api, historyController, redisClient)
}
