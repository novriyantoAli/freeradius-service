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
	radcheckService service.RadcheckService
	logger          *zap.Logger
}

func NewRadcheckHandler(radcheckService service.RadcheckService, logger *zap.Logger) *RadcheckHandler {
	return &RadcheckHandler{
		radcheckService: radcheckService,
		logger:          logger,
	}
}

func (h *RadcheckHandler) RegisterRoutes(r *gin.Engine) {
	radcheckGroup := r.Group("/api/v1/radcheck")
	{
		radcheckGroup.POST("", h.CreateRadcheck)
		radcheckGroup.GET("", h.ListRadcheck)
		radcheckGroup.GET("/:id", h.GetRadcheck)
		radcheckGroup.PUT("/:id", h.UpdateRadcheck)
		radcheckGroup.DELETE("/:id", h.DeleteRadcheck)
	}
}

// CreateRadcheck godoc
// @Summary Create a new Radcheck entry
// @Description Create a new RADIUS check entry for user authentication
// @Tags Radcheck
// @Accept json
// @Produce json
// @Param request body dto.CreateRadcheckRequest true "Create Radcheck Request"
// @Success 201 {object} dto.RadcheckResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/radcheck [post]
func (h *RadcheckHandler) CreateRadcheck(c *gin.Context) {
	var req dto.CreateRadcheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.radcheckService.CreateRadcheck(&req)
	if err != nil {
		h.logger.Error("Failed to create radcheck", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetRadcheck godoc
// @Summary Get Radcheck by ID
// @Description Get a RADIUS check entry by ID
// @Tags Radcheck
// @Accept json
// @Produce json
// @Param id path int true "Radcheck ID"
// @Success 200 {object} dto.RadcheckResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/radcheck/{id} [get]
func (h *RadcheckHandler) GetRadcheck(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error("Invalid ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	resp, err := h.radcheckService.GetRadcheckByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get radcheck", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListRadcheck godoc
// @Summary List Radcheck entries
// @Description List RADIUS check entries with pagination and filtering
// @Tags Radcheck
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param username query string false "Filter by username"
// @Param attribute query string false "Filter by attribute"
// @Success 200 {object} dto.ListRadcheckResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/radcheck [get]
func (h *RadcheckHandler) ListRadcheck(c *gin.Context) {
	var filter dto.RadcheckFilter
	filter.Page = 1
	filter.PageSize = 10

	if err := c.ShouldBindQuery(&filter); err != nil {
		h.logger.Error("Invalid query parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.radcheckService.ListRadcheck(&filter)
	if err != nil {
		h.logger.Error("Failed to list radcheck", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateRadcheck godoc
// @Summary Update Radcheck entry
// @Description Update a RADIUS check entry
// @Tags Radcheck
// @Accept json
// @Produce json
// @Param id path int true "Radcheck ID"
// @Param request body dto.UpdateRadcheckRequest true "Update Radcheck Request"
// @Success 200 {object} dto.RadcheckResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/radcheck/{id} [put]
func (h *RadcheckHandler) UpdateRadcheck(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error("Invalid ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateRadcheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.radcheckService.UpdateRadcheck(uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update radcheck", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteRadcheck godoc
// @Summary Delete Radcheck entry
// @Description Delete a RADIUS check entry
// @Tags Radcheck
// @Accept json
// @Produce json
// @Param id path int true "Radcheck ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/radcheck/{id} [delete]
func (h *RadcheckHandler) DeleteRadcheck(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error("Invalid ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.radcheckService.DeleteRadcheck(uint(id))
	if err != nil {
		h.logger.Error("Failed to delete radcheck", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
