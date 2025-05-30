package main

import (
	"fmt"
	"pelita/config"
	"pelita/entity"
	"pelita/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// Load Env
	err := godotenv.Load()
	if err != nil {
		panic("error loading ENV")
	}

	// Connect DB
	db := config.ConnectDatabase()
	MigrateAll(db)

	// Setup Gin
	router := gin.Default()
	redisClient := config.InitRedis()
	routes.SetUpRoutes(router, db, redisClient)

	// Run server
	router.Run(":9000")
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
	)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Migrate Success!")
}
