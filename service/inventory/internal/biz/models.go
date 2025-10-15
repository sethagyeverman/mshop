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
	Version    int32          `gorm:"column:version;not null" json:"version"`
	AddTime    time.Time      `gorm:"column:add_time;not null" json:"add_time"`
	IsDeleted  bool           `gorm:"column:is_deleted" json:"is_deleted"`
	UpdateTime time.Time      `gorm:"column:update_time;not null" json:"update_time"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (Inventory) TableName() string {
	return "inventory"
}

// InventoryHistory 库存历史模型
type InventoryHistory struct {
	ID         int32          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId     int32          `gorm:"column:user_id;not null" json:"user_id"`
	GoodsId    int32          `gorm:"column:goods_id;not null" json:"goods_id"`
	Num        int32          `gorm:"column:num;not null" json:"num"`
	OrderSn    string         `gorm:"column:order_sn;not null" json:"order_sn"`
	Status     int32          `gorm:"column:status;not null" json:"status"` // 订单的状态 1.表示库存是与扣减， 幂等性 2. 表示已经支付
	AddTime    time.Time      `gorm:"column:add_time;not null" json:"add_time"`
	IsDeleted  bool           `gorm:"column:is_deleted" json:"is_deleted"`
	UpdateTime time.Time      `gorm:"column:update_time;not null" json:"update_time"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (InventoryHistory) TableName() string {
	return "inventory_history"
}
