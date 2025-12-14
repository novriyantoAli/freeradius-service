package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/dto"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRadreplyHandler_CreateRadreply(t *testing.T) {
	t.Run("should create radreply successfully", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		req := testutil.CreateRadreplyRequestFixture()
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("POST", "/api/v1/radreply", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("should return bad request for invalid json", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("POST", "/api/v1/radreply", bytes.NewBuffer([]byte("invalid")))
		httpReq.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return internal error when service fails", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		service.CreateRadreplyFn = func(ctx context.Context, req *dto.CreateRadreplyRequest) (*dto.RadreplyResponse, error) {
			return nil, gorm.ErrInvalidDB
		}
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		req := testutil.CreateRadreplyRequestFixture()
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("POST", "/api/v1/radreply", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRadreplyHandler_GetRadreply(t *testing.T) {
	t.Run("should get radreply successfully", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		service.GetRadreplyByIDFn = func(ctx context.Context, id uint) (*dto.RadreplyResponse, error) {
			return &dto.RadreplyResponse{
				ID:        1,
				Username:  "john",
				Attribute: "Reply-Message",
				Op:        "=",
				Value:     "Welcome",
			}, nil
		}
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("GET", "/api/v1/radreply/1", nil)
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return not found when radreply not found", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		service.GetRadreplyByIDFn = func(ctx context.Context, id uint) (*dto.RadreplyResponse, error) {
			return nil, errors.New("radreply not found")
		}
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("GET", "/api/v1/radreply/9999", nil)
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should return bad request for invalid id", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("GET", "/api/v1/radreply/invalid", nil)
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRadreplyHandler_ListRadreply(t *testing.T) {
	t.Run("should list radreply successfully", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		service.ListRadreplyFn = func(ctx context.Context, filter *dto.RadreplyFilter) (*dto.ListRadreplyResponse, error) {
			return &dto.ListRadreplyResponse{
				Data: []dto.RadreplyResponse{
					{
						ID:        1,
						Username:  "john",
						Attribute: "Reply-Message",
						Op:        "=",
						Value:     "Welcome",
					},
				},
				Total:     1,
				Page:      1,
				PageSize:  10,
				TotalPage: 1,
			}, nil
		}
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("GET", "/api/v1/radreply?page=1&page_size=10", nil)
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return bad request for invalid query", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("GET", "/api/v1/radreply?page=invalid", nil)
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return internal error when service fails", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		service.ListRadreplyFn = func(ctx context.Context, filter *dto.RadreplyFilter) (*dto.ListRadreplyResponse, error) {
			return nil, gorm.ErrInvalidDB
		}
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("GET", "/api/v1/radreply?page=1&page_size=10", nil)
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRadreplyHandler_UpdateRadreply(t *testing.T) {
	t.Run("should update radreply successfully", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		service.UpdateRadreplyFn = func(ctx context.Context, id uint, req *dto.UpdateRadreplyRequest) (*dto.RadreplyResponse, error) {
			return &dto.RadreplyResponse{
				ID:        1,
				Username:  "john",
				Attribute: "Reply-Message",
				Op:        "=",
				Value:     "Updated",
			}, nil
		}
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		req := testutil.CreateUpdateRadreplyRequestFixture()
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("PUT", "/api/v1/radreply/1", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return not found when radreply not found", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		service.UpdateRadreplyFn = func(ctx context.Context, id uint, req *dto.UpdateRadreplyRequest) (*dto.RadreplyResponse, error) {
			return nil, errors.New("radreply not found")
		}
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		req := testutil.CreateUpdateRadreplyRequestFixture()
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("PUT", "/api/v1/radreply/9999", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should return bad request for invalid id", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		req := testutil.CreateUpdateRadreplyRequestFixture()
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("PUT", "/api/v1/radreply/invalid", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return bad request for invalid json", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("PUT", "/api/v1/radreply/1", bytes.NewBuffer([]byte("invalid")))
		httpReq.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRadreplyHandler_DeleteRadreply(t *testing.T) {
	t.Run("should delete radreply successfully", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		service.DeleteRadreplyFn = func(ctx context.Context, id uint) error {
			return nil
		}
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("DELETE", "/api/v1/radreply/1", nil)
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("should return bad request for invalid id", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("DELETE", "/api/v1/radreply/invalid", nil)
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return internal error when service fails", func(t *testing.T) {
		service := testutil.NewMockRadreplyService()
		service.DeleteRadreplyFn = func(ctx context.Context, id uint) error {
			return gorm.ErrInvalidDB
		}
		handler := NewRadreplyHandler(service, testutil.NewSilentLogger())

		w := httptest.NewRecorder()
		r := gin.New()
		apiGroup := r.Group("/api/v1")
		handler.RegisterRoutes(apiGroup)

		httpReq, _ := http.NewRequest("DELETE", "/api/v1/radreply/1", nil)
		r.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
