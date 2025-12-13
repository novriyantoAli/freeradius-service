package repository

import (
	"errors"

	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RadcheckRepository interface {
	Create(req *dto.CreateRadcheckRequest) (*entity.Radcheck, error)
	GetByID(id uint) (*entity.Radcheck, error)
	GetByUsernameAndAttribute(username, attribute string) (*entity.Radcheck, error)
	GetAll(filter *dto.RadcheckFilter) (*dto.ListRadcheckResponse, error)
	Update(id uint, req *dto.UpdateRadcheckRequest) (*entity.Radcheck, error)
	Delete(id uint) error
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

func (r *radcheckRepository) Create(req *dto.CreateRadcheckRequest) (*entity.Radcheck, error) {
	radcheck := &entity.Radcheck{
		Username:  req.Username,
		Attribute: req.Attribute,
		Op:        req.Op,
		Value:     req.Value,
	}

	if radcheck.Op == "" {
		radcheck.Op = ":="
	}

	if err := r.db.Create(radcheck).Error; err != nil {
		r.logger.Error("Failed to create radcheck", zap.Error(err))
		return nil, err
	}

	r.logger.Info("Radcheck created successfully", zap.Uint("id", radcheck.ID), zap.String("username", radcheck.Username))
	return radcheck, nil
}

func (r *radcheckRepository) GetByID(id uint) (*entity.Radcheck, error) {
	var radcheck entity.Radcheck
	if err := r.db.First(&radcheck, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warn("Radcheck not found", zap.Uint("id", id))
			return nil, err
		}
		r.logger.Error("Failed to get radcheck by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	return &radcheck, nil
}

func (r *radcheckRepository) GetByUsernameAndAttribute(username, attribute string) (*entity.Radcheck, error) {
	var radcheck entity.Radcheck
	if err := r.db.Where("username = ? AND attribute = ?", username, attribute).First(&radcheck).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warn("Radcheck not found", zap.String("username", username), zap.String("attribute", attribute))
			return nil, err
		}
		r.logger.Error("Failed to get radcheck", zap.String("username", username), zap.String("attribute", attribute), zap.Error(err))
		return nil, err
	}
	return &radcheck, nil
}

func (r *radcheckRepository) GetAll(filter *dto.RadcheckFilter) (*dto.ListRadcheckResponse, error) {
	var radchecks []entity.Radcheck
	var total int64

	query := r.db

	if filter.Username != "" {
		query = query.Where("username LIKE ?", "%"+filter.Username+"%")
	}
	if filter.Attribute != "" {
		query = query.Where("attribute LIKE ?", "%"+filter.Attribute+"%")
	}

	if err := query.Model(&entity.Radcheck{}).Count(&total).Error; err != nil {
		r.logger.Error("Failed to count radcheck", zap.Error(err))
		return nil, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Offset(offset).Limit(filter.PageSize).Find(&radchecks).Error; err != nil {
		r.logger.Error("Failed to get all radcheck", zap.Error(err))
		return nil, err
	}

	responses := make([]dto.RadcheckResponse, len(radchecks))
	for i, rc := range radchecks {
		responses[i] = dto.RadcheckResponse{
			ID:        rc.ID,
			Username:  rc.Username,
			Attribute: rc.Attribute,
			Op:        rc.Op,
			Value:     rc.Value,
		}
	}

	totalPage := int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize))

	return &dto.ListRadcheckResponse{
		Data:      responses,
		Total:     total,
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		TotalPage: totalPage,
	}, nil
}

func (r *radcheckRepository) Update(id uint, req *dto.UpdateRadcheckRequest) (*entity.Radcheck, error) {
	radcheck, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Attribute != "" {
		updates["attribute"] = req.Attribute
	}
	if req.Op != "" {
		updates["op"] = req.Op
	}
	if req.Value != "" {
		updates["value"] = req.Value
	}

	if err := r.db.Model(radcheck).Updates(updates).Error; err != nil {
		r.logger.Error("Failed to update radcheck", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	r.logger.Info("Radcheck updated successfully", zap.Uint("id", id))
	return radcheck, nil
}

func (r *radcheckRepository) Delete(id uint) error {
	if err := r.db.Delete(&entity.Radcheck{}, id).Error; err != nil {
		r.logger.Error("Failed to delete radcheck", zap.Uint("id", id), zap.Error(err))
		return err
	}

	r.logger.Info("Radcheck deleted successfully", zap.Uint("id", id))
	return nil
}
