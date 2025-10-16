package biz

import (
	"context"
	"fmt"
	"time"

	"mshop/pkg/errx"
	goodsV1 "mshop/service/goods/api/goods/v1"
	pb "mshop/service/inventory/api/inventory/v1"

	"github.com/go-redsync/redsync/v4"
)

func (uc *InventoryUsecase) SetInv(ctx context.Context, req *pb.GoodsInvInfo) (_ *pb.Empty, err error) {

	var inventory Inventory
	now := time.Now()
	// 已经存在
	if result := uc.db.Where("goods_id = ?", req.GoodsId).First(&inventory); result.RowsAffected != 0 {
		inventory.Stock = req.Num
		inventory.UpdateTime = now
		uc.db.Save(&inventory)
		return
	}

	// 不存在
	if _, err := uc.goodsClient.GetGoodsDetail(ctx, &goodsV1.GoodInfoRequest{
		Id: req.GoodsId,
	}); err != nil {
		return nil, err
	}

	inventory.GoodsId = req.GoodsId
	inventory.Stock = req.Num
	inventory.UpdateTime = now
	inventory.AddTime = now
	uc.db.Create(&inventory)

	return
}

func (uc *InventoryUsecase) InvDetail(ctx context.Context, req *pb.GoodsInvInfo) (resp *pb.GoodsInvInfo, err error) {
	var inventory Inventory
	if result := uc.db.Where("goods_id = ?", req.GoodsId).First(&inventory); result.RowsAffected == 0 {
		return nil, errx.ErrorInventoryNotFound("inventory not found")
	}

	return &pb.GoodsInvInfo{
		GoodsId: inventory.GoodsId,
		Num:     inventory.Stock,
	}, nil
}

func (uc *InventoryUsecase) Sell(ctx context.Context, req *pb.SellInfo) (_ *pb.Empty, err error) {

	for _, good := range req.GoodsInfo {
		// 创建分布式锁，key 为 goods_id
		lockKey := fmt.Sprintf("inventory:lock:goods:%d", good.GoodsId)
		mutex := uc.rs.NewMutex(lockKey,
			redsync.WithExpiry(10*time.Second),           // 锁过期时间
			redsync.WithTries(3),                         // 重试次数
			redsync.WithRetryDelay(100*time.Millisecond), // 重试间隔
		)

		// 获取锁
		if err := mutex.LockContext(ctx); err != nil {
			uc.log.Errorf("Failed to acquire lock for goods %d: %v", good.GoodsId, err)
			return nil, errx.ErrorInventoryLockFailed("failed to acquire lock for goods %d", good.GoodsId)
		}

		// 在锁保护下执行库存扣减
		var inv Inventory
		if result := uc.db.Where("goods_id = ?", good.GoodsId).First(&inv); result.RowsAffected == 0 {
			return nil, errx.ErrorInventoryNotFound("goods id %d not found", good.GoodsId)
		}

		if inv.Stock < good.Num {
			return nil, errx.ErrorInventoryInsufficient("goods id %d inventory insufficient", good.GoodsId)
		}

		// 直接扣减库存（有分布式锁保护，不需要乐观锁）
		inv.Stock -= good.Num
		if err := uc.db.Save(&inv).Error; err != nil {
			return nil, err
		}

		if _, err := mutex.UnlockContext(ctx); err != nil {
			uc.log.Errorf("Failed to release lock: %v", err)
		}
	}

	return &pb.Empty{}, nil
}

func (uc *InventoryUsecase) Reback(ctx context.Context, req *pb.SellInfo) (_ *pb.Empty, err error) {
	tx := uc.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	for _, good := range req.GoodsInfo {
		lockKey := fmt.Sprintf("inventory:lock:goods:%d", good.GoodsId)
		mutex := uc.rs.NewMutex(lockKey,
			redsync.WithExpiry(10*time.Second),           // 锁过期时间
			redsync.WithTries(3),                         // 重试次数
			redsync.WithRetryDelay(100*time.Millisecond), // 重试间隔
		)

		if err := mutex.LockContext(ctx); err != nil {
			var inv Inventory
			if result := tx.Where("goods_id = ?", good.GoodsId).First(&inv); result.RowsAffected == 0 {
				return nil, errx.ErrorInventoryNotFound("goods id %d not found", good.GoodsId)
			}

			tx.Where("goods_id = ?", good.GoodsId).Updates(&Inventory{
				Stock: inv.Stock + good.Num,
			})

			if _, err := mutex.UnlockContext(ctx); err != nil {
				uc.log.Errorf("Failed to release lock: %v", err)
			}
		}
	}
	return

}
