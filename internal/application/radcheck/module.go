package radcheck

import (
	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/handler"
	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/repository"
	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck/service"

	"go.uber.org/fx"
)

// Module provides all radcheck domain dependencies
var Module = fx.Options(
	fx.Provide(
		repository.NewRadcheckRepository,
		service.NewRadcheckService,
		handler.NewRadcheckHandler,
	),
)

// WorkerModule provides only worker dependencies for worker api
var WorkerModule = fx.Options(
	fx.Provide(
		repository.NewRadcheckRepository,
		service.NewRadcheckService,
	),
)
