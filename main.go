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

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
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

	// Task Scheduler
	// Initialize Repositories
	assetMaintenanceRepo := repository.NewAssetMaintenanceRepository(db)
	adminRepo := repository.NewAdminRepository(db)

	// Initialize Services
	assetMaintenanceService := service.NewAssetMaintenanceService(assetMaintenanceRepo)
	adminService := service.NewAdminService(adminRepo)

	// Initialize Scheduler
	maintenanceScheduler := scheduler.NewAssetMaintenanceScheduler(assetMaintenanceService, adminService)

	// Init Scheduler
	c := cron.New()
	Scheduler(c, maintenanceScheduler)
	c.Start()
	defer c.Stop()

	// Run server
	router.Run(":9000")
}

func Scheduler(c *cron.Cron, maintenanceScheduler *scheduler.AssetMaintenanceScheduler) {
	// Production (e.g., every day at 00:10 AM)
	c.AddFunc("10 0 * * *", maintenanceScheduler.ReminderSchedulerTodayMaintenance)

	// Development (after 5 sec)
	go func() {
		time.Sleep(5 * time.Second)
		maintenanceScheduler.ReminderSchedulerTodayMaintenance()
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
	)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Migrate Success!")
}
