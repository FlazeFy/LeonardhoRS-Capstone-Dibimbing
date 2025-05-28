package controller

import (
	"net/http"
	"pelita/entity"
	"pelita/service"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	RoomService service.RoomService
}

func NewRoomRepository(roomService service.RoomService) *RoomController {
	return &RoomController{RoomService: roomService}
}

func (rc *RoomController) GetAllRoom(c *gin.Context) {
	// Service: Get All Room
	room, err := rc.RoomService.GetAllRoom()
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
