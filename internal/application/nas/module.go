package nas

import (
	"github.com/novriyantoAli/freeradius-service/internal/application/nas/handler"
	"github.com/novriyantoAli/freeradius-service/internal/application/nas/repository"
	"github.com/novriyantoAli/freeradius-service/internal/application/nas/service"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var NASModule = fx.Module(
	"nas",
	fx.Provide(
		NewNASRepository,
		service.NewNASService,
		handler.NewNASHandler,
	),
)

func NewNASRepository(db *gorm.DB, logger *zap.Logger) repository.NASRepository {
	return repository.NewNASRepository(db, logger)
}
