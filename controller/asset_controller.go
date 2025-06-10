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

type AssetController struct {
	AssetService service.AssetService
}

func NewAssetRepository(assetService service.AssetService) *AssetController {
	return &AssetController{AssetService: assetService}
}

type Config struct {
	MaxSizeFile     int64
	AllowedFileType []string
}

var config = Config{
	MaxSizeFile:     10000000, // 10 MB
	AllowedFileType: []string{"jpg", "jpeg"},
}

// @Summary      Get All Asset
// @Description  Returns a paginated list of assets available
// @Tags         Asset
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetAllAsset
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/asset [get]
func (rc *AssetController) GetAllAsset(c *gin.Context) {
	// Pagination
	pagination := utils.GetPagination(c)

	// Service: Get All Asset
	asset, total, err := rc.AssetService.GetAllAsset(pagination)
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
		"message": "asset fetched",
		"status":  "success",
		"data":    asset,
		"metadata": gin.H{
			"total":       total,
			"page":        pagination.Page,
			"limit":       pagination.Limit,
			"total_pages": totalPages,
		},
	})
}

// @Summary      Get Deleted Asset
// @Description  Returns a list of deleted assets
// @Tags         Asset
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetDeletedAsset
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/asset/deleted [get]
func (rc *AssetController) GetDeletedAsset(c *gin.Context) {
	// Service: Get All Deleted Asset
	asset, err := rc.AssetService.GetDeleted()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "deleted asset fetched",
		"status":  "success",
		"data":    asset,
	})
}

func (rc *AssetController) Create(c *gin.Context) {
	// Model
	var req entity.Asset

	// Multipart Form
	if err := c.Request.ParseMultipartForm(20 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed to parse form data",
			"status":  "failed",
		})
		return
	}

	// Parse Form
	req.AssetName = c.PostForm("asset_name")
	req.AssetDesc = utils.OptionalString(c.PostForm("asset_desc"))
	req.AssetMerk = utils.OptionalString(c.PostForm("asset_merk"))
	req.AssetCategory = c.PostForm("asset_category")
	req.AssetPrice = utils.OptionalString(c.PostForm("asset_price"))
	req.AssetStatus = c.PostForm("asset_status")

	// Get User Id
	adminId, err := utils.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
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

	// Service : Create Asset
	if err := rc.AssetService.Create(&req, adminId, fileHeader, fileExt, fileSize); err != nil {
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

func (rc *AssetController) SoftDeleteById(c *gin.Context) {
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

	// Service : Soft Delete Asset By Id
	if err := rc.AssetService.SoftDeleteById(assetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "asset deleted",
		"status":  "success",
	})
}

func (rc *AssetController) RecoverDeletedById(c *gin.Context) {
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

	// Service : Recover Delete Asset By Id
	if err := rc.AssetService.RecoverDeletedById(assetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "asset recovered",
		"status":  "success",
	})
}

// @Summary      Get Most Context Asset
// @Description  Returns a list of most appear item in asset by given field
// @Tags         Asset
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetMostContext
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/asset/most_context/{targe_col} [get]
// @Param        target_col  path  string  true  "Target Column to Analyze (such as: asset_merk, asset_category, or asset_status)"
func (rc *AssetController) GetMostContext(c *gin.Context) {
	// Param
	targetCol := c.Param("target_col")

	// Validator : Target Column Validator
	validTarget := []string{"asset_merk", "asset_category", "asset_status"}
	if !utils.Contains(validTarget, targetCol) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "target_col is not valid",
			"status":  "failed",
		})
		return
	}

	// Service: Get Most Context
	asset, err := rc.AssetService.GetMostContext(targetCol)
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
