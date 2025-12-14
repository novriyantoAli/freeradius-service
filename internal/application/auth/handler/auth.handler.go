package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup) {
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("", h.CreateAuth)
	}
}

// CreateAuth godoc
// @Summary Create authentication credentials
// @Description Create authentication credentials with radcheck and radreply entries in a transaction
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.CreateAuthRequest true "Create Auth Request"
// @Success 201 {object} dto.CreateAuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/auth [post]
func (h *AuthHandler) CreateAuth(ctx *gin.Context) {
	var req dto.CreateAuthRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, err := h.service.CreateAuth(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": result})
}
