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

// @Summary      Get All Asset Placement
// @Description  Returns a paginated list of assets placement
// @Tags         Asset
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetAllAssetPlacement
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/asset/placement [get]
func (rc *AssetPlacementController) GetAllAssetPlacement(c *gin.Context) {
	// Pagination
	pagination := utils.GetPagination(c)

	// Service: Get All Asset Placement
	assetPlacement, total, err := rc.AssetPlacementService.GetAllAssetPlacement(pagination)
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
	utils.BuildResponseMessage(c, "success", "asset placement", "get", http.StatusOK, assetPlacement, metadata)
}

// @Summary      Post Create Asset Placement By Id
// @Description  Create an asset maintenance by Id
// @Tags         Asset
// @Accept       application/json
// @Produce      json
// @Param        request  body  entity.RequestCreateUpdateAssetPlacement  true  "Create Asset Placement Request Body"
// @Success      201  {object}  entity.ResponseCreateAssetPlacement
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/asset/placement [post]
func (rc *AssetPlacementController) Create(c *gin.Context) {
	// Model
	var req entity.AssetPlacement

	// Validator
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get User Id
	adminId, err := utils.GetCurrentUserID(c)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Service : Create Asset Placement
	if err := rc.AssetPlacementService.Create(&req, adminId); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "asset placement", "post", http.StatusCreated, &req, nil)
}

// @Summary      Put Update Asset Placement By Id
// @Description  Create an asset placement by Id
// @Tags         Asset
// @Accept       application/json
// @Produce      json
// @Param        request  body  entity.RequestCreateUpdateAssetPlacement  true  "Put Update Asset Placement Request Body"
// @Success      200  {object}  entity.ResponsePutUpdateAssetPlacement
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/asset/placement/{id} [put]
// @Param        id  path  string  true  "Id of asset placement"
func (rc *AssetPlacementController) UpdateById(c *gin.Context) {
	// Param
	id := c.Param("id")

	// Model
	var req entity.AssetPlacement

	// Validator
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Parse Id
	assetPlacementID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	// Service : Update Asset Placement
	if err := rc.AssetPlacementService.UpdateById(&req, assetPlacementID); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "asset placement", "put", http.StatusOK, &req, nil)
}

// @Summary      Delete Asset Placement By Id
// @Description  Permanentally delete asset placement by id
// @Tags         Asset
// @Success      200  {object}  entity.ResponseDeleteAssetPlacementById
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/asset/placement/{id} [delete]
// @Param        id  path  string  true  "Id of asset placement"
func (rc *AssetPlacementController) DeleteById(c *gin.Context) {
	// Param
	id := c.Param("id")

	// Parse Id
	assetPlacementID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	// Service : Delete Asset Placement By Id
	if err := rc.AssetPlacementService.DeleteById(assetPlacementID); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "asset placement", "soft delete", http.StatusOK, nil, nil)
}
