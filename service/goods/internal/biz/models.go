package biz

import (
	"time"

	"gorm.io/gorm"
)

// Banner 轮播图模型
type Banner struct {
	ID         int32          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AddTime    time.Time      `gorm:"column:add_time;not null" json:"add_time"`
	IsDeleted  bool           `gorm:"column:is_deleted" json:"is_deleted"`
	UpdateTime time.Time      `gorm:"column:update_time;not null" json:"update_time"`
	Image      string         `gorm:"column:image;type:varchar(200);not null" json:"image"`
	URL        string         `gorm:"column:url;type:varchar(200);not null" json:"url"`
	Index      int32          `gorm:"column:index;not null" json:"index"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

// TableName 指定表名
func (Banner) TableName() string {
	return "banner"
}

// Brands 品牌模型
type Brands struct {
	ID         int32          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string         `gorm:"column:name;type:varchar(50);not null;uniqueIndex:brands_name" json:"name"`
	Logo       string         `gorm:"column:logo;type:varchar(200)" json:"logo"`
	AddTime    time.Time      `gorm:"column:add_time;not null" json:"add_time"`
	IsDeleted  bool           `gorm:"column:is_deleted" json:"is_deleted"`
	UpdateTime time.Time      `gorm:"column:update_time;not null" json:"update_time"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

// TableName 指定表名
func (Brands) TableName() string {
	return "brands"
}

// Category 商品分类模型
type Category struct {
	ID               int32          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name             string         `gorm:"column:name;type:varchar(50);not null;index:category_name" json:"name"`
	ParentCategoryID int32          `gorm:"column:parent_category_id;index:category_parent_category_id" json:"parent_category_id"`
	Level            int32          `gorm:"column:level;not null" json:"level"`
	IsTab            bool           `gorm:"column:is_tab;not null" json:"is_tab"`
	URL              string         `gorm:"column:url;type:varchar(300);not null;index:category_url" json:"url"`
	AddTime          time.Time      `gorm:"column:add_time" json:"add_time"`
	IsDeleted        bool           `gorm:"column:is_deleted" json:"is_deleted"`
	UpdateTime       time.Time      `gorm:"column:update_time" json:"update_time"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`

	// 自引用关联
	ParentCategory *Category   `gorm:"foreignKey:ParentCategoryID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"parent_category,omitempty"`
	SubCategories  []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_categories,omitempty"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "category"
}

// Goods 商品模型
type Goods struct {
	ID              int32          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
	AddTime         time.Time      `gorm:"column:add_time;not null" json:"add_time"`
	IsDeleted       bool           `gorm:"column:is_deleted" json:"is_deleted"`
	UpdateTime      time.Time      `gorm:"column:update_time;not null" json:"update_time"`
	CategoryID      int32          `gorm:"column:category_id;not null;index:goods2_category_id" json:"category_id"`
	BrandID         int32          `gorm:"column:brand_id;not null;index:goods2_brand_id" json:"brand_id"`
	OnSale          bool           `gorm:"column:on_sale;not null" json:"on_sale"`
	GoodsSn         string         `gorm:"column:goods_sn;type:varchar(50);not null" json:"goods_sn"`
	Name            string         `gorm:"column:name;type:varchar(100);not null" json:"name"`
	ClickNum        int32          `gorm:"column:click_num;not null;default:0" json:"click_num"`
	SoldNum         int32          `gorm:"column:sold_num;not null;default:0" json:"sold_num"`
	FavNum          int32          `gorm:"column:fav_num;not null;default:0" json:"fav_num"`
	Stocks          int32          `gorm:"column:stocks;not null;default:0" json:"stocks"`
	MarketPrice     float32        `gorm:"column:market_price;not null" json:"market_price"`
	ShopPrice       float32        `gorm:"column:shop_price;not null" json:"shop_price"`
	GoodsBrief      string         `gorm:"column:goods_brief;type:varchar(200);not null" json:"goods_brief"`
	ShipFree        bool           `gorm:"column:ship_free;not null" json:"ship_free"`
	Images          GormList       `gorm:"column:images;type:json;not null;serializer:json" json:"images"`
	DescImages      GormList       `gorm:"column:desc_images;type:json;not null;serializer:json" json:"desc_images"`
	GoodsFrontImage string         `gorm:"column:goods_front_image;type:varchar(200);not null" json:"goods_front_image"`
	IsNew           bool           `gorm:"column:is_new;not null" json:"is_new"`
	IsHot           bool           `gorm:"column:is_hot;not null" json:"is_hot"`

	// 外键关联
	Category *Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE" json:"category,omitempty"`
	Brand    *Brands   `gorm:"foreignKey:BrandID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE" json:"brand,omitempty"`
}

// TableName 指定表名
func (Goods) TableName() string {
	return "goods"
}

// es中商品数据模型
type EsGoods struct {
	ID         int32 `json:"id"`
	CategoryID int32 `json:"category_id"`
	OnSale     bool  `json:"on_sale"`
	ShipFree   bool  `json:"ship_free"`
	IsNew      bool  `json:"is_new"`
	IsHot      bool  `json:"is_hot"`

	Name     string `json:"name"`
	ClickNum int32  `json:"click_num"`
	SoldNum  int32  `json:"sold_num"`
	FavNum   int32  `json:"fav_num"`

	MarketPrice float32 `json:"market_price"`
	GoodsBrief  string  `json:"goods_brief"`
	ShopPrice   float32 `json:"shop_price"`
}

// GoodsCategoryBrand 商品分类品牌关联模型
type GoodsCategoryBrand struct {
	ID         int32          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CategoryID int32          `gorm:"column:category_id;not null;uniqueIndex:goodscategorybrand_category_id_brand_id,priority:1;index:goodscategorybrand_category_id" json:"category_id"`
	BrandsID   int32          `gorm:"column:brands_id;not null;uniqueIndex:goodscategorybrand_category_id_brand_id,priority:2;index:goodscategorybrand_brand_id" json:"brands_id"`
	AddTime    time.Time      `gorm:"column:add_time;not null" json:"add_time"`
	IsDeleted  bool           `gorm:"column:is_deleted" json:"is_deleted"`
	UpdateTime time.Time      `gorm:"column:update_time;not null" json:"update_time"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`

	// 外键关联
	Category *Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"category,omitempty"`
	Brand    *Brands   `gorm:"foreignKey:BrandsID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"brand,omitempty"`
}

// TableName 指定表名
func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

// GormList 自定义类型，用于处理 JSON 数组字段
type GormList []string
