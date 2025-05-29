package controller

import (
	"net/http"
	"pelita/entity"
	"pelita/service"
	"pelita/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (rc *AssetController) UpdateById(c *gin.Context) {
	// Param
	id := c.Param("id")

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

	// Parse Id
	assetID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid UUID format",
			"status":  "failed",
		})
		return
	}

	// Service : Update Asset
	if err := rc.AssetService.UpdateById(&req, assetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "asset update successfully",
		"status":  "success",
	})
}

func (rc *AssetController) HardDeleteById(c *gin.Context) {
	// Param
	id := c.Param("id")

	// Parse Id
	assetID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid UUID format",
			"status":  "failed",
		})
		return
	}

	// Service : Hard Delete Asset By Id
	if err := rc.AssetService.HardDeleteById(assetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "asset permanentally deleted",
		"status":  "success",
	})
}
