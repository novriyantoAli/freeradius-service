package repository

import (
	"context"

	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/entity"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/database"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RadreplyRepository interface {
	Create(ctx context.Context, radreply *entity.Radreply) error
	GetByID(ctx context.Context, id uint) (*entity.Radreply, error)
	GetByUsernameAndAttribute(ctx context.Context, username, attribute string) (*entity.Radreply, error)
	GetAll(ctx context.Context, filter *dto.RadreplyFilter) ([]entity.Radreply, int64, error)
	Update(ctx context.Context, radreply *entity.Radreply) error
	Delete(ctx context.Context, id uint) error
}

type radreplyRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewRadreplyRepository(db *gorm.DB, logger *zap.Logger) RadreplyRepository {
	return &radreplyRepository{
		db:     db,
		logger: logger,
	}
}

func (r *radreplyRepository) Create(ctx context.Context, radreply *entity.Radreply) error {
	r.logger.Info("Creating radreply", zap.String("username", radreply.Username))
	db := database.GetDB(ctx, r.db).(*gorm.DB)
	return db.Create(radreply).Error
}

func (r *radreplyRepository) GetByID(ctx context.Context, id uint) (*entity.Radreply, error) {
	var radreply entity.Radreply
	db := database.GetDB(ctx, r.db).(*gorm.DB)
	err := db.Where("id = ?", id).First(&radreply).Error
	if err != nil {
		r.logger.Error("Failed to get radreply by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	return &radreply, nil
}

func (r *radreplyRepository) GetByUsernameAndAttribute(ctx context.Context, username, attribute string) (*entity.Radreply, error) {
	var radreply entity.Radreply
	db := database.GetDB(ctx, r.db).(*gorm.DB)
	err := db.Where("username = ? AND attribute = ?", username, attribute).First(&radreply).Error
	if err != nil {
		r.logger.Error("Failed to get radreply", zap.String("username", username), zap.String("attribute", attribute), zap.Error(err))
		return nil, err
	}
	return &radreply, nil
}

func (r *radreplyRepository) GetAll(ctx context.Context, filter *dto.RadreplyFilter) ([]entity.Radreply, int64, error) {
	var radreply []entity.Radreply
	var total int64

	query := database.GetDB(ctx, r.db).(*gorm.DB)

	if filter.Username != "" {
		query = query.Where("username = ?", filter.Username)
	}

	if filter.Attribute != "" {
		query = query.Where("attribute = ?", filter.Attribute)
	}

	// Get total count
	if err := query.Model(&entity.Radreply{}).Count(&total).Error; err != nil {
		r.logger.Error("Failed to get radreply count", zap.Error(err))
		return nil, 0, err
	}

	// Get paginated results
	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Offset(offset).Limit(filter.PageSize).Find(&radreply).Error; err != nil {
		r.logger.Error("Failed to get radreply list", zap.Error(err))
		return nil, 0, err
	}

	return radreply, total, nil
}

func (r *radreplyRepository) Update(ctx context.Context, radreply *entity.Radreply) error {
	r.logger.Info("Updating radreply", zap.Uint("id", radreply.ID))
	db := database.GetDB(ctx, r.db).(*gorm.DB)
	return db.Save(radreply).Error
}

func (r *radreplyRepository) Delete(ctx context.Context, id uint) error {
	r.logger.Info("Deleting radreply", zap.Uint("id", id))
	db := database.GetDB(ctx, r.db).(*gorm.DB)
	return db.Delete(&entity.Radreply{}, id).Error
}
