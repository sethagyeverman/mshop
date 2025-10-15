package biz

import (
	"context"
	"mshop/pkg/errx"
	pb "mshop/service/goods/api/goods/v1"
)

func (s *GoodsUsecase) GoodsList(ctx context.Context, req *pb.GoodsFilterRequest) (resp *pb.GoodsListResponse, err error) {
	resp = &pb.GoodsListResponse{
		Data: make([]*pb.GoodsInfoResponse, 0),
	}

	// 构建查询
	query := s.db.Model(&Goods{}).Preload("Category").Preload("Brand")

	// 价格过滤
	if req.PriceMin > 0 {
		query = query.Where("shop_price >= ?", req.PriceMin)
	}
	if req.PriceMax > 0 {
		query = query.Where("shop_price <= ?", req.PriceMax)
	}

	// 热销/新品过滤
	if req.IsHot {
		query = query.Where("is_hot = ?", true)
	}
	if req.IsNew {
		query = query.Where("is_new = ?", true)
	}

	// 品牌过滤
	if req.Brand > 0 {
		query = query.Where("brand_id = ?", req.Brand)
	}

	// 分类过滤
	if req.TopCategory > 0 {
		query = query.Where("category_id = ?", req.TopCategory)
	}

	// 关键词搜索
	if req.KeyWords != "" {
		query = query.Where("name LIKE ?", "%"+req.KeyWords+"%")
	}

	// 获取总数
	var count int64
	query.Count(&count)
	resp.Total = int32(count)

	// 分页查询
	var goods []Goods
	if result := query.Scopes(s.Paginate(req.Pages, req.PagePerNums)).Find(&goods); result.Error != nil {
		return nil, result.Error
	}

	// 构建响应数据
	for _, good := range goods {
		goodsInfo := &pb.GoodsInfoResponse{
			Id:              good.ID,
			Name:            good.Name,
			ShopPrice:       good.ShopPrice,
			MarketPrice:     good.MarketPrice,
			GoodsBrief:      good.GoodsBrief,
			GoodsDesc:       good.GoodsSn,
			GoodsSn:         good.GoodsSn,
			Images:          good.Images,
			DescImages:      good.DescImages,
			GoodsFrontImage: good.GoodsFrontImage,
			IsNew:           good.IsNew,
			IsHot:           good.IsHot,
			OnSale:          good.OnSale,
			AddTime:         good.AddTime.Unix(),
			ShipFree:        good.ShipFree,
			ClickNum:        good.ClickNum,
			SoldNum:         good.SoldNum,
			FavNum:          good.FavNum,
		}

		if good.Category != nil {
			goodsInfo.Category = &pb.CategoryBriefInfoResponse{
				Id:   good.Category.ID,
				Name: good.Category.Name,
			}
		}

		if good.Brand != nil {
			goodsInfo.Brand = &pb.BrandInfoResponse{
				Id:   good.Brand.ID,
				Name: good.Brand.Name,
				Logo: good.Brand.Logo,
			}
		}

		resp.Data = append(resp.Data, goodsInfo)
	}

	return
}

func (s *GoodsUsecase) BatchGetGoods(ctx context.Context, req *pb.BatchGoodsIdInfo) (resp *pb.GoodsListResponse, err error) {

	goods := make([]Goods, 0)
	if result := s.db.Preload("Category").Preload("Brand").Find(&goods, req.Id); result.Error != nil {
		return nil, result.Error
	}

	resp = &pb.GoodsListResponse{
		Data: make([]*pb.GoodsInfoResponse, 0),
	}

	for _, good := range goods {
		resp.Data = append(resp.Data, &pb.GoodsInfoResponse{
			Id:        good.ID,
			Name:      good.Name,
			ShopPrice: good.ShopPrice,
			Category: &pb.CategoryBriefInfoResponse{
				Id:   good.Category.ID,
				Name: good.Category.Name,
			},
			Brand: &pb.BrandInfoResponse{
				Id:   good.Brand.ID,
				Name: good.Brand.Name,
			},
			Images:          good.Images,
			DescImages:      good.DescImages,
			GoodsFrontImage: good.GoodsFrontImage,
			IsNew:           good.IsNew,
			IsHot:           good.IsHot,
			OnSale:          good.OnSale,
			AddTime:         good.AddTime.Unix(),
			GoodsBrief:      good.GoodsBrief,
			GoodsDesc:       good.GoodsSn,
			ShipFree:        good.ShipFree,
			ClickNum:        good.ClickNum,
			SoldNum:         good.SoldNum,
			FavNum:          good.FavNum,
			MarketPrice:     good.MarketPrice,
			GoodsSn:         good.GoodsSn,
		})
	}

	return
}
func (s *GoodsUsecase) CreateGoods(ctx context.Context, req *pb.CreateGoodsInfo) (resp *pb.GoodsInfoResponse, err error) {
	// 检查分类是否存在
	var category Category
	if result := s.db.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, errx.ErrorCategoryNotFound("category not found")
	}

	// 检查品牌是否存在
	var brand Brands
	if result := s.db.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, errx.ErrorBrandNotFound("brand not found")
	}

	// 创建商品
	goods := &Goods{
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		Stocks:          req.Stocks,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
		CategoryID:      req.CategoryId,
		BrandID:         req.BrandId,
	}

	if result := s.db.Create(goods); result.Error != nil {
		return nil, result.Error
	}

	// 预加载关联数据
	s.db.Preload("Category").Preload("Brand").First(goods, goods.ID)

	// 构建响应
	resp = &pb.GoodsInfoResponse{
		Id:              goods.ID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ShopPrice:       goods.ShopPrice,
		MarketPrice:     goods.MarketPrice,
		GoodsBrief:      goods.GoodsBrief,
		GoodsDesc:       goods.GoodsSn,
		ShipFree:        goods.ShipFree,
		Images:          goods.Images,
		DescImages:      goods.DescImages,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		AddTime:         goods.AddTime.Unix(),
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
	}

	if goods.Category != nil {
		resp.Category = &pb.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		}
	}

	if goods.Brand != nil {
		resp.Brand = &pb.BrandInfoResponse{
			Id:   goods.Brand.ID,
			Name: goods.Brand.Name,
			Logo: goods.Brand.Logo,
		}
	}

	return
}
func (s *GoodsUsecase) DeleteGoods(ctx context.Context, req *pb.DeleteGoodsInfo) (resp *pb.Empty, err error) {
	if result := s.db.Delete(&Goods{}, req.Id); result.RowsAffected == 0 {
		return nil, errx.ErrorGoodsNotFound("goods not found")
	}
	return &pb.Empty{}, nil
}
func (s *GoodsUsecase) UpdateGoods(ctx context.Context, req *pb.CreateGoodsInfo) (resp *pb.Empty, err error) {
	// 查找商品
	var goods Goods
	if result := s.db.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, errx.ErrorGoodsNotFound("goods not found")
	}

	// 如果更新分类，检查分类是否存在
	if req.CategoryId > 0 {
		var category Category
		if result := s.db.First(&category, req.CategoryId); result.RowsAffected == 0 {
			return nil, errx.ErrorCategoryNotFound("category not found")
		}
		goods.CategoryID = req.CategoryId
	}

	// 如果更新品牌，检查品牌是否存在
	if req.BrandId > 0 {
		var brand Brands
		if result := s.db.First(&brand, req.BrandId); result.RowsAffected == 0 {
			return nil, errx.ErrorBrandNotFound("brand not found")
		}
		goods.BrandID = req.BrandId
	}

	// 更新字段
	if req.Name != "" {
		goods.Name = req.Name
	}
	if req.GoodsSn != "" {
		goods.GoodsSn = req.GoodsSn
	}
	if req.Stocks > 0 {
		goods.Stocks = req.Stocks
	}
	if req.MarketPrice > 0 {
		goods.MarketPrice = req.MarketPrice
	}
	if req.ShopPrice > 0 {
		goods.ShopPrice = req.ShopPrice
	}
	if req.GoodsBrief != "" {
		goods.GoodsBrief = req.GoodsBrief
	}
	if req.GoodsFrontImage != "" {
		goods.GoodsFrontImage = req.GoodsFrontImage
	}
	if len(req.Images) > 0 {
		goods.Images = req.Images
	}
	if len(req.DescImages) > 0 {
		goods.DescImages = req.DescImages
	}

	// 布尔值字段直接赋值
	goods.ShipFree = req.ShipFree
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	// 保存更新
	if result := s.db.Save(&goods); result.Error != nil {
		return nil, result.Error
	}

	return &pb.Empty{}, nil
}
func (s *GoodsUsecase) GetGoodsDetail(ctx context.Context, req *pb.GoodInfoRequest) (resp *pb.GoodsInfoResponse, err error) {
	var goods Goods
	if result := s.db.Preload("Category").Preload("Brand").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, errx.ErrorGoodsNotFound("goods not found")
	}

	// 构建响应
	resp = &pb.GoodsInfoResponse{
		Id:              goods.ID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ShopPrice:       goods.ShopPrice,
		MarketPrice:     goods.MarketPrice,
		GoodsBrief:      goods.GoodsBrief,
		GoodsDesc:       goods.GoodsSn,
		ShipFree:        goods.ShipFree,
		Images:          goods.Images,
		DescImages:      goods.DescImages,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		AddTime:         goods.AddTime.Unix(),
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
	}

	if goods.Category != nil {
		resp.Category = &pb.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		}
	}

	if goods.Brand != nil {
		resp.Brand = &pb.BrandInfoResponse{
			Id:   goods.Brand.ID,
			Name: goods.Brand.Name,
			Logo: goods.Brand.Logo,
		}
	}

	return
}
