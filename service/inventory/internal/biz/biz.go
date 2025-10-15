package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewInventoryUsecase)

type InventoryUsecase struct {
	db  *gorm.DB
	log *log.Helper
}

func NewInventoryUsecase(db *gorm.DB, logger log.Logger) *InventoryUsecase {
	return &InventoryUsecase{
		db:  db,
		log: log.NewHelper(logger),
	}
}
