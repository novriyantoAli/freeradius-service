package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/dto"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/testutil"
)

func setupRadcheckHandler() (*RadcheckHandler, *testutil.MockRadcheckService) {
	gin.SetMode(gin.TestMode)
	mockService := &testutil.MockRadcheckService{}
	logger := testutil.NewSilentLogger()
	handler := NewRadcheckHandler(mockService, logger)
	return handler, mockService
}

func TestRadcheckHandler_CreateRadcheck(t *testing.T) {
	t.Run("should create radcheck successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		req := testutil.CreateRadcheckRequestFixture()
		response := &dto.RadcheckResponse{
			ID:        1,
			Username:  req.Username,
			Attribute: req.Attribute,
			Op:        req.Op,
			Value:     req.Value,
		}

		mockService.On("CreateRadcheck", mock.Anything, mock.AnythingOfType("*dto.CreateRadcheckRequest")).Return(response, nil)

		// Prepare request
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/api/v1/radcheck", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// When
		handler.CreateRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Contains(t, result, "data")
		data := result["data"].(map[string]interface{})
		assert.Equal(t, float64(1), data["id"])
		assert.Equal(t, req.Username, data["username"])
	})

	t.Run("should return bad request for invalid JSON", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/api/v1/radcheck", bytes.NewBuffer([]byte("invalid json")))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// When
		handler.CreateRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		req := testutil.CreateRadcheckRequestFixture()
		mockService.On("CreateRadcheck", mock.Anything, mock.AnythingOfType("*dto.CreateRadcheckRequest")).Return(nil, errors.New("database error"))

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("POST", "/api/v1/radcheck", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// When
		handler.CreateRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestRadcheckHandler_GetRadcheck(t *testing.T) {
	t.Run("should get radcheck successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		radcheckID := uint(1)
		response := &dto.RadcheckResponse{
			ID:        radcheckID,
			Username:  "testuser",
			Attribute: "User-Password",
			Op:        ":=",
			Value:     "testpassword",
		}

		mockService.On("GetRadcheckByID", mock.Anything, radcheckID).Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/v1/radcheck/1", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		// When
		handler.GetRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Contains(t, result, "data")
		data := result["data"].(map[string]interface{})
		assert.Equal(t, float64(1), data["id"])
		assert.Equal(t, "testuser", data["username"])
	})

	t.Run("should return bad request for invalid ID", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/v1/radcheck/invalid", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		// When
		handler.GetRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return not found when radcheck not found", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		radcheckID := uint(999)
		mockService.On("GetRadcheckByID", mock.Anything, radcheckID).Return(nil, errors.New("radcheck not found"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/v1/radcheck/999", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "999"},
		}

		// When
		handler.GetRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestRadcheckHandler_ListRadcheck(t *testing.T) {
	t.Run("should list radchecks successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		response := &dto.ListRadcheckResponse{
			Data: []dto.RadcheckResponse{
				{ID: 1, Username: "user1", Attribute: "User-Password", Op: ":=", Value: "pass1"},
				{ID: 2, Username: "user2", Attribute: "User-Password", Op: ":=", Value: "pass2"},
			},
			Total:     2,
			Page:      1,
			PageSize:  10,
			TotalPage: 1,
		}

		mockService.On("ListRadcheck", mock.Anything, mock.AnythingOfType("*dto.RadcheckFilter")).Return(response, nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/v1/radcheck?page=1&page_size=10", nil)

		// When
		handler.ListRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result dto.ListRadcheckResponse
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, int64(2), result.Total)
	})

	t.Run("should return bad request for invalid query parameters", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/v1/radcheck?page=invalid", nil)

		// When
		handler.ListRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		mockService.On("ListRadcheck", mock.Anything, mock.AnythingOfType("*dto.RadcheckFilter")).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("GET", "/api/v1/radcheck?page=1&page_size=10", nil)

		// When
		handler.ListRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestRadcheckHandler_UpdateRadcheck(t *testing.T) {
	t.Run("should update radcheck successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		radcheckID := uint(1)
		req := testutil.CreateUpdateRadcheckRequestFixture()
		response := &dto.RadcheckResponse{
			ID:        radcheckID,
			Username:  "testuser",
			Attribute: "User-Password",
			Op:        ":=",
			Value:     req.Value,
		}

		mockService.On("UpdateRadcheck", mock.Anything, radcheckID, mock.AnythingOfType("*dto.UpdateRadcheckRequest")).Return(response, nil)

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/api/v1/radcheck/1", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		// When
		handler.UpdateRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Contains(t, result, "data")
		data := result["data"].(map[string]interface{})
		assert.Equal(t, float64(1), data["id"])
	})

	t.Run("should return bad request for invalid ID", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		req := testutil.CreateUpdateRadcheckRequestFixture()
		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/api/v1/radcheck/invalid", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		// When
		handler.UpdateRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return bad request for invalid JSON", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/api/v1/radcheck/1", bytes.NewBuffer([]byte("invalid json")))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		// When
		handler.UpdateRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return not found when radcheck not found", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		radcheckID := uint(999)
		req := testutil.CreateUpdateRadcheckRequestFixture()
		mockService.On("UpdateRadcheck", mock.Anything, radcheckID, mock.AnythingOfType("*dto.UpdateRadcheckRequest")).Return(nil, errors.New("radcheck not found"))

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/api/v1/radcheck/999", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "999"},
		}

		// When
		handler.UpdateRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal server error for other errors", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		radcheckID := uint(1)
		req := testutil.CreateUpdateRadcheckRequestFixture()
		mockService.On("UpdateRadcheck", mock.Anything, radcheckID, mock.AnythingOfType("*dto.UpdateRadcheckRequest")).Return(nil, errors.New("database error"))

		reqBody, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("PUT", "/api/v1/radcheck/1", bytes.NewBuffer(reqBody))
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		// When
		handler.UpdateRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestRadcheckHandler_DeleteRadcheck(t *testing.T) {
	t.Run("should delete radcheck successfully", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		radcheckID := uint(1)
		mockService.On("DeleteRadcheck", mock.Anything, radcheckID).Return(nil)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/api/v1/radcheck/1", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		// When
		handler.DeleteRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

		var result map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &result)
		assert.Contains(t, result, "message")
		assert.Equal(t, "Radcheck deleted successfully", result["message"])
	})

	t.Run("should return bad request for invalid ID", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/api/v1/radcheck/invalid", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "invalid"},
		}

		// When
		handler.DeleteRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return not found when radcheck not found", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		radcheckID := uint(999)
		mockService.On("DeleteRadcheck", mock.Anything, radcheckID).Return(errors.New("radcheck not found"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/api/v1/radcheck/999", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "999"},
		}

		// When
		handler.DeleteRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("should return internal server error when delete fails", func(t *testing.T) {
		// Setup
		handler, mockService := setupRadcheckHandler()

		radcheckID := uint(1)
		mockService.On("DeleteRadcheck", mock.Anything, radcheckID).Return(errors.New("database error"))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = httptest.NewRequest("DELETE", "/api/v1/radcheck/1", nil)
		ctx.Params = gin.Params{
			{Key: "id", Value: "1"},
		}

		// When
		handler.DeleteRadcheck(ctx)

		// Then
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
