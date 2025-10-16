package service

import (
	"context"

	pb "mshop/service/inventory/api/inventory/v1"
	"mshop/service/inventory/internal/biz"
)

type InventoryService struct {
	pb.UnimplementedInventoryServer
	inventoryUsecase *biz.InventoryUsecase
}

func NewInventoryService(inventoryUsecase *biz.InventoryUsecase) *InventoryService {
	return &InventoryService{
		inventoryUsecase: inventoryUsecase,
	}
}

func (s *InventoryService) SetInv(ctx context.Context, req *pb.GoodsInvInfo) (*pb.Empty, error) {
	return s.inventoryUsecase.SetInv(ctx, req)
}
func (s *InventoryService) InvDetail(ctx context.Context, req *pb.GoodsInvInfo) (*pb.GoodsInvInfo, error) {
	return s.inventoryUsecase.InvDetail(ctx, req)
}
func (s *InventoryService) Sell(ctx context.Context, req *pb.SellInfo) (*pb.Empty, error) {
	return s.inventoryUsecase.Sell(ctx, req)
}
func (s *InventoryService) Reback(ctx context.Context, req *pb.SellInfo) (*pb.Empty, error) {
	return s.inventoryUsecase.Reback(ctx, req)
}
