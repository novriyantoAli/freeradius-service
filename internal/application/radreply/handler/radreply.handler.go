package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/service"
	"gorm.io/gorm"
)

type RadreplyHandler struct {
	service service.RadreplyService
}

func NewRadreplyHandler(service service.RadreplyService) *RadreplyHandler {
	return &RadreplyHandler{service: service}
}

func (h *RadreplyHandler) RegisterRoutes(r *gin.RouterGroup) {
	radreplyRoutes := r.Group("/radreply")
	{
		radreplyRoutes.POST("", h.CreateRadreply)
		radreplyRoutes.GET("/:id", h.GetRadreply)
		radreplyRoutes.GET("", h.ListRadreply)
		radreplyRoutes.PUT("/:id", h.UpdateRadreply)
		radreplyRoutes.DELETE("/:id", h.DeleteRadreply)
	}
}

// CreateRadreply godoc
// @Summary Create a new radreply entry
// @Description Create a new RADIUS reply entry for user replies
// @Tags Radreply
// @Accept json
// @Produce json
// @Param request body dto.CreateRadreplyRequest true "Create Radreply Request"
// @Success 201 {object} dto.RadreplyResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/radreply [post]
func (h *RadreplyHandler) CreateRadreply(ctx *gin.Context) {
	var req dto.CreateRadreplyRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, err := h.service.CreateRadreply(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": result})
}

// GetRadreply godoc
// @Summary Get a radreply by ID
// @Description Get a RADIUS reply entry by ID
// @Tags Radreply
// @Accept json
// @Produce json
// @Param id path int true "Radreply ID"
// @Success 200 {object} dto.RadreplyResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/radreply/{id} [get]
func (h *RadreplyHandler) GetRadreply(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	result, err := h.service.GetRadreplyByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "radreply not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// ListRadreply godoc
// @Summary List radreply entries
// @Description List RADIUS reply entries with pagination and filtering
// @Tags Radreply
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param username query string false "Filter by username"
// @Param attribute query string false "Filter by attribute"
// @Success 200 {object} dto.ListRadreplyResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/radreply [get]
func (h *RadreplyHandler) ListRadreply(ctx *gin.Context) {
	var filter dto.RadreplyFilter

	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, err := h.service.ListRadreply(&filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// UpdateRadreply godoc
// @Summary Update radreply entry
// @Description Update a RADIUS reply entry
// @Tags Radreply
// @Accept json
// @Produce json
// @Param id path int true "Radreply ID"
// @Param request body dto.UpdateRadreplyRequest true "Update Radreply Request"
// @Success 200 {object} dto.RadreplyResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/radreply/{id} [put]
func (h *RadreplyHandler) UpdateRadreply(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	var req dto.UpdateRadreplyRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, err := h.service.UpdateRadreply(uint(id), &req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "radreply not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

// DeleteRadreply godoc
// @Summary Delete radreply entry
// @Description Delete a RADIUS reply entry
// @Tags Radreply
// @Accept json
// @Produce json
// @Param id path int true "Radreply ID"
// @Success 204
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/radreply/{id} [delete]
func (h *RadreplyHandler) DeleteRadreply(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
		return
	}

	if err := h.service.DeleteRadreply(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "radreply deleted successfully"})
}
