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
