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
// @Failure      404  {object}  entity.ResponseNotFound
// @Router       /api/v1/assets/maintenances [get]
func (rc *AssetMaintenanceController) GetAllAssetMaintenance(c *gin.Context) {
	// Pagination
	pagination := utils.GetPagination(c)

	// Service: Get All Asset Maintenance
	assetMaintenance, total, err := rc.AssetMaintenanceService.GetAllAssetMaintenance(pagination)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
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
// @Failure      404  {object}  entity.ResponseNotFound
// @Router       /api/v1/assets/maintenances/schedule [get]
func (rc *AssetMaintenanceController) GetAllAssetMaintenanceSchedule(c *gin.Context) {
	// Service: Get All Asset Maintenance
	assetMaintenance, err := rc.AssetMaintenanceService.GetAllAssetMaintenanceSchedule()
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
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
// @Failure      400  {object}  entity.ResponseBadRequest
// @Router       /api/v1/assets/maintenances [post]
func (rc *AssetMaintenanceController) Create(c *gin.Context) {
	// Model
	var req entity.AssetMaintenance

	// Validator JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validator Rules
	validDays := map[string]bool{"Sun": true, "Mon": true, "Tue": true, "Wed": true, "Thu": true, "Fri": true, "Sat": true}
	if !validDays[req.MaintenanceDay] {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "maintenance day must be one of: Sun, Mon, Tue, Wed, Thu, Fri, Sat")
		return
	}

	// Get User Id
	adminId, err := utils.GetCurrentUserID(c)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Validator Field
	if req.AssetPlacementId == uuid.Nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "asset placement id is required")
		return
	}
	if req.MaintenanceBy == uuid.Nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "asset maintenance by is required")
		return
	}
	if req.MaintenanceDay == "" {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "asset maintenance day is required")
		return
	}

	// Service : Create Asset Maintenance
	if err := rc.AssetMaintenanceService.Create(&req, adminId); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
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
// @Failure      400  {object}  entity.ResponseBadRequest
// @Router       /api/v1/assets/maintenances/{id} [put]
// @Param        id  path  string  true  "Id of asset maintenance"
func (rc *AssetMaintenanceController) UpdateById(c *gin.Context) {
	// Param
	id := c.Param("id")

	// Model
	var req entity.AssetMaintenance

	// Validator JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Parse Id
	assetMaintenanceID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	// Validator Field
	if req.AssetPlacementId == uuid.Nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "asset placement id is required")
		return
	}
	if req.MaintenanceBy == uuid.Nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "asset maintenance by is required")
		return
	}
	if req.MaintenanceDay == "" {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "asset maintenance day is required")
		return
	}

	// Service : Update Asset Maintenance
	if err := rc.AssetMaintenanceService.UpdateById(&req, assetMaintenanceID); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "asset maintenance", "put", http.StatusOK, &req, nil)
}

// @Summary      Delete Asset Maintenance By Id
// @Description  Permanentally delete asset maintenance by id
// @Tags         Asset
// @Success      200  {object}  entity.ResponseDeleteAssetMaintenanceById
// @Failure      400  {object}  entity.ResponseBadRequest
// @Router       /api/v1/assets/maintenances/{id} [delete]
// @Param        id  path  string  true  "Id of asset maintenance"
func (rc *AssetMaintenanceController) DeleteById(c *gin.Context) {
	// Param
	id := c.Param("id")

	// Parse Id
	assetMaintenanceID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	// Service : Delete Asset Maintenance By Id
	if err := rc.AssetMaintenanceService.DeleteById(assetMaintenanceID); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
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
// @Failure      404  {object}  entity.ResponseNotFound
// @Router       /api/v1/assets/most-context/{targetCol} [get]
// @Param        targetCol  path  string  true  "Target Column to Analyze (such as: maintenance_day)"
func (rc *AssetMaintenanceController) GetMostContext(c *gin.Context) {
	// Param
	targetCol := c.Param("targetCol")

	// Validator : Target Column Validator
	if targetCol != "maintenance_day" {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "targetCol is not valid")
		return
	}

	// Service: Get Most Context
	assetMaintenance, err := rc.AssetMaintenanceService.GetMostContext(targetCol)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "asset maintenance", "get", http.StatusOK, assetMaintenance, nil)
}
