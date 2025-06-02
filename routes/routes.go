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
	roomService := service.NewRoomService(roomRepo)
	roomController := controller.NewRoomRepository(roomService)

	// Asset Module
	assetRepo := repository.NewAssetRepository(db)
	assetService := service.NewAssetService(assetRepo)
	assetController := controller.NewAssetRepository(assetService)

	// Asset Placement Module
	assetPlacementRepo := repository.NewAssetPlacementRepository(db)
	assetPlacementService := service.NewAssetPlacementService(assetPlacementRepo)
	assetPlacementController := controller.NewAssetPlacementRepository(assetPlacementService)

	// Asset Maintenance Module
	assetMaintenanceRepo := repository.NewAssetMaintenanceRepository(db)
	assetMaintenanceService := service.NewAssetMaintenanceService(assetMaintenanceRepo)
	assetMaintenanceController := controller.NewAssetMaintenanceRepository(assetMaintenanceService)

	// Asset Finding Module
	assetFindingRepo := repository.NewAssetFindingRepository(db)
	assetFindingService := service.NewAssetFindingService(assetFindingRepo)
	assetFindingController := controller.NewAssetFindingRepository(assetFindingService)

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
	}

	// Admin Only
	protected_admin := api.Group("/")
	protected_admin.Use(middleware.AuthMiddleware(redisClient, "admin"))
	{
		room := protected_admin.Group("/room")
		{
			room.POST("/", roomController.Create)
			room.DELETE("/:id", roomController.DeleteById)
			room.PUT("/:id", roomController.UpdateById)
		}
		technician := protected_admin.Group("/technician")
		{
			technician.POST("/", technicianController.Create)
			technician.PUT("/:id", technicianController.UpdateById)
			technician.DELETE("/:id", technicianController.DeleteById)
		}
		asset := protected_admin.Group("/asset")
		{
			asset.POST("/", assetController.Create)
			asset.GET("/", assetController.GetAllAsset)
			asset.DELETE("/destroy/:id", assetController.HardDeleteById)
			asset.DELETE("/:id", assetController.SoftDeleteById)
			asset.PUT("/:id", assetController.UpdateById)
			asset.PUT("/recover/:id", assetController.RecoverDeletedById)

			asset_placement := asset.Group("/placement")
			{
				asset_placement.POST("/", assetPlacementController.Create)
				asset_placement.DELETE("/:id", assetPlacementController.DeleteById)
			}
			asset_maintenance := asset.Group("/maintenance")
			{
				asset_maintenance.POST("/", assetMaintenanceController.Create)
				asset_maintenance.PUT("/:id", assetMaintenanceController.UpdateById)
				asset_maintenance.DELETE("/:id", assetMaintenanceController.DeleteById)
			}
			asset_finding := asset.Group("/finding")
			{
				asset_finding.DELETE("/:id", assetFindingController.DeleteById)
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
				asset_placement.PUT("/:id", assetPlacementController.UpdateById)
			}
			asset_maintenance := asset.Group("/maintenance")
			{
				asset_maintenance.GET("/", assetMaintenanceController.GetAllAssetMaintenance)
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
				asset_finding.POST("/", assetFindingController.Create)
			}
		}
	}
}
