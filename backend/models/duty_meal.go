package models

import (
	"time"
)

// DutyMealSetting 值班餐设置表
type DutyMealSetting struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:100;not null" json:"name"`        // 设置名称
	StartTime    string    `gorm:"size:10;not null" json:"start_time"`   // 开始时间 HH:MM
	EndTime      string    `gorm:"size:10;not null" json:"end_time"`     // 结束时间 HH:MM
	Subsidy      float64   `gorm:"type:decimal(10,2);not null" json:"subsidy"` // 餐补标准
	Status       int       `gorm:"default:1" json:"status"`              // 状态：1-启用，0-禁用
	Description  string    `gorm:"type:text" json:"description"`         // 描述
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// DutyMealOrder 值班餐订单表
type DutyMealOrder struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	SettingID    uint      `gorm:"not null" json:"setting_id"`          // 值班餐设置ID
	OrderDate    time.Time `gorm:"not null" json:"order_date"`          // 订餐日期
	DeliveryTime time.Time `gorm:"not null" json:"delivery_time"`       // 配送时间
	TotalAmount  float64   `gorm:"type:decimal(10,2);not null" json:"total_amount"` // 订单总额
	SubsidyUsed  float64   `gorm:"type:decimal(10,2);not null" json:"subsidy_used"` // 使用的餐补
	ActualPaid   float64   `gorm:"type:decimal(10,2);not null" json:"actual_paid"`  // 实际支付
	Status       int       `gorm:"default:1" json:"status"`             // 状态：1-待处理，2-制作中，3-配送中，4-已完成，5-已取消
	Remark       string    `gorm:"type:text" json:"remark"`             // 备注
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联
	User    User            `gorm:"foreignKey:UserID" json:"user"`
	Setting DutyMealSetting `gorm:"foreignKey:SettingID" json:"setting"`
	Items   []DutyMealOrderItem `gorm:"foreignKey:OrderID" json:"items"`
}

// DutyMealOrderItem 值班餐订单详情表
type DutyMealOrderItem struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	OrderID  uint    `gorm:"not null" json:"order_id"`
	DishID   uint    `gorm:"not null" json:"dish_id"`
	Quantity int     `gorm:"not null" json:"quantity"`
	Price    float64 `gorm:"type:decimal(10,2);not null" json:"price"`

	// 关联
	Dish Dish `gorm:"foreignKey:DishID" json:"dish"`
}

// 值班餐订单状态常量
const (
	DutyMealOrderStatusPending    = 1 // 待处理
	DutyMealOrderStatusCooking    = 2 // 制作中
	DutyMealOrderStatusDelivering = 3 // 配送中
	DutyMealOrderStatusCompleted  = 4 // 已完成
	DutyMealOrderStatusCancelled  = 5 // 已取消
)

// GetStatusText 获取值班餐订单状态文本
func (dmo *DutyMealOrder) GetStatusText() string {
	switch dmo.Status {
	case DutyMealOrderStatusPending:
		return "待处理"
	case DutyMealOrderStatusCooking:
		return "制作中"
	case DutyMealOrderStatusDelivering:
		return "配送中"
	case DutyMealOrderStatusCompleted:
		return "已完成"
	case DutyMealOrderStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

// IsInDutyTime 检查当前时间是否在值班时间内
func (dms *DutyMealSetting) IsInDutyTime() bool {
	now := time.Now()
	currentTime := now.Format("15:04")
	
	return currentTime >= dms.StartTime && currentTime <= dms.EndTime
}

// DutyMealStats 值班餐统计
type DutyMealStats struct {
	TotalOrders    int64   `json:"total_orders"`
	TotalAmount    float64 `json:"total_amount"`
	TotalSubsidy   float64 `json:"total_subsidy"`
	MonthlyOrders  int64   `json:"monthly_orders"`
	MonthlyAmount  float64 `json:"monthly_amount"`
	MonthlySubsidy float64 `json:"monthly_subsidy"`
}