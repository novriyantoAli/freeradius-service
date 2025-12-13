package repository

import (
	"fmt"
	"testing"

	nasDto "github.com/novriyantoAli/freeradius-service/internal/application/nas/dto"
	nasEntity "github.com/novriyantoAli/freeradius-service/internal/application/nas/entity"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNASRepository_Create(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewNASRepository(db, logger)

	t.Run("should create NAS successfully", func(t *testing.T) {
		// Given
		nas := testutil.CreateNASFixture()
		nas.ID = 0 // Reset ID for creation

		// When
		err := repo.Create(nas)

		// Then
		assert.NoError(t, err)
		assert.NotZero(t, nas.ID)

		// Verify NAS was created in database
		var dbNAS nasEntity.NAS
		err = db.First(&dbNAS, nas.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, nas.NASName, dbNAS.NASName)
		assert.Equal(t, nas.ShortName, dbNAS.ShortName)
		assert.Equal(t, nas.Type, dbNAS.Type)
	})

	t.Run("should fail to create NAS with duplicate nasname", func(t *testing.T) {
		// Given
		nas1 := testutil.CreateNASFixture()
		nas1.ID = 0
		nas1.NASName = "duplicate-nas"

		nas2 := testutil.CreateNASFixture()
		nas2.ID = 0
		nas2.NASName = "duplicate-nas"

		// When
		err1 := repo.Create(nas1)
		err2 := repo.Create(nas2)

		// Then
		assert.NoError(t, err1)
		assert.Error(t, err2) // Should fail due to unique constraint
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestNASRepository_GetByID(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewNASRepository(db, logger)

	t.Run("should get NAS by ID successfully", func(t *testing.T) {
		// Given
		nas := testutil.CreateNASFixture()
		nas.ID = 0
		err := repo.Create(nas)
		require.NoError(t, err)

		// When
		foundNAS, err := repo.GetByID(nas.ID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, nas.ID, foundNAS.ID)
		assert.Equal(t, nas.NASName, foundNAS.NASName)
		assert.Equal(t, nas.ShortName, foundNAS.ShortName)
		assert.Equal(t, nas.Type, foundNAS.Type)
	})

	t.Run("should return error when NAS not found", func(t *testing.T) {
		// When
		_, err := repo.GetByID(999)

		// Then
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestNASRepository_GetByNASName(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewNASRepository(db, logger)

	t.Run("should get NAS by NASName successfully", func(t *testing.T) {
		// Given
		nas := testutil.CreateNASFixture()
		nas.ID = 0
		nas.NASName = "unique-nas-name"
		err := repo.Create(nas)
		require.NoError(t, err)

		// When
		foundNAS, err := repo.GetByNASName("unique-nas-name")

		// Then
		assert.NoError(t, err)
		assert.Equal(t, nas.ID, foundNAS.ID)
		assert.Equal(t, "unique-nas-name", foundNAS.NASName)
		assert.Equal(t, nas.ShortName, foundNAS.ShortName)
	})

	t.Run("should return error when NAS NASName not found", func(t *testing.T) {
		// When
		_, err := repo.GetByNASName("nonexistent-nas")

		// Then
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestNASRepository_GetAll(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewNASRepository(db, logger)

	t.Run("should get all NAS with pagination", func(t *testing.T) {
		// Given - Create multiple NAS
		for i := 0; i < 5; i++ {
			nas := testutil.CreateNASFixture()
			nas.ID = 0
			nas.NASName = fmt.Sprintf("nas-%d", i)
			nas.ShortName = fmt.Sprintf("n%d", i)
			err := repo.Create(nas)
			require.NoError(t, err)
		}

		filter := &nasDto.NASFilter{
			Page:     1,
			PageSize: 3,
		}

		// When
		nasServers, totalCount, err := repo.GetAll(filter)

		// Then
		assert.NoError(t, err)
		assert.Len(t, nasServers, 3)      // Should return 3 NAS due to page size
		assert.Equal(t, int64(5), totalCount) // Total count should be 5
	})

	t.Run("should filter NAS by NASName", func(t *testing.T) {
		// Given
		nas1 := testutil.CreateNASFixture()
		nas1.ID = 0
		nas1.NASName = "radius-server-1"
		nas1.ShortName = "rs1"
		err := repo.Create(nas1)
		require.NoError(t, err)

		nas2 := testutil.CreateNASFixture()
		nas2.ID = 0
		nas2.NASName = "radius-server-2"
		nas2.ShortName = "rs2"
		err = repo.Create(nas2)
		require.NoError(t, err)

		filter := &nasDto.NASFilter{
			NASName: "radius-server-1",
		}

		// When
		nasServers, totalCount, err := repo.GetAll(filter)

		// Then
		assert.NoError(t, err)
		assert.Len(t, nasServers, 1)
		assert.Equal(t, int64(1), totalCount)
		assert.Equal(t, "radius-server-1", nasServers[0].NASName)
	})

	t.Run("should filter NAS by Type", func(t *testing.T) {
		// Given
		nas1 := testutil.CreateNASFixture()
		nas1.ID = 0
		nas1.NASName = "cisco-nas"
		nas1.Type = "cisco"
		err := repo.Create(nas1)
		require.NoError(t, err)

		nas2 := testutil.CreateNASFixture()
		nas2.ID = 0
		nas2.NASName = "juniper-nas"
		nas2.Type = "juniper"
		err = repo.Create(nas2)
		require.NoError(t, err)

		filter := &nasDto.NASFilter{
			Type: "cisco",
		}

		// When
		nasServers, totalCount, err := repo.GetAll(filter)

		// Then
		assert.NoError(t, err)
		assert.Len(t, nasServers, 1)
		assert.Equal(t, int64(1), totalCount)
		assert.Equal(t, "cisco", nasServers[0].Type)
	})

	t.Run("should filter NAS by Description", func(t *testing.T) {
		// Given
		nas1 := testutil.CreateNASFixture()
		nas1.ID = 0
		nas1.NASName = "prod-nas"
		nas1.Description = "Production RADIUS Server"
		err := repo.Create(nas1)
		require.NoError(t, err)

		nas2 := testutil.CreateNASFixture()
		nas2.ID = 0
		nas2.NASName = "test-nas"
		nas2.Description = "Testing RADIUS Server"
		err = repo.Create(nas2)
		require.NoError(t, err)

		filter := &nasDto.NASFilter{
			Description: "Production",
		}

		// When
		nasServers, totalCount, err := repo.GetAll(filter)

		// Then
		assert.NoError(t, err)
		assert.Len(t, nasServers, 1)
		assert.Equal(t, int64(1), totalCount)
		assert.Equal(t, "Production RADIUS Server", nasServers[0].Description)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestNASRepository_Update(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewNASRepository(db, logger)

	t.Run("should update NAS successfully", func(t *testing.T) {
		// Given
		nas := testutil.CreateNASFixture()
		nas.ID = 0
		err := repo.Create(nas)
		require.NoError(t, err)

		// When
		nas.ShortName = "updated-short"
		nas.Description = "Updated Description"
		ports := 1813
		nas.Ports = ports
		err = repo.Update(nas)

		// Then
		assert.NoError(t, err)

		// Verify update in database
		var dbNAS nasEntity.NAS
		err = db.First(&dbNAS, nas.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "updated-short", dbNAS.ShortName)
		assert.Equal(t, "Updated Description", dbNAS.Description)
		assert.Equal(t, 1813, dbNAS.Ports)
	})

	// Cleanup
	testutil.CleanDB(db)
}

func TestNASRepository_Delete(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	logger := testutil.NewTestLogger(t)
	repo := NewNASRepository(db, logger)

	t.Run("should delete NAS successfully", func(t *testing.T) {
		// Given
		nas := testutil.CreateNASFixture()
		nas.ID = 0
		err := repo.Create(nas)
		require.NoError(t, err)

		// When
		err = repo.Delete(nas.ID)

		// Then
		assert.NoError(t, err)

		// Verify NAS is deleted (soft delete)
		var dbNAS nasEntity.NAS
		err = db.First(&dbNAS, nas.ID).Error
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	// Cleanup
	testutil.CleanDB(db)
}
