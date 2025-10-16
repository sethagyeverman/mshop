package biz

import (
	"context"
	"time"

	"mshop/pkg/errx"
	goodsV1 "mshop/service/goods/api/goods/v1"
	pb "mshop/service/order/api/order/v1"
)

// CartItemList 获取用户的购物车列表
func (uc *OrderUsecase) CartItemList(ctx context.Context, req *pb.UserInfo) (resp *pb.CartItemListResponse, err error) {
	resp = &pb.CartItemListResponse{
		Data: make([]*pb.ShopCartInfoResponse, 0),
	}

	var shopCarts []*ShoppingCart
	if result := uc.db.WithContext(ctx).Where("user_id = ?", req.Id).Find(&shopCarts); result.Error != nil {
		return nil, result.Error
	}

	resp.Total = int32(len(resp.Data))

	for _, shopCart := range shopCarts {
		resp.Data = append(resp.Data, &pb.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.UserId,
			GoodsId: shopCart.GoodsId,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}

	return resp, nil

}

// CreateCartItem 添加商品到购物车
func (uc *OrderUsecase) CreateCartItem(ctx context.Context, req *pb.CartItemRequest) (resp *pb.ShopCartInfoResponse, err error) {
	var shoppingCarts ShoppingCart

	if _, err = uc.goodsClient.GetGoodsDetail(ctx, &goodsV1.GoodInfoRequest{Id: req.GoodsId}); err != nil {
		return nil, err
	}

	if result := uc.db.WithContext(ctx).Where("user_id = ? AND goods_id = ?", req.UserId, req.GoodsId).First(&shoppingCarts); result.RowsAffected == 0 {
		shoppingCarts.GoodsId = req.GoodsId
		shoppingCarts.UserId = req.UserId

		shoppingCarts.AddTime = time.Now()
	}

	shoppingCarts.Nums += req.Nums
	shoppingCarts.UpdateTime = time.Now()
	uc.db.WithContext(ctx).Save(&shoppingCarts)

	resp = &pb.ShopCartInfoResponse{
		Id:      shoppingCarts.ID,
		UserId:  shoppingCarts.UserId,
		GoodsId: shoppingCarts.GoodsId,
		Nums:    shoppingCarts.Nums,
		Checked: shoppingCarts.Checked,
	}
	return resp, nil
}

// UpdateCartItem 更新购物车
func (uc *OrderUsecase) UpdateCartItem(ctx context.Context, req *pb.CartItemRequest) (resp *pb.Empty, err error) {

	var shopCart ShoppingCart

	if result := uc.db.WithContext(ctx).Where("id = ?", req.Id).First(&shopCart); result.RowsAffected == 0 {
		return nil, errx.ErrorCartNotFound("cart item not found")
	}

	shopCart.Checked = req.Checked
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	shopCart.UpdateTime = time.Now()
	uc.db.WithContext(ctx).Save(&shopCart)

	return &pb.Empty{}, nil
}

// DeleteCartItem 删除购物车条目
func (uc *OrderUsecase) DeleteCartItem(ctx context.Context, req *pb.CartItemRequest) (resp *pb.Empty, err error) {
	if result := uc.db.WithContext(ctx).Where("id = ?", req.Id).Delete(&ShoppingCart{}); result.RowsAffected == 0 {
		return nil, errx.ErrorCartNotFound("cart item not found")
	}

	return &pb.Empty{}, nil
}
