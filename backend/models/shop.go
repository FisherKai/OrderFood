package models

import (
	"time"
)

// ShopCategory 商城商品分类表
type ShopCategory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	Icon      string    `gorm:"size:255" json:"icon"`
	Sort      int       `gorm:"default:0" json:"sort"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Products []ShopProduct `gorm:"foreignKey:CategoryID" json:"products"`
}

// ShopProduct 商城商品表
type ShopProduct struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CategoryID  uint      `gorm:"not null;index" json:"category_id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginPrice float64   `gorm:"type:decimal(10,2)" json:"origin_price"` // 原价
	Description string    `gorm:"type:text" json:"description"`
	Stock       int       `gorm:"not null;default:0" json:"stock"`
	SalesCount  int       `gorm:"default:0" json:"sales_count"`           // 销量
	Status      int       `gorm:"default:1" json:"status"`                // 状态：1-上架，0-下架
	IsPromotion bool      `gorm:"default:false" json:"is_promotion"`      // 是否促销
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	Category ShopCategory      `gorm:"foreignKey:CategoryID" json:"category"`
	Images   []ShopProductImage `gorm:"foreignKey:ProductID" json:"images"`
}

// ShopProductImage 商城商品图片表
type ShopProductImage struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ProductID uint   `gorm:"not null" json:"product_id"`
	ImageURL  string `gorm:"size:500;not null" json:"image_url"`
	IsMain    bool   `gorm:"default:false" json:"is_main"` // 是否主图
	Sort      int    `gorm:"default:0" json:"sort"`        // 排序
}

// ShopOrder 商城订单表
type ShopOrder struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderNo      string    `gorm:"size:50;not null;unique" json:"order_no"`  // 订单号
	UserID       uint      `gorm:"not null" json:"user_id"`
	TotalAmount  float64   `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status       int       `gorm:"default:1" json:"status"`                  // 状态：1-待支付，2-待发货，3-待收货，4-已完成，5-已取消
	PaymentMethod string   `gorm:"size:20" json:"payment_method"`            // 支付方式
	DeliveryType int       `gorm:"default:1" json:"delivery_type"`           // 配送方式：1-自提，2-配送
	DeliveryAddress string `gorm:"type:text" json:"delivery_address"`        // 配送地址
	Remark       string    `gorm:"type:text" json:"remark"`                  // 备注
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联
	User  User             `gorm:"foreignKey:UserID" json:"user"`
	Items []ShopOrderItem  `gorm:"foreignKey:OrderID" json:"items"`
}

// ShopOrderItem 商城订单详情表
type ShopOrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"not null" json:"order_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"type:decimal(10,2);not null" json:"price"`

	// 关联
	Product ShopProduct `gorm:"foreignKey:ProductID" json:"product"`
}

// InventoryRecord 库存记录表
type InventoryRecord struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ProductID     uint      `gorm:"not null" json:"product_id"`
	Type          int       `gorm:"not null" json:"type"`                     // 类型：1-入库，2-出库，3-调整
	Quantity      int       `gorm:"not null" json:"quantity"`                 // 变动数量
	BeforeStock   int       `gorm:"not null" json:"before_stock"`             // 变动前库存
	AfterStock    int       `gorm:"not null" json:"after_stock"`              // 变动后库存
	Reason        string    `gorm:"size:200" json:"reason"`                   // 变动原因
	OperatorID    uint      `json:"operator_id"`                              // 操作人ID
	OperatorType  string    `gorm:"size:20;default:admin" json:"operator_type"` // 操作人类型：admin/system
	CreatedAt     time.Time `json:"created_at"`

	// 关联
	Product ShopProduct `gorm:"foreignKey:ProductID" json:"product"`
}

// 商城订单状态常量
const (
	ShopOrderStatusPending   = 1 // 待支付
	ShopOrderStatusPaid      = 2 // 待发货
	ShopOrderStatusShipping  = 3 // 待收货
	ShopOrderStatusCompleted = 4 // 已完成
	ShopOrderStatusCancelled = 5 // 已取消
)

// 库存变动类型常量
const (
	InventoryTypeIn     = 1 // 入库
	InventoryTypeOut    = 2 // 出库
	InventoryTypeAdjust = 3 // 调整
)

// GetStatusText 获取商城订单状态文本
func (so *ShopOrder) GetStatusText() string {
	switch so.Status {
	case ShopOrderStatusPending:
		return "待支付"
	case ShopOrderStatusPaid:
		return "待发货"
	case ShopOrderStatusShipping:
		return "待收货"
	case ShopOrderStatusCompleted:
		return "已完成"
	case ShopOrderStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

// GetTypeText 获取库存变动类型文本
func (ir *InventoryRecord) GetTypeText() string {
	switch ir.Type {
	case InventoryTypeIn:
		return "入库"
	case InventoryTypeOut:
		return "出库"
	case InventoryTypeAdjust:
		return "调整"
	default:
		return "未知"
	}
}