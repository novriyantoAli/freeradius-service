package service

import (
	"context"
	"testing"

	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/entity"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRadreplyService_CreateRadreply(t *testing.T) {
	t.Run("should create radreply successfully", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		req := testutil.CreateRadreplyRequestFixture()

		result, err := service.CreateRadreply(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, req.Username, result.Username)
		assert.Equal(t, req.Attribute, result.Attribute)
	})

	t.Run("should fail when repository fails", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		repo.CreateFn = func(ctx context.Context, radreply *entity.Radreply) error {
			return gorm.ErrInvalidDB
		}
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		req := testutil.CreateRadreplyRequestFixture()

		result, err := service.CreateRadreply(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should fail on validation error - empty username", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		req := &dto.CreateRadreplyRequest{
			Username:  "",
			Attribute: "Reply-Message",
			Op:        "=",
			Value:     "Welcome",
		}

		result, err := service.CreateRadreply(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "username is required")
	})

	t.Run("should fail on validation error - username too long", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		req := &dto.CreateRadreplyRequest{
			Username:  "a" + string(make([]byte, 65)),
			Attribute: "Reply-Message",
			Op:        "=",
			Value:     "Welcome",
		}

		result, err := service.CreateRadreply(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should fail on validation error - empty value", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		req := &dto.CreateRadreplyRequest{
			Username:  "john",
			Attribute: "Reply-Message",
			Op:        "=",
			Value:     "",
		}

		result, err := service.CreateRadreply(context.Background(), req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "value is required")
	})
}

func TestRadreplyService_GetRadreplyByID(t *testing.T) {
	t.Run("should get radreply by id successfully", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		fixture := testutil.CreateRadreplyFixture()
		repo.GetByIDFn = func(ctx context.Context, id uint) (*entity.Radreply, error) {
			return fixture, nil
		}

		result, err := service.GetRadreplyByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, fixture.Username, result.Username)
	})

	t.Run("should fail when radreply not found", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		repo.GetByIDFn = func(ctx context.Context, id uint) (*entity.Radreply, error) {
			return nil, gorm.ErrRecordNotFound
		}
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		result, err := service.GetRadreplyByID(context.Background(), 9999)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestRadreplyService_GetRadreplyByUsernameAndAttribute(t *testing.T) {
	t.Run("should get radreply by username and attribute successfully", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		fixture := testutil.CreateRadreplyFixture()
		repo.GetByUsernameAndAttributeFn = func(ctx context.Context, username, attribute string) (*entity.Radreply, error) {
			return fixture, nil
		}

		result, err := service.GetRadreplyByUsernameAndAttribute(context.Background(), "john", "Reply-Message")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, fixture.Username, result.Username)
	})

	t.Run("should fail when record not found", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		repo.GetByUsernameAndAttributeFn = func(ctx context.Context, username, attribute string) (*entity.Radreply, error) {
			return nil, gorm.ErrRecordNotFound
		}
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		result, err := service.GetRadreplyByUsernameAndAttribute(context.Background(), "nonexistent", "nonexistent")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestRadreplyService_ListRadreply(t *testing.T) {
	t.Run("should list radreply with pagination", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		fixtures := []entity.Radreply{
			*testutil.CreateRadreplyFixture(),
			*testutil.CreateRadreplyFixture(),
		}
		repo.GetAllFn = func(ctx context.Context, filter *dto.RadreplyFilter) ([]entity.Radreply, int64, error) {
			return fixtures, int64(len(fixtures)), nil
		}

		filter := &dto.RadreplyFilter{Page: 1, PageSize: 10}
		result, err := service.ListRadreply(context.Background(), filter)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(2), result.Total)
		assert.Len(t, result.Data, 2)
	})

	t.Run("should apply default pagination when not provided", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		repo.GetAllFn = func(ctx context.Context, filter *dto.RadreplyFilter) ([]entity.Radreply, int64, error) {
			assert.Equal(t, 1, filter.Page)
			assert.Equal(t, 10, filter.PageSize)
			return []entity.Radreply{}, 0, nil
		}

		result, err := service.ListRadreply(context.Background(), &dto.RadreplyFilter{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("should cap page size at 100", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		repo.GetAllFn = func(ctx context.Context, filter *dto.RadreplyFilter) ([]entity.Radreply, int64, error) {
			assert.Equal(t, 100, filter.PageSize)
			return []entity.Radreply{}, 0, nil
		}

		result, err := service.ListRadreply(context.Background(), &dto.RadreplyFilter{PageSize: 200})

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("should fail when repository fails", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		repo.GetAllFn = func(ctx context.Context, filter *dto.RadreplyFilter) ([]entity.Radreply, int64, error) {
			return nil, 0, gorm.ErrInvalidDB
		}
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		result, err := service.ListRadreply(context.Background(), &dto.RadreplyFilter{})

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestRadreplyService_UpdateRadreply(t *testing.T) {
	t.Run("should update radreply successfully", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		fixture := testutil.CreateRadreplyFixture()
		repo.GetByIDFn = func(ctx context.Context, id uint) (*entity.Radreply, error) {
			return fixture, nil
		}
		repo.UpdateFn = func(ctx context.Context, radreply *entity.Radreply) error {
			return nil
		}

		req := testutil.CreateUpdateRadreplyRequestFixture()
		result, err := service.UpdateRadreply(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("should fail when radreply not found", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		repo.GetByIDFn = func(ctx context.Context, id uint) (*entity.Radreply, error) {
			return nil, gorm.ErrRecordNotFound
		}
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		req := testutil.CreateUpdateRadreplyRequestFixture()
		result, err := service.UpdateRadreply(context.Background(), 9999, req)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("should fail when repository update fails", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		fixture := testutil.CreateRadreplyFixture()
		repo.GetByIDFn = func(ctx context.Context, id uint) (*entity.Radreply, error) {
			return fixture, nil
		}
		repo.UpdateFn = func(ctx context.Context, radreply *entity.Radreply) error {
			return gorm.ErrInvalidDB
		}

		req := testutil.CreateUpdateRadreplyRequestFixture()
		result, err := service.UpdateRadreply(context.Background(), 1, req)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestRadreplyService_DeleteRadreply(t *testing.T) {
	t.Run("should delete radreply successfully", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		repo.DeleteFn = func(ctx context.Context, id uint) error {
			return nil
		}

		err := service.DeleteRadreply(context.Background(), 1)

		assert.NoError(t, err)
	})

	t.Run("should fail when repository fails", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		repo.DeleteFn = func(ctx context.Context, id uint) error {
			return gorm.ErrInvalidDB
		}
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		err := service.DeleteRadreply(context.Background(), 1)

		assert.Error(t, err)
	})
}

func TestRadreplyService_entityToResponse(t *testing.T) {
	t.Run("should convert entity to response correctly", func(t *testing.T) {
		repo := testutil.NewMockRadreplyRepository()
		service := NewRadreplyService(repo, testutil.NewSilentLogger())

		entity := &entity.Radreply{
			ID:        1,
			Username:  "john",
			Attribute: "Reply-Message",
			Op:        "=",
			Value:     "Welcome",
		}

		response := service.(*radreplyService).entityToResponse(entity)

		assert.Equal(t, entity.ID, response.ID)
		assert.Equal(t, entity.Username, response.Username)
		assert.Equal(t, entity.Attribute, response.Attribute)
		assert.Equal(t, entity.Op, response.Op)
		assert.Equal(t, entity.Value, response.Value)
	})
}
