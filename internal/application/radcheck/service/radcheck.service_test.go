package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/entity"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/testutil"
)

func TestRadcheckService_CreateRadcheck(t *testing.T) {
	t.Run("should create radcheck successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		req := testutil.CreateRadcheckRequestFixture()

		// Mock expectations
		mockRepo.On("Create", mock.AnythingOfType("*entity.Radcheck")).Return(nil).Run(func(args mock.Arguments) {
			radcheck := args.Get(0).(*entity.Radcheck)
			radcheck.ID = 1
		})

		// When
		response, err := service.CreateRadcheck(req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, req.Username, response.Username)
		assert.Equal(t, req.Attribute, response.Attribute)
		assert.Equal(t, req.Op, response.Op)
		assert.Equal(t, req.Value, response.Value)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when create fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		req := testutil.CreateRadcheckRequestFixture()

		// Mock expectations
		mockRepo.On("Create", mock.AnythingOfType("*entity.Radcheck")).Return(errors.New("create failed"))

		// When
		response, err := service.CreateRadcheck(req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "create failed")
		mockRepo.AssertExpectations(t)
	})
}

func TestRadcheckService_GetRadcheckByID(t *testing.T) {
	t.Run("should get radcheck by ID successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		radcheckID := uint(1)
		radcheck := testutil.CreateRadcheckFixture()
		radcheck.ID = radcheckID

		// Mock expectations
		mockRepo.On("GetByID", radcheckID).Return(radcheck, nil)

		// When
		response, err := service.GetRadcheckByID(radcheckID)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, radcheckID, response.ID)
		assert.Equal(t, radcheck.Username, response.Username)
		assert.Equal(t, radcheck.Attribute, response.Attribute)
		assert.Equal(t, radcheck.Op, response.Op)
		assert.Equal(t, radcheck.Value, response.Value)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when radcheck not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		radcheckID := uint(999)

		// Mock expectations
		mockRepo.On("GetByID", radcheckID).Return(nil, gorm.ErrRecordNotFound)

		// When
		response, err := service.GetRadcheckByID(radcheckID)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "radcheck not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		radcheckID := uint(1)

		// Mock expectations
		mockRepo.On("GetByID", radcheckID).Return(nil, errors.New("database error"))

		// When
		response, err := service.GetRadcheckByID(radcheckID)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestRadcheckService_GetRadcheckByUsernameAndAttribute(t *testing.T) {
	t.Run("should get radcheck by username and attribute successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		username := "testuser"
		attribute := "User-Password"
		radcheck := testutil.CreateRadcheckFixture()
		radcheck.Username = username
		radcheck.Attribute = attribute

		// Mock expectations
		mockRepo.On("GetByUsernameAndAttribute", username, attribute).Return(radcheck, nil)

		// When
		response, err := service.GetRadcheckByUsernameAndAttribute(username, attribute)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, username, response.Username)
		assert.Equal(t, attribute, response.Attribute)
		assert.Equal(t, radcheck.Op, response.Op)
		assert.Equal(t, radcheck.Value, response.Value)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when radcheck not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		username := "nonexistent"
		attribute := "User-Password"

		// Mock expectations
		mockRepo.On("GetByUsernameAndAttribute", username, attribute).Return(nil, gorm.ErrRecordNotFound)

		// When
		response, err := service.GetRadcheckByUsernameAndAttribute(username, attribute)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "radcheck not found")
		mockRepo.AssertExpectations(t)
	})
}

func TestRadcheckService_ListRadcheck(t *testing.T) {
	t.Run("should list radchecks with pagination successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		filter := &dto.RadcheckFilter{
			Page:     1,
			PageSize: 10,
		}

		radchecks := []entity.Radcheck{
			*testutil.CreateRadcheckFixture(),
			*testutil.CreateRadcheckFixture(),
		}
		radchecks[0].ID = 1
		radchecks[1].ID = 2
		radchecks[1].Username = "user2"

		// Mock expectations
		mockRepo.On("GetAll", filter).Return(radchecks, int64(2), nil)

		// When
		response, err := service.ListRadcheck(filter)

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
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		filter := &dto.RadcheckFilter{
			Page:     0,
			PageSize: 0,
		}

		expectedFilter := &dto.RadcheckFilter{
			Page:     1,
			PageSize: 10,
		}

		// Mock expectations
		mockRepo.On("GetAll", expectedFilter).Return([]entity.Radcheck{}, int64(0), nil)

		// When
		response, err := service.ListRadcheck(filter)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, response.Page)
		assert.Equal(t, 10, response.PageSize)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should cap page size at 100", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		filter := &dto.RadcheckFilter{
			Page:     1,
			PageSize: 200,
		}

		expectedFilter := &dto.RadcheckFilter{
			Page:     1,
			PageSize: 100,
		}

		// Mock expectations
		mockRepo.On("GetAll", expectedFilter).Return([]entity.Radcheck{}, int64(0), nil)

		// When
		response, err := service.ListRadcheck(filter)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 100, response.PageSize)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		filter := &dto.RadcheckFilter{
			Page:     1,
			PageSize: 10,
		}

		// Mock expectations
		mockRepo.On("GetAll", filter).Return(nil, int64(0), errors.New("database error"))

		// When
		response, err := service.ListRadcheck(filter)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "database error")
		mockRepo.AssertExpectations(t)
	})
}

func TestRadcheckService_UpdateRadcheck(t *testing.T) {
	t.Run("should update radcheck successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		radcheckID := uint(1)
		existingRadcheck := testutil.CreateRadcheckFixture()
		existingRadcheck.ID = radcheckID
		existingRadcheck.Value = "oldvalue"

		req := testutil.CreateUpdateRadcheckRequestFixture()
		req.Value = "newvalue"

		// Mock expectations
		mockRepo.On("GetByID", radcheckID).Return(existingRadcheck, nil)
		mockRepo.On("Update", mock.AnythingOfType("*entity.Radcheck")).Return(nil)

		// When
		response, err := service.UpdateRadcheck(radcheckID, req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, radcheckID, response.ID)
		assert.Equal(t, "newvalue", response.Value)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when radcheck not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		radcheckID := uint(999)
		req := testutil.CreateUpdateRadcheckRequestFixture()

		// Mock expectations
		mockRepo.On("GetByID", radcheckID).Return(nil, gorm.ErrRecordNotFound)

		// When
		response, err := service.UpdateRadcheck(radcheckID, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "radcheck not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle partial update with empty fields", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		radcheckID := uint(1)
		existingRadcheck := testutil.CreateRadcheckFixture()
		existingRadcheck.ID = radcheckID

		req := &dto.UpdateRadcheckRequest{
			Value: "newvalue",
		}

		// Mock expectations
		mockRepo.On("GetByID", radcheckID).Return(existingRadcheck, nil)
		mockRepo.On("Update", mock.AnythingOfType("*entity.Radcheck")).Return(nil)

		// When
		response, err := service.UpdateRadcheck(radcheckID, req)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "newvalue", response.Value)
		// Username, Attribute, Op should remain unchanged
		assert.Equal(t, existingRadcheck.Username, response.Username)
		assert.Equal(t, existingRadcheck.Attribute, response.Attribute)
		assert.Equal(t, existingRadcheck.Op, response.Op)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		radcheckID := uint(1)
		existingRadcheck := testutil.CreateRadcheckFixture()
		existingRadcheck.ID = radcheckID

		req := testutil.CreateUpdateRadcheckRequestFixture()

		// Mock expectations
		mockRepo.On("GetByID", radcheckID).Return(existingRadcheck, nil)
		mockRepo.On("Update", mock.AnythingOfType("*entity.Radcheck")).Return(errors.New("update failed"))

		// When
		response, err := service.UpdateRadcheck(radcheckID, req)

		// Then
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "update failed")
		mockRepo.AssertExpectations(t)
	})
}

func TestRadcheckService_DeleteRadcheck(t *testing.T) {
	t.Run("should delete radcheck successfully", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		radcheckID := uint(1)
		radcheck := testutil.CreateRadcheckFixture()
		radcheck.ID = radcheckID

		// Mock expectations
		mockRepo.On("GetByID", radcheckID).Return(radcheck, nil)
		mockRepo.On("Delete", radcheckID).Return(nil)

		// When
		err := service.DeleteRadcheck(radcheckID)

		// Then
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when radcheck not found", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		radcheckID := uint(999)

		// Mock expectations
		mockRepo.On("GetByID", radcheckID).Return(nil, gorm.ErrRecordNotFound)

		// When
		err := service.DeleteRadcheck(radcheckID)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "radcheck not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when delete fails", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger)

		radcheckID := uint(1)
		radcheck := testutil.CreateRadcheckFixture()
		radcheck.ID = radcheckID

		// Mock expectations
		mockRepo.On("GetByID", radcheckID).Return(radcheck, nil)
		mockRepo.On("Delete", radcheckID).Return(errors.New("delete failed"))

		// When
		err := service.DeleteRadcheck(radcheckID)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
		mockRepo.AssertExpectations(t)
	})
}

func TestRadcheckService_entityToResponse(t *testing.T) {
	t.Run("should convert entity to response correctly", func(t *testing.T) {
		// Setup
		mockRepo := &testutil.MockRadcheckRepository{}
		logger := testutil.NewSilentLogger()
		service := NewRadcheckService(mockRepo, logger).(*radcheckService)

		radcheck := testutil.CreateRadcheckFixture()
		radcheck.ID = 1
		radcheck.Username = "testuser"
		radcheck.Attribute = "User-Password"
		radcheck.Op = ":="
		radcheck.Value = "testpassword"

		// When
		response := service.entityToResponse(radcheck)

		// Then
		assert.Equal(t, radcheck.ID, response.ID)
		assert.Equal(t, radcheck.Username, response.Username)
		assert.Equal(t, radcheck.Attribute, response.Attribute)
		assert.Equal(t, radcheck.Op, response.Op)
		assert.Equal(t, radcheck.Value, response.Value)
	})
}
