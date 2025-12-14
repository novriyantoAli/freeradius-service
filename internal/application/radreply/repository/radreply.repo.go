package repository

import (
	"context"

	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/entity"
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
	db *gorm.DB
}

func NewRadreplyRepository(db *gorm.DB) RadreplyRepository {
	return &radreplyRepository{db: db}
}

func (r *radreplyRepository) Create(ctx context.Context, radreply *entity.Radreply) error {
	return r.db.WithContext(ctx).Create(radreply).Error
}

func (r *radreplyRepository) GetByID(ctx context.Context, id uint) (*entity.Radreply, error) {
	var radreply entity.Radreply
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&radreply).Error
	if err != nil {
		return nil, err
	}
	return &radreply, nil
}

func (r *radreplyRepository) GetByUsernameAndAttribute(ctx context.Context, username, attribute string) (*entity.Radreply, error) {
	var radreply entity.Radreply
	err := r.db.WithContext(ctx).Where("username = ? AND attribute = ?", username, attribute).First(&radreply).Error
	if err != nil {
		return nil, err
	}
	return &radreply, nil
}

func (r *radreplyRepository) GetAll(ctx context.Context, filter *dto.RadreplyFilter) ([]entity.Radreply, int64, error) {
	var radreply []entity.Radreply
	var total int64

	query := r.db.WithContext(ctx)

	if filter.Username != "" {
		query = query.Where("username = ?", filter.Username)
	}

	if filter.Attribute != "" {
		query = query.Where("attribute = ?", filter.Attribute)
	}

	// Get total count
	if err := query.Model(&entity.Radreply{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Offset(offset).Limit(filter.PageSize).Find(&radreply).Error; err != nil {
		return nil, 0, err
	}

	return radreply, total, nil
}

func (r *radreplyRepository) Update(ctx context.Context, radreply *entity.Radreply) error {
	return r.db.WithContext(ctx).Save(radreply).Error
}

func (r *radreplyRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Radreply{}, id).Error
}
