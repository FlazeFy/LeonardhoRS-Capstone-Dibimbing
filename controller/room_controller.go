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
	utils.BuildResponseMessage(c, "success", "room", "get", http.StatusOK, room, metadata)
}

// @Summary      Get Room Asset By Floor And Room Name (Detail)
// @Description  Returns a paginated list of asset that found in a room
// @Tags         Room
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetRoomAssetByFloorAndRoomName
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/room/asset/detail/{floor}/{roomName} [get]
// @Param        roomName  path  string  true  "In which Room you want to find the asset. Type 'all' to search in all room"
// @Param        floor  path  string  true  "In which Floor you want to find the asset."
func (rc *RoomController) GetRoomAssetByFloorAndRoomName(c *gin.Context) {
	// Params
	roomName := c.Param("roomName")
	floor := c.Param("floor")

	// Service: Get Find Room Asset By Floor And Room Name
	room, err := rc.RoomService.GetRoomAssetByFloorAndRoomName(floor, roomName)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "room asset", "get", http.StatusOK, room, nil)
}

// @Summary      Get Room Asset By Floor And Room Name (Short)
// @Description  Returns a paginated list of asset that found in a room. in short format
// @Tags         Room
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetRoomAssetShortByFloorAndRoomName
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/room/asset/short/{floor}/{roomName} [get]
// @Param        roomName  path  string  true  "In which Room you want to find the asset. Type 'all' to search in all room"
// @Param        floor  path  string  true  "In which Floor you want to find the asset."
func (rc *RoomController) GetRoomAssetShortByFloorAndRoomName(c *gin.Context) {
	// Params
	roomName := c.Param("roomName")
	floor := c.Param("floor")

	// Service: Get Find Room Asset Short By Floor And Room Name
	room, err := rc.RoomService.GetRoomAssetShortByFloorAndRoomName(floor, roomName)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "room asset", "get", http.StatusOK, room, nil)
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
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Service : Create Room
	err := rc.RoomService.Create(&req)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "room", "post", http.StatusCreated, &req, nil)
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
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Parse Id
	roomID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "Invalid UUID Format")
		return
	}

	// Service : Update Room
	if err := rc.RoomService.UpdateById(&req, roomID); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "room", "put", http.StatusOK, &req, nil)
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
		utils.BuildErrorMessage(c, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	// Service : Delete Room By Id
	if err := rc.RoomService.DeleteById(roomID); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "room", "delete", http.StatusOK, nil, nil)
}

// @Summary      Get Most Context Room
// @Description  Returns a list of most appear item in room by given field
// @Tags         Room
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetMostContext
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/room/most-context/{targetCol} [get]
// @Param        targetCol  path  string  true  "Target Column to Analyze (such as: asset_merk, asset_category, or asset_status)"
func (rc *RoomController) GetMostContext(c *gin.Context) {
	// Param
	targetCol := c.Param("targetCol")

	// Validator : Target Column Validator
	validTarget := []string{"floor", "room_dept"}
	if !utils.Contains(validTarget, targetCol) {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "targetCol is not valid")
		return
	}

	// Service: Get My Room
	room, err := rc.RoomService.GetMostContext(targetCol)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "room", "get", http.StatusOK, room, nil)
}
