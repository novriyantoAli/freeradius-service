package auth

import (
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/handler"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/service"
	radcheckrepo "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/repository"
	radreplyrepo "github.com/novriyantoAli/freeradius-service/internal/application/radreply/repository"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/database"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides authentication dependencies
var Module = fx.Module("auth",
	fx.Provide(
		provideAuthService,
		provideAuthHandler,
		provideAuthGrpcHandler,
	),
)

func provideAuthService(
	radcheckRepo radcheckrepo.RadcheckRepository,
	radreplyRepo radreplyrepo.RadreplyRepository,
	txManager database.TransactionManagerI,
) service.AuthService {
	return service.NewAuthService(radcheckRepo, radreplyRepo, txManager)
}

func provideAuthHandler(authService service.AuthService) *handler.AuthHandler {
	return handler.NewAuthHandler(authService)
}

func provideAuthGrpcHandler(authService service.AuthService, logger *zap.Logger) *handler.AuthGrpcHandler {
	return handler.NewAuthGrpcHandler(authService, logger)
}
