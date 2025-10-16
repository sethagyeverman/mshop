package biz

import "time"

// ShoppingCart 购物车表
type ShoppingCart struct {
	ID         int32     `gorm:"primarykey;type:int" json:"id"`
	UserID     int32     `gorm:"type:int;index" json:"user_id"`
	GoodsID    int32     `gorm:"type:int;index" json:"goods_id"`
	Nums       int32     `gorm:"type:int" json:"nums"`
	Checked    bool      `gorm:"type:boolean;default:false" json:"checked"`
	AddTime    time.Time `gorm:"type:datetime" json:"add_time"`
	UpdateTime time.Time `gorm:"type:datetime" json:"update_time"`
}

// OrderInfo 订单表
type OrderInfo struct {
	ID      int32  `gorm:"primarykey;type:int" json:"id"`
	UserID  int32  `gorm:"type:int;index" json:"user_id"`
	OrderSn string `gorm:"type:varchar(30);index" json:"order_sn"` // 订单号
	PayType string `gorm:"type:varchar(20);comment:'支付方式 alipay(支付宝)、wechat(微信)'" json:"pay_type"`

	// 支付状态
	Status string `gorm:"type:varchar(20);comment:'订单状态 PAYING(待支付), TRADE_SUCCESS(成功), TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)'" json:"status"`

	// 交易号
	TradeNo string `gorm:"type:varchar(100);comment:'交易号'" json:"trade_no"`

	// 订单留言
	OrderMount string `gorm:"type:varchar(200)" json:"order_mount"`

	// 支付时间
	PayTime *time.Time `gorm:"type:datetime" json:"pay_time"`

	// 收货人信息
	Address      string `gorm:"type:varchar(200)" json:"address"`
	SignerName   string `gorm:"type:varchar(20)" json:"signer_name"`
	SignerMobile string `gorm:"type:varchar(11)" json:"signer_mobile"`
	Post         string `gorm:"type:varchar(20)" json:"post"`

	AddTime    time.Time `gorm:"type:datetime" json:"add_time"`
	UpdateTime time.Time `gorm:"type:datetime" json:"update_time"`
}

// OrderGoods 订单商品表
type OrderGoods struct {
	ID int32 `gorm:"primarykey;type:int" json:"id"`

	OrderID int32 `gorm:"type:int;index" json:"order_id"` // 订单ID
	GoodsID int32 `gorm:"type:int;index" json:"goods_id"` // 商品ID

	// 冗余字段（避免商品信息变化后订单信息不准确）
	GoodsName  string  `gorm:"type:varchar(100)" json:"goods_name"`
	GoodsImage string  `gorm:"type:varchar(200)" json:"goods_image"`
	GoodsPrice float32 `gorm:"type:float" json:"goods_price"`
	Nums       int32   `gorm:"type:int" json:"nums"` // 购买数量

	AddTime    time.Time `gorm:"type:datetime" json:"add_time"`
	UpdateTime time.Time `gorm:"type:datetime" json:"update_time"`
}

func (ShoppingCart) TableName() string {
	return "shopping_cart"
}

func (OrderInfo) TableName() string {
	return "order_info"
}

func (OrderGoods) TableName() string {
	return "order_goods"
}
