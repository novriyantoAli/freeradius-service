package service

import (
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/entity"
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/repository"
)

type RadreplyService interface {
	CreateRadreply(req *dto.CreateRadreplyRequest) (*dto.RadreplyResponse, error)
	GetRadreplyByID(id uint) (*dto.RadreplyResponse, error)
	GetRadreplyByUsernameAndAttribute(username, attribute string) (*dto.RadreplyResponse, error)
	ListRadreply(filter *dto.RadreplyFilter) (*dto.ListRadreplyResponse, error)
	UpdateRadreply(id uint, req *dto.UpdateRadreplyRequest) (*dto.RadreplyResponse, error)
	DeleteRadreply(id uint) error
}

type radreplyService struct {
	repository repository.RadreplyRepository
}

func NewRadreplyService(repository repository.RadreplyRepository) RadreplyService {
	return &radreplyService{repository: repository}
}

func (s *radreplyService) CreateRadreply(req *dto.CreateRadreplyRequest) (*dto.RadreplyResponse, error) {
	radreply := &entity.Radreply{
		Username:  req.Username,
		Attribute: req.Attribute,
		Op:        req.Op,
		Value:     req.Value,
	}

	if err := s.repository.Create(radreply); err != nil {
		return nil, err
	}

	return s.entityToResponse(radreply), nil
}

func (s *radreplyService) GetRadreplyByID(id uint) (*dto.RadreplyResponse, error) {
	radreply, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(radreply), nil
}

func (s *radreplyService) GetRadreplyByUsernameAndAttribute(username, attribute string) (*dto.RadreplyResponse, error) {
	radreply, err := s.repository.GetByUsernameAndAttribute(username, attribute)
	if err != nil {
		return nil, err
	}

	return s.entityToResponse(radreply), nil
}

func (s *radreplyService) ListRadreply(filter *dto.RadreplyFilter) (*dto.ListRadreplyResponse, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	radreply, total, err := s.repository.GetAll(filter)
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

func (s *radreplyService) UpdateRadreply(id uint, req *dto.UpdateRadreplyRequest) (*dto.RadreplyResponse, error) {
	radreply, err := s.repository.GetByID(id)
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

	if err := s.repository.Update(radreply); err != nil {
		return nil, err
	}

	return s.entityToResponse(radreply), nil
}

func (s *radreplyService) DeleteRadreply(id uint) error {
	return s.repository.Delete(id)
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
