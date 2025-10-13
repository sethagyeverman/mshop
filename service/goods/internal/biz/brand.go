package biz

import (
	"context"
	"mshop/pkg/errx"
	pb "mshop/service/goods/api/goods/v1"
)

func (uc *GoodsUsecase) BrandList(ctx context.Context, req *pb.BrandFilterRequest) (resp *pb.BrandListResponse, err error) {
	resp = &pb.BrandListResponse{}

	var brands []*Brands
	if result := uc.db.Scopes(uc.Paginate(req.Pages, req.PagePerNums)).Find(&brands); result.Error != nil {
		resp.Data = make([]*pb.BrandInfoResponse, 0)
		return
	}

	var count int64
	uc.db.Model(&Brands{}).Count(&count)
	resp.Total = int32(count)

	for _, b := range brands {
		resp.Data = append(resp.Data, &pb.BrandInfoResponse{
			Id:   b.ID,
			Name: b.Name,
			Logo: b.Logo,
		})
	}
	return
}

func (uc *GoodsUsecase) CreateBrand(ctx context.Context, req *pb.BrandRequest) (resp *pb.BrandInfoResponse, err error) {
	var brand *Brands
	if result := uc.db.Where("name = ?", req.Name).First(&brand); result.RowsAffected != 0 {
		return nil, errx.ErrorBrandNameExists("brand name already exists")
	}

	brand = &Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	uc.db.Save(&brand)

	return &pb.BrandInfoResponse{
		Id: brand.ID,
	}, nil
}

func (uc *GoodsUsecase) DeleteBrand(ctx context.Context, req *pb.BrandRequest) (_ *pb.Empty, err error) {
	if result := uc.db.Delete(&Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, errx.ErrorBrandNotFound("brand not found")
	}
	return &pb.Empty{}, nil
}

func (uc *GoodsUsecase) UpdateBrand(ctx context.Context, req *pb.BrandRequest) (resp *pb.Empty, err error) {
	var brand *Brands

	if result := uc.db.Where("id = ?", req.Id).First(&brand); result.RowsAffected == 0 {
		return nil, errx.ErrorBrandNotFound("brand not found")
	}

	if req.Name != "" {
		if result := uc.db.Where("name = ?", req.Name).First(&brand); result.RowsAffected != 0 {
			return nil, errx.ErrorBrandNameExists("brand name already exists")
		}
		brand.Name = req.Name
	}
	if req.Logo != "" {
		brand.Logo = req.Logo
	}
	uc.db.Save(&brand)

	return &pb.Empty{}, nil
}
