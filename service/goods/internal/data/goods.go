package data

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	pb "mshop/service/goods/api/goods/v1"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kratos/kratos/v2/log"
)

// GoodsRepo 商品数据仓库
type GoodsRepo struct {
	esClient *elasticsearch.Client
	log      *log.Helper
}

// NewGoodsRepo 创建商品数据仓库
func NewGoodsRepo(data *Data, logger log.Logger) *GoodsRepo {
	return &GoodsRepo{
		esClient: data.es,
		log:      log.NewHelper(log.With(logger, "module", "data/goods")),
	}
}

// SearchGoodsIDs 在 ES 中搜索商品，返回商品 ID 列表和总数
func (r *GoodsRepo) SearchGoodsIDs(ctx context.Context, req *pb.GoodsFilterRequest) (ids []int32, total int64, err error) {
	if r.esClient == nil {
		r.log.Warn("elasticsearch client is nil, skip search")
		return nil, 0, fmt.Errorf("elasticsearch client is not initialized")
	}

	// 构建查询条件
	must := make([]map[string]interface{}, 0)
	filter := make([]map[string]interface{}, 0)

	// 关键词搜索 - 使用 multi_match 搜索多个字段
	if req.KeyWords != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.KeyWords,
				"fields": []string{"name", "goods_brief"},
			},
		})
	}

	// 价格过滤
	if req.PriceMin > 0 || req.PriceMax > 0 {
		priceRange := make(map[string]interface{})
		if req.PriceMin > 0 {
			priceRange["gte"] = req.PriceMin
		}
		if req.PriceMax > 0 {
			priceRange["lte"] = req.PriceMax
		}
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"shop_price": priceRange,
			},
		})
	}

	// 布尔值过滤
	if req.IsHot {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"is_hot": true,
			},
		})
	}
	if req.IsNew {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"is_new": true,
			},
		})
	}

	// 分类过滤
	if req.TopCategory > 0 {
		filter = append(filter, map[string]interface{}{
			"term": map[string]interface{}{
				"category_id": req.TopCategory,
			},
		})
	}

	// 构建完整查询
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   must,
				"filter": filter,
			},
		},
		// 只返回 ID 字段
		"_source": []string{"id"},
	}

	// 添加分页
	if req.Pages > 0 && req.PagePerNums > 0 {
		query["from"] = (req.Pages - 1) * req.PagePerNums
		query["size"] = req.PagePerNums
	}

	// 序列化查询
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		r.log.Errorf("failed to encode query: %v", err)
		return nil, 0, err
	}

	// 执行搜索
	res, err := r.esClient.Search(
		r.esClient.Search.WithContext(ctx),
		r.esClient.Search.WithIndex(GoodsIndexName),
		r.esClient.Search.WithBody(&buf),
		r.esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		r.log.Errorf("failed to search goods: %v", err)
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		r.log.Errorf("elasticsearch search error: %s", res.String())
		return nil, 0, fmt.Errorf("elasticsearch search error: %s", res.String())
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		r.log.Errorf("failed to decode response: %v", err)
		return nil, 0, err
	}

	// 提取总数
	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("invalid response format: hits not found")
	}

	totalObj, ok := hits["total"].(map[string]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("invalid response format: total not found")
	}

	totalValue, ok := totalObj["value"].(float64)
	if !ok {
		return nil, 0, fmt.Errorf("invalid response format: total value not found")
	}
	total = int64(totalValue)

	// 提取商品 ID
	hitsArray, ok := hits["hits"].([]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("invalid response format: hits array not found")
	}

	ids = make([]int32, 0, len(hitsArray))
	for _, hit := range hitsArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		id, ok := source["id"].(float64)
		if ok {
			ids = append(ids, int32(id))
		}
	}

	r.log.Infof("search goods found %d results, total: %d", len(ids), total)
	return ids, total, nil
}
