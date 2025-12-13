package handler

import (
	"context"
	"time"

	"github.com/novriyantoAli/freeradius-service/api/proto/nas"
	"github.com/novriyantoAli/freeradius-service/internal/application/nas/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/nas/service"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NASGrpcHandler struct {
	nas.UnimplementedNASServiceServer
	nasService service.NASService
	logger     *zap.Logger
}

func NewNASGrpcHandler(nasService service.NASService, logger *zap.Logger) *NASGrpcHandler {
	return &NASGrpcHandler{
		nasService: nasService,
		logger:     logger,
	}
}

func (h *NASGrpcHandler) CreateNAS(
	ctx context.Context,
	req *nas.CreateNASRequest,
) (*nas.CreateNASResponse, error) {
	ports := int(req.Ports)
	createReq := &dto.CreateNASRequest{
		NASName:         req.Nasname,
		ShortName:       req.Shortname,
		Type:            req.Type,
		Ports:           &ports,
		Secret:          req.Secret,
		Server:          req.Server,
		Community:       req.Community,
		Description:     req.Description,
		RequireMa:       req.RequireMa,
		LimitProxyState: req.LimitProxyState,
	}

	nasResponse, err := h.nasService.CreateNAS(createReq)
	if err != nil {
		h.logger.Error("Failed to create NAS via gRPC", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create NAS: %v", err)
	}

	return &nas.CreateNASResponse{
		Nas: h.toProtoNAS(nasResponse),
	}, nil
}

func (h *NASGrpcHandler) GetNAS(ctx context.Context, req *nas.GetNASRequest) (*nas.GetNASResponse, error) {
	nasResponse, err := h.nasService.GetNASByID(uint(req.Id))
	if err != nil {
		h.logger.Error("Failed to get NAS via gRPC", zap.Uint32("id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "NAS not found: %v", err)
	}

	return &nas.GetNASResponse{
		Nas: h.toProtoNAS(nasResponse),
	}, nil
}

func (h *NASGrpcHandler) ListNAS(ctx context.Context, req *nas.ListNASRequest) (*nas.ListNASResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	filter := &dto.NASFilter{
		NASName:   req.Filter.Nasname,
		ShortName: req.Filter.Shortname,
		Type:      req.Filter.Type,
		Page:      page,
		PageSize:  pageSize,
	}

	listResponse, err := h.nasService.ListNAS(filter)
	if err != nil {
		h.logger.Error("Failed to list NAS via gRPC", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to list NAS: %v", err)
	}

	protoNAS := make([]*nas.NAS, len(listResponse.Data))
	for i, n := range listResponse.Data {
		protoNAS[i] = h.toProtoNAS(&n)
	}

	return &nas.ListNASResponse{
		Nas:      protoNAS,
		Total:    listResponse.Total,
		Page:     int32(listResponse.Page),
		PageSize: int32(listResponse.PageSize),
	}, nil
}

func (h *NASGrpcHandler) UpdateNAS(
	ctx context.Context,
	req *nas.UpdateNASRequest,
) (*nas.UpdateNASResponse, error) {
	ports := int(req.Ports)
	updateReq := &dto.UpdateNASRequest{
		NASName:         req.Nasname,
		ShortName:       req.Shortname,
		Type:            req.Type,
		Ports:           &ports,
		Secret:          req.Secret,
		Server:          req.Server,
		Community:       req.Community,
		Description:     req.Description,
		RequireMa:       req.RequireMa,
		LimitProxyState: req.LimitProxyState,
	}

	nasResponse, err := h.nasService.UpdateNAS(uint(req.Id), updateReq)
	if err != nil {
		h.logger.Error("Failed to update NAS via gRPC", zap.Uint32("id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to update NAS: %v", err)
	}

	return &nas.UpdateNASResponse{
		Nas: h.toProtoNAS(nasResponse),
	}, nil
}

func (h *NASGrpcHandler) DeleteNAS(
	ctx context.Context,
	req *nas.DeleteNASRequest,
) (*nas.DeleteNASResponse, error) {
	err := h.nasService.DeleteNAS(uint(req.Id))
	if err != nil {
		h.logger.Error("Failed to delete NAS via gRPC", zap.Uint32("id", req.Id), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to delete NAS: %v", err)
	}

	return &nas.DeleteNASResponse{
		Success: true,
	}, nil
}

func (h *NASGrpcHandler) toProtoNAS(n *dto.NASResponse) *nas.NAS {
	var ports int32
	if n.Ports != nil {
		ports = int32(*n.Ports)
	}

	// Parse CreatedAt string to time.Time
	createdAt, _ := time.Parse("2006-01-02 15:04:05", n.CreatedAt)
	updatedAt, _ := time.Parse("2006-01-02 15:04:05", n.UpdatedAt)

	return &nas.NAS{
		Id:              uint32(n.ID),
		Nasname:         n.NASName,
		Shortname:       n.ShortName,
		Type:            n.Type,
		Ports:           ports,
		Secret:          n.Secret,
		Server:          n.Server,
		Community:       n.Community,
		Description:     n.Description,
		RequireMa:       n.RequireMa,
		LimitProxyState: n.LimitProxyState,
		CreatedAt:       timestamppb.New(createdAt),
		UpdatedAt:       timestamppb.New(updatedAt),
	}
}
