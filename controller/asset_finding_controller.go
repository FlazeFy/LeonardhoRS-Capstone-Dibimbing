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

// @Summary      Get All Asset Finding
// @Description  Returns a paginated list of assets finding
// @Tags         Asset
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetAllAssetFinding
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/asset/finding [get]
func (rc *AssetFindingController) GetAllAssetFinding(c *gin.Context) {
	// Pagination
	pagination := utils.GetPagination(c)

	// Service: Get All Asset Finding
	assetFinding, total, err := rc.AssetFindingService.GetAllAssetFinding(pagination)
	if err != nil {
		utils.BuildErrorMessage(c, err.Error())
		return
	}

	// Response
	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))
	metadata := gin.H{
		"total":       total,
		"page":        pagination.Page,
		"limit":       pagination.Limit,
		"total_pages": totalPages,
	}
	utils.BuildResponseMessage(c, "success", "asset finding", "get", http.StatusOK, assetFinding, metadata)
}

// @Summary      Get All Asset Finding Hour Total
// @Description  Returns a paginated list of assets finding total per hour
// @Tags         Asset
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetFindingHourTotal
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/asset/finding/hour_total [get]
func (rc *AssetFindingController) GetFindingHourTotal(c *gin.Context) {
	// Service: Get All Asset Finding
	assetFinding, err := rc.AssetFindingService.GetFindingHourTotal()
	if err != nil {
		utils.BuildErrorMessage(c, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "asset finding", "get", http.StatusOK, assetFinding, nil)
}

// @Summary      Get Most Context Asset Finding
// @Description  Returns a list of most appear item in asset finding by given field
// @Tags         Asset
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetMostContext
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/asset/most_context/{targe_col} [get]
// @Param        target_col  path  string  true  "Target Column to Analyze (such as: finding_category)"
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
		utils.BuildErrorMessage(c, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "asset finding", "get", http.StatusOK, assetFinding, nil)
}

// @Summary      Post Create Asset Finding
// @Description  Create an asset finding
// @Tags         Asset
// @Accept       multipart/form-data
// @Produce      json
// @Param        finding_category     	formData  string  true  "Finding Category"
// @Param        finding_notes     		formData  string  true  "Finding Notes"
// @Param        finding_image     		formData  file    true  "Finding Image (JPG,PNG,JPEG)"
// @Param        asset_placement_id 	formData  string  true  "Asset Placement Id"
// @Success      201  {object}  entity.ResponseCreateAssetFinding
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/asset/finding [post]
func (rc *AssetFindingController) Create(c *gin.Context) {
	// Model
	var req entity.AssetFinding

	// Validator
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorMessage(c, err.Error())
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
		utils.BuildErrorMessage(c, err.Error())
		return
	}

	// Get Role
	role, err := utils.GetCurrentRole(c)
	if err != nil {
		utils.BuildErrorMessage(c, err.Error())
		return
	}

	// Define The Role Id
	var technicianId, userId *uuid.UUID
	if role == "technician" {
		technicianId = &technicianOrUserId
	} else {
		userId = &technicianOrUserId
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
		utils.BuildErrorMessage(c, err.Error())
		return
	}

	// Response
	cleanedData := utils.CleanResponse(req, "users", "technicians", "asset_placements")
	utils.BuildResponseMessage(c, "success", "asset finding", "post", http.StatusCreated, cleanedData, nil)
}

// @Summary      Delete Asset Finding By Id
// @Description  Permanentally delete asset finding by id
// @Tags         Asset
// @Success      200  {object}  entity.ResponseDeleteAssetFindingById
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/asset/finding/{id} [delete]
// @Param        id  path  string  true  "Id of asset finding"
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
		utils.BuildErrorMessage(c, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "asset finding", "soft delete", http.StatusOK, nil, nil)
}
