package controller

import (
	"math"
	"net/http"
	"pelita/service"
	"pelita/utils"

	"github.com/gin-gonic/gin"
)

type HistoryController struct {
	HistoryService service.HistoryService
}

func NewHistoryRepository(historyService service.HistoryService) *HistoryController {
	return &HistoryController{HistoryService: historyService}
}

// @Summary      Get All History
// @Description  Returns a paginated list of all users histories
// @Tags         History
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetAllHistory
// @Failure      404  {object}  entity.ResponseNotFound
// @Router       /api/v1/histories/all [get]
func (rc *HistoryController) GetAllHistory(c *gin.Context) {
	// Pagination
	pagination := utils.GetPagination(c)

	// Service: Get All History
	history, total, err := rc.HistoryService.GetAllHistory(pagination)
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
	utils.BuildResponseMessage(c, "success", "history", "get", http.StatusOK, history, metadata)
}

// @Summary      Get My History
// @Description  Returns a paginated list of my histories
// @Tags         History
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetMyHistory
// @Failure      404  {object}  entity.ResponseNotFound
// @Router       /api/v1/histories/my [get]
func (rc *HistoryController) GetMyHistory(c *gin.Context) {
	// Pagination
	pagination := utils.GetPagination(c)

	// Get User Id
	userId, err := utils.GetCurrentUserID(c)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Get Role
	role, err := utils.GetCurrentRole(c)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Service: Get My History
	history, total, err := rc.HistoryService.GetMyHistory(pagination, userId, role)
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
	utils.BuildResponseMessage(c, "success", "history", "get", http.StatusOK, history, metadata)
}

// @Summary      Get Most Context History
// @Description  Returns a list of most appear item in history by given field
// @Tags         History
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetMostContext
// @Failure      404  {object}  entity.ResponseNotFound
// @Router       /api/v1/histories/most-context/{targetCol} [get]
// @Param        targetCol  path  string  true  "Target Column to Analyze (such as: type_user, type_history)"
func (rc *HistoryController) GetMostContext(c *gin.Context) {
	// Param
	targetCol := c.Param("targetCol")

	// Validator : Target Column Validator
	validTarget := []string{"type_user", "type_history"}
	if !utils.Contains(validTarget, targetCol) {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "targetCol is not valid")
		return
	}

	// Service: Get My History
	history, err := rc.HistoryService.GetMostContext(targetCol)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "history", "get", http.StatusOK, history, nil)
}
