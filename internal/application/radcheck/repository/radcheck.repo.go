package repository

import (
	"context"

	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/entity"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/database"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RadcheckRepository interface {
	Create(ctx context.Context, radcheck *entity.Radcheck) error
	GetByID(ctx context.Context, id uint) (*entity.Radcheck, error)
	GetByUsernameAndAttribute(ctx context.Context, username, attribute string) (*entity.Radcheck, error)
	GetAll(ctx context.Context, filter *dto.RadcheckFilter) ([]entity.Radcheck, int64, error)
	Update(ctx context.Context, radcheck *entity.Radcheck) error
	Delete(ctx context.Context, id uint) error
}

type radcheckRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewRadcheckRepository(db *gorm.DB, logger *zap.Logger) RadcheckRepository {
	return &radcheckRepository{
		db:     db,
		logger: logger,
	}
}

func (r *radcheckRepository) Create(ctx context.Context, radcheck *entity.Radcheck) error {
	if radcheck.Op == "" {
		radcheck.Op = ":="
	}
	r.logger.Info("Creating radcheck", zap.String("username", radcheck.Username))
	db := database.GetDB(ctx, r.db).(*gorm.DB)
	return db.Create(radcheck).Error
}

func (r *radcheckRepository) GetByID(ctx context.Context, id uint) (*entity.Radcheck, error) {
	var radcheck entity.Radcheck
	db := database.GetDB(ctx, r.db).(*gorm.DB)
	err := db.First(&radcheck, id).Error
	if err != nil {
		r.logger.Error("Failed to get radcheck by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	return &radcheck, nil
}

func (r *radcheckRepository) GetByUsernameAndAttribute(ctx context.Context, username, attribute string) (*entity.Radcheck, error) {
	var radcheck entity.Radcheck
	db := database.GetDB(ctx, r.db).(*gorm.DB)
	err := db.Where("username = ? AND attribute = ?", username, attribute).First(&radcheck).Error
	if err != nil {
		r.logger.Error("Failed to get radcheck", zap.String("username", username), zap.String("attribute", attribute), zap.Error(err))
		return nil, err
	}
	return &radcheck, nil
}

func (r *radcheckRepository) GetAll(ctx context.Context, filter *dto.RadcheckFilter) ([]entity.Radcheck, int64, error) {
	var radchecks []entity.Radcheck
	var totalCount int64

	query := database.GetDB(ctx, r.db).(*gorm.DB).Model(&entity.Radcheck{})

	if filter.Username != "" {
		query = query.Where("username LIKE ?", "%"+filter.Username+"%")
	}
	if filter.Attribute != "" {
		query = query.Where("attribute LIKE ?", "%"+filter.Attribute+"%")
	}

	if err := query.Count(&totalCount).Error; err != nil {
		r.logger.Error("Failed to count radchecks", zap.Error(err))
		return nil, 0, err
	}

	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	err := query.Find(&radchecks).Error
	if err != nil {
		r.logger.Error("Failed to get radchecks", zap.Error(err))
		return nil, 0, err
	}

	return radchecks, totalCount, nil
}

func (r *radcheckRepository) Update(ctx context.Context, radcheck *entity.Radcheck) error {
	r.logger.Info("Updating radcheck", zap.Uint("id", radcheck.ID))
	db := database.GetDB(ctx, r.db).(*gorm.DB)
	return db.Save(radcheck).Error
}

func (r *radcheckRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting radcheck", zap.Uint("id", id))
	db := database.GetDB(ctx, r.db).(*gorm.DB)
	return db.Delete(&entity.Radcheck{}, id).Error
}
