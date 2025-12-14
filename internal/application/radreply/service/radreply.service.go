package service

import (
	"context"
	"errors"

	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/entity"
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/repository"
)

type RadreplyService interface {
	CreateRadreply(ctx context.Context, req *dto.CreateRadreplyRequest) (*dto.RadreplyResponse, error)
	GetRadreplyByID(ctx context.Context, id uint) (*dto.RadreplyResponse, error)
	GetRadreplyByUsernameAndAttribute(ctx context.Context, username, attribute string) (*dto.RadreplyResponse, error)
	ListRadreply(ctx context.Context, filter *dto.RadreplyFilter) (*dto.ListRadreplyResponse, error)
	UpdateRadreply(ctx context.Context, id uint, req *dto.UpdateRadreplyRequest) (*dto.RadreplyResponse, error)
	DeleteRadreply(ctx context.Context, id uint) error
}

type radreplyService struct {
	repository repository.RadreplyRepository
}

func NewRadreplyService(repository repository.RadreplyRepository) RadreplyService {
	return &radreplyService{repository: repository}
}

func (s *radreplyService) CreateRadreply(ctx context.Context, req *dto.CreateRadreplyRequest) (*dto.RadreplyResponse, error) {
	// Validate constraints before transaction (fail-fast principle)
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	radreply := &entity.Radreply{
		Username:  req.Username,
		Attribute: req.Attribute,
		Op:        req.Op,
		Value:     req.Value,
	}

	if err := s.repository.Create(ctx, radreply); err != nil {
		return nil, err
	}

	return s.entityToResponse(radreply), nil
}

func (s *radreplyService) GetRadreplyByID(ctx context.Context, id uint) (*dto.RadreplyResponse, error) {
	radreply, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(radreply), nil
}

func (s *radreplyService) GetRadreplyByUsernameAndAttribute(ctx context.Context, username, attribute string) (*dto.RadreplyResponse, error) {
	radreply, err := s.repository.GetByUsernameAndAttribute(ctx, username, attribute)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(radreply), nil
}

func (s *radreplyService) ListRadreply(ctx context.Context, filter *dto.RadreplyFilter) (*dto.ListRadreplyResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	radreply, total, err := s.repository.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.RadreplyResponse, 0, len(radreply))
	for _, r := range radreply {
		responses = append(responses, *s.entityToResponse(&r))
	}

	totalPage := int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize))

	return &dto.ListRadreplyResponse{
		Data:      responses,
		Total:     total,
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		TotalPage: totalPage,
	}, nil
}

func (s *radreplyService) UpdateRadreply(ctx context.Context, id uint, req *dto.UpdateRadreplyRequest) (*dto.RadreplyResponse, error) {
	// Validate constraints before fetching (fail-fast principle)
	if err := s.validateUpdateRequest(req); err != nil {
		return nil, err
	}

	radreply, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Username != "" {
		radreply.Username = req.Username
	}
	if req.Attribute != "" {
		radreply.Attribute = req.Attribute
	}
	if req.Op != "" {
		radreply.Op = req.Op
	}
	if req.Value != "" {
		radreply.Value = req.Value
	}

	if err := s.repository.Update(ctx, radreply); err != nil {
		return nil, err
	}

	return s.entityToResponse(radreply), nil
}

func (s *radreplyService) DeleteRadreply(ctx context.Context, id uint) error {
	return s.repository.Delete(ctx, id)
}

func (s *radreplyService) entityToResponse(radreply *entity.Radreply) *dto.RadreplyResponse {
	return &dto.RadreplyResponse{
		ID:        radreply.ID,
		Username:  radreply.Username,
		Attribute: radreply.Attribute,
		Op:        radreply.Op,
		Value:     radreply.Value,
	}
}

// DB Constraint Validation Methods

func (s *radreplyService) validateCreateRequest(req *dto.CreateRadreplyRequest) error {
	// Username constraint: NOT NULL, size 1-64
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) > 64 {
		return errors.New("username must not exceed 64 characters")
	}

	// Attribute constraint: NOT NULL, size 1-64
	if req.Attribute == "" {
		return errors.New("attribute is required")
	}
	if len(req.Attribute) > 64 {
		return errors.New("attribute must not exceed 64 characters")
	}

	// Op constraint: NOT NULL, size 1-2, default '='
	if req.Op == "" {
		return errors.New("op is required")
	}
	if len(req.Op) > 2 {
		return errors.New("op must not exceed 2 characters")
	}

	// Value constraint: NOT NULL, size 1-253
	if req.Value == "" {
		return errors.New("value is required")
	}
	if len(req.Value) > 253 {
		return errors.New("value must not exceed 253 characters")
	}

	return nil
}

func (s *radreplyService) validateUpdateRequest(req *dto.UpdateRadreplyRequest) error {
	// Only validate fields that are being updated
	if req.Username != "" && len(req.Username) > 64 {
		return errors.New("username must not exceed 64 characters")
	}

	if req.Attribute != "" && len(req.Attribute) > 64 {
		return errors.New("attribute must not exceed 64 characters")
	}

	if req.Op != "" && len(req.Op) > 2 {
		return errors.New("op must not exceed 2 characters")
	}

	if req.Value != "" && len(req.Value) > 253 {
		return errors.New("value must not exceed 253 characters")
	}

	return nil
}
