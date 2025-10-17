package biz

import (
	"mshop/service/goods/internal/data"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGoodsUsecase)

type GoodsUsecase struct {
	db        *gorm.DB
	log       *log.Helper
	goodsRepo *data.GoodsRepo
}

func NewGoodsUsecase(db *gorm.DB, logger log.Logger, goodsRepo *data.GoodsRepo) *GoodsUsecase {
	return &GoodsUsecase{
		db:        db,
		log:       log.NewHelper(logger),
		goodsRepo: goodsRepo,
	}
}
