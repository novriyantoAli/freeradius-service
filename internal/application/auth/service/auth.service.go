package service

import (
	"context"
	"errors"

	"github.com/novriyantoAli/freeradius-service/internal/application/auth/dto"
	radcheckentity "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/entity"
	radcheckrepo "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/repository"
	radreplyentity "github.com/novriyantoAli/freeradius-service/internal/application/radreply/entity"
	radreplyrepo "github.com/novriyantoAli/freeradius-service/internal/application/radreply/repository"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/database"
)

// AuthService defines authentication business logic
type AuthService interface {
	CreateAuth(ctx context.Context, req *dto.CreateAuthRequest) (*dto.CreateAuthResponse, error)
}

type authService struct {
	radcheckRepo radcheckrepo.RadcheckRepository
	radreplyRepo radreplyrepo.RadreplyRepository
	txManager    database.TransactionManagerI
}

// NewAuthService creates a new authentication service
func NewAuthService(
	radcheckRepo radcheckrepo.RadcheckRepository,
	radreplyRepo radreplyrepo.RadreplyRepository,
	txManager database.TransactionManagerI,
) AuthService {
	return &authService{
		radcheckRepo: radcheckRepo,
		radreplyRepo: radreplyRepo,
		txManager:    txManager,
	}
}

// CreateAuth creates authentication credentials with radcheck and radreply entries in a transaction
func (s *authService) CreateAuth(ctx context.Context, req *dto.CreateAuthRequest) (*dto.CreateAuthResponse, error) {
	if req.Username == "" {
		return nil, errors.New("username is required")
	}
	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	var response dto.CreateAuthResponse
	response.Username = req.Username
	response.Password = req.Password

	// Execute in transaction
	err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// Create User-Password radcheck entry
		passwordRadcheck := &radcheckentity.Radcheck{
			Username:  req.Username,
			Attribute: "User-Password",
			Op:        ":=",
			Value:     req.Password,
		}

		if err := s.radcheckRepo.Create(txCtx, passwordRadcheck); err != nil {
			return err
		}

		response.Attributes = append(response.Attributes, dto.AuthCreateAttrResponse{
			ID:        passwordRadcheck.ID,
			Attribute: passwordRadcheck.Attribute,
			Value:     "***",
			Op:        passwordRadcheck.Op,
		})

		// Create additional radcheck attributes
		for _, attr := range req.Attributes {
			if attr.Attribute == "User-Password" {
				continue // Skip, already created
			}

			op := attr.Op
			if op == "" {
				op = ":="
			}

			radcheck := &radcheckentity.Radcheck{
				Username:  req.Username,
				Attribute: attr.Attribute,
				Op:        op,
				Value:     attr.Value,
			}

			if err := s.radcheckRepo.Create(txCtx, radcheck); err != nil {
				return err
			}

			response.Attributes = append(response.Attributes, dto.AuthCreateAttrResponse{
				ID:        radcheck.ID,
				Attribute: radcheck.Attribute,
				Value:     radcheck.Value,
				Op:        radcheck.Op,
			})
		}

		// Create radreply attributes
		for _, attr := range req.ReplyAttrs {
			op := attr.Op
			if op == "" {
				op = "+=" // Default for radreply
			}

			radreply := &radreplyentity.Radreply{
				Username:  req.Username,
				Attribute: attr.Attribute,
				Op:        op,
				Value:     attr.Value,
			}

			if err := s.radreplyRepo.Create(txCtx, radreply); err != nil {
				return err
			}

			response.ReplyAttrs = append(response.ReplyAttrs, dto.AuthCreateAttrResponse{
				ID:        radreply.ID,
				Attribute: radreply.Attribute,
				Value:     radreply.Value,
				Op:        radreply.Op,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &response, nil
}
