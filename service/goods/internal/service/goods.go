package service

import (
	"context"

	pb "mshop/service/goods/api/goods/v1"
	"mshop/service/goods/internal/biz"
)

type GoodsService struct {
	pb.UnimplementedGoodsServer
	goodsUsecase *biz.GoodsUsecase
}

func NewGoodsService(goodsUsecase *biz.GoodsUsecase) *GoodsService {
	return &GoodsService{
		goodsUsecase: goodsUsecase,
	}
}

func (s *GoodsService) GoodsList(ctx context.Context, req *pb.GoodsFilterRequest) (*pb.GoodsListResponse, error) {
	return &pb.GoodsListResponse{}, nil
}
func (s *GoodsService) BatchGetGoods(ctx context.Context, req *pb.BatchGoodsIdInfo) (*pb.GoodsListResponse, error) {
	return &pb.GoodsListResponse{}, nil
}
func (s *GoodsService) CreateGoods(ctx context.Context, req *pb.CreateGoodsInfo) (*pb.GoodsInfoResponse, error) {
	return &pb.GoodsInfoResponse{}, nil
}
func (s *GoodsService) DeleteGoods(ctx context.Context, req *pb.DeleteGoodsInfo) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
func (s *GoodsService) UpdateGoods(ctx context.Context, req *pb.CreateGoodsInfo) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
func (s *GoodsService) GetGoodsDetail(ctx context.Context, req *pb.GoodInfoRequest) (*pb.GoodsInfoResponse, error) {
	return &pb.GoodsInfoResponse{}, nil
}

func (s *GoodsService) GetAllCategorysList(ctx context.Context, req *pb.Empty) (*pb.CategoryListResponse, error) {
	return s.goodsUsecase.GetAllCategorysList(ctx, req)
}
func (s *GoodsService) GetSubCategory(ctx context.Context, req *pb.CategoryListRequest) (*pb.SubCategoryListResponse, error) {
	return s.goodsUsecase.GetSubCategory(ctx, req)
}
func (s *GoodsService) CreateCategory(ctx context.Context, req *pb.CategoryInfoRequest) (*pb.CategoryInfoResponse, error) {
	return s.goodsUsecase.CreateCategory(ctx, req)
}
func (s *GoodsService) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.Empty, error) {
	return s.goodsUsecase.DeleteCategory(ctx, req)
}
func (s *GoodsService) UpdateCategory(ctx context.Context, req *pb.CategoryInfoRequest) (*pb.Empty, error) {
	return s.goodsUsecase.UpdateCategory(ctx, req)
}

func (s *GoodsService) BrandList(ctx context.Context, req *pb.BrandFilterRequest) (*pb.BrandListResponse, error) {
	return s.goodsUsecase.BrandList(ctx, req)
}
func (s *GoodsService) CreateBrand(ctx context.Context, req *pb.BrandRequest) (*pb.BrandInfoResponse, error) {
	return s.goodsUsecase.CreateBrand(ctx, req)
}
func (s *GoodsService) DeleteBrand(ctx context.Context, req *pb.BrandRequest) (*pb.Empty, error) {
	return s.goodsUsecase.DeleteBrand(ctx, req)
}
func (s *GoodsService) UpdateBrand(ctx context.Context, req *pb.BrandRequest) (*pb.Empty, error) {
	return s.goodsUsecase.UpdateBrand(ctx, req)
}

func (s *GoodsService) BannerList(ctx context.Context, req *pb.Empty) (*pb.BannerListResponse, error) {
	return s.goodsUsecase.BannerList(ctx)
}
func (s *GoodsService) CreateBanner(ctx context.Context, req *pb.BannerRequest) (*pb.BannerResponse, error) {
	return s.goodsUsecase.CreateBanner(ctx, req)
}
func (s *GoodsService) DeleteBanner(ctx context.Context, req *pb.BannerRequest) (*pb.Empty, error) {
	return s.goodsUsecase.DeleteBanner(ctx, req)
}
func (s *GoodsService) UpdateBanner(ctx context.Context, req *pb.BannerRequest) (*pb.Empty, error) {
	return s.goodsUsecase.UpdateBanner(ctx, req)
}

func (s *GoodsService) CategoryBrandList(ctx context.Context, req *pb.CategoryBrandFilterRequest) (*pb.CategoryBrandListResponse, error) {
	return &pb.CategoryBrandListResponse{}, nil
}
func (s *GoodsService) GetCategoryBrandList(ctx context.Context, req *pb.CategoryInfoRequest) (*pb.BrandListResponse, error) {
	return &pb.BrandListResponse{}, nil
}
func (s *GoodsService) CreateCategoryBrand(ctx context.Context, req *pb.CategoryBrandRequest) (*pb.CategoryBrandResponse, error) {
	return &pb.CategoryBrandResponse{}, nil
}
func (s *GoodsService) DeleteCategoryBrand(ctx context.Context, req *pb.CategoryBrandRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
func (s *GoodsService) UpdateCategoryBrand(ctx context.Context, req *pb.CategoryBrandRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
