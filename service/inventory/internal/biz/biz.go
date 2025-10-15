package biz

import (
	goodsV1 "mshop/service/goods/api/goods/v1"
	"mshop/service/inventory/internal/data"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redsync/redsync/v4"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewInventoryUsecase)

type InventoryUsecase struct {
	db          *gorm.DB
	log         *log.Helper
	goodsClient goodsV1.GoodsClient
	rs          *redsync.Redsync // 分布式锁管理器
}

func NewInventoryUsecase(db *gorm.DB, data *data.Data, logger log.Logger) *InventoryUsecase {
	return &InventoryUsecase{
		db:          db,
		log:         log.NewHelper(logger),
		goodsClient: data.GoodsClient,
		rs:          data.RS,
	}
}
