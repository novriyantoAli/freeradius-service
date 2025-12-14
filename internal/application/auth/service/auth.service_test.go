package service_test

import (
	"context"
	"testing"

	"github.com/novriyantoAli/freeradius-service/internal/application/auth/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/service"
	radcheckEntity "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/entity"
	radreplyEntity "github.com/novriyantoAli/freeradius-service/internal/application/radreply/entity"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestAuthService_CreateAuth_Success(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.CreateFn = func(ctx context.Context, radcheck *radcheckEntity.Radcheck) error {
		radcheck.ID = uint(len([]radcheckEntity.Radcheck{}) + 1)
		return nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()
	mockRadreplyRepo.CreateFn = func(ctx context.Context, radreply *radreplyEntity.Radreply) error {
		radreply.ID = uint(len([]radreplyEntity.Radreply{}) + 1)
		return nil
	}

	mockTxManager := &testutil.MockTransactionManager{}
	mockTxManager.WithinTransactionFn = func(ctx context.Context, fn func(ctx context.Context) error) error {
		return fn(ctx)
	}

	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo, mockTxManager)

	req := &dto.CreateAuthRequest{
		Username: "newuser",
		Password: "password123",
		Attributes: []dto.CreateAuthAttribute{
			{Attribute: "Framed-IP-Address", Value: "192.168.1.100", Op: "="},
		},
		ReplyAttrs: []dto.CreateAuthAttribute{
			{Attribute: "Reply-Message", Value: "Welcome", Op: "="},
		},
	}

	result, err := authService.CreateAuth(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, "newuser", result.Username)
	require.Equal(t, "password123", result.Password)
	require.Greater(t, len(result.Attributes), 0)
	require.Greater(t, len(result.ReplyAttrs), 0)
}

func TestAuthService_CreateAuth_MissingUsername(t *testing.T) {
	mockTxManager := &testutil.MockTransactionManager{}
	authService := service.NewAuthService(nil, nil, mockTxManager)

	req := &dto.CreateAuthRequest{
		Username: "",
		Password: "password123",
	}

	result, err := authService.CreateAuth(context.Background(), req)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, "username is required", err.Error())
}

func TestAuthService_CreateAuth_MissingPassword(t *testing.T) {
	mockTxManager := &testutil.MockTransactionManager{}
	authService := service.NewAuthService(nil, nil, mockTxManager)

	req := &dto.CreateAuthRequest{
		Username: "newuser",
		Password: "",
	}

	result, err := authService.CreateAuth(context.Background(), req)

	require.Error(t, err)
	require.Nil(t, result)
	require.Equal(t, "password is required", err.Error())
}
