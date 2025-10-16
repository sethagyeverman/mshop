package biz

import (
	goodsV1 "mshop/service/goods/api/goods/v1"
	inventoryV1 "mshop/service/inventory/api/inventory/v1"
	"mshop/service/order/internal/data"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewOrderUsecase)

type OrderUsecase struct {
	db              *gorm.DB
	rdb             *redis.Client
	log             *log.Helper
	goodsClient     goodsV1.GoodsClient
	inventoryClient inventoryV1.InventoryClient
}

func NewOrderUsecase(db *gorm.DB, data *data.Data, logger log.Logger) *OrderUsecase {
	return &OrderUsecase{
		db:              db,
		rdb:             data.RDB,
		log:             log.NewHelper(logger),
		goodsClient:     data.GoodsClient,
		inventoryClient: data.InventoryClient,
	}
}
