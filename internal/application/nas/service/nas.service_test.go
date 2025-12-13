package service

import (
	"errors"
	"testing"

	nasDto "github.com/novriyantoAli/freeradius-service/internal/application/nas/dto"
	nasEntity "github.com/novriyantoAli/freeradius-service/internal/application/nas/entity"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestNASService_CreateNAS(t *testing.T) {
	t.Run("should create NAS successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		req := testutil.CreateNASRequestFixture()

		// Mock expectations
		mockRepo.On("GetByNASName", req.NASName).Return(nil, gorm.ErrRecordNotFound)
		mockRepo.On("Create", mock.AnythingOfType("*entity.NAS")).Return(nil).Run(func(args mock.Arguments) {
			nas := args.Get(0).(*nasEntity.NAS)
			nas.ID = 1
		})

		// When
		response, err := service.CreateNAS(req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, req.NASName, response.NASName)
		assert.Equal(t, req.ShortName, response.ShortName)
		assert.Equal(t, req.Type, response.Type)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when nasname already exists", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		req := testutil.CreateNASRequestFixture()
		existingNAS := testutil.CreateNASFixture()

		// Mock expectations
		mockRepo.On("GetByNASName", req.NASName).Return(existingNAS, nil)

		// When
		response, err := service.CreateNAS(req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "nasname already exists")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when nasname check fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		req := testutil.CreateNASRequestFixture()

		// Mock expectations
		mockRepo.On("GetByNASName", req.NASName).Return(nil, errors.New("database error"))

		// When
		response, err := service.CreateNAS(req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when NAS creation fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		req := testutil.CreateNASRequestFixture()

		// Mock expectations
		mockRepo.On("GetByNASName", req.NASName).Return(nil, gorm.ErrRecordNotFound)
		mockRepo.On("Create", mock.AnythingOfType("*entity.NAS")).Return(errors.New("create failed"))

		// When
		response, err := service.CreateNAS(req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "create failed")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle nil ports pointer", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		req := testutil.CreateNASRequestFixture()
		req.Ports = nil // Explicitly set to nil

		// Mock expectations
		mockRepo.On("GetByNASName", req.NASName).Return(nil, gorm.ErrRecordNotFound)
		mockRepo.On("Create", mock.AnythingOfType("*entity.NAS")).Return(nil).Run(func(args mock.Arguments) {
			nas := args.Get(0).(*nasEntity.NAS)
			nas.ID = 1
			assert.Equal(t, 0, nas.Ports) // Should be 0 when nil pointer
		})

		// When
		response, err := service.CreateNAS(req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		mockRepo.AssertExpectations(t)
	})
}

func TestNASService_GetNASByID(t *testing.T) {
	t.Run("should get NAS by ID successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		nasID := uint(1)
		nas := testutil.CreateNASFixture()
		nas.ID = nasID

		// Mock expectations
		mockRepo.On("GetByID", nasID).Return(nas, nil)

		// When
		response, err := service.GetNASByID(nasID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, nasID, response.ID)
		assert.Equal(t, nas.NASName, response.NASName)
		assert.Equal(t, nas.ShortName, response.ShortName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when NAS not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		nasID := uint(999)

		// Mock expectations
		mockRepo.On("GetByID", nasID).Return(nil, gorm.ErrRecordNotFound)

		// When
		response, err := service.GetNASByID(nasID)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "nas not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		nasID := uint(1)

		// Mock expectations
		mockRepo.On("GetByID", nasID).Return(nil, errors.New("database error"))

		// When
		response, err := service.GetNASByID(nasID)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestNASService_ListNAS(t *testing.T) {
	t.Run("should list NAS with pagination successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		filter := &nasDto.NASFilter{
			Page:     1,
			PageSize: 10,
		}

		nasList := []nasEntity.NAS{
			*testutil.CreateNASFixture(),
			*testutil.CreateNASFixture(),
		}
		nasList[0].ID = 1
		nasList[1].ID = 2
		nasList[1].NASName = "nas-02"

		// Mock expectations
		mockRepo.On("GetAll", filter).Return(nasList, int64(2), nil)

		// When
		response, err := service.ListNAS(filter)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Data, 2)
		assert.Equal(t, int64(2), response.Total)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.PageSize)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should set default pagination values", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		filter := &nasDto.NASFilter{
			Page:     0,
			PageSize: 0,
		}

		expectedFilter := &nasDto.NASFilter{
			Page:     1,
			PageSize: 10,
		}

		// Mock expectations
		mockRepo.On("GetAll", expectedFilter).Return([]nasEntity.NAS{}, int64(0), nil)

		// When
		response, err := service.ListNAS(filter)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.PageSize)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should cap page size at 100", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		filter := &nasDto.NASFilter{
			Page:     1,
			PageSize: 200,
		}

		expectedFilter := &nasDto.NASFilter{
			Page:     1,
			PageSize: 10,
		}

		// Mock expectations
		mockRepo.On("GetAll", expectedFilter).Return([]nasEntity.NAS{}, int64(0), nil)

		// When
		response, err := service.ListNAS(filter)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 10, response.PageSize)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should calculate total pages correctly", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		filter := &nasDto.NASFilter{
			Page:     1,
			PageSize: 3,
		}

		nasList := []nasEntity.NAS{
			*testutil.CreateNASFixture(),
			*testutil.CreateNASFixture(),
			*testutil.CreateNASFixture(),
		}

		// Mock expectations - total 10 items, page size 3 = 4 total pages (3+3+3+1)
		mockRepo.On("GetAll", filter).Return(nasList, int64(10), nil)

		// When
		response, err := service.ListNAS(filter)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 4, response.TotalPage)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		filter := &nasDto.NASFilter{
			Page:     1,
			PageSize: 10,
		}

		// Mock expectations
		mockRepo.On("GetAll", filter).Return(nil, int64(0), errors.New("database error"))

		// When
		response, err := service.ListNAS(filter)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestNASService_UpdateNAS(t *testing.T) {
	t.Run("should update NAS successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		nasID := uint(1)
		existingNAS := testutil.CreateNASFixture()
		existingNAS.ID = nasID
		existingNAS.ShortName = "old-short"
		existingNAS.Description = "Old description"

		req := testutil.CreateUpdateNASRequestFixture()
		req.ShortName = "new-short"
		req.Description = "New description"

		// Mock expectations
		mockRepo.On("GetByID", nasID).Return(existingNAS, nil)
		mockRepo.On("Update", mock.AnythingOfType("*entity.NAS")).Return(nil)

		// When
		response, err := service.UpdateNAS(nasID, req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, nasID, response.ID)
		assert.Equal(t, "new-short", response.ShortName)
		assert.Equal(t, "New description", response.Description)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when NAS not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		nasID := uint(999)
		req := testutil.CreateUpdateNASRequestFixture()

		// Mock expectations
		mockRepo.On("GetByID", nasID).Return(nil, gorm.ErrRecordNotFound)

		// When
		response, err := service.UpdateNAS(nasID, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "nas not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should only update provided fields", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		nasID := uint(1)
		existingNAS := testutil.CreateNASFixture()
		existingNAS.ID = nasID
		existingNAS.ShortName = "original-short"
		existingNAS.Type = "original-type"

		req := &nasDto.UpdateNASRequest{
			ShortName: "updated-short",
		}

		// Mock expectations
		mockRepo.On("GetByID", nasID).Return(existingNAS, nil)
		mockRepo.On("Update", mock.MatchedBy(func(nas *nasEntity.NAS) bool {
			return nas.ShortName == "updated-short" && nas.Type == "original-type"
		})).Return(nil)

		// When
		response, err := service.UpdateNAS(nasID, req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "updated-short", response.ShortName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails to get NAS", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		nasID := uint(1)
		req := testutil.CreateUpdateNASRequestFixture()

		// Mock expectations
		mockRepo.On("GetByID", nasID).Return(nil, errors.New("database error"))

		// When
		response, err := service.UpdateNAS(nasID, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		nasID := uint(1)
		existingNAS := testutil.CreateNASFixture()
		existingNAS.ID = nasID

		req := testutil.CreateUpdateNASRequestFixture()

		// Mock expectations
		mockRepo.On("GetByID", nasID).Return(existingNAS, nil)
		mockRepo.On("Update", mock.AnythingOfType("*entity.NAS")).Return(errors.New("update failed"))

		// When
		response, err := service.UpdateNAS(nasID, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "update failed")
		mockRepo.AssertExpectations(t)
	})
}

func TestNASService_DeleteNAS(t *testing.T) {
	t.Run("should delete NAS successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		nasID := uint(1)
		nas := testutil.CreateNASFixture()
		nas.ID = nasID

		// Mock expectations
		mockRepo.On("GetByID", nasID).Return(nas, nil)
		mockRepo.On("Delete", nasID).Return(nil)

		// When
		err := service.DeleteNAS(nasID)

		// Then
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when NAS not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		nasID := uint(999)

		// Mock expectations
		mockRepo.On("GetByID", nasID).Return(nil, gorm.ErrRecordNotFound)

		// When
		err := service.DeleteNAS(nasID)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nas not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails to get NAS", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		nasID := uint(1)

		// Mock expectations
		mockRepo.On("GetByID", nasID).Return(nil, errors.New("database error"))

		// When
		err := service.DeleteNAS(nasID)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when delete fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockNASRepository{}
		logger := testutil.NewSilentLogger()
		service := NewNASService(mockRepo, logger)

		nasID := uint(1)
		nas := testutil.CreateNASFixture()
		nas.ID = nasID

		// Mock expectations
		mockRepo.On("GetByID", nasID).Return(nas, nil)
		mockRepo.On("Delete", nasID).Return(errors.New("delete failed"))

		// When
		err := service.DeleteNAS(nasID)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
		mockRepo.AssertExpectations(t)
	})
}
