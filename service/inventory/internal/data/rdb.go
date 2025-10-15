package data

import (
	"context"
	"mshop/service/inventory/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redsync/redsync/v4"
	goredislib "github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

// NewRedisClient 创建 Redis 客户端
func NewRedisClient(conf *conf.Data, logger log.Logger) (*redis.Client, func(), error) {
	l := log.NewHelper(logger)

	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
	})

	// 测试连接
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		l.Errorf("Failed to connect to Redis: %v", err)
		return nil, nil, err
	}

	cleanup := func() {
		l.Info("Closing Redis connection")
		rdb.Close()
	}

	l.Infof("Connected to Redis at: %s", conf.Redis.Addr)
	return rdb, cleanup, nil
}

// NewRedsync 创建分布式锁管理器
func NewRedsync(rdb *redis.Client) *redsync.Redsync {
	pool := goredislib.NewPool(rdb)
	return redsync.New(pool)
}
