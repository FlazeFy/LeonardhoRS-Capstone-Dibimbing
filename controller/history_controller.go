package controller

import (
	"net/http"
	"pelita/service"

	"github.com/gin-gonic/gin"
)

type HistoryController struct {
	HistoryService service.HistoryService
}

func NewHistoryRepository(historyService service.HistoryService) *HistoryController {
	return &HistoryController{HistoryService: historyService}
}

func (rc *HistoryController) GetAllHistory(c *gin.Context) {
	// Service: Get All History
	history, err := rc.HistoryService.GetAllHistory()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "history fetched",
		"status":  "success",
		"data":    history,
	})
}
