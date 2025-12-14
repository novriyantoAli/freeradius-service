package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/entity"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestRadreplyRepository_Create(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	repo := NewRadreplyRepository(db)

	radreply := &entity.Radreply{
		Username:  "john",
		Attribute: "Reply-Message",
		Op:        "=",
		Value:     "Welcome",
	}

	err = repo.Create(context.Background(), radreply)

	assert.NoError(t, err)
	assert.NotZero(t, radreply.ID)
}

func TestRadreplyRepository_GetByID(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	repo := NewRadreplyRepository(db)

	radreply := testutil.CreateRadreplyFixture()
	db.Create(radreply)

	result, err := repo.GetByID(context.Background(), radreply.ID)

	assert.NoError(t, err)
	assert.Equal(t, radreply.ID, result.ID)
	assert.Equal(t, radreply.Username, result.Username)
}

func TestRadreplyRepository_GetByID_NotFound(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	repo := NewRadreplyRepository(db)

	result, err := repo.GetByID(context.Background(), 9999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestRadreplyRepository_GetByUsernameAndAttribute(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	repo := NewRadreplyRepository(db)

	radreply := testutil.CreateRadreplyFixture()
	db.Create(radreply)

	result, err := repo.GetByUsernameAndAttribute(context.Background(), radreply.Username, radreply.Attribute)

	assert.NoError(t, err)
	assert.Equal(t, radreply.Username, result.Username)
	assert.Equal(t, radreply.Attribute, result.Attribute)
}

func TestRadreplyRepository_GetByUsernameAndAttribute_NotFound(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	repo := NewRadreplyRepository(db)

	result, err := repo.GetByUsernameAndAttribute(context.Background(), "nonexistent", "nonexistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestRadreplyRepository_GetAll(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	repo := NewRadreplyRepository(db)

	radreply1 := &entity.Radreply{
		Username:  "user1",
		Attribute: "Reply-Message",
		Op:        "=",
		Value:     "Welcome 1",
	}
	radreply2 := &entity.Radreply{
		Username:  "user2",
		Attribute: "Reply-Message",
		Op:        "=",
		Value:     "Welcome 2",
	}
	db.Create(radreply1)
	db.Create(radreply2)

	filter := &dto.RadreplyFilter{
		Page:     1,
		PageSize: 10,
	}

	result, total, err := repo.GetAll(context.Background(), filter)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
}

func TestRadreplyRepository_GetAll_WithFilter(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	repo := NewRadreplyRepository(db)

	radreply1 := testutil.CreateRadreplyFixture()
	radreply2 := testutil.CreateRadreplyFixture()
	radreply2.Username = "jane"
	db.Create(radreply1)
	db.Create(radreply2)

	filter := &dto.RadreplyFilter{
		Username: "john",
		Page:     1,
		PageSize: 10,
	}

	result, total, err := repo.GetAll(context.Background(), filter)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
}

func TestRadreplyRepository_GetAll_Pagination(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	repo := NewRadreplyRepository(db)

	for i := 0; i < 15; i++ {
		radreply := &entity.Radreply{
			Username:  fmt.Sprintf("user%d", i),
			Attribute: "Reply-Message",
			Op:        "=",
			Value:     "Welcome",
		}
		db.Create(radreply)
	}

	filter := &dto.RadreplyFilter{
		Page:     2,
		PageSize: 10,
	}

	result, total, err := repo.GetAll(context.Background(), filter)

	assert.NoError(t, err)
	assert.Equal(t, int64(15), total)
	assert.Len(t, result, 5)
}

func TestRadreplyRepository_GetAll_Empty(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	repo := NewRadreplyRepository(db)

	filter := &dto.RadreplyFilter{
		Page:     1,
		PageSize: 10,
	}

	result, total, err := repo.GetAll(context.Background(), filter)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Len(t, result, 0)
}

func TestRadreplyRepository_Update(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	repo := NewRadreplyRepository(db)

	radreply := testutil.CreateRadreplyFixture()
	db.Create(radreply)

	radreply.Value = "Updated Value"

	err = repo.Update(context.Background(), radreply)

	assert.NoError(t, err)

	updated, _ := repo.GetByID(context.Background(), radreply.ID)
	assert.Equal(t, "Updated Value", updated.Value)
}

func TestRadreplyRepository_Delete(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	repo := NewRadreplyRepository(db)

	radreply := testutil.CreateRadreplyFixture()
	db.Create(radreply)

	err = repo.Delete(context.Background(), radreply.ID)

	assert.NoError(t, err)

	result, err := repo.GetByID(context.Background(), radreply.ID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
