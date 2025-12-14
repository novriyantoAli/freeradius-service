package service

import (
	"context"
	"errors"

	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/entity"
	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RadcheckService interface {
	CreateRadcheck(ctx context.Context, req *dto.CreateRadcheckRequest) (*dto.RadcheckResponse, error)
	GetRadcheckByID(ctx context.Context, id uint) (*dto.RadcheckResponse, error)
	GetRadcheckByUsernameAndAttribute(ctx context.Context, username, attribute string) (*dto.RadcheckResponse, error)
	ListRadcheck(ctx context.Context, filter *dto.RadcheckFilter) (*dto.ListRadcheckResponse, error)
	UpdateRadcheck(ctx context.Context, id uint, req *dto.UpdateRadcheckRequest) (*dto.RadcheckResponse, error)
	DeleteRadcheck(ctx context.Context, id uint) error
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

func (s *radcheckService) CreateRadcheck(ctx context.Context, req *dto.CreateRadcheckRequest) (*dto.RadcheckResponse, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	radcheck := &entity.Radcheck{
		Username:  req.Username,
		Attribute: req.Attribute,
		Op:        req.Op,
		Value:     req.Value,
	}

	err := s.repo.Create(ctx, radcheck)
	if err != nil {
		s.logger.Error("Failed to create radcheck", zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(radcheck), nil
}

func (s *radcheckService) GetRadcheckByID(ctx context.Context, id uint) (*dto.RadcheckResponse, error) {
	radcheck, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("radcheck not found")
		}
		s.logger.Error("Failed to get radcheck by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(radcheck), nil
}

func (s *radcheckService) GetRadcheckByUsernameAndAttribute(ctx context.Context, username, attribute string) (*dto.RadcheckResponse, error) {
	radcheck, err := s.repo.GetByUsernameAndAttribute(ctx, username, attribute)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("radcheck not found")
		}
		s.logger.Error("Failed to get radcheck", zap.String("username", username), zap.String("attribute", attribute), zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(radcheck), nil
}

func (s *radcheckService) ListRadcheck(ctx context.Context, filter *dto.RadcheckFilter) (*dto.ListRadcheckResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	radchecks, totalCount, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to list radchecks", zap.Error(err))
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

func (s *radcheckService) UpdateRadcheck(ctx context.Context, id uint, req *dto.UpdateRadcheckRequest) (*dto.RadcheckResponse, error) {
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	radcheck, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("radcheck not found")
		}
		s.logger.Error("Failed to get radcheck by ID", zap.Uint("id", id), zap.Error(err))
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

	err = s.repo.Update(ctx, radcheck)
	if err != nil {
		s.logger.Error("Failed to update radcheck", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	return s.entityToResponse(radcheck), nil
}

func (s *radcheckService) DeleteRadcheck(ctx context.Context, id uint) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("radcheck not found")
		}
		s.logger.Error("Failed to get radcheck by ID", zap.Uint("id", id), zap.Error(err))
		return err
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Failed to delete radcheck", zap.Uint("id", id), zap.Error(err))
		return err
	}

	return nil
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

func (s *radcheckService) validateCreateRequest(req *dto.CreateRadcheckRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) < 1 || len(req.Username) > 64 {
		return errors.New("username must be between 1 and 64 characters")
	}

	if req.Attribute == "" {
		return errors.New("attribute is required")
	}
	if len(req.Attribute) < 1 || len(req.Attribute) > 64 {
		return errors.New("attribute must be between 1 and 64 characters")
	}

	if req.Value == "" {
		return errors.New("value is required")
	}
	if len(req.Value) < 1 || len(req.Value) > 253 {
		return errors.New("value must be between 1 and 253 characters")
	}

	return nil
}

func (s *radcheckService) validateUpdateRequest(req *dto.UpdateRadcheckRequest) error {
	if req.Username != "" && (len(req.Username) < 1 || len(req.Username) > 64) {
		return errors.New("username must be between 1 and 64 characters")
	}

	if req.Attribute != "" && (len(req.Attribute) < 1 || len(req.Attribute) > 64) {
		return errors.New("attribute must be between 1 and 64 characters")
	}

	if req.Value != "" && (len(req.Value) < 1 || len(req.Value) > 253) {
		return errors.New("value must be between 1 and 253 characters")
	}

	return nil
}