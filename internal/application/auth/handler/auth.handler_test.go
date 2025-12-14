package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/handler"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/service"
	radcheckEntity "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/entity"
	radreplyEntity "github.com/novriyantoAli/freeradius-service/internal/application/radreply/entity"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestAuthHandler_CreateAuth_Success(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.CreateFn = func(ctx context.Context, radcheck *radcheckEntity.Radcheck) error {
		radcheck.ID = 1
		return nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()
	mockRadreplyRepo.CreateFn = func(ctx context.Context, radreply *radreplyEntity.Radreply) error {
		radreply.ID = 1
		return nil
	}

	mockTxManager := &testutil.MockTransactionManager{}
	mockTxManager.WithinTransactionFn = func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	}

	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo, mockTxManager)
	authHandler := handler.NewAuthHandler(authService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	authHandler.RegisterRoutes(api)

	req := dto.CreateAuthRequest{
		Username: "newuser",
		Password: "password123",
		Attributes: []dto.CreateAuthAttribute{
			{Attribute: "Framed-IP-Address", Value: "192.168.1.100", Op: "="},
		},
		ReplyAttrs: []dto.CreateAuthAttribute{
			{Attribute: "Reply-Message", Value: "Welcome", Op: "="},
		},
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/auth", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, httpReq)

	require.Equal(t, http.StatusCreated, writer.Code)

	var response map[string]interface{}
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	require.Equal(t, "newuser", data["username"].(string))
	require.Equal(t, "password123", data["password"].(string))
}

func TestAuthHandler_CreateAuth_MissingUsername(t *testing.T) {
	mockTxManager := &testutil.MockTransactionManager{}
	authService := service.NewAuthService(nil, nil, mockTxManager)
	authHandler := handler.NewAuthHandler(authService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	authHandler.RegisterRoutes(api)

	req := dto.CreateAuthRequest{
		Username: "",
		Password: "password123",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/auth", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, httpReq)

	require.Equal(t, http.StatusBadRequest, writer.Code)
}
