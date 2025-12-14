package service_test

import (
	"context"
	"testing"

	"github.com/novriyantoAli/freeradius-service/internal/application/auth/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/service"
	radcheckdto "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/dto"
	radcheckEntity "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/entity"
	radrepldto "github.com/novriyantoAli/freeradius-service/internal/application/radreply/dto"
	radreplyEntity "github.com/novriyantoAli/freeradius-service/internal/application/radreply/entity"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Authenticate_Success(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.GetAllFn = func(ctx context.Context, filter *radcheckdto.RadcheckFilter) ([]radcheckEntity.Radcheck, int64, error) {
		return []radcheckEntity.Radcheck{*testutil.CreateRadcheckFixture()}, 1, nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()
	mockRadreplyRepo.GetAllFn = func(ctx context.Context, filter *radrepldto.RadreplyFilter) ([]radreplyEntity.Radreply, int64, error) {
		return []radreplyEntity.Radreply{*testutil.CreateRadreplyFixture()}, 1, nil
	}

	mockTxManager := &testutil.MockTransactionManager{}

	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo, mockTxManager)

	req := &dto.AuthenticateRequest{
		Username: "testuser",
		Password: "testing123",
	}

	result, err := authService.Authenticate(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.Success)
	require.Equal(t, "Authentication successful", result.Message)
	require.Equal(t, "testuser", result.User.Username)
}

func TestAuthService_Authenticate_InvalidPassword(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.GetAllFn = func(ctx context.Context, filter *radcheckdto.RadcheckFilter) ([]radcheckEntity.Radcheck, int64, error) {
		return []radcheckEntity.Radcheck{*testutil.CreateRadcheckFixture()}, 1, nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()
	mockRadreplyRepo.GetAllFn = func(ctx context.Context, filter *radrepldto.RadreplyFilter) ([]radreplyEntity.Radreply, int64, error) {
		return []radreplyEntity.Radreply{}, 0, nil
	}

	mockTxManager := &testutil.MockTransactionManager{}
	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo, mockTxManager)

	req := &dto.AuthenticateRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	result, err := authService.Authenticate(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.Success)
	require.Equal(t, "Invalid credentials", result.Message)
}

func TestAuthService_Authenticate_UserNotFound(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.GetAllFn = func(ctx context.Context, filter *radcheckdto.RadcheckFilter) ([]radcheckEntity.Radcheck, int64, error) {
		return []radcheckEntity.Radcheck{}, 0, nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()
	mockRadreplyRepo.GetAllFn = func(ctx context.Context, filter *radrepldto.RadreplyFilter) ([]radreplyEntity.Radreply, int64, error) {
		return []radreplyEntity.Radreply{}, 0, nil
	}

	mockTxManager := &testutil.MockTransactionManager{}
	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo, mockTxManager)

	req := &dto.AuthenticateRequest{
		Username: "nonexistent",
		Password: "anypassword",
	}

	result, err := authService.Authenticate(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.Success)
	require.Equal(t, "User not found", result.Message)
}

func TestAuthService_Authenticate_MultipleAttributes(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.GetAllFn = func(ctx context.Context, filter *radcheckdto.RadcheckFilter) ([]radcheckEntity.Radcheck, int64, error) {
		return []radcheckEntity.Radcheck{
			*testutil.CreateRadcheckFixture(),
			{
				ID:        2,
				Username:  "testuser",
				Attribute: "NAS-IP-Address",
				Op:        "==",
				Value:     "192.168.1.1",
			},
		}, 2, nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()
	mockRadreplyRepo.GetAllFn = func(ctx context.Context, filter *radrepldto.RadreplyFilter) ([]radreplyEntity.Radreply, int64, error) {
		return []radreplyEntity.Radreply{}, 0, nil
	}

	mockTxManager := &testutil.MockTransactionManager{}
	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo, mockTxManager)

	req := &dto.AuthenticateRequest{
		Username: "testuser",
		Password: "testing123",
	}

	result, err := authService.Authenticate(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.Success)
	require.Greater(t, len(result.User.Attributes), 0)
}

func TestAuthService_Authenticate_WithReplyAttributes(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.GetAllFn = func(ctx context.Context, filter *radcheckdto.RadcheckFilter) ([]radcheckEntity.Radcheck, int64, error) {
		return []radcheckEntity.Radcheck{*testutil.CreateRadcheckFixture()}, 1, nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()
	mockRadreplyRepo.GetAllFn = func(ctx context.Context, filter *radrepldto.RadreplyFilter) ([]radreplyEntity.Radreply, int64, error) {
		return []radreplyEntity.Radreply{
			{
				ID:        1,
				Username:  "testuser",
				Attribute: "Reply-Message",
				Op:        "=",
				Value:     "Welcome",
			},
		}, 1, nil
	}

	mockTxManager := &testutil.MockTransactionManager{}
	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo, mockTxManager)

	req := &dto.AuthenticateRequest{
		Username: "testuser",
		Password: "testing123",
	}

	result, err := authService.Authenticate(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.Success)
}

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

