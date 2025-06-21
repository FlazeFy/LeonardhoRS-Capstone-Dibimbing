package main

import (
	"fmt"
	"log"
	"os"
	"pelita/config"
	"pelita/entity"
	"pelita/routes"

	_ "pelita/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func initLogging() {
	f, err := os.OpenFile("pelita.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	initLogging()
	log.Println("Pelita service is starting...")

	// Load Env
	err := godotenv.Load()
	if err != nil {
		panic("error loading ENV")
	}

	config.InitFirebase()

	// Connect DB
	db := config.ConnectDatabase()
	MigrateAll(db)

	// Setup Gin & Redis
	router := gin.Default()
	redisClient := config.InitRedis()

	// Setup Dependecy, Scheduler, and Seeder
	routes.SetUpDependency(router, db, redisClient)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run server
	port := os.Getenv("PORT")
	router.Run(":" + port)

	log.Printf("Pelita is running on port %s\n", port)
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
