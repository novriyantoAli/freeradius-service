package handler

import (
	"context"

	"github.com/novriyantoAli/freeradius-service/api/proto/auth"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/service"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthGrpcHandler struct {
	auth.UnimplementedAuthServiceServer
	authService service.AuthService
	logger      *zap.Logger
}

func NewAuthGrpcHandler(authService service.AuthService, logger *zap.Logger) *AuthGrpcHandler {
	return &AuthGrpcHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthGrpcHandler) CreateAuth(
	ctx context.Context,
	req *auth.CreateAuthRequest,
) (*auth.CreateAuthResponse, error) {
	// Validate username and password
	if req.Username == "" {
		h.logger.Warn("CreateAuth gRPC request missing username")
		return nil, status.Errorf(codes.InvalidArgument, "username is required")
	}
	if req.Password == "" {
		h.logger.Warn("CreateAuth gRPC request missing password")
		return nil, status.Errorf(codes.InvalidArgument, "password is required")
	}

	// Convert proto attributes to DTO
	var attributes []dto.CreateAuthAttribute
	for _, attr := range req.Attributes {
		attributes = append(attributes, dto.CreateAuthAttribute{
			Attribute: attr.Attribute,
			Value:     attr.Value,
			Op:        attr.Op,
		})
	}

	// Convert proto reply attributes to DTO
	var replyAttrs []dto.CreateAuthAttribute
	for _, attr := range req.ReplyAttributes {
		replyAttrs = append(replyAttrs, dto.CreateAuthAttribute{
			Attribute: attr.Attribute,
			Value:     attr.Value,
			Op:        attr.Op,
		})
	}

	// Create auth request
	createReq := &dto.CreateAuthRequest{
		Username:   req.Username,
		Password:   req.Password,
		Attributes: attributes,
		ReplyAttrs: replyAttrs,
	}

	// Call service
	authResponse, err := h.authService.CreateAuth(ctx, createReq)
	if err != nil {
		h.logger.Error("Failed to create auth via gRPC", zap.String("username", req.Username), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create auth: %v", err)
	}

	// Convert response to proto
	return h.toProtoCreateAuthResponse(authResponse), nil
}

func (h *AuthGrpcHandler) toProtoCreateAuthResponse(resp *dto.CreateAuthResponse) *auth.CreateAuthResponse {
	// Convert attributes
	attributes := make([]*auth.AuthCreateAttrResponse, len(resp.Attributes))
	for i, attr := range resp.Attributes {
		attributes[i] = &auth.AuthCreateAttrResponse{
			Id:        uint32(attr.ID),
			Attribute: attr.Attribute,
			Value:     attr.Value,
			Op:        attr.Op,
		}
	}

	// Convert reply attributes
	replyAttrs := make([]*auth.AuthCreateAttrResponse, len(resp.ReplyAttrs))
	for i, attr := range resp.ReplyAttrs {
		replyAttrs[i] = &auth.AuthCreateAttrResponse{
			Id:        uint32(attr.ID),
			Attribute: attr.Attribute,
			Value:     attr.Value,
			Op:        attr.Op,
		}
	}

	return &auth.CreateAuthResponse{
		Username:        resp.Username,
		Password:        resp.Password,
		Attributes:      attributes,
		ReplyAttributes: replyAttrs,
	}
}
