package controller

import (
	"math"
	"net/http"
	"pelita/entity"
	"pelita/service"
	"pelita/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssetMaintenanceController struct {
	AssetMaintenanceService service.AssetMaintenanceService
}

func NewAssetMaintenanceRepository(assetMaintenanceService service.AssetMaintenanceService) *AssetMaintenanceController {
	return &AssetMaintenanceController{AssetMaintenanceService: assetMaintenanceService}
}

func (rc *AssetMaintenanceController) GetAllAssetMaintenance(c *gin.Context) {
	// Pagination
	pagination := utils.GetPagination(c)

	// Service: Get All Asset Maintenance
	assetMaintenance, total, err := rc.AssetMaintenanceService.GetAllAssetMaintenance(pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))
	c.JSON(http.StatusOK, gin.H{
		"message": "asset maintenance fetched",
		"status":  "success",
		"data":    assetMaintenance,
		"metadata": gin.H{
			"total":       total,
			"page":        pagination.Page,
			"limit":       pagination.Limit,
			"total_pages": totalPages,
		},
	})
}

func (rc *AssetMaintenanceController) GetAllAssetMaintenanceSchedule(c *gin.Context) {
	// Service: Get All Asset Maintenance
	assetMaintenance, err := rc.AssetMaintenanceService.GetAllAssetMaintenanceSchedule()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "asset maintenance schedule fetched",
		"status":  "success",
		"data":    assetMaintenance,
	})
}

func (rc *AssetMaintenanceController) Create(c *gin.Context) {
	// Model
	var req entity.AssetMaintenance

	// Validator
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Validator Rules
	validDays := map[string]bool{"Sun": true, "Mon": true, "Tue": true, "Wed": true, "Thu": true, "Fri": true, "Sat": true}
	if !validDays[req.MaintenanceDay] {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "maintenance day must be one of: Sun, Mon, Tue, Wed, Thu, Fri, Sat",
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

	// Service : Create Asset Maintenance
	if err := rc.AssetMaintenanceService.Create(&req, adminId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "asset maintenance created successfully",
		"status":  "success",
		"data":    &req,
	})
}

func (rc *AssetMaintenanceController) UpdateById(c *gin.Context) {
	// Param
	id := c.Param("id")

	// Model
	var req entity.AssetMaintenance

	// Validator
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Parse Id
	assetMaintenanceID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid UUID format",
			"status":  "failed",
		})
		return
	}

	// Service : Update Asset Maintenance
	if err := rc.AssetMaintenanceService.UpdateById(&req, assetMaintenanceID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "asset maintenance update successfully",
		"status":  "success",
	})
}

func (rc *AssetMaintenanceController) DeleteById(c *gin.Context) {
	// Param
	id := c.Param("id")

	// Parse Id
	assetMaintenanceID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid UUID format",
			"status":  "failed",
		})
		return
	}

	// Service : Delete Asset Maintenance By Id
	if err := rc.AssetMaintenanceService.DeleteById(assetMaintenanceID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "asset maintenance deleted",
		"status":  "success",
	})
}

func (rc *AssetMaintenanceController) GetMostContext(c *gin.Context) {
	// Param
	targetCol := c.Param("target_col")

	// Validator : Target Column Validator
	if targetCol != "maintenance_day" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "target_col is not valid",
			"status":  "failed",
		})
		return
	}

	// Service: Get Most Context
	assetMaintenance, err := rc.AssetMaintenanceService.GetMostContext(targetCol)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "asset maintenance fetched",
		"status":  "success",
		"data":    assetMaintenance,
	})
}
