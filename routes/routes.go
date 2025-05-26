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
	authService := service.NewAuthService(userRepo, redisClient)
	authController := controller.NewAuthController(authService)

	// User Module
	userService := service.NewUserService(userRepo, redisClient)
	userController := controller.NewUserController(userService)

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
	protected.Use(middleware.AuthMiddleware(redisClient))
	{
		protected.GET("/profile", userController.GetMyProfile)
	}
}
