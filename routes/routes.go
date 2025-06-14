package routes

import (
	"pelita/controller"
	"pelita/middleware"
	"pelita/repository"
	"pelita/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRoutes(r *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
	// Stats Module
	statsRepo := repository.NewStatsRepository(db)

	// Auth Module
	userRepo := repository.NewUserRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	technicianRepo := repository.NewTechnicianRepository(db)
	authService := service.NewAuthService(userRepo, adminRepo, technicianRepo, redisClient)
	authController := controller.NewAuthController(authService)

	// Technician Module
	technicianService := service.NewTechnicianService(technicianRepo)
	technicianController := controller.NewTechnicianController(technicianService)

	// User Module
	userService := service.NewUserService(userRepo, redisClient)
	userController := controller.NewUserController(userService)

	// Room Module
	roomRepo := repository.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo, statsRepo)
	roomController := controller.NewRoomRepository(roomService)

	// Asset Module
	assetRepo := repository.NewAssetRepository(db)
	assetService := service.NewAssetService(assetRepo, statsRepo)
	assetController := controller.NewAssetRepository(assetService)

	// Asset Placement Module
	assetPlacementRepo := repository.NewAssetPlacementRepository(db)
	assetPlacementService := service.NewAssetPlacementService(assetPlacementRepo)
	assetPlacementController := controller.NewAssetPlacementRepository(assetPlacementService)

	// Asset Maintenance Module
	assetMaintenanceRepo := repository.NewAssetMaintenanceRepository(db)
	assetMaintenanceService := service.NewAssetMaintenanceService(assetMaintenanceRepo, technicianRepo, assetRepo, statsRepo)
	assetMaintenanceController := controller.NewAssetMaintenanceRepository(assetMaintenanceService)

	// Asset Finding Module
	assetFindingRepo := repository.NewAssetFindingRepository(db)
	assetFindingService := service.NewAssetFindingService(assetFindingRepo, statsRepo)
	assetFindingController := controller.NewAssetFindingRepository(assetFindingService)

	// History Module
	historyRepo := repository.NewHistoryRepository(db)
	historyService := service.NewHistoryService(historyRepo, statsRepo)
	historyController := controller.NewHistoryRepository(historyService)

	api := r.Group("/api/v1")
	{
		// Public Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/signout", authController.SignOut)
		}
	}

	// All Role
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(redisClient, "admin", "technician", "guest"))
	{
		protected.GET("/profile", userController.GetMyProfile)

		room := protected.Group("/room")
		{
			room.GET("/", roomController.GetAllRoom)
		}
		history := protected.Group("/history")
		{
			history.GET("/my", historyController.GetMyHistory)
		}
	}

	// Admin Only
	protected_admin := api.Group("/")
	protected_admin.Use(middleware.AuthMiddleware(redisClient, "admin"))
	{
		room := protected_admin.Group("/room")
		{
			room.GET("/most_context/:target_col", roomController.GetMostContext)
			room.POST("/", roomController.Create, middleware.AuditTrailMiddleware(db, "create_room"))
			room.DELETE("/:id", roomController.DeleteById, middleware.AuditTrailMiddleware(db, "delete_room_by_id"))
			room.PUT("/:id", roomController.UpdateById, middleware.AuditTrailMiddleware(db, "update_room_by_id"))
		}
		history := protected_admin.Group("/history")
		{
			history.GET("/all", historyController.GetAllHistory)
			history.GET("/:target_col", historyController.GetMostContext)
		}
		technician := protected_admin.Group("/technician")
		{
			technician.POST("/", technicianController.Create, middleware.AuditTrailMiddleware(db, "create_technician"))
			technician.PUT("/:id", technicianController.UpdateById, middleware.AuditTrailMiddleware(db, "update_technician_by_id"))
			technician.DELETE("/:id", technicianController.DeleteById, middleware.AuditTrailMiddleware(db, "delete_technician_by_id"))
		}
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
		technician := protected_admin_technician.Group("/technician")
		{
			technician.GET("/", technicianController.GetAllTechnician)
		}
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
		room := protected_admin_technician.Group("/room/asset")
		{
			room.GET("/detail/:floor/:room_name", roomController.GetRoomAssetByFloorAndRoomName)
			room.GET("/short/:floor/:room_name", roomController.GetRoomAssetShortByFloorAndRoomName)
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
