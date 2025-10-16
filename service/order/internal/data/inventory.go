package data

import (
	"context"
	"time"

	inventoryV1 "mshop/service/inventory/api/inventory/v1"
	"mshop/service/order/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewInventoryServiceClient 创建 Inventory 服务客户端
func NewInventoryServiceClient(conf *conf.Services, logger log.Logger) (inventoryV1.InventoryClient, error) {
	l := log.NewHelper(logger)

	// 设置超时时间
	timeout := 5 * time.Second
	if conf.Inventory != nil && conf.Inventory.Timeout != nil {
		timeout = conf.Inventory.Timeout.AsDuration()
	}

	// 创建 gRPC 连接
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(conf.Inventory.Endpoint),
		grpc.WithTimeout(timeout),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		l.Errorf("Failed to connect to inventory service: %v", err)
		return nil, err
	}

	l.Infof("Connected to inventory service at: %s", conf.Inventory.Endpoint)

	// 创建 Inventory 客户端
	return inventoryV1.NewInventoryClient(conn), nil
}
