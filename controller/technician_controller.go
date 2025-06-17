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

type TechnicianController struct {
	TechnicianService service.TechnicianService
}

func NewTechnicianController(technicianService service.TechnicianService) *TechnicianController {
	return &TechnicianController{
		TechnicianService: technicianService,
	}
}

// @Summary      Get All Technician
// @Description  Returns a paginated list of technician
// @Tags         Technician
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.ResponseGetAllTechnician
// @Failure      404  {object}  entity.ResponseNotFound
// @Router       /api/v1/technicians [get]
func (rc *TechnicianController) GetAllTechnician(c *gin.Context) {
	// Pagination
	pagination := utils.GetPagination(c)

	// Service: Get All Technician
	technician, total, err := rc.TechnicianService.GetAllTechnician(pagination)
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
	utils.BuildResponseMessage(c, "success", "technician", "get", http.StatusOK, technician, metadata)
}

// @Summary      Post Create Technician
// @Description  Create a technician
// @Tags         Technician
// @Accept       application/json
// @Produce      json
// @Param        request  body  entity.RequestPostUpdateTechnicianById  true  "Post Technician Request Body"
// @Success      201  {object}  entity.ResponsePostTechnician
// @Failure      400  {object}  entity.ResponseBadRequest
// @Router       /api/v1/technicians [post]
// @Param        id  path  string  true  "Id of asset technician"
func (rc *TechnicianController) Create(c *gin.Context) {
	// Model
	var req entity.Technician

	// Validator JSON
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

	// Validator Field
	if req.Username == "" {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "username is required")
		return
	}
	if req.Password == "" {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "password is required")
		return
	}
	if req.Email == "" {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "email is required")
		return
	}

	// Service : Create Technician
	if err := rc.TechnicianService.Create(&req, adminId); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "technician", "post", http.StatusCreated, &req, nil)
}

// @Summary      Put Update Technician By Id
// @Description  Update a Technician by Id
// @Tags         Technician
// @Accept       application/json
// @Produce      json
// @Param        request  body  entity.RequestPostUpdateTechnicianById  true  "Update Technician Request Body"
// @Success      200  {object}  entity.ResponseUpdateTechnicianById
// @Failure      400  {object}  entity.ResponseBadRequest
// @Router       /api/v1/technicians/{id} [put]
// @Param        id  path  string  true  "Id of technician"
func (rc *TechnicianController) UpdateById(c *gin.Context) {
	// Param
	id := c.Param("id")

	// Model
	var req entity.Technician

	// Validator JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Parse Id
	technicianID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	// Validator Field
	if req.Username == "" {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "username is required")
		return
	}
	if req.Password == "" {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "password is required")
		return
	}
	if req.Email == "" {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "email is required")
		return
	}

	// Service : Update Technician
	if err := rc.TechnicianService.UpdateById(&req, technicianID); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "technician", "put", http.StatusOK, nil, nil)
}

// @Summary      Delete Technician By Id
// @Description  Permanentally delete technician by id
// @Tags         Technician
// @Success      200  {object}  entity.ResponseDeleteTechnicianById
// @Failure      400  {object}  entity.ResponseBadRequest
// @Router       /api/v1/technicians/{id} [delete]
// @Param        id  path  string  true  "Id of technician"
func (rc *TechnicianController) DeleteById(c *gin.Context) {
	// Param
	id := c.Param("id")

	// Parse Id
	technicianID, err := uuid.Parse(id)
	if err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	// Service : Delete Technician By Id
	if err := rc.TechnicianService.DeleteById(technicianID); err != nil {
		utils.BuildErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	// Response
	utils.BuildResponseMessage(c, "success", "technician", "soft delete", http.StatusOK, nil, nil)
}
