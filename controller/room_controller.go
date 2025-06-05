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

type RoomController struct {
	RoomService service.RoomService
}

func NewRoomRepository(roomService service.RoomService) *RoomController {
	return &RoomController{RoomService: roomService}
}

func (rc *RoomController) GetAllRoom(c *gin.Context) {
	// Pagination
	pagination := utils.GetPagination(c)

	// Service: Get All Room
	room, total, err := rc.RoomService.GetAllRoom(pagination)
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
		"message": "room fetched",
		"status":  "success",
		"data":    room,
		"metadata": gin.H{
			"total":       total,
			"page":        pagination.Page,
			"limit":       pagination.Limit,
			"total_pages": totalPages,
		},
	})
}

func (rc *RoomController) GetRoomAssetByFloorAndRoomName(c *gin.Context) {
	// Params
	roomName := c.Param("room_name")
	floor := c.Param("floor")

	// Service: Get Find Room Asset By Floor And Room Name
	room, err := rc.RoomService.GetRoomAssetByFloorAndRoomName(floor, roomName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "room asset fetched",
		"status":  "success",
		"data":    room,
	})
}

func (rc *RoomController) GetRoomAssetShortByFloorAndRoomName(c *gin.Context) {
	// Params
	roomName := c.Param("room_name")
	floor := c.Param("floor")

	// Service: Get Find Room Asset Short By Floor And Room Name
	room, err := rc.RoomService.GetRoomAssetShortByFloorAndRoomName(floor, roomName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "room asset fetched",
		"status":  "success",
		"data":    room,
	})
}

func (rc *RoomController) Create(c *gin.Context) {
	// Model
	var req entity.Room

	// Validator
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Service : Create Room
	err := rc.RoomService.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "room created successfully",
		"status":  "success",
		"data":    &req,
	})
}

func (rc *RoomController) UpdateById(c *gin.Context) {
	// Param
	id := c.Param("id")

	// Model
	var req entity.Room

	// Validator
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Parse Id
	roomID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid UUID format",
			"status":  "failed",
		})
		return
	}

	// Service : Update Room
	if err := rc.RoomService.UpdateById(&req, roomID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "room update successfully",
		"status":  "success",
	})
}

func (rc *RoomController) DeleteById(c *gin.Context) {
	// Param
	id := c.Param("id")

	// Parse Id
	roomID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid UUID format",
			"status":  "failed",
		})
		return
	}

	// Service : Delete Room By Id
	if err := rc.RoomService.DeleteById(roomID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "room deleted",
		"status":  "success",
	})
}

func (rc *RoomController) GetMostContext(c *gin.Context) {
	// Param
	targetCol := c.Param("target_col")

	// Validator : Target Column Validator
	validTarget := []string{"floor", "room_dept"}
	if !utils.Contains(validTarget, targetCol) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "target_col is not valid",
			"status":  "failed",
		})
		return
	}

	// Service: Get My Room
	room, err := rc.RoomService.GetMostContext(targetCol)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"status":  "failed",
		})
		return
	}

	// Response
	c.JSON(http.StatusOK, gin.H{
		"message": "room fetched",
		"status":  "success",
		"data":    room,
	})
}
