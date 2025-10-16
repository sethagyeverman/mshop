package data

import (
	"context"
	"time"

	goodsV1 "mshop/service/goods/api/goods/v1"
	"mshop/service/inventory/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGoodsServiceClient 创建 Goods 服务客户端
func NewGoodsServiceClient(conf *conf.Services, logger log.Logger) (goodsV1.GoodsClient, error) {
	l := log.NewHelper(logger)

	// 设置超时时间
	timeout := 5 * time.Second
	if conf.Goods != nil && conf.Goods.Timeout != nil {
		timeout = conf.Goods.Timeout.AsDuration()
	}

	// 创建 gRPC 连接
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(conf.Goods.Endpoint),
		grpc.WithTimeout(timeout),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		l.Errorf("Failed to connect to goods service: %v", err)
		return nil, err
	}

	l.Infof("Connected to goods service at: %s", conf.Goods.Endpoint)

	// 创建 Goods 客户端
	return goodsV1.NewGoodsClient(conn), nil
}
