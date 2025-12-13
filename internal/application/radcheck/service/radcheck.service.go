package service

import (
	"errors"
	"fmt"

	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RadcheckService interface {
	CreateRadcheck(req *dto.CreateRadcheckRequest) (*dto.RadcheckResponse, error)
	GetRadcheckByID(id uint) (*dto.RadcheckResponse, error)
	ListRadcheck(filter *dto.RadcheckFilter) (*dto.ListRadcheckResponse, error)
	UpdateRadcheck(id uint, req *dto.UpdateRadcheckRequest) (*dto.RadcheckResponse, error)
	DeleteRadcheck(id uint) error
}

type radcheckService struct {
	repo   repository.RadcheckRepository
	logger *zap.Logger
}

func NewRadcheckService(repo repository.RadcheckRepository, logger *zap.Logger) RadcheckService {
	return &radcheckService{
		repo:   repo,
		logger: logger,
	}
}

func (s *radcheckService) CreateRadcheck(req *dto.CreateRadcheckRequest) (*dto.RadcheckResponse, error) {
	if req.Username == "" {
		err := fmt.Errorf("username is required")
		s.logger.Error("Invalid request", zap.Error(err))
		return nil, err
	}

	if req.Attribute == "" {
		err := fmt.Errorf("attribute is required")
		s.logger.Error("Invalid request", zap.Error(err))
		return nil, err
	}

	if req.Value == "" {
		err := fmt.Errorf("value is required")
		s.logger.Error("Invalid request", zap.Error(err))
		return nil, err
	}

	radcheck, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return &dto.RadcheckResponse{
		ID:        radcheck.ID,
		Username:  radcheck.Username,
		Attribute: radcheck.Attribute,
		Op:        radcheck.Op,
		Value:     radcheck.Value,
		CreatedAt: radcheck.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: radcheck.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *radcheckService) GetRadcheckByID(id uint) (*dto.RadcheckResponse, error) {
	radcheck, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Radcheck not found", zap.Uint("id", id))
			return nil, fmt.Errorf("radcheck not found")
		}
		return nil, err
	}

	return &dto.RadcheckResponse{
		ID:        radcheck.ID,
		Username:  radcheck.Username,
		Attribute: radcheck.Attribute,
		Op:        radcheck.Op,
		Value:     radcheck.Value,
		CreatedAt: radcheck.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: radcheck.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *radcheckService) ListRadcheck(filter *dto.RadcheckFilter) (*dto.ListRadcheckResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	return s.repo.GetAll(filter)
}

func (s *radcheckService) UpdateRadcheck(id uint, req *dto.UpdateRadcheckRequest) (*dto.RadcheckResponse, error) {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Radcheck not found", zap.Uint("id", id))
			return nil, fmt.Errorf("radcheck not found")
		}
		return nil, err
	}

	radcheck, err := s.repo.Update(id, req)
	if err != nil {
		return nil, err
	}

	return &dto.RadcheckResponse{
		ID:        radcheck.ID,
		Username:  radcheck.Username,
		Attribute: radcheck.Attribute,
		Op:        radcheck.Op,
		Value:     radcheck.Value,
		CreatedAt: radcheck.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: radcheck.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *radcheckService) DeleteRadcheck(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Radcheck not found", zap.Uint("id", id))
			return fmt.Errorf("radcheck not found")
		}
		return err
	}

	return s.repo.Delete(id)
}
