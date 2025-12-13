package repository

import (
	"github.com/novriyantoAli/freeradius-service/internal/application/nas/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/nas/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NASRepository interface {
	Create(nas *entity.NAS) error
	GetByID(id uint) (*entity.NAS, error)
	GetByNASName(nasname string) (*entity.NAS, error)
	GetAll(filter *dto.NASFilter) ([]entity.NAS, int64, error)
	Update(nas *entity.NAS) error
	Delete(id uint) error
}

type nasRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewNASRepository(db *gorm.DB, logger *zap.Logger) NASRepository {
	return &nasRepository{
		db:     db,
		logger: logger,
	}
}

func (r *nasRepository) Create(nas *entity.NAS) error {
	r.logger.Info("Creating NAS", zap.String("nasname", nas.NASName))
	return r.db.Create(nas).Error
}

func (r *nasRepository) GetByID(id uint) (*entity.NAS, error) {
	var nas entity.NAS
	err := r.db.First(&nas, id).Error
	if err != nil {
		r.logger.Error("Failed to get NAS by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	return &nas, nil
}

func (r *nasRepository) GetByNASName(nasname string) (*entity.NAS, error) {
	var nas entity.NAS
	err := r.db.Where("nas_name = ?", nasname).First(&nas).Error
	if err != nil {
		r.logger.Error("Failed to get NAS by name", zap.String("nasname", nasname), zap.Error(err))
		return nil, err
	}
	return &nas, nil
}

func (r *nasRepository) GetAll(filter *dto.NASFilter) ([]entity.NAS, int64, error) {
	var nasList []entity.NAS
	var totalCount int64

	query := r.db.Model(&entity.NAS{})

	if filter.NASName != "" {
		query = query.Where("nas_name LIKE ?", "%"+filter.NASName+"%")
	}
	if filter.ShortName != "" {
		query = query.Where("short_name LIKE ?", "%"+filter.ShortName+"%")
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Description != "" {
		query = query.Where("description LIKE ?", "%"+filter.Description+"%")
	}

	query.Count(&totalCount)

	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(offset).Limit(filter.PageSize)
	}

	err := query.Find(&nasList).Error
	if err != nil {
		r.logger.Error("Failed to get NAS list", zap.Error(err))
		return nil, 0, err
	}

	return nasList, totalCount, nil
}

func (r *nasRepository) Update(nas *entity.NAS) error {
	r.logger.Info("Updating NAS", zap.Uint("id", nas.ID))
	return r.db.Save(nas).Error
}

func (r *nasRepository) Delete(id uint) error {
	r.logger.Info("Deleting NAS", zap.Uint("id", id))
	return r.db.Delete(&entity.NAS{}, id).Error
}
