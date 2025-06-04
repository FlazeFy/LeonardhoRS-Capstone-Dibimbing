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

func (rc *HistoryController) GetAllHistory(c *gin.Context) {
	// Pagination
	pagination := utils.GetPagination(c)

	// Service: Get All History
	history, total, err := rc.HistoryService.GetAllHistory(pagination)
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
		"message": "history fetched",
		"status":  "success",
		"data":    history,
		"metadata": gin.H{
			"total":       total,
			"page":        pagination.Page,
			"limit":       pagination.Limit,
			"total_pages": totalPages,
		},
	})
}

func (rc *HistoryController) GetMyHistory(c *gin.Context) {
	// Pagination
	pagination := utils.GetPagination(c)

	// Get User Id
	userId, err := utils.GetCurrentUserID(c)
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

	// Service: Get My History
	history, total, err := rc.HistoryService.GetMyHistory(pagination, userId, role)
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
		"message": "history fetched",
		"status":  "success",
		"data":    history,
		"metadata": gin.H{
			"total":       total,
			"page":        pagination.Page,
			"limit":       pagination.Limit,
			"total_pages": totalPages,
		},
	})
}
