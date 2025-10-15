package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGoodsUsecase)

type GoodsUsecase struct {
	db  *gorm.DB
	log *log.Helper
}

func NewGoodsUsecase(db *gorm.DB, logger log.Logger) *GoodsUsecase {
	return &GoodsUsecase{
		db:  db,
		log: log.NewHelper(logger),
	}
}
