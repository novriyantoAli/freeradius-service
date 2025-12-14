package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/service"
	"go.uber.org/zap"
)

type RadcheckHandler struct {
	service service.RadcheckService
	logger  *zap.Logger
}

func NewRadcheckHandler(service service.RadcheckService, logger *zap.Logger) *RadcheckHandler {
	return &RadcheckHandler{
		service: service,
		logger:  logger,
	}
}

// CreateRadcheck godoc
// @Summary Create a new radcheck entry
// @Description Create a new RADIUS check entry for user authentication
// @Tags radcheck
// @Accept json
// @Produce json
// @Param request body dto.CreateRadcheckRequest true "Radcheck creation request"
// @Success 201 {object} map[string]interface{} "Created radcheck"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/radcheck [post]
func (h *RadcheckHandler) CreateRadcheck(ctx *gin.Context) {
	var req dto.CreateRadcheckRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	radcheck, err := h.service.CreateRadcheck(ctx.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create radcheck", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create radcheck"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": radcheck})
}

// GetRadcheck godoc
// @Summary Get a radcheck by ID
// @Description Get a single radcheck by its ID
// @Tags radcheck
// @Accept json
// @Produce json
// @Param id path int true "Radcheck ID"
// @Success 200 {object} map[string]interface{} "Radcheck details"
// @Failure 400 {object} map[string]interface{} "Invalid radcheck ID"
// @Failure 404 {object} map[string]interface{} "Radcheck not found"
// @Router /api/v1/radcheck/{id} [get]
func (h *RadcheckHandler) GetRadcheck(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid radcheck ID"})
		return
	}

	radcheck, err := h.service.GetRadcheckByID(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get radcheck", zap.Error(err))
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Radcheck not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": radcheck})
}

// ListRadcheck godoc
// @Summary List all radchecks
// @Description Get a list of radchecks with optional filtering and pagination
// @Tags radcheck
// @Accept json
// @Produce json
// @Param username query string false "Filter by username"
// @Param attribute query string false "Filter by attribute"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Number of items per page" default(10)
// @Success 200 {object} dto.ListRadcheckResponse "List of radchecks"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/radcheck [get]
func (h *RadcheckHandler) ListRadcheck(ctx *gin.Context) {
	var filter dto.RadcheckFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		h.logger.Error("Invalid query parameters", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	radchecks, err := h.service.ListRadcheck(ctx.Request.Context(), &filter)
	if err != nil {
		h.logger.Error("Failed to list radcheck", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list radcheck"})
		return
	}

	ctx.JSON(http.StatusOK, radchecks)
}

// UpdateRadcheck godoc
// @Summary Update a radcheck entry
// @Description Update a radcheck entry by ID
// @Tags radcheck
// @Accept json
// @Produce json
// @Param id path int true "Radcheck ID"
// @Param request body dto.UpdateRadcheckRequest true "Radcheck update request"
// @Success 200 {object} map[string]interface{} "Updated radcheck"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 404 {object} map[string]interface{} "Radcheck not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/radcheck/{id} [put]
func (h *RadcheckHandler) UpdateRadcheck(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid radcheck ID"})
		return
	}

	var req dto.UpdateRadcheckRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	radcheck, err := h.service.UpdateRadcheck(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update radcheck", zap.Error(err))
		if err.Error() == "radcheck not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update radcheck"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": radcheck})
}

// DeleteRadcheck godoc
// @Summary Delete a radcheck entry
// @Description Delete a radcheck entry by ID
// @Tags radcheck
// @Accept json
// @Produce json
// @Param id path int true "Radcheck ID"
// @Success 200 {object} map[string]interface{} "Radcheck deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid radcheck ID"
// @Failure 404 {object} map[string]interface{} "Radcheck not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/radcheck/{id} [delete]
func (h *RadcheckHandler) DeleteRadcheck(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid radcheck ID"})
		return
	}

	err = h.service.DeleteRadcheck(ctx.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to delete radcheck", zap.Error(err))
		if err.Error() == "radcheck not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete radcheck"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Radcheck deleted successfully"})
}

func (h *RadcheckHandler) RegisterRoutes(api *gin.RouterGroup) {
	radcheck := api.Group("/radcheck")
	{
		radcheck.POST("", h.CreateRadcheck)
		radcheck.GET("", h.ListRadcheck)
		radcheck.GET("/:id", h.GetRadcheck)
		radcheck.PUT("/:id", h.UpdateRadcheck)
		radcheck.DELETE("/:id", h.DeleteRadcheck)
	}
}
