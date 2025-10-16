package data

import (
	goodsV1 "mshop/service/goods/api/goods/v1"
	inventoryV1 "mshop/service/inventory/api/inventory/v1"
	"mshop/service/order/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewDB,
	NewGoodsServiceClient,
	NewInventoryServiceClient,
	NewRedisClient,
)

// Data .
type Data struct {
	GoodsClient     goodsV1.GoodsClient
	InventoryClient inventoryV1.InventoryClient
	RDB             *redis.Client
}

// NewData .
func NewData(
	c *conf.Data,
	goodsClient goodsV1.GoodsClient,
	inventoryClient inventoryV1.InventoryClient,
	rdb *redis.Client,
	logger log.Logger,
) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	d := &Data{
		GoodsClient:     goodsClient,
		InventoryClient: inventoryClient,
		RDB:             rdb,
	}

	return d, cleanup, nil
}
