package controller

import (
	"net/http"
	"pelita/entity"
	"pelita/service"
	"pelita/utils"

	"github.com/gin-gonic/gin"
)

type AssetPlacementController struct {
	AssetPlacementService service.AssetPlacementService
}

func NewAssetPlacementRepository(assetPlacementService service.AssetPlacementService) *AssetPlacementController {
	return &AssetPlacementController{AssetPlacementService: assetPlacementService}
}

func (rc *AssetPlacementController) GetAllAssetPlacement(c *gin.Context) {
	// Service: Get All AssetPlacement
	assetPlacement, err := rc.AssetPlacementService.GetAllAssetPlacement()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "asset placement fetched",
		"status":  "success",
		"data":    assetPlacement,
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
