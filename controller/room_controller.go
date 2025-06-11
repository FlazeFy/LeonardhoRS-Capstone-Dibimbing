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

// @Summary      Get All Room
// @Description  Returns a paginated list of room
// @Tags         Room
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetAllRoom
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/room [get]
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

// @Summary      Get Room Asset By Floor And Room Name (Detail)
// @Description  Returns a paginated list of asset that found in a room
// @Tags         Room
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetRoomAssetByFloorAndRoomName
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/room/asset/detail/{floor}/{room_name} [get]
// @Param        room_name  path  string  true  "In which Room you want to find the asset. Type 'all' to search in all room"
// @Param        floor  path  string  true  "In which Floor you want to find the asset."
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

// @Summary      Get Room Asset By Floor And Room Name (Short)
// @Description  Returns a paginated list of asset that found in a room. in short format
// @Tags         Room
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetRoomAssetShortByFloorAndRoomName
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/room/asset/short/{floor}/{room_name} [get]
// @Param        room_name  path  string  true  "In which Room you want to find the asset. Type 'all' to search in all room"
// @Param        floor  path  string  true  "In which Floor you want to find the asset."
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

// @Summary      Post Create Room
// @Description  Create an room
// @Tags         Room
// @Accept       application/json
// @Produce      json
// @Param        request  body  entity.RequestPostCreateUpdateRoom true  "Post Create Room Request Body"
// @Success      200  {object}  entity.ResponsePutUpdateRoom
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/room [post]
// @Param        id  path  string  true  "Id of asset placement"
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

// @Summary      Put Update Room
// @Description  Update an room by id
// @Tags         Room
// @Accept       application/json
// @Produce      json
// @Param        request  body  entity.RequestPostCreateUpdateRoom true  "Put Update Room Request Body"
// @Success      200  {object}  entity.ResponsePutUpdateRoom
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/room [post]
// @Param        id  path  string  true  "Id of room"
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

// @Summary      Delete Room By Id
// @Description  Permanentally delete room by id
// @Tags         Room
// @Success      200  {object}  entity.ResponseDeleteRoomById
// @Failure      400  {object}  map[string]string
// @Router       /api/v1/room/{id} [delete]
// @Param        id  path  string  true  "Id of room"
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

// @Summary      Get Most Context Room
// @Description  Returns a list of most appear item in room by given field
// @Tags         Room
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetMostContext
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/room/most_context/{targe_col} [get]
// @Param        target_col  path  string  true  "Target Column to Analyze (such as: asset_merk, asset_category, or asset_status)"
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
