package controller

import (
	"net/http"
	"pelita/entity"
	"pelita/service"

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
