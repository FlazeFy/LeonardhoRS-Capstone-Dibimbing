package routes

import (
	"pelita/controller"
	"pelita/repository"
	"pelita/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpDependency(r *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
	// Dependency Repositories
	statsRepo := repository.NewStatsRepository(db)
	userRepo := repository.NewUserRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	technicianRepo := repository.NewTechnicianRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	assetRepo := repository.NewAssetRepository(db)
	assetPlacementRepo := repository.NewAssetPlacementRepository(db)
	assetMaintenanceRepo := repository.NewAssetMaintenanceRepository(db)
	assetFindingRepo := repository.NewAssetFindingRepository(db)
	historyRepo := repository.NewHistoryRepository(db)

	// Dependency Services
	authService := service.NewAuthService(userRepo, adminRepo, technicianRepo, redisClient)
	technicianService := service.NewTechnicianService(technicianRepo)
	userService := service.NewUserService(userRepo, redisClient)
	roomService := service.NewRoomService(roomRepo, statsRepo)
	assetService := service.NewAssetService(assetRepo, statsRepo)
	assetPlacementService := service.NewAssetPlacementService(assetPlacementRepo)
	assetMaintenanceService := service.NewAssetMaintenanceService(assetMaintenanceRepo, technicianRepo, assetRepo, statsRepo)
	assetFindingService := service.NewAssetFindingService(assetFindingRepo, statsRepo)
	historyService := service.NewHistoryService(historyRepo, statsRepo)
	adminService := service.NewAdminService(adminRepo)

	// Dependency Controllers
	authController := controller.NewAuthController(authService)
	technicianController := controller.NewTechnicianController(technicianService)
	userController := controller.NewUserController(userService)
	roomController := controller.NewRoomRepository(roomService)
	assetController := controller.NewAssetRepository(assetService)
	assetPlacementController := controller.NewAssetPlacementRepository(assetPlacementService)
	assetMaintenanceController := controller.NewAssetMaintenanceRepository(assetMaintenanceService)
	assetFindingController := controller.NewAssetFindingRepository(assetFindingService)
	historyController := controller.NewHistoryRepository(historyService)

	// Routes Endpoint
	SetUpRoutes(r, db, redisClient,
		authController,
		technicianController,
		userController,
		roomController,
		assetController,
		assetPlacementController,
		assetMaintenanceController,
		assetFindingController,
		historyController,
	)

	// Task Scheduler
	SetUpScheduler(assetMaintenanceService, assetFindingService, adminService)

	// Seeder & Factories
	SetUpSeeder(db, roomRepo, adminRepo, technicianRepo, userRepo, assetRepo)
}
