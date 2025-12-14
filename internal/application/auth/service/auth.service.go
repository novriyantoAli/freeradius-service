package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/novriyantoAli/freeradius-service/internal/application/auth/dto"
	radcheckdto "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/dto"
	radcheckentity "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/entity"
	radcheckrepo "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/repository"
	radrepldto "github.com/novriyantoAli/freeradius-service/internal/application/radreply/dto"
	radreplyentity "github.com/novriyantoAli/freeradius-service/internal/application/radreply/entity"
	radreplyrepo "github.com/novriyantoAli/freeradius-service/internal/application/radreply/repository"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/database"
)

// AuthService defines authentication business logic
type AuthService interface {
	Authenticate(ctx context.Context, req *dto.AuthenticateRequest) (*dto.AuthenticateResponse, error)
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

// Authenticate verifies user credentials against radcheck entries
func (s *authService) Authenticate(ctx context.Context, req *dto.AuthenticateRequest) (*dto.AuthenticateResponse, error) {
	// Get all radcheck entries for the user
	filter := &radcheckdto.RadcheckFilter{
		Username: req.Username,
	}
	radchecks, _, err := s.radcheckRepo.GetAll(ctx, filter)
	if err != nil {
		return &dto.AuthenticateResponse{
			Success: false,
			Message: "Failed to retrieve user authentication attributes",
		}, nil
	}

	if len(radchecks) == 0 {
		return &dto.AuthenticateResponse{
			Success: false,
			Message: "User not found",
		}, nil
	}

	// Verify password against User-Password attribute
	authenticated := false
	for _, radcheck := range radchecks {
		if radcheck.Attribute == "User-Password" {
			if s.verifyPassword(req.Password, radcheck.Value) {
				authenticated = true
				break
			}
		}
	}

	if !authenticated {
		return &dto.AuthenticateResponse{
			Success: false,
			Message: "Invalid credentials",
		}, nil
	}

	// Collect user attributes from radcheck
	var userAttrs []dto.AttrValue
	for _, radcheck := range radchecks {
		if radcheck.Attribute != "User-Password" {
			userAttrs = append(userAttrs, dto.AttrValue{
				Attribute: radcheck.Attribute,
				Value:     radcheck.Value,
			})
		}
	}

	// Get reply attributes for the user
	replyFilter := &radrepldto.RadreplyFilter{
		Username: req.Username,
	}
	radeplies, _, err := s.radreplyRepo.GetAll(ctx, replyFilter)
	if err != nil {
		radeplies = []radreplyentity.Radreply{}
	}

	// Convert reply attributes
	var replies []dto.ReplyAttr
	for _, radReply := range radeplies {
		replies = append(replies, dto.ReplyAttr{
			Attribute: radReply.Attribute,
			Value:     radReply.Value,
		})
	}

	return &dto.AuthenticateResponse{
		Success: true,
		Message: "Authentication successful",
		User: dto.UserAuth{
			Username:   req.Username,
			Attributes: userAttrs,
		},
		Replies: replies,
	}, nil
}

// verifyPassword compares plaintext password with stored password
// Supports MD5 and plaintext passwords
func (s *authService) verifyPassword(plaintext, stored string) bool {
	// Check plaintext
	if plaintext == stored {
		return true
	}

	// Check MD5 hash
	hash := sha1.Sum([]byte(plaintext))
	hashedPlaintext := fmt.Sprintf("{SHA}%s", hex.EncodeToString(hash[:]))
	if hashedPlaintext == stored {
		return true
	}

	// Check unsalted MD5
	if fmt.Sprintf("%x", sha1.Sum([]byte(plaintext))) == stored {
		return true
	}

	return false
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