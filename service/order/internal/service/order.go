package service

import (
	"context"

	pb "mshop/service/order/api/order/v1"
	"mshop/service/order/internal/biz"
)

type OrderService struct {
	pb.UnimplementedOrderServer
	orderUsecase *biz.OrderUsecase
}

func NewOrderService(orderUsecase *biz.OrderUsecase) *OrderService {
	return &OrderService{
		orderUsecase: orderUsecase,
	}
}

// 购物车相关方法
func (s *OrderService) CartItemList(ctx context.Context, req *pb.UserInfo) (*pb.CartItemListResponse, error) {
	return s.orderUsecase.CartItemList(ctx, req)
}

func (s *OrderService) CreateCartItem(ctx context.Context, req *pb.CartItemRequest) (*pb.ShopCartInfoResponse, error) {
	return s.orderUsecase.CreateCartItem(ctx, req)
}

func (s *OrderService) UpdateCartItem(ctx context.Context, req *pb.CartItemRequest) (*pb.Empty, error) {
	return s.orderUsecase.UpdateCartItem(ctx, req)
}

func (s *OrderService) DeleteCartItem(ctx context.Context, req *pb.CartItemRequest) (*pb.Empty, error) {
	return s.orderUsecase.DeleteCartItem(ctx, req)
}

// 订单相关方法
func (s *OrderService) CreateOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderInfoResponse, error) {
	return s.orderUsecase.CreateOrder(ctx, req)
}

func (s *OrderService) OrderList(ctx context.Context, req *pb.OrderFilterRequest) (*pb.OrderListResponse, error) {
	return s.orderUsecase.OrderList(ctx, req)
}

func (s *OrderService) OrderDetail(ctx context.Context, req *pb.OrderRequest) (*pb.OrderInfoDetailResponse, error) {
	return s.orderUsecase.OrderDetail(ctx, req)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, req *pb.OrderStatus) (*pb.Empty, error) {
	return s.orderUsecase.UpdateOrderStatus(ctx, req)
}
