package biz

import (
	"time"

	"gorm.io/gorm"
)

// Inventory 库存模型
type Inventory struct {
	ID         int32          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GoodsId    int32          `gorm:"column:goods_id;not null" json:"goods_id"`
	Stock      int32          `gorm:"column:stock;not null" json:"stock"`
	AddTime    time.Time      `gorm:"column:add_time;not null" json:"add_time"`
	IsDeleted  bool           `gorm:"column:is_deleted" json:"is_deleted"`
	UpdateTime time.Time      `gorm:"column:update_time;not null" json:"update_time"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

// TableName 指定表名
func (Inventory) TableName() string {
	return "inventory"
}
