package service

import (
	"errors"

	"github.com/novriyantoAli/freeradius-service/internal/application/nas/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/nas/entity"
	"github.com/novriyantoAli/freeradius-service/internal/application/nas/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NASService interface {
	CreateNAS(req *dto.CreateNASRequest) (*dto.NASResponse, error)
	GetNASByID(id uint) (*dto.NASResponse, error)
	ListNAS(filter *dto.NASFilter) (*dto.ListNASResponse, error)
	UpdateNAS(id uint, req *dto.UpdateNASRequest) (*dto.NASResponse, error)
	DeleteNAS(id uint) error
}

type nasService struct {
	nasRepo repository.NASRepository
	logger  *zap.Logger
}

func NewNASService(nasRepo repository.NASRepository, logger *zap.Logger) NASService {
	return &nasService{
		nasRepo: nasRepo,
		logger:  logger,
	}
}

func (s *nasService) CreateNAS(req *dto.CreateNASRequest) (*dto.NASResponse, error) {
	s.logger.Info("Creating NAS", zap.String("nasname", req.NASName))

	// Check if NASName already exists
	existing, err := s.nasRepo.GetByNASName(req.NASName)
	if err == nil && existing != nil {
		s.logger.Warn("NASName already exists", zap.String("nasname", req.NASName))
		return nil, errors.New("nasname already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to check NASName existence", zap.Error(err))
		return nil, err
	}

	nas := &entity.NAS{
		NASName:         req.NASName,
		ShortName:       req.ShortName,
		Type:            req.Type,
		Ports:           0,
		Secret:          req.Secret,
		Server:          req.Server,
		Community:       req.Community,
		Description:     req.Description,
		RequireMa:       req.RequireMa,
		LimitProxyState: req.LimitProxyState,
	}

	if req.Ports != nil {
		nas.Ports = *req.Ports
	}

	if err := s.nasRepo.Create(nas); err != nil {
		s.logger.Error("Failed to create NAS", zap.Error(err))
		return nil, err
	}

	s.logger.Info("NAS created successfully", zap.Uint("id", nas.ID))
	return entityToResponse(nas), nil
}

func (s *nasService) GetNASByID(id uint) (*dto.NASResponse, error) {
	s.logger.Info("Getting NAS by ID", zap.Uint("id", id))

	nas, err := s.nasRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("NAS not found", zap.Uint("id", id))
			return nil, errors.New("nas not found")
		}
		s.logger.Error("Failed to get NAS", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	return entityToResponse(nas), nil
}

func (s *nasService) ListNAS(filter *dto.NASFilter) (*dto.ListNASResponse, error) {
	s.logger.Info("Listing NAS", zap.Int("page", filter.Page), zap.Int("page_size", filter.PageSize))

	// Set defaults
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 10
	}

	nasList, total, err := s.nasRepo.GetAll(filter)
	if err != nil {
		s.logger.Error("Failed to list NAS", zap.Error(err))
		return nil, err
	}

	responses := make([]dto.NASResponse, 0, len(nasList))
	for _, nas := range nasList {
		responses = append(responses, *entityToResponse(&nas))
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	s.logger.Info("NAS list retrieved", zap.Int64("total", total), zap.Int("page", filter.Page))
	return &dto.ListNASResponse{
		Data:      responses,
		Total:     total,
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		TotalPage: totalPages,
	}, nil
}

func (s *nasService) UpdateNAS(id uint, req *dto.UpdateNASRequest) (*dto.NASResponse, error) {
	s.logger.Info("Updating NAS", zap.Uint("id", id))

	nas, err := s.nasRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("NAS not found", zap.Uint("id", id))
			return nil, errors.New("nas not found")
		}
		s.logger.Error("Failed to get NAS", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	// Update fields if provided
	if req.NASName != "" {
		nas.NASName = req.NASName
	}
	if req.ShortName != "" {
		nas.ShortName = req.ShortName
	}
	if req.Type != "" {
		nas.Type = req.Type
	}
	if req.Ports != nil {
		nas.Ports = *req.Ports
	}
	if req.Secret != "" {
		nas.Secret = req.Secret
	}
	if req.Server != "" {
		nas.Server = req.Server
	}
	if req.Community != "" {
		nas.Community = req.Community
	}
	if req.Description != "" {
		nas.Description = req.Description
	}
	if req.RequireMa != "" {
		nas.RequireMa = req.RequireMa
	}
	if req.LimitProxyState != "" {
		nas.LimitProxyState = req.LimitProxyState
	}

	if err := s.nasRepo.Update(nas); err != nil {
		s.logger.Error("Failed to update NAS", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	s.logger.Info("NAS updated successfully", zap.Uint("id", id))
	return entityToResponse(nas), nil
}

func (s *nasService) DeleteNAS(id uint) error {
	s.logger.Info("Deleting NAS", zap.Uint("id", id))

	nas, err := s.nasRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("NAS not found", zap.Uint("id", id))
			return errors.New("nas not found")
		}
		s.logger.Error("Failed to get NAS", zap.Uint("id", id), zap.Error(err))
		return err
	}

	if err := s.nasRepo.Delete(nas.ID); err != nil {
		s.logger.Error("Failed to delete NAS", zap.Uint("id", id), zap.Error(err))
		return err
	}

	s.logger.Info("NAS deleted successfully", zap.Uint("id", id))
	return nil
}

// Helper function to convert entity to response
func entityToResponse(nas *entity.NAS) *dto.NASResponse {
	ports := nas.Ports
	return &dto.NASResponse{
		ID:              nas.ID,
		NASName:         nas.NASName,
		ShortName:       nas.ShortName,
		Type:            nas.Type,
		Ports:           &ports,
		Secret:          nas.Secret,
		Server:          nas.Server,
		Community:       nas.Community,
		Description:     nas.Description,
		RequireMa:       nas.RequireMa,
		LimitProxyState: nas.LimitProxyState,
		CreatedAt:       nas.CreatedAt.String(),
		UpdatedAt:       nas.UpdatedAt.String(),
	}
}
