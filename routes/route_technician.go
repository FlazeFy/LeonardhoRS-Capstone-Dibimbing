package routes

import (
	"pelita/controller"
	"pelita/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouteTechnician(api *gin.RouterGroup, technicianController *controller.TechnicianController, redisClient *redis.Client, db *gorm.DB) {
	// Admin Only
	protected_admin := api.Group("/")
	protected_admin.Use(middleware.AuthMiddleware(redisClient, "admin"))
	{
		technician := protected_admin.Group("/technicians")
		{
			technician.POST("/", technicianController.Create, middleware.AuditTrailMiddleware(db, "create_technician"))
			technician.PUT("/:id", technicianController.UpdateById, middleware.AuditTrailMiddleware(db, "update_technician_by_id"))
			technician.DELETE("/:id", technicianController.DeleteById, middleware.AuditTrailMiddleware(db, "delete_technician_by_id"))
		}
	}
	// Admin & Technician Only
	protected_admin_technician := api.Group("/")
	protected_admin_technician.Use(middleware.AuthMiddleware(redisClient, "admin", "technician"))
	{
		technician := protected_admin_technician.Group("/technicians")
		{
			technician.GET("/", technicianController.GetAllTechnician)
		}
	}
}
