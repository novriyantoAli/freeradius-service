package radreply

import (
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/handler"
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/repository"
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply/service"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(repository.NewRadreplyRepository),
	fx.Provide(service.NewRadreplyService),
	fx.Provide(handler.NewRadreplyHandler),
)
