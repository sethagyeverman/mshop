package biz

import (
	"context"
	"time"

	"mshop/pkg/errx"
	"mshop/pkg/utils"
	goodsV1 "mshop/service/goods/api/goods/v1"
	inventoryV1 "mshop/service/inventory/api/inventory/v1"
	pb "mshop/service/order/api/order/v1"

	"github.com/google/uuid"
)

// CreateOrder 创建订单
func (uc *OrderUsecase) CreateOrder(ctx context.Context, req *pb.OrderRequest) (resp *pb.OrderInfoResponse, err error) {
	/*
		从购物车中获取选中的商品
		计算商品总金额
		商品库存扣减
		创建订单表项
		从购物车删除选中的商品
	*/

	// 获取购物车选中的商品
	cartItems, err := uc.CartItemList(ctx, &pb.UserInfo{Id: req.UserId})
	if err != nil {
		return nil, err
	}

	goodsIds := make([]int32, 0, len(cartItems.Data))
	goodsId2Num := make(map[int32]int32)
	goodsInfos := make([]*inventoryV1.GoodsInvInfo, 0, len(cartItems.Data))

	for _, cartItem := range cartItems.Data {
		if cartItem.Checked {
			goodsIds = append(goodsIds, cartItem.GoodsId)
			goodsId2Num[cartItem.GoodsId] = cartItem.Nums
			goodsInfos = append(goodsInfos, &inventoryV1.GoodsInvInfo{
				GoodsId: cartItem.GoodsId,
				Num:     cartItem.Nums,
			})
		}
	}

	if len(goodsIds) == 0 {
		return nil, errx.ErrorOrderCreateFailed("no goods selected")
	}

	// 查询商品服务
	getGoodsResp, err := uc.goodsClient.BatchGetGoods(ctx, &goodsV1.BatchGoodsIdInfo{
		Id: goodsIds,
	})
	if err != nil {
		return nil, err
	}

	// 计算本次消费金额
	amount := float32(0)
	for _, good := range getGoodsResp.Data {
		amount += good.ShopPrice * float32(goodsId2Num[good.Id])
	}

	// 库存扣减
	if _, err := uc.inventoryClient.Sell(ctx, &inventoryV1.SellInfo{
		GoodsInfo: goodsInfos,
	}); err != nil {
		return nil, err
	}

	// 购物车扣减
	for _, cartItem := range cartItems.Data {
		if cartItem.Checked {
			if _, err := uc.DeleteCartItem(ctx, &pb.CartItemRequest{Id: cartItem.Id}); err != nil {
				return nil, err
			}
		}
	}

	// 保存订单信息
	tx := uc.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	orderInfo := &OrderInfo{
		UserId:       req.UserId,
		OrderSn:      uuid.New().String()[:15],
		Status:       "PAYING",
		OrderAmount:  float64(amount),
		PayTime:      time.Now(),
		Address:      req.Address,
		SignerName:   req.Name,
		SignerMobile: req.Mobile,
		Post:         req.Post,

		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	}
	if result := tx.Save(orderInfo); result.Error != nil {
		tx.Rollback()
		return nil, errx.ErrorOrderCreateFailed("create order info failed")
	}

	for _, good := range getGoodsResp.Data {
		orderGoods := &OrderGoods{
			OrderId:    orderInfo.ID,
			GoodsId:    good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsId2Num[good.Id],
			AddTime:    time.Now(),
			UpdateTime: time.Now(),
		}
		if result := tx.Save(orderGoods); result.Error != nil {
			tx.Rollback()
			return nil, errx.ErrorOrderCreateFailed("create order goods failed")
		}
	}

	return &pb.OrderInfoResponse{
		Id:      orderInfo.ID,
		UserId:  orderInfo.UserId,
		OrderSn: orderInfo.OrderSn,
		PayType: orderInfo.PayType,
		Status:  orderInfo.Status,
		Post:    orderInfo.Post,
		Total:   float32(orderInfo.OrderAmount),
	}, nil
}

// OrderList 获取订单列表
func (uc *OrderUsecase) OrderList(ctx context.Context, req *pb.OrderFilterRequest) (resp *pb.OrderListResponse, err error) {

	var orderInfos []*OrderInfo
	if result := uc.db.Scopes(utils.Paginate(req.Pages, req.PagePerNums)).Where("user_id = ?", req.UserId).Find(&orderInfos); result.Error != nil {
		return nil, result.Error
	}

	var count int64
	uc.db.Model(&OrderInfo{}).Where("user_id = ?", req.UserId).Count(&count)

	resp = &pb.OrderListResponse{
		Data:  make([]*pb.OrderInfoResponse, 0, len(orderInfos)),
		Total: int32(count),
	}
	for _, orderInfo := range orderInfos {
		resp.Data = append(resp.Data, &pb.OrderInfoResponse{
			Id:      orderInfo.ID,
			UserId:  orderInfo.UserId,
			OrderSn: orderInfo.OrderSn,
			PayType: orderInfo.PayType,
			Status:  orderInfo.Status,
			Post:    orderInfo.Post,
			Total:   float32(orderInfo.OrderAmount),
			Address: orderInfo.Address,
			Name:    orderInfo.SignerName,
			Mobile:  orderInfo.SignerMobile,
			AddTime: orderInfo.AddTime.Format(time.DateTime),
		})
	}
	return
}

// OrderDetail 获取订单详情
func (uc *OrderUsecase) OrderDetail(ctx context.Context, req *pb.OrderRequest) (resp *pb.OrderInfoDetailResponse, err error) {
	var orderInfo OrderInfo
	if result := uc.db.First(&orderInfo, req.Id); result.Error != nil {
		return nil, result.Error
	}

	resp = &pb.OrderInfoDetailResponse{
		OrderInfo: &pb.OrderInfoResponse{
			Id:      orderInfo.ID,
			UserId:  orderInfo.UserId,
			OrderSn: orderInfo.OrderSn,
			PayType: orderInfo.PayType,
			Status:  orderInfo.Status,
			Post:    orderInfo.Post,
			Total:   float32(orderInfo.OrderAmount),
			Address: orderInfo.Address,
			Name:    orderInfo.SignerName,
			Mobile:  orderInfo.SignerMobile,
			AddTime: orderInfo.AddTime.Format(time.DateTime),
		},
	}

	var goods []*OrderGoods
	if result := uc.db.Where("order_id = ?", orderInfo.ID).Find(&goods); result.Error != nil {
		return nil, errx.ErrorOrderGoodsEmpty("order has no goods")
	}

	resp.Goods = make([]*pb.OrderItemResponse, 0, len(goods))
	for _, good := range goods {
		resp.Goods = append(resp.Goods, &pb.OrderItemResponse{
			Id:         good.ID,
			OrderId:    good.OrderId,
			GoodsId:    good.GoodsId,
			GoodsName:  good.GoodsName,
			GoodsImage: good.GoodsImage,
			GoodsPrice: float32(good.GoodsPrice),
			Nums:       good.Nums,
		})
	}

	return
}

// UpdateOrderStatus 更新订单状态
func (uc *OrderUsecase) UpdateOrderStatus(ctx context.Context, req *pb.OrderStatus) (resp *pb.Empty, err error) {

	orderInfo := &OrderInfo{
		Status:     req.Status,
		UpdateTime: time.Now(),
	}

	if result := uc.db.Where("order_sn = ?", req.OrderSn).Updates(orderInfo); result.Error != nil {
		return nil, errx.ErrorOrderUpdateFailed("update order status failed: %v", result.Error)
	}

	return
}
