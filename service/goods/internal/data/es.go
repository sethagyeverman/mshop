package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"mshop/service/goods/internal/conf"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kratos/kratos/v2/log"
)

const (
	// GoodsIndexName ES 商品索引名称
	GoodsIndexName = "goods"
)

// GoodsMapping 商品 ES mapping 定义
var GoodsMapping = map[string]interface{}{
	"mappings": map[string]interface{}{
		"properties": map[string]interface{}{
			"id": map[string]interface{}{
				"type": "integer",
			},
			"category_id": map[string]interface{}{
				"type": "integer",
			},
			"on_sale": map[string]interface{}{
				"type": "boolean",
			},
			"ship_free": map[string]interface{}{
				"type": "boolean",
			},
			"is_new": map[string]interface{}{
				"type": "boolean",
			},
			"is_hot": map[string]interface{}{
				"type": "boolean",
			},
			"name": map[string]interface{}{
				"type":     "text",
				"analyzer": "ik_max_word",
				"fields": map[string]interface{}{
					"keyword": map[string]interface{}{
						"type":         "keyword",
						"ignore_above": 256,
					},
				},
			},
			"click_num": map[string]interface{}{
				"type": "integer",
			},
			"sold_num": map[string]interface{}{
				"type": "integer",
			},
			"fav_num": map[string]interface{}{
				"type": "integer",
			},
			"market_price": map[string]interface{}{
				"type": "float",
			},
			"shop_price": map[string]interface{}{
				"type": "float",
			},
			"goods_brief": map[string]interface{}{
				"type":     "text",
				"analyzer": "ik_max_word",
			},
		},
	},
	"settings": map[string]interface{}{
		"number_of_shards":   3,
		"number_of_replicas": 1,
	},
}

// NewElasticsearch 创建 Elasticsearch 客户端
func NewElasticsearch(c *conf.Data, logger log.Logger) (*elasticsearch.Client, error) {
	helper := log.NewHelper(log.With(logger, "module", "data/elasticsearch"))

	if c.Elasticsearch == nil {
		helper.Warn("elasticsearch config is nil, skip initialization")
		return nil, nil
	}

	cfg := elasticsearch.Config{
		Addresses: c.Elasticsearch.Addresses,
	}

	// 设置认证信息
	if c.Elasticsearch.Username != "" && c.Elasticsearch.Password != "" {
		cfg.Username = c.Elasticsearch.Username
		cfg.Password = c.Elasticsearch.Password
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		helper.Errorf("failed to create elasticsearch client: %v", err)
		return nil, err
	}
 
	// 检查连接
	res, err := client.Info()
	if err != nil {
		helper.Errorf("failed to get elasticsearch info: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		helper.Errorf("elasticsearch returned error: %s", res.String())
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	helper.Info("elasticsearch client initialized successfully")

	// 初始化索引
	if err := initGoodsIndex(client, helper); err != nil {
		helper.Warnf("failed to initialize goods index: %v", err)
		// 不返回错误，允许应用启动
	}

	return client, nil
}

// initGoodsIndex 初始化商品索引
func initGoodsIndex(client *elasticsearch.Client, helper *log.Helper) error {
	ctx := context.Background()

	// 检查索引是否存在
	res, err := client.Indices.Exists([]string{GoodsIndexName})
	if err != nil {
		return fmt.Errorf("check index exists error: %w", err)
	}
	defer res.Body.Close()

	// 索引已存在
	if res.StatusCode == 200 {
		helper.Infof("index [%s] already exists", GoodsIndexName)
		return nil
	}

	// 创建索引
	mappingJSON, err := json.Marshal(GoodsMapping)
	if err != nil {
		return fmt.Errorf("marshal mapping error: %w", err)
	}

	res, err = client.Indices.Create(
		GoodsIndexName,
		client.Indices.Create.WithBody(strings.NewReader(string(mappingJSON))),
		client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("create index error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("create index failed: %s", res.String())
	}

	helper.Infof("index [%s] created successfully", GoodsIndexName)
	return nil
}
