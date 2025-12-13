package nas

import (
	"github.com/novriyantoAli/freeradius-service/internal/application/nas/handler"
	"github.com/novriyantoAli/freeradius-service/internal/application/nas/repository"
	"github.com/novriyantoAli/freeradius-service/internal/application/nas/service"

	"go.uber.org/fx"
)

// Module provides all NAS domain dependencies
var Module = fx.Options(
	fx.Provide(
		repository.NewNASRepository,
		service.NewNASService,
		handler.NewNASHandler,
	),
)

// WorkerModule provides only worker dependencies for worker api
var WorkerModule = fx.Options(
	fx.Provide(
		repository.NewNASRepository,
		service.NewNASService,
	),
)
