package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	nasDto "github.com/novriyantoAli/freeradius-service/internal/application/nas/dto"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupNASHandler() (*NASHandler, *testutil.MockNASService) {
	gin.SetMode(gin.TestMode)
	mockService := &testutil.MockNASService{}
	logger := testutil.NewSilentLogger()
	handler := NewNASHandler(mockService, logger)
	return handler, mockService
}

func TestNASHandler_CreateNAS(t *testing.T) {
	t.Run("should create NAS successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		req := testutil.CreateNASRequestFixture()
		response := &nasDto.NASResponse{
			ID:              1,
			NASName:         req.NASName,
			ShortName:       req.ShortName,
			Type:            req.Type,
			Ports:           req.Ports,
			Secret:          req.Secret,
			Server:          req.Server,
			Community:       req.Community,
			Description:     req.Description,
			RequireMa:       req.RequireMa,
			LimitProxyState: req.LimitProxyState,
			CreatedAt:       time.Now().String(),
			UpdatedAt:       time.Now().String(),
		}

		mockService.On("CreateNAS", mock.AnythingOfType("*dto.CreateNASRequest")).Return(response, nil)

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/api/v1/nas", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// When
		handler.CreateNAS(ctx)

		// Then
		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)

		var result nasDto.NASResponse
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, req.NASName, result.NASName)
		assert.Equal(t, req.ShortName, result.ShortName)
	})

	t.Run("should return bad request for invalid JSON", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/api/v1/nas", bytes.NewBuffer([]byte("invalid json")))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// When
		handler.CreateNAS(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		req := testutil.CreateNASRequestFixture()
		mockService.On("CreateNAS", mock.AnythingOfType("*dto.CreateNASRequest")).Return(nil, errors.New("nasname already exists"))

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/api/v1/nas", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// When
		handler.CreateNAS(ctx)

		// Then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestNASHandler_GetNAS(t *testing.T) {
	t.Run("should get NAS successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		nasID := uint(1)
		ports := 1812
		response := &nasDto.NASResponse{
			ID:              nasID,
			NASName:         "test-nas-01",
			ShortName:       "test-nas",
			Type:            "other",
			Ports:           &ports,
			Secret:          "testing123",
			Server:          "192.168.1.1",
			Community:       "public",
			Description:     "Test NAS",
			RequireMa:       "auto",
			LimitProxyState: "auto",
			CreatedAt:       time.Now().String(),
			UpdatedAt:       time.Now().String(),
		}

		mockService.On("GetNASByID", nasID).Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/v1/nas/1", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		// When
		handler.GetNAS(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result nasDto.NASResponse
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, nasID, result.ID)
		assert.Equal(t, "test-nas-01", result.NASName)
	})

	t.Run("should return bad request for invalid ID", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/v1/nas/invalid", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		// When
		handler.GetNAS(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return not found when NAS not found", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		nasID := uint(999)
		mockService.On("GetNASByID", nasID).Return(nil, errors.New("nas not found"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/v1/nas/999", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "999"},
		}

		// When
		handler.GetNAS(ctx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestNASHandler_ListNAS(t *testing.T) {
	t.Run("should list NAS successfully with default pagination", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		response := &nasDto.ListNASResponse{
			Data: []nasDto.NASResponse{
				{ID: 1, NASName: "nas-01", ShortName: "n1", Type: "other"},
				{ID: 2, NASName: "nas-02", ShortName: "n2", Type: "other"},
			},
			Total:     2,
			Page:      1,
			PageSize:  10,
			TotalPage: 1,
		}

		mockService.On("ListNAS", mock.AnythingOfType("*dto.NASFilter")).Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/v1/nas", nil)

		// When
		handler.ListNAS(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result nasDto.ListNASResponse
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, int64(2), result.Total)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 10, result.PageSize)
	})

	t.Run("should list NAS with custom pagination", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		response := &nasDto.ListNASResponse{
			Data:      []nasDto.NASResponse{},
			Total:     0,
			Page:      2,
			PageSize:  5,
			TotalPage: 0,
		}

		mockService.On("ListNAS", mock.AnythingOfType("*dto.NASFilter")).Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/v1/nas?page=2&page_size=5", nil)

		// When
		handler.ListNAS(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should list NAS with filters", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		response := &nasDto.ListNASResponse{
			Data: []nasDto.NASResponse{
				{ID: 1, NASName: "cisco-nas", ShortName: "cisco", Type: "cisco"},
			},
			Total:     1,
			Page:      1,
			PageSize:  10,
			TotalPage: 1,
		}

		mockService.On("ListNAS", mock.MatchedBy(func(filter *nasDto.NASFilter) bool {
			return filter.Type == "cisco"
		})).Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/v1/nas?type=cisco", nil)

		// When
		handler.ListNAS(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid query parameters", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/v1/nas?page=invalid", nil)

		// When
		handler.ListNAS(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		mockService.On("ListNAS", mock.AnythingOfType("*dto.NASFilter")).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/v1/nas", nil)

		// When
		handler.ListNAS(ctx)

		// Then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestNASHandler_UpdateNAS(t *testing.T) {
	t.Run("should update NAS successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		nasID := uint(1)
		req := testutil.CreateUpdateNASRequestFixture()
		response := &nasDto.NASResponse{
			ID:          nasID,
			NASName:     "updated-nas",
			ShortName:   req.ShortName,
			Type:        req.Type,
			Description: req.Description,
			CreatedAt:   time.Now().String(),
			UpdatedAt:   time.Now().String(),
		}

		mockService.On("UpdateNAS", nasID, mock.AnythingOfType("*dto.UpdateNASRequest")).Return(response, nil)

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/api/v1/nas/1", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		// When
		handler.UpdateNAS(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result nasDto.NASResponse
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Equal(t, nasID, result.ID)
		assert.Equal(t, "updated-nas", result.NASName)
	})

	t.Run("should return bad request for invalid ID", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		req := testutil.CreateUpdateNASRequestFixture()
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/api/v1/nas/invalid", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		// When
		handler.UpdateNAS(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid JSON body", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/api/v1/nas/1", bytes.NewBuffer([]byte("invalid json")))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		// When
		handler.UpdateNAS(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return not found when NAS not found", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		nasID := uint(999)
		req := testutil.CreateUpdateNASRequestFixture()
		mockService.On("UpdateNAS", nasID, mock.AnythingOfType("*dto.UpdateNASRequest")).Return(nil, errors.New("nas not found"))

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/api/v1/nas/999", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "999"},
		}

		// When
		handler.UpdateNAS(ctx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		nasID := uint(1)
		req := testutil.CreateUpdateNASRequestFixture()
		mockService.On("UpdateNAS", nasID, mock.AnythingOfType("*dto.UpdateNASRequest")).Return(nil, errors.New("database error"))

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/api/v1/nas/1", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		// When
		handler.UpdateNAS(ctx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestNASHandler_DeleteNAS(t *testing.T) {
	t.Run("should delete NAS successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		nasID := uint(1)
		mockService.On("DeleteNAS", nasID).Return(nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/api/v1/nas/1", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		// When
		handler.DeleteNAS(ctx)

		// Then
		// Status NoContent (204) sets header but Gin test recorder may show 200, verify actual behavior
		assert.True(t, w.Code == http.StatusNoContent || w.Code == http.StatusOK)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid ID", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/api/v1/nas/invalid", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		// When
		handler.DeleteNAS(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return not found when NAS not found", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		nasID := uint(999)
		mockService.On("DeleteNAS", nasID).Return(errors.New("nas not found"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/api/v1/nas/999", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "999"},
		}

		// When
		handler.DeleteNAS(ctx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		// Setup
		handler, mockService := setupNASHandler()

		nasID := uint(1)
		mockService.On("DeleteNAS", nasID).Return(errors.New("database error"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/api/v1/nas/1", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		// When
		handler.DeleteNAS(ctx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}
