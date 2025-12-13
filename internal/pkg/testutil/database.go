package testutil

import (
	nasEntity "github.com/novriyantoAli/freeradius-service/internal/application/nas/entity"
	paymentEntity "github.com/novriyantoAli/freeradius-service/internal/application/payment/entity"
	userEntity "github.com/novriyantoAli/freeradius-service/internal/application/user/entity"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate all entities
	err = db.AutoMigrate(
		&userEntity.User{},
		&paymentEntity.Payment{},
		&nasEntity.NAS{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CleanDB cleans all data from test database
func CleanDB(db *gorm.DB) error {
	// Delete in reverse order of dependencies
	if err := db.Exec("DELETE FROM payments").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM nas").Error; err != nil {
		return err
	}
	return nil
}
