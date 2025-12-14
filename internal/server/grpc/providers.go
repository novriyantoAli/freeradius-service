package grpc

import (
	"github.com/novriyantoAli/freeradius-service/internal/application/auth"
	authHandler "github.com/novriyantoAli/freeradius-service/internal/application/auth/handler"
	"github.com/novriyantoAli/freeradius-service/internal/application/payment"
	paymentHandler "github.com/novriyantoAli/freeradius-service/internal/application/payment/handler"
	"github.com/novriyantoAli/freeradius-service/internal/application/user"
	userHandler "github.com/novriyantoAli/freeradius-service/internal/application/user/handler"

	"go.uber.org/fx"
)

var Module = fx.Options(
	// Include domain modules
	auth.Module,
	user.Module,
	payment.Module,

	// gRPC handlers
	fx.Provide(
		authHandler.NewAuthGrpcHandler,
		userHandler.NewUserGrpcHandler,
		paymentHandler.NewPaymentGrpcHandler,
		NewServer,
	),
)
