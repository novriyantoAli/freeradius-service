package service_test

import (
	"testing"

	"github.com/novriyantoAli/freeradius-service/internal/application/auth/dto"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/service"
	radcheck "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/entity"
	radcheckdto "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/dto"
	radreply "github.com/novriyantoAli/freeradius-service/internal/application/radreply/entity"
	radrepldto "github.com/novriyantoAli/freeradius-service/internal/application/radreply/dto"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/testutil"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Authenticate_Success(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.GetAllFn = func(filter *radcheckdto.RadcheckFilter) ([]radcheck.Radcheck, int64, error) {
		return []radcheck.Radcheck{*testutil.CreateRadcheckFixture()}, 1, nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()
	mockRadreplyRepo.GetAllFn = func(filter *radrepldto.RadreplyFilter) ([]radreply.Radreply, int64, error) {
		return []radreply.Radreply{*testutil.CreateRadreplyFixture()}, 1, nil
	}

	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo)

	req := &dto.AuthenticateRequest{
		Username: "testuser",
		Password: "testing123",
	}

	result, err := authService.Authenticate(req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.Success)
	require.Equal(t, "Authentication successful", result.Message)
	require.Equal(t, "testuser", result.User.Username)
}

func TestAuthService_Authenticate_InvalidPassword(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.GetAllFn = func(filter *radcheckdto.RadcheckFilter) ([]radcheck.Radcheck, int64, error) {
		return []radcheck.Radcheck{*testutil.CreateRadcheckFixture()}, 1, nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()
	mockRadreplyRepo.GetAllFn = func(filter *radrepldto.RadreplyFilter) ([]radreply.Radreply, int64, error) {
		return []radreply.Radreply{}, 0, nil
	}

	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo)

	req := &dto.AuthenticateRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	result, err := authService.Authenticate(req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.Success)
	require.Equal(t, "Invalid credentials", result.Message)
}

func TestAuthService_Authenticate_UserNotFound(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.GetAllFn = func(filter *radcheckdto.RadcheckFilter) ([]radcheck.Radcheck, int64, error) {
		return []radcheck.Radcheck{}, 0, nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()
	mockRadreplyRepo.GetAllFn = func(filter *radrepldto.RadreplyFilter) ([]radreply.Radreply, int64, error) {
		return []radreply.Radreply{}, 0, nil
	}

	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo)

	req := &dto.AuthenticateRequest{
		Username: "nonexistent",
		Password: "anypassword",
	}

	result, err := authService.Authenticate(req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.False(t, result.Success)
	require.Equal(t, "User not found", result.Message)
}

func TestAuthService_Authenticate_MultipleAttributes(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.GetAllFn = func(filter *radcheckdto.RadcheckFilter) ([]radcheck.Radcheck, int64, error) {
		return []radcheck.Radcheck{
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
	mockRadreplyRepo.GetAllFn = func(filter *radrepldto.RadreplyFilter) ([]radreply.Radreply, int64, error) {
		return []radreply.Radreply{}, 0, nil
	}

	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo)

	req := &dto.AuthenticateRequest{
		Username: "testuser",
		Password: "testing123",
	}

	result, err := authService.Authenticate(req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.Success)
	require.Greater(t, len(result.User.Attributes), 0)
}

func TestAuthService_Authenticate_WithReplyAttributes(t *testing.T) {
	mockRadcheckRepo := testutil.NewMockRadcheckRepositoryWithFn()
	mockRadcheckRepo.GetAllFn = func(filter *radcheckdto.RadcheckFilter) ([]radcheck.Radcheck, int64, error) {
		return []radcheck.Radcheck{*testutil.CreateRadcheckFixture()}, 1, nil
	}

	mockRadreplyRepo := testutil.NewMockRadreplyRepository()
	mockRadreplyRepo.GetAllFn = func(filter *radrepldto.RadreplyFilter) ([]radreply.Radreply, int64, error) {
		return []radreply.Radreply{
			{
				ID:        1,
				Username:  "testuser",
				Attribute: "Reply-Message",
				Op:        "=",
				Value:     "Welcome",
			},
		}, 1, nil
	}

	authService := service.NewAuthService(mockRadcheckRepo, mockRadreplyRepo)

	req := &dto.AuthenticateRequest{
		Username: "testuser",
		Password: "testing123",
	}

	result, err := authService.Authenticate(req)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.True(t, result.Success)
}
