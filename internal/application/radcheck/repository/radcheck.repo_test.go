package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/entity"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/testutil"
)

func TestRadcheckRepository_Create(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewRadcheckRepository(db, logger)

	t.Run("should create radcheck successfully", func(t *testing.T) {
		// Given
		radcheck := testutil.CreateRadcheckFixture()
		radcheck.ID = 0

		// When
		err := repo.Create(radcheck)

		// Then
		require.NoError(t, err)
		assert.NotZero(t, radcheck.ID)

		// Verify in database
		saved := &entity.Radcheck{}
		result := db.First(saved, radcheck.ID)
		require.NoError(t, result.Error)
		assert.Equal(t, radcheck.Username, saved.Username)
		assert.Equal(t, radcheck.Attribute, saved.Attribute)
		assert.Equal(t, radcheck.Op, saved.Op)
		assert.Equal(t, radcheck.Value, saved.Value)
	})

	t.Run("should create radcheck with different values", func(t *testing.T) {
		// Given
		radcheck := &entity.Radcheck{
			Username:  "anotheruser",
			Attribute: "Cleartext-Password",
			Op:        ":=",
			Value:     "secret123",
		}

		// When
		err := repo.Create(radcheck)

		// Then
		require.NoError(t, err)
		assert.NotZero(t, radcheck.ID)

		// Verify in database
		saved := &entity.Radcheck{}
		result := db.First(saved, radcheck.ID)
		require.NoError(t, result.Error)
		assert.Equal(t, "anotheruser", saved.Username)
		assert.Equal(t, "Cleartext-Password", saved.Attribute)
	})
}

func TestRadcheckRepository_GetByID(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewRadcheckRepository(db, logger)

	t.Run("should get radcheck by id", func(t *testing.T) {
		// Given
		fixture := testutil.CreateRadcheckFixture()
		fixture.ID = 0
		_ = repo.Create(fixture)

		// When
		radcheck, err := repo.GetByID(fixture.ID)

		// Then
		require.NoError(t, err)
		assert.NotNil(t, radcheck)
		assert.Equal(t, fixture.Username, radcheck.Username)
		assert.Equal(t, fixture.Attribute, radcheck.Attribute)
	})

	t.Run("should return error when radcheck not found", func(t *testing.T) {
		// When
		radcheck, err := repo.GetByID(9999)

		// Then
		assert.Error(t, err)
		assert.Nil(t, radcheck)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestRadcheckRepository_GetByUsernameAndAttribute(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewRadcheckRepository(db, logger)

	t.Run("should get radcheck by username and attribute", func(t *testing.T) {
		// Given
		fixture := testutil.CreateRadcheckFixture()
		fixture.ID = 0
		_ = repo.Create(fixture)

		// When
		radcheck, err := repo.GetByUsernameAndAttribute(fixture.Username, fixture.Attribute)

		// Then
		require.NoError(t, err)
		assert.NotNil(t, radcheck)
		assert.Equal(t, fixture.Username, radcheck.Username)
		assert.Equal(t, fixture.Attribute, radcheck.Attribute)
		assert.Equal(t, fixture.Value, radcheck.Value)
	})

	t.Run("should return error when radcheck not found by username and attribute", func(t *testing.T) {
		// When
		radcheck, err := repo.GetByUsernameAndAttribute("nonexistent", "Attribute")

		// Then
		assert.Error(t, err)
		assert.Nil(t, radcheck)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}

func TestRadcheckRepository_GetAll(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewRadcheckRepository(db, logger)

	t.Run("should get all radchecks with pagination", func(t *testing.T) {
		// Given
		for i := 0; i < 15; i++ {
			radcheck := &entity.Radcheck{
				Username:  "user" + string(rune(i)),
				Attribute: "User-Password",
				Op:        ":=",
				Value:     "password" + string(rune(i)),
			}
			_ = repo.Create(radcheck)
		}

		filter := &dto.RadcheckFilter{
			Page:     1,
			PageSize: 10,
		}

		// When
		radchecks, total, err := repo.GetAll(filter)

		// Then
		require.NoError(t, err)
		assert.NotNil(t, radchecks)
		assert.Equal(t, 10, len(radchecks))
		assert.Equal(t, int64(15), total)
	})

	t.Run("should filter radchecks by username", func(t *testing.T) {
		// Given - clean and create test data
		testutil.CleanDB(db)

		radcheck1 := &entity.Radcheck{
			Username:  "alice",
			Attribute: "User-Password",
			Op:        ":=",
			Value:     "pass123",
		}
		radcheck2 := &entity.Radcheck{
			Username:  "bob",
			Attribute: "User-Password",
			Op:        ":=",
			Value:     "pass456",
		}
		_ = repo.Create(radcheck1)
		_ = repo.Create(radcheck2)

		filter := &dto.RadcheckFilter{
			Username: "alice",
			Page:     1,
			PageSize: 10,
		}

		// When
		radchecks, total, err := repo.GetAll(filter)

		// Then
		require.NoError(t, err)
		assert.Equal(t, 1, len(radchecks))
		assert.Equal(t, int64(1), total)
		assert.Equal(t, "alice", radchecks[0].Username)
	})

	t.Run("should filter radchecks by attribute", func(t *testing.T) {
		// Given - clean and create test data
		testutil.CleanDB(db)

		radcheck1 := &entity.Radcheck{
			Username:  "testuser",
			Attribute: "User-Password",
			Op:        ":=",
			Value:     "pass123",
		}
		radcheck2 := &entity.Radcheck{
			Username:  "testuser",
			Attribute: "Cleartext-Password",
			Op:        ":=",
			Value:     "pass456",
		}
		_ = repo.Create(radcheck1)
		_ = repo.Create(radcheck2)

		filter := &dto.RadcheckFilter{
			Attribute: "User-Password",
			Page:      1,
			PageSize:  10,
		}

		// When
		radchecks, total, err := repo.GetAll(filter)

		// Then
		require.NoError(t, err)
		assert.Equal(t, 1, len(radchecks))
		assert.Equal(t, int64(1), total)
		assert.Equal(t, "User-Password", radchecks[0].Attribute)
	})

	t.Run("should handle empty results", func(t *testing.T) {
		// Given - clean database
		testutil.CleanDB(db)

		filter := &dto.RadcheckFilter{
			Page:     1,
			PageSize: 10,
		}

		// When
		radchecks, total, err := repo.GetAll(filter)

		// Then
		require.NoError(t, err)
		assert.Equal(t, 0, len(radchecks))
		assert.Equal(t, int64(0), total)
	})
}

func TestRadcheckRepository_Update(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewRadcheckRepository(db, logger)

	t.Run("should update radcheck", func(t *testing.T) {
		// Given
		fixture := testutil.CreateRadcheckFixture()
		fixture.ID = 0
		_ = repo.Create(fixture)

		// Update values
		fixture.Value = "updatedpassword"
		fixture.Op = "+="

		// When
		err := repo.Update(fixture)

		// Then
		require.NoError(t, err)

		// Verify update in database
		updated := &entity.Radcheck{}
		result := db.First(updated, fixture.ID)
		require.NoError(t, result.Error)
		assert.Equal(t, "updatedpassword", updated.Value)
		assert.Equal(t, "+=", updated.Op)
	})

	t.Run("should handle update of non-existent radcheck", func(t *testing.T) {
		// Given
		radcheck := &entity.Radcheck{
			ID:        9999,
			Username:  "nonexistent",
			Attribute: "User-Password",
			Op:        ":=",
			Value:     "test",
		}

		// When
		err := repo.Update(radcheck)

		// Then
		// GORM doesn't error on update of non-existent record, so we just verify no error
		assert.NoError(t, err)
	})
}

func TestRadcheckRepository_Delete(t *testing.T) {
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	defer testutil.CleanDB(db)

	logger := testutil.NewTestLogger(t)
	repo := NewRadcheckRepository(db, logger)

	t.Run("should delete radcheck", func(t *testing.T) {
		// Given
		fixture := testutil.CreateRadcheckFixture()
		fixture.ID = 0
		_ = repo.Create(fixture)
		createdID := fixture.ID

		// When
		err := repo.Delete(createdID)

		// Then
		require.NoError(t, err)

		// Verify deletion in database
		deleted := &entity.Radcheck{}
		result := db.First(deleted, createdID)
		assert.Error(t, result.Error)
		assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
	})

	t.Run("should handle delete of non-existent radcheck", func(t *testing.T) {
		// When
		err := repo.Delete(9999)

		// Then
		// GORM doesn't error on delete of non-existent record, so we just verify no error
		assert.NoError(t, err)
	})
}
