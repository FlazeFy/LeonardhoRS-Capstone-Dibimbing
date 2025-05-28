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

	// User Module
	userService := service.NewUserService(userRepo, redisClient)
	userController := controller.NewUserController(userService)

	// Room Module
	roomRepo := repository.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo)
	roomController := controller.NewRoomRepository(roomService)

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

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(redisClient, "admin", "technician", "guest"))
	{
		protected.GET("/profile", userController.GetMyProfile)

		room := protected.Group("/room")
		{
			room.GET("/", roomController.GetAllRoom)
		}
	}

	protected_admin := api.Group("/")
	protected_admin.Use(middleware.AuthMiddleware(redisClient, "admin"))
	{
		room := protected_admin.Group("/room")
		{
			room.POST("/", roomController.Create)
			room.DELETE("/:id", roomController.DeleteById)
			room.PUT("/:id", roomController.UpdateById)
		}
	}
}
