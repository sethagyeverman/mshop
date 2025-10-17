package data

import (
	"mshop/service/goods/internal/conf"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewDB, NewElasticsearch, NewGoodsRepo)

// Data .
type Data struct {
	db *gorm.DB
	es *elasticsearch.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB, es *elasticsearch.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db: db,
		es: es,
	}, cleanup, nil
}
