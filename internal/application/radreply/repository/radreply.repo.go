package repository

import (
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/entity"
	"gorm.io/gorm"
)

type RadreplyRepository interface {
	Create(radreply *entity.Radreply) error
	GetByID(id uint) (*entity.Radreply, error)
	GetByUsernameAndAttribute(username, attribute string) (*entity.Radreply, error)
	GetAll(filter *dto.RadreplyFilter) ([]entity.Radreply, int64, error)
	Update(radreply *entity.Radreply) error
	Delete(id uint) error
}

type radreplyRepository struct {
	db *gorm.DB
}

func NewRadreplyRepository(db *gorm.DB) RadreplyRepository {
	return &radreplyRepository{db: db}
}

func (r *radreplyRepository) Create(radreply *entity.Radreply) error {
	return r.db.Create(radreply).Error
}

func (r *radreplyRepository) GetByID(id uint) (*entity.Radreply, error) {
	var radreply entity.Radreply
	err := r.db.Where("id = ?", id).First(&radreply).Error
	if err != nil {
		return nil, err
	}
	return &radreply, nil
}

func (r *radreplyRepository) GetByUsernameAndAttribute(username, attribute string) (*entity.Radreply, error) {
	var radreply entity.Radreply
	err := r.db.Where("username = ? AND attribute = ?", username, attribute).First(&radreply).Error
	if err != nil {
		return nil, err
	}
	return &radreply, nil
}

func (r *radreplyRepository) GetAll(filter *dto.RadreplyFilter) ([]entity.Radreply, int64, error) {
	var radreply []entity.Radreply
	var total int64

	query := r.db

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

func (r *radreplyRepository) Update(radreply *entity.Radreply) error {
	return r.db.Save(radreply).Error
}

func (r *radreplyRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Radreply{}, id).Error
}
