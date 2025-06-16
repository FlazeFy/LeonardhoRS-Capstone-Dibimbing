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

// @Summary      Get All Asset Maintenance
// @Description  Returns a paginated list of assets maintenance
// @Tags         Asset
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetAllAssetMaintenance
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/asset/maintenance [get]
func (rc *AssetMaintenanceController) GetAllAssetMaintenance(c *gin.Context) {
	// Pagination
	pagination := utils.GetPagination(c)

	// Service: Get All Asset Maintenance
	assetMaintenance, total, err := rc.AssetMaintenanceService.GetAllAssetMaintenance(pagination)
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
	utils.BuildResponseMessage(c, "success", "asset maintenance", "get", http.StatusOK, assetMaintenance, metadata)
}

// @Summary      Get All Asset Maintenance Schedule
// @Description  Returns a list of assets maintenance schedule
// @Tags         Asset
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetAllAssetMaintenanceSchedule
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/asset/maintenance/schedule [get]
func (rc *AssetMaintenanceController) GetAllAssetMaintenanceSchedule(c *gin.Context) {
	// Service: Get All Asset Maintenance
	assetMaintenance, err := rc.AssetMaintenanceService.GetAllAssetMaintenanceSchedule()
	if err != nil {
		utils.BuildErrorMessage(c, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "asset maintenance", "get", http.StatusOK, assetMaintenance, nil)
}

// @Summary      Post Create Asset Maintenance By Id
// @Description  Create an asset maintenance by Id
// @Tags         Asset
// @Accept       application/json
// @Produce      json
// @Param        request  body  entity.RequestCreateUpdateAssetMaintenance  true  "Create Asset Maintenance Request Body"
// @Success      201  {object}  entity.ResponseCreateAssetMaintenance
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/asset/maintenance [post]
func (rc *AssetMaintenanceController) Create(c *gin.Context) {
	// Model
	var req entity.AssetMaintenance

	// Validator
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorMessage(c, err.Error())
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
		utils.BuildErrorMessage(c, err.Error())
		return
	}

	// Service : Create Asset Maintenance
	if err := rc.AssetMaintenanceService.Create(&req, adminId); err != nil {
		utils.BuildErrorMessage(c, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "asset maintenance", "post", http.StatusCreated, &req, nil)
}

// @Summary      Put Update Asset Maintenance By Id
// @Description  Create an asset maintenance by Id
// @Tags         Asset
// @Accept       application/json
// @Produce      json
// @Param        request  body  entity.RequestCreateUpdateAssetMaintenance  true  "Put Update Asset Maintenance Request Body"
// @Success      200  {object}  entity.ResponsePutUpdateAssetMaintenance
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/asset/maintenance/{id} [put]
// @Param        id  path  string  true  "Id of asset maintenance"
func (rc *AssetMaintenanceController) UpdateById(c *gin.Context) {
	// Param
	id := c.Param("id")

	// Model
	var req entity.AssetMaintenance

	// Validator
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorMessage(c, err.Error())
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
		utils.BuildErrorMessage(c, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "asset maintenance", "put", http.StatusOK, &req, nil)
}

// @Summary      Delete Asset Maintenance By Id
// @Description  Permanentally delete asset maintenance by id
// @Tags         Asset
// @Success      200  {object}  entity.ResponseDeleteAssetMaintenanceById
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/asset/maintenance/{id} [delete]
// @Param        id  path  string  true  "Id of asset maintenance"
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
		utils.BuildErrorMessage(c, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "asset maintenance", "soft delete", http.StatusOK, nil, nil)
}

// @Summary      Get Most Context Asset Maintenance
// @Description  Returns a list of most appear item in asset maintenance by given field
// @Tags         Asset
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetMostContext
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/asset/most_context/{targe_col} [get]
// @Param        target_col  path  string  true  "Target Column to Analyze (such as: maintenance_day)"
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
		utils.BuildErrorMessage(c, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "asset maintenance", "get", http.StatusOK, assetMaintenance, nil)
}
