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

type AssetPlacementController struct {
	AssetPlacementService service.AssetPlacementService
}

func NewAssetPlacementRepository(assetPlacementService service.AssetPlacementService) *AssetPlacementController {
	return &AssetPlacementController{AssetPlacementService: assetPlacementService}
}

func (rc *AssetPlacementController) GetAllAssetPlacement(c *gin.Context) {
	// Pagination
	pagination := utils.GetPagination(c)

	// Service: Get All Asset Placement
	assetPlacement, total, err := rc.AssetPlacementService.GetAllAssetPlacement(pagination)
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
		"message": "asset placement fetched",
		"status":  "success",
		"data":    assetPlacement,
		"metadata": gin.H{
			"total":       total,
			"page":        pagination.Page,
			"limit":       pagination.Limit,
			"total_pages": totalPages,
		},
	})
}

func (rc *AssetPlacementController) Create(c *gin.Context) {
	// Model
	var req entity.AssetPlacement

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

	// Service : Create Asset Placement
	if err := rc.AssetPlacementService.Create(&req, adminId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "asset placement created successfully",
		"status":  "success",
		"data":    &req,
	})
}

func (rc *AssetPlacementController) UpdateById(c *gin.Context) {
	// Param
	id := c.Param("id")

	// Model
	var req entity.AssetPlacement

	// Validator
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Parse Id
	assetPlacementID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid UUID format",
			"status":  "failed",
		})
		return
	}

	// Service : Update Asset Placement
	if err := rc.AssetPlacementService.UpdateById(&req, assetPlacementID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "asset placement update successfully",
		"status":  "success",
	})
}

func (rc *AssetPlacementController) DeleteById(c *gin.Context) {
	// Param
	id := c.Param("id")

	// Parse Id
	assetPlacementID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid UUID format",
			"status":  "failed",
		})
		return
	}

	// Service : Delete Asset Placement By Id
	if err := rc.AssetPlacementService.DeleteById(assetPlacementID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "asset placement deleted",
		"status":  "success",
	})
}
