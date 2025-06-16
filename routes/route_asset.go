package routes

import (
	"pelita/controller"
	"pelita/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteAsset(api *gin.RouterGroup, assetController *controller.AssetController, assetFindingController *controller.AssetFindingController, assetMaintenanceController *controller.AssetMaintenanceController, assetPlacementController *controller.AssetPlacementController, redisClient *redis.Client, db *gorm.DB) {
	// Admin Only
	protected_admin := api.Group("/")
	protected_admin.Use(middleware.AuthMiddleware(redisClient, "admin"))
	{
		asset := protected_admin.Group("/asset")
		{
			asset.POST("/", assetController.Create, middleware.AuditTrailMiddleware(db, "create_asset"))
			asset.GET("/most_context/:target_col", assetController.GetMostContext)
			asset.GET("/", assetController.GetAllAsset)
			asset.DELETE("/destroy/:id", assetController.HardDeleteById, middleware.AuditTrailMiddleware(db, "hard_delete_asset_by_id"))
			asset.DELETE("/:id", assetController.SoftDeleteById, middleware.AuditTrailMiddleware(db, "soft_delete_asset_by_id"))
			asset.PUT("/:id", assetController.UpdateById, middleware.AuditTrailMiddleware(db, "update_asset_by_id"))
			asset.PUT("/recover/:id", assetController.RecoverDeletedById, middleware.AuditTrailMiddleware(db, "recover_delete_asset_by_id"))

			asset_placement := asset.Group("/placement")
			{
				asset_placement.POST("/", assetPlacementController.Create)
				asset_placement.DELETE("/:id", assetPlacementController.DeleteById)
			}
			asset_maintenance := asset.Group("/maintenance")
			{
				asset_maintenance.GET("/most_context/:target_col", assetMaintenanceController.GetMostContext)
				asset_maintenance.POST("/", assetMaintenanceController.Create, middleware.AuditTrailMiddleware(db, "create_asset_maintenance_by_id"))
				asset_maintenance.PUT("/:id", assetMaintenanceController.UpdateById, middleware.AuditTrailMiddleware(db, "update_asset_maintenance_by_id"))
				asset_maintenance.DELETE("/:id", assetMaintenanceController.DeleteById, middleware.AuditTrailMiddleware(db, "delete_asset_maintenance_by_id"))
			}
			asset_finding := asset.Group("/finding")
			{
				asset_finding.DELETE("/:id", assetFindingController.DeleteById, middleware.AuditTrailMiddleware(db, "delete_asset_finding_by_id"))
				asset_finding.GET("/most_context/:target_col", assetFindingController.GetMostContext)
				asset_finding.GET("/hour_total", assetFindingController.GetFindingHourTotal)
			}
		}
	}

	// Admin & Technician Only
	protected_admin_technician := api.Group("/")
	protected_admin_technician.Use(middleware.AuthMiddleware(redisClient, "admin", "technician"))
	{
		asset := protected_admin_technician.Group("/asset")
		{
			asset.GET("/deleted", assetController.GetDeletedAsset)

			asset_placement := asset.Group("/placement")
			{
				asset_placement.GET("/", assetPlacementController.GetAllAssetPlacement)
				asset_placement.PUT("/:id", assetPlacementController.UpdateById, middleware.AuditTrailMiddleware(db, "update_asset_placement_by_id"))
			}
			asset_maintenance := asset.Group("/maintenance")
			{
				asset_maintenance.GET("/", assetMaintenanceController.GetAllAssetMaintenance)
				asset_maintenance.GET("/schedule", assetMaintenanceController.GetAllAssetMaintenanceSchedule)
			}
			asset_finding := asset.Group("/finding")
			{
				asset_finding.GET("/", assetFindingController.GetAllAssetFinding)
			}
		}
	}

	// User / Guest & Technician Only
	protected_user_technician := api.Group("/")
	protected_user_technician.Use(middleware.AuthMiddleware(redisClient, "guest", "technician"))
	{
		asset := protected_user_technician.Group("/asset")
		{
			asset_finding := asset.Group("/finding")
			{
				asset_finding.POST("/", assetFindingController.Create, middleware.AuditTrailMiddleware(db, "create_asset_finding"))
			}
		}
	}
}
