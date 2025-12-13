package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/novriyantoAli/freeradius-service/internal/application/nas/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/nas/service"
	"go.uber.org/zap"
)

type NASHandler struct {
	nasService service.NASService
	logger     *zap.Logger
}

func NewNASHandler(nasService service.NASService, logger *zap.Logger) *NASHandler {
	return &NASHandler{
		nasService: nasService,
		logger:     logger,
	}
}

func (h *NASHandler) RegisterRoutes(r *gin.Engine) {
	nasGroup := r.Group("/api/v1/nas")
	{
		nasGroup.POST("", h.CreateNAS)
		nasGroup.GET("", h.ListNAS)
		nasGroup.GET("/:id", h.GetNAS)
		nasGroup.PUT("/:id", h.UpdateNAS)
		nasGroup.DELETE("/:id", h.DeleteNAS)
	}
}

// CreateNAS godoc
// @Summary Create a new NAS
// @Description Create a new Network Access Server
// @Tags NAS
// @Accept json
// @Produce json
// @Param request body dto.CreateNASRequest true "Create NAS Request"
// @Success 201 {object} dto.NASResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/nas [post]
func (h *NASHandler) CreateNAS(c *gin.Context) {
	var req dto.CreateNASRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.nasService.CreateNAS(&req)
	if err != nil {
		h.logger.Error("Failed to create NAS", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetNAS godoc
// @Summary Get NAS by ID
// @Description Get a Network Access Server by ID
// @Tags NAS
// @Accept json
// @Produce json
// @Param id path int true "NAS ID"
// @Success 200 {object} dto.NASResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/nas/{id} [get]
func (h *NASHandler) GetNAS(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error("Invalid ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	resp, err := h.nasService.GetNASByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get NAS", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListNAS godoc
// @Summary List NAS
// @Description List Network Access Servers with pagination and filtering
// @Tags NAS
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param nasname query string false "Filter by NAS name"
// @Param shortname query string false "Filter by short name"
// @Param type query string false "Filter by type"
// @Param description query string false "Filter by description"
// @Success 200 {object} dto.ListNASResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/nas [get]
func (h *NASHandler) ListNAS(c *gin.Context) {
	var filter dto.NASFilter
	filter.Page = 1
	filter.PageSize = 10

	if err := c.ShouldBindQuery(&filter); err != nil {
		h.logger.Error("Invalid query parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.nasService.ListNAS(&filter)
	if err != nil {
		h.logger.Error("Failed to list NAS", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateNAS godoc
// @Summary Update NAS
// @Description Update a Network Access Server
// @Tags NAS
// @Accept json
// @Produce json
// @Param id path int true "NAS ID"
// @Param request body dto.UpdateNASRequest true "Update NAS Request"
// @Success 200 {object} dto.NASResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/nas/{id} [put]
func (h *NASHandler) UpdateNAS(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error("Invalid ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateNASRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.nasService.UpdateNAS(uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update NAS", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteNAS godoc
// @Summary Delete NAS
// @Description Delete a Network Access Server
// @Tags NAS
// @Accept json
// @Produce json
// @Param id path int true "NAS ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/nas/{id} [delete]
func (h *NASHandler) DeleteNAS(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		h.logger.Error("Invalid ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.nasService.DeleteNAS(uint(id))
	if err != nil {
		h.logger.Error("Failed to delete NAS", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
