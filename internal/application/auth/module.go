package auth

import (
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/handler"
	"github.com/novriyantoAli/freeradius-service/internal/application/auth/service"
	radcheckrepo "github.com/novriyantoAli/freeradius-service/internal/application/radcheck/repository"
	radreplyrepo "github.com/novriyantoAli/freeradius-service/internal/application/radreply/repository"
	"github.com/novriyantoAli/freeradius-service/internal/pkg/database"
	"go.uber.org/fx"
)

// Module provides authentication dependencies
var Module = fx.Module("auth",
	fx.Provide(
		provideAuthService,
		provideAuthHandler,
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
