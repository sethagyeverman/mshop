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
	return &pb.CategoryListResponse{}, nil
}
func (s *GoodsService) GetSubCategory(ctx context.Context, req *pb.CategoryListRequest) (*pb.SubCategoryListResponse, error) {
	return &pb.SubCategoryListResponse{}, nil
}
func (s *GoodsService) CreateCategory(ctx context.Context, req *pb.CategoryInfoRequest) (*pb.CategoryInfoResponse, error) {
	return &pb.CategoryInfoResponse{}, nil
}
func (s *GoodsService) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
func (s *GoodsService) UpdateCategory(ctx context.Context, req *pb.CategoryInfoRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
func (s *GoodsService) BrandList(ctx context.Context, req *pb.BrandFilterRequest) (*pb.BrandListResponse, error) {
	return &pb.BrandListResponse{}, nil
}
func (s *GoodsService) CreateBrand(ctx context.Context, req *pb.BrandRequest) (*pb.BrandInfoResponse, error) {
	return &pb.BrandInfoResponse{}, nil
}
func (s *GoodsService) DeleteBrand(ctx context.Context, req *pb.BrandRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
func (s *GoodsService) UpdateBrand(ctx context.Context, req *pb.BrandRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
func (s *GoodsService) BannerList(ctx context.Context, req *pb.Empty) (*pb.BannerListResponse, error) {
	return &pb.BannerListResponse{}, nil
}
func (s *GoodsService) CreateBanner(ctx context.Context, req *pb.BannerRequest) (*pb.BannerResponse, error) {
	return &pb.BannerResponse{}, nil
}
func (s *GoodsService) DeleteBanner(ctx context.Context, req *pb.BannerRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
func (s *GoodsService) UpdateBanner(ctx context.Context, req *pb.BannerRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
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
