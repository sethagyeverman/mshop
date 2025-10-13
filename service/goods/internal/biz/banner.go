package biz

import (
	"context"
	"mshop/pkg/errx"
	pb "mshop/service/goods/api/goods/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
)

func (uc *GoodsUsecase) BannerList(ctx context.Context) (resp *pb.BannerListResponse, err error) {
	resp = &pb.BannerListResponse{}
	var banners []*Banner
	if result := uc.db.Find(&banners); result.Error != nil {
		if !errors.As(result.Error, &gorm.ErrRecordNotFound) {
			return nil, errx.ErrorDatabaseError("db error: %v", result.Error)
		}
		resp.Data = make([]*pb.BannerResponse, 0)
		return
	}

	for _, b := range banners {
		resp.Data = append(resp.Data, &pb.BannerResponse{
			Id:    b.ID,
			Image: b.Image,
			Url:   b.URL,
			Index: b.Index,
		})
	}

	return
}

func (uc *GoodsUsecase) CreateBanner(ctx context.Context, req *pb.BannerRequest) (resp *pb.BannerResponse, err error) {
	banner := &Banner{
		Image: req.Image,
		URL:   req.Url,
		Index: req.Index,
	}
	uc.db.Save(&banner)

	return &pb.BannerResponse{
		Id:    banner.ID,
		Image: banner.Image,
		Url:   banner.URL,
		Index: banner.Index,
	}, nil
}

func (uc *GoodsUsecase) DeleteBanner(ctx context.Context, req *pb.BannerRequest) (_ *pb.Empty, err error) {
	if result := uc.db.Delete(&Banner{}, req.Id); result.RowsAffected == 0 {
		return nil, errx.ErrorBannerNotFound("banner not found")
	}
	return &pb.Empty{}, nil
}

func (uc *GoodsUsecase) UpdateBanner(ctx context.Context, req *pb.BannerRequest) (_ *pb.Empty, err error) {
	var banner *Banner
	if result := uc.db.Where("id = ?", req.Id).First(&banner); result.RowsAffected == 0 {
		return nil, errx.ErrorBannerNotFound("banner not found")
	}
	if req.Image != "" {
		banner.Image = req.Image
	}
	if req.Url != "" {
		banner.URL = req.Url
	}
	if req.Index != 0 {
		banner.Index = req.Index
	}
	uc.db.Save(&banner)

	return &pb.Empty{}, nil
}
