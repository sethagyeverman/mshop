package data

import (
	goodsV1 "mshop/service/goods/api/goods/v1"
	"mshop/service/inventory/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redsync/redsync/v4"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewDB,
	NewGoodsServiceClient,
	NewRedisClient,
	NewRedsync,
)

// Data .
type Data struct {
	GoodsClient goodsV1.GoodsClient
	RDB         *redis.Client
	RS          *redsync.Redsync
}

// NewData .
func NewData(
	c *conf.Data,
	goodsClient goodsV1.GoodsClient,
	rdb *redis.Client,
	rs *redsync.Redsync,
	logger log.Logger,
) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	d := &Data{
		GoodsClient: goodsClient,
		RDB:         rdb,
		RS:          rs,
	}

	return d, cleanup, nil
}
