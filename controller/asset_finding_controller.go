package controller

import (
	"fmt"
	"math"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"pelita/entity"
	"pelita/service"
	"pelita/utils"
	"strings"

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
	// Pagination
	pagination := utils.GetPagination(c)

	// Service: Get All Asset Finding
	assetFinding, total, err := rc.AssetFindingService.GetAllAssetFinding(pagination)
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
		"message": "asset finding fetched",
		"status":  "success",
		"data":    assetFinding,
		"metadata": gin.H{
			"total":       total,
			"page":        pagination.Page,
			"limit":       pagination.Limit,
			"total_pages": totalPages,
		},
	})
}

func (rc *AssetFindingController) GetMostContext(c *gin.Context) {
	// Param
	targetCol := c.Param("target_col")

	// Validator : Target Column Validator
	if targetCol != "finding_category" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "target_col is not valid",
			"status":  "failed",
		})
		return
	}

	// Service: Get Most Context
	assetFinding, err := rc.AssetFindingService.GetMostContext(targetCol)
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

	// Default values
	var fileExt string
	var fileSize int64
	var fileHeader *multipart.FileHeader = nil

	file, err := c.FormFile("asset_image")
	if file != nil {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "failed to retrieve the file",
				"status":  "failed",
			})
			return
		}

		fileExt = strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Filename), "."))
		fileSize = file.Size
		fileHeader = file

		// Validate file size
		if fileSize > config.MaxSizeFile {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("The file size must be under %.2f MB", float64(config.MaxSizeFile)/1000000),
				"status":  "failed",
			})
			return
		}

		// Optional: open file to validate it can be read
		fileReader, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to open the file",
				"status":  "failed",
			})
			return
		}
		defer fileReader.Close()
	}

	// Service : Create Asset Finding
	if err := rc.AssetFindingService.Create(&req, technicianId, userId, fileHeader, fileExt, fileSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	cleanedRes := utils.CleanResponse(req, "users", "technicians", "asset_placements")

	c.JSON(http.StatusCreated, gin.H{
		"message": "asset finding created successfully",
		"status":  "success",
		"data":    cleanedRes,
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
