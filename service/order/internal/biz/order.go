package biz

import (
	"context"
	"fmt"
	"time"

	"mshop/pkg/errx"
	goodsV1 "mshop/service/goods/api/goods/v1"
	inventoryV1 "mshop/service/inventory/api/inventory/v1"
	pb "mshop/service/order/api/order/v1"
)

// generateOrderSn 生成订单号
func (uc *OrderUsecase) generateOrderSn(userId int32) string {
	// 生成订单号：时间戳 + 用户ID
	return fmt.Sprintf("%d%d%d", time.Now().Unix(), time.Now().UnixNano()%1000000, userId)
}

// CreateOrder 创建订单
func (uc *OrderUsecase) CreateOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderInfoResponse, error) {
	// 1. 从购物车中获取选中的商品
	var shopCarts []ShoppingCart
	if result := uc.db.Where("user_id = ? and checked = ?", req.UserId, true).Find(&shopCarts); result.Error != nil {
		return nil, result.Error
	}

	if len(shopCarts) == 0 {
		return nil, errx.ErrorOrderCreateFailed("no items selected in cart")
	}

	// 2. 批量查询商品信息
	var goodsIds []int32
	goodsNumsMap := make(map[int32]int32) // 商品ID -> 购买数量
	for _, cart := range shopCarts {
		goodsIds = append(goodsIds, cart.GoodsID)
		goodsNumsMap[cart.GoodsID] = cart.Nums
	}

	goodsListRsp, err := uc.goodsClient.BatchGetGoods(ctx, &goodsV1.BatchGoodsIdInfo{
		Id: goodsIds,
	})
	if err != nil {
		return nil, err
	}

	// 3. 计算订单总价
	var totalPrice float32
	var orderGoods []OrderGoods
	var sellInfo []*inventoryV1.GoodsInvInfo

	for _, goodsInfo := range goodsListRsp.Data {
		nums := goodsNumsMap[goodsInfo.Id]
		totalPrice += goodsInfo.ShopPrice * float32(nums)

		orderGoods = append(orderGoods, OrderGoods{
			GoodsID:    goodsInfo.Id,
			GoodsName:  goodsInfo.Name,
			GoodsImage: goodsInfo.GoodsFrontImage,
			GoodsPrice: goodsInfo.ShopPrice,
			Nums:       nums,
		})

		sellInfo = append(sellInfo, &inventoryV1.GoodsInvInfo{
			GoodsId: goodsInfo.Id,
			Num:     nums,
		})
	}

	// 4. 开始事务
	tx := uc.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 5. 创建订单
	orderInfo := OrderInfo{
		UserID:       req.UserId,
		OrderSn:      uc.generateOrderSn(req.UserId),
		Address:      req.Address,
		SignerName:   req.Name,
		SignerMobile: req.Mobile,
		Post:         req.Post,
		Status:       "WAIT_BUYER_PAY",
		AddTime:      time.Now(),
		UpdateTime:   time.Now(),
	}

	if err := tx.Create(&orderInfo).Error; err != nil {
		tx.Rollback()
		return nil, errx.ErrorOrderCreateFailed("failed to create order: %v", err)
	}

	// 6. 创建订单商品明细
	for i := range orderGoods {
		orderGoods[i].OrderID = orderInfo.ID
		orderGoods[i].AddTime = time.Now()
		orderGoods[i].UpdateTime = time.Now()
	}

	if err := tx.Create(&orderGoods).Error; err != nil {
		tx.Rollback()
		return nil, errx.ErrorOrderCreateFailed("failed to create order goods: %v", err)
	}

	// 7. 扣减库存
	if _, err := uc.inventoryClient.Sell(ctx, &inventoryV1.SellInfo{
		GoodsInfo: sellInfo,
	}); err != nil {
		tx.Rollback()
		return nil, errx.ErrorOrderCreateFailed("failed to sell inventory: %v", err)
	}

	// 8. 删除购物车中已选中的商品
	if err := tx.Where("user_id = ? and checked = ?", req.UserId, true).Delete(&ShoppingCart{}).Error; err != nil {
		tx.Rollback()
		// 库存已扣减，需要归还
		uc.inventoryClient.Reback(ctx, &inventoryV1.SellInfo{
			GoodsInfo: sellInfo,
		})
		return nil, errx.ErrorOrderCreateFailed("failed to delete cart items: %v", err)
	}

	// 9. 提交事务
	if err := tx.Commit().Error; err != nil {
		// 归还库存
		uc.inventoryClient.Reback(ctx, &inventoryV1.SellInfo{
			GoodsInfo: sellInfo,
		})
		return nil, errx.ErrorOrderCreateFailed("failed to commit transaction: %v", err)
	}

	return &pb.OrderInfoResponse{
		Id:      orderInfo.ID,
		UserId:  orderInfo.UserID,
		OrderSn: orderInfo.OrderSn,
		PayType: orderInfo.PayType,
		Status:  orderInfo.Status,
		Post:    orderInfo.Post,
		Total:   totalPrice,
		Address: orderInfo.Address,
		Name:    orderInfo.SignerName,
		Mobile:  orderInfo.SignerMobile,
		AddTime: orderInfo.AddTime.Format("2006-01-02 15:04:05"),
	}, nil
}

// OrderList 获取订单列表
func (uc *OrderUsecase) OrderList(ctx context.Context, req *pb.OrderFilterRequest) (*pb.OrderListResponse, error) {
	var orders []OrderInfo
	var total int64

	// 构建查询
	query := uc.db.Model(&OrderInfo{})
	if req.UserId > 0 {
		query = query.Where("user_id = ?", req.UserId)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (req.Pages - 1) * req.PagePerNums
	if err := query.Offset(int(offset)).Limit(int(req.PagePerNums)).Find(&orders).Error; err != nil {
		return nil, err
	}

	rsp := &pb.OrderListResponse{
		Total: int32(total),
	}

	for _, order := range orders {
		rsp.Data = append(rsp.Data, &pb.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.UserID,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SignerMobile,
			AddTime: order.AddTime.Format("2006-01-02 15:04:05"),
		})
	}

	return rsp, nil
}

// OrderDetail 获取订单详情
func (uc *OrderUsecase) OrderDetail(ctx context.Context, req *pb.OrderRequest) (*pb.OrderInfoDetailResponse, error) {
	var order OrderInfo
	if result := uc.db.Where("id = ?", req.Id).First(&order); result.RowsAffected == 0 {
		return nil, errx.ErrorOrderNotFound("order not found")
	}

	// 如果指定了用户ID，验证订单是否属于该用户
	if req.UserId > 0 && order.UserID != req.UserId {
		return nil, errx.ErrorOrderNotFound("order not found")
	}

	// 查询订单商品
	var orderGoods []OrderGoods
	if err := uc.db.Where("order_id = ?", order.ID).Find(&orderGoods).Error; err != nil {
		return nil, err
	}

	rsp := &pb.OrderInfoDetailResponse{
		OrderInfo: &pb.OrderInfoResponse{
			Id:      order.ID,
			UserId:  order.UserID,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SignerMobile,
			AddTime: order.AddTime.Format("2006-01-02 15:04:05"),
		},
	}

	for _, goods := range orderGoods {
		rsp.Goods = append(rsp.Goods, &pb.OrderItemResponse{
			Id:         goods.ID,
			OrderId:    goods.OrderID,
			GoodsId:    goods.GoodsID,
			GoodsName:  goods.GoodsName,
			GoodsImage: goods.GoodsImage,
			GoodsPrice: goods.GoodsPrice,
			Nums:       goods.Nums,
		})
	}

	return rsp, nil
}

// UpdateOrderStatus 更新订单状态
func (uc *OrderUsecase) UpdateOrderStatus(ctx context.Context, req *pb.OrderStatus) (*pb.Empty, error) {
	var order OrderInfo
	if result := uc.db.Where("id = ? and order_sn = ?", req.Id, req.OrderSn).First(&order); result.RowsAffected == 0 {
		return nil, errx.ErrorOrderNotFound("order not found")
	}

	order.Status = req.Status
	order.UpdateTime = time.Now()

	if err := uc.db.Save(&order).Error; err != nil {
		return nil, errx.ErrorOrderUpdateFailed("failed to update order status: %v", err)
	}

	return &pb.Empty{}, nil
}
