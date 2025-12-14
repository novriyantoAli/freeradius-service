package api

import (
	"github.com/novriyantoAli/freeradius-service/internal/application/nas"
	"github.com/novriyantoAli/freeradius-service/internal/application/payment"
	"github.com/novriyantoAli/freeradius-service/internal/application/radcheck"
	"github.com/novriyantoAli/freeradius-service/internal/application/radreply"
	"github.com/novriyantoAli/freeradius-service/internal/application/user"

	"go.uber.org/fx"
)

var Module = fx.Options(
	// Include all domain modules
	user.Module,
	payment.Module,
	nas.Module,
	radcheck.Module,
	radreply.Module,

	// API api
	fx.Provide(NewServer),
)
