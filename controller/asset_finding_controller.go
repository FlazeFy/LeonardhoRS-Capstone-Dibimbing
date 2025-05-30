package controller

import (
	"net/http"
	"pelita/entity"
	"pelita/service"
	"pelita/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssetFindingController struct {
	AssetFindingService service.AssetFindingService
}

func NewAssetFindingRepository(assetFindingService service.AssetFindingService) *AssetFindingController {
	return &AssetFindingController{AssetFindingService: assetFindingService}
}

func (rc *AssetFindingController) GetAllAssetFinding(c *gin.Context) {
	// Service: Get All Asset Finding
	assetFinding, err := rc.AssetFindingService.GetAllAssetFinding()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "asset finding fetched",
		"status":  "success",
		"data":    assetFinding,
	})
}

func (rc *AssetFindingController) Create(c *gin.Context) {
	// Model
	var req entity.AssetFinding

	// Validator
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Validator Rules
	validDays := map[string]bool{"Missing": true, "Broken": true, "Empty": true, "Dirty": true}
	if !validDays[req.FindingCategory] {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "finding category must be one of: Missing, Broken, Empty, Dirty",
			"status":  "failed",
		})
		return
	}

	// Get User Id / Technician Id
	technicianOrUserId, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Get Role
	role, err := utils.GetCurrentRole(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Define The Role Id
	var technicianId uuid.NullUUID
	var userId uuid.NullUUID
	if role == "technician" {
		technicianId = uuid.NullUUID{UUID: technicianOrUserId, Valid: true}
	} else {
		userId = uuid.NullUUID{UUID: technicianOrUserId, Valid: true}
	}

	// Service : Create Asset Finding
	if err := rc.AssetFindingService.Create(&req, technicianId, userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "asset finding created successfully",
		"status":  "success",
		"data":    &req,
	})
}

func (rc *AssetFindingController) DeleteById(c *gin.Context) {
	// Param
	id := c.Param("id")

	// Parse Id
	assetFindingID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid UUID format",
			"status":  "failed",
		})
		return
	}

	// Service : Delete Asset Finding By Id
	if err := rc.AssetFindingService.DeleteById(assetFindingID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "asset finding deleted",
		"status":  "success",
	})
}
