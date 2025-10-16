package biz

import (
	"context"
	"time"

	"mshop/pkg/errx"
	pb "mshop/service/order/api/order/v1"
)

// CartItemList 获取用户的购物车列表
func (uc *OrderUsecase) CartItemList(ctx context.Context, req *pb.UserInfo) (*pb.CartItemListResponse, error) {
	var shopCarts []ShoppingCart

	if result := uc.db.Where("user_id = ?", req.Id).Find(&shopCarts); result.Error != nil {
		return nil, result.Error
	}

	rsp := &pb.CartItemListResponse{
		Total: int32(len(shopCarts)),
	}

	for _, cart := range shopCarts {
		rsp.Data = append(rsp.Data, &pb.ShopCartInfoResponse{
			Id:      cart.ID,
			UserId:  cart.UserID,
			GoodsId: cart.GoodsID,
			Nums:    cart.Nums,
			Checked: cart.Checked,
		})
	}

	return rsp, nil
}

// CreateCartItem 添加商品到购物车
func (uc *OrderUsecase) CreateCartItem(ctx context.Context, req *pb.CartItemRequest) (*pb.ShopCartInfoResponse, error) {
	var shopCart ShoppingCart

	// 检查商品是否已经在购物车中
	result := uc.db.Where("user_id = ? and goods_id = ?", req.UserId, req.GoodsId).First(&shopCart)

	if result.RowsAffected == 1 {
		// 如果已存在，更新数量
		shopCart.Nums += req.Nums
		shopCart.UpdateTime = time.Now()
	} else {
		// 不存在，创建新记录
		shopCart.UserID = req.UserId
		shopCart.GoodsID = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
		shopCart.AddTime = time.Now()
		shopCart.UpdateTime = time.Now()
	}

	if err := uc.db.Save(&shopCart).Error; err != nil {
		return nil, errx.ErrorCartCreateFailed("failed to create cart item: %v", err)
	}

	return &pb.ShopCartInfoResponse{
		Id:      shopCart.ID,
		UserId:  shopCart.UserID,
		GoodsId: shopCart.GoodsID,
		Nums:    shopCart.Nums,
		Checked: shopCart.Checked,
	}, nil
}

// UpdateCartItem 更新购物车
func (uc *OrderUsecase) UpdateCartItem(ctx context.Context, req *pb.CartItemRequest) (*pb.Empty, error) {
	var shopCart ShoppingCart

	if result := uc.db.Where("id = ? and user_id = ?", req.Id, req.UserId).First(&shopCart); result.RowsAffected == 0 {
		return nil, errx.ErrorCartNotFound("cart item not found")
	}

	// 更新字段
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	shopCart.Checked = req.Checked
	shopCart.UpdateTime = time.Now()

	if err := uc.db.Save(&shopCart).Error; err != nil {
		return nil, errx.ErrorCartUpdateFailed("failed to update cart item: %v", err)
	}

	return &pb.Empty{}, nil
}

// DeleteCartItem 删除购物车条目
func (uc *OrderUsecase) DeleteCartItem(ctx context.Context, req *pb.CartItemRequest) (*pb.Empty, error) {
	if result := uc.db.Where("id = ? and user_id = ?", req.Id, req.UserId).Delete(&ShoppingCart{}); result.RowsAffected == 0 {
		return nil, errx.ErrorCartNotFound("cart item not found")
	}

	return &pb.Empty{}, nil
}
