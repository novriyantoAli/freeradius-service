package service

import (
	"errors"

	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/entity"
	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RadcheckService interface {
	CreateRadcheck(req *dto.CreateRadcheckRequest) (*dto.RadcheckResponse, error)
	GetRadcheckByID(id uint) (*dto.RadcheckResponse, error)
	GetRadcheckByUsernameAndAttribute(username, attribute string) (*dto.RadcheckResponse, error)
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
	radcheck := &entity.Radcheck{
		Username:  req.Username,
		Attribute: req.Attribute,
		Op:        req.Op,
		Value:     req.Value,
	}

	err := s.repo.Create(radcheck)
	if err != nil {
		s.logger.Error("Failed to create radcheck", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(radcheck), nil
}

func (s *radcheckService) GetRadcheckByID(id uint) (*dto.RadcheckResponse, error) {
	radcheck, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("radcheck not found")
		}
		return nil, err
	}

	return s.entityToResponse(radcheck), nil
}

func (s *radcheckService) GetRadcheckByUsernameAndAttribute(username, attribute string) (*dto.RadcheckResponse, error) {
	radcheck, err := s.repo.GetByUsernameAndAttribute(username, attribute)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("radcheck not found")
		}
		return nil, err
	}

	return s.entityToResponse(radcheck), nil
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

	radchecks, totalCount, err := s.repo.GetAll(filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.RadcheckResponse, 0, len(radchecks))
	for _, radcheck := range radchecks {
		responses = append(responses, *s.entityToResponse(&radcheck))
	}

	return &dto.ListRadcheckResponse{
		Data:      responses,
		Total:     totalCount,
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		TotalPage: (int(totalCount) + filter.PageSize - 1) / filter.PageSize,
	}, nil
}

func (s *radcheckService) UpdateRadcheck(id uint, req *dto.UpdateRadcheckRequest) (*dto.RadcheckResponse, error) {
	radcheck, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("radcheck not found")
		}
		return nil, err
	}

	if req.Username != "" {
		radcheck.Username = req.Username
	}
	if req.Attribute != "" {
		radcheck.Attribute = req.Attribute
	}
	if req.Op != "" {
		radcheck.Op = req.Op
	}
	if req.Value != "" {
		radcheck.Value = req.Value
	}

	err = s.repo.Update(radcheck)
	if err != nil {
		s.logger.Error("Failed to update radcheck", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(radcheck), nil
}

func (s *radcheckService) DeleteRadcheck(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("radcheck not found")
		}
		return err
	}

	return s.repo.Delete(id)
}

func (s *radcheckService) entityToResponse(radcheck *entity.Radcheck) *dto.RadcheckResponse {
	return &dto.RadcheckResponse{
		ID:        radcheck.ID,
		Username:  radcheck.Username,
		Attribute: radcheck.Attribute,
		Op:        radcheck.Op,
		Value:     radcheck.Value,
	}
}
