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
	radcheckdto "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/dto"
	radcheckEntity "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/entity"
	radrepldto "github.com/novriyantoAli/freeradius-service/internal/application/radreply/dto"
	radreplyEntity "github.com/novriyantoAli/freeradius-service/internal/application/radreply/entity"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestAuthHandler_Authenticate_Success(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.GetAllFn = func(ctx context.Context, filter *radcheckdto.RadcheckFilter) ([]radcheckEntity.Radcheck, int64, error) {
		return []radcheckEntity.Radcheck{*testutil.CreateRadcheckFixture()}, 1, nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()
	mockRadreplyRepo.GetAllFn = func(ctx context.Context, filter *radrepldto.RadreplyFilter) ([]radreplyEntity.Radreply, int64, error) {
		return []radreplyEntity.Radreply{}, 0, nil
	}

	mockTxManager := &testutil.MockTransactionManager{}
	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo, mockTxManager)
	authHandler := handler.NewAuthHandler(authService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	authHandler.RegisterRoutes(api)

	req := dto.AuthenticateRequest{
		Username: "testuser",
		Password: "testing123",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/auth/authenticate", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, httpReq)

	require.Equal(t, http.StatusOK, writer.Code)

	var response map[string]interface{}
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	require.True(t, data["success"].(bool))
}

func TestAuthHandler_Authenticate_MissingUsername(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.GetAllFn = func(ctx context.Context, filter *radcheckdto.RadcheckFilter) ([]radcheckEntity.Radcheck, int64, error) {
		return []radcheckEntity.Radcheck{}, 0, nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()

	mockTxManager := &testutil.MockTransactionManager{}
	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo, mockTxManager)
	authHandler := handler.NewAuthHandler(authService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	authHandler.RegisterRoutes(api)

	req := dto.AuthenticateRequest{
		Password: "test123",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/auth/authenticate", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, httpReq)

	require.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestAuthHandler_Authenticate_MissingPassword(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.GetAllFn = func(ctx context.Context, filter *radcheckdto.RadcheckFilter) ([]radcheckEntity.Radcheck, int64, error) {
		return []radcheckEntity.Radcheck{}, 0, nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()

	mockTxManager := &testutil.MockTransactionManager{}
	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo, mockTxManager)
	authHandler := handler.NewAuthHandler(authService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	authHandler.RegisterRoutes(api)

	req := dto.AuthenticateRequest{
		Username: "testuser",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/auth/authenticate", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, httpReq)

	require.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestAuthHandler_Authenticate_InvalidCredentials(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.GetAllFn = func(ctx context.Context, filter *radcheckdto.RadcheckFilter) ([]radcheckEntity.Radcheck, int64, error) {
		return []radcheckEntity.Radcheck{*testutil.CreateRadcheckFixture()}, 1, nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()

	mockTxManager := &testutil.MockTransactionManager{}
	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo, mockTxManager)
	authHandler := handler.NewAuthHandler(authService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	authHandler.RegisterRoutes(api)

	req := dto.AuthenticateRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/auth/authenticate", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, httpReq)

	require.Equal(t, http.StatusOK, writer.Code)

	var response map[string]interface{}
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	require.False(t, data["success"].(bool))
}

func TestAuthHandler_Authenticate_UserNotFound(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.GetAllFn = func(ctx context.Context, filter *radcheckdto.RadcheckFilter) ([]radcheckEntity.Radcheck, int64, error) {
		return []radcheckEntity.Radcheck{}, 0, nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()

	mockTxManager := &testutil.MockTransactionManager{}
	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo, mockTxManager)
	authHandler := handler.NewAuthHandler(authService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api/v1")
	authHandler.RegisterRoutes(api)

	req := dto.AuthenticateRequest{
		Username: "nonexistent",
		Password: "test123",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/api/v1/auth/authenticate", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	writer := httptest.NewRecorder()
	router.ServeHTTP(writer, httpReq)

	require.Equal(t, http.StatusOK, writer.Code)

	var response map[string]interface{}
	err := json.Unmarshal(writer.Body.Bytes(), &response)
	require.NoError(t, err)

	data := response["data"].(map[string]interface{})
	require.False(t, data["success"].(bool))
	require.Equal(t, "User not found", data["message"].(string))
}

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

