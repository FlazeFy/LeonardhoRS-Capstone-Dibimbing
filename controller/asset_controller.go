package controller

import (
	"net/http"
	"pelita/entity"
	"pelita/service"
	"pelita/utils"

	"github.com/gin-gonic/gin"
)

type AssetController struct {
	AssetService service.AssetService
}

func NewAssetRepository(assetService service.AssetService) *AssetController {
	return &AssetController{AssetService: assetService}
}

func (rc *AssetController) GetAllAsset(c *gin.Context) {
	// Service: Get All Asset
	asset, err := rc.AssetService.GetAllAsset()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "asset fetched",
		"status":  "success",
		"data":    asset,
	})
}

func (rc *AssetController) Create(c *gin.Context) {
	// Model
	var req entity.Asset

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

	// Service : Create Asset
	if err := rc.AssetService.Create(&req, adminId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "asset created successfully",
		"status":  "success",
		"data":    &req,
	})
}
