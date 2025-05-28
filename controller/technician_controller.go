package controller

import (
	"net/http"
	"pelita/entity"
	"pelita/service"
	"pelita/utils"

	"github.com/gin-gonic/gin"
)

type TechnicianController struct {
	TechnicianService service.TechnicianService
}

func NewTechnicianController(technicianService service.TechnicianService) *TechnicianController {
	return &TechnicianController{
		TechnicianService: technicianService,
	}
}

func (rc *TechnicianController) GetAllTechnician(c *gin.Context) {
	// Service: Get All Technician
	technician, err := rc.TechnicianService.GetAllTechnician()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "technician fetched",
		"status":  "success",
		"data":    technician,
	})
}

func (rc *TechnicianController) Create(c *gin.Context) {
	// Model
	var req entity.Technician

	// Validator
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Get User Id
	adminId, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Service : Create Technician
	if err := rc.TechnicianService.Create(&req, adminId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "technician created successfully",
		"status":  "success",
		"data":    &req,
	})
}
