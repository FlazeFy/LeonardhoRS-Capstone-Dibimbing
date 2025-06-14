package main

import (
	"fmt"
	"pelita/config"
	"pelita/entity"
	"pelita/repository"
	"pelita/routes"
	"pelita/scheduler"
	"pelita/service"
	"time"

	_ "pelita/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func main() {
	// Load Env
	err := godotenv.Load()
	if err != nil {
		panic("error loading ENV")
	}

	config.InitFirebase()

	// Connect DB
	db := config.ConnectDatabase()
	MigrateAll(db)

	// Setup Gin
	router := gin.Default()
	redisClient := config.InitRedis()
	routes.SetUpRoutes(router, db, redisClient)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Task Scheduler
	// Initialize Repositories
	assetMaintenanceRepo := repository.NewAssetMaintenanceRepository(db)
	assetFindingRepo := repository.NewAssetFindingRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	assetRepo := repository.NewAssetRepository(db)
	statsRepo := repository.NewStatsRepository(db)
	technicianRepo := repository.NewTechnicianRepository(db)

	// Initialize Services
	assetMaintenanceService := service.NewAssetMaintenanceService(assetMaintenanceRepo, technicianRepo, assetRepo, statsRepo)
	assetFindingService := service.NewAssetFindingService(assetFindingRepo, statsRepo)
	adminService := service.NewAdminService(adminRepo)

	// Initialize Scheduler
	maintenanceScheduler := scheduler.NewAssetMaintenanceScheduler(assetMaintenanceService, assetFindingService, adminService)

	// Init Scheduler
	c := cron.New()
	Scheduler(c, maintenanceScheduler)
	c.Start()
	defer c.Stop()

	// Run server
	router.Run(":9000")
}

func Scheduler(c *cron.Cron, maintenanceScheduler *scheduler.AssetMaintenanceScheduler) {
	// Production
	// Every day at 00:10 AM)
	c.AddFunc("10 0 * * *", maintenanceScheduler.ReminderSchedulerTodayMaintenance)
	// Every day at 00:20 AM)
	c.AddFunc("20 0 * * *", maintenanceScheduler.AuditSchedulerAssetFindingReport)

	// Development (after 5 sec)
	go func() {
		time.Sleep(5 * time.Second)
		maintenanceScheduler.ReminderSchedulerTodayMaintenance()
		maintenanceScheduler.AuditSchedulerAssetFindingReport()
	}()
}

func MigrateAll(db *gorm.DB) {
	err := db.AutoMigrate(
		&entity.User{},
		&entity.Admin{},
		&entity.Technician{},
		&entity.Room{},
		&entity.Asset{},
		&entity.AssetPlacement{},
		&entity.AssetMaintenance{},
		&entity.AssetFinding{},
		&entity.History{},
	)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Migrate Success!")
}
