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
	return &pb.Empty{}, nil
}
func (s *InventoryService) InvDetail(ctx context.Context, req *pb.GoodsInvInfo) (*pb.GoodsInvInfo, error) {
	return &pb.GoodsInvInfo{}, nil
}
func (s *InventoryService) Sell(ctx context.Context, req *pb.SellInfo) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
func (s *InventoryService) Reback(ctx context.Context, req *pb.SellInfo) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
