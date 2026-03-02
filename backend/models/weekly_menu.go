package models

import (
	"time"
	"gorm.io/gorm"
)

// WeeklyMenu 一周菜谱
type WeeklyMenu struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	WeekStart   time.Time      `gorm:"not null;index" json:"week_start"`    // 周开始日期
	WeekEnd     time.Time      `gorm:"not null;index" json:"week_end"`      // 周结束日期
	Title       string         `gorm:"size:100;not null" json:"title"`      // 菜谱标题
	Status      int            `gorm:"default:0" json:"status"`             // 状态：0-草稿，1-已发布
	IsCycle     bool           `gorm:"default:false" json:"is_cycle"`       // 是否为循环菜谱
	CreatedBy   uint           `gorm:"not null" json:"created_by"`          // 创建人ID
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	MenuItems []MenuItem `gorm:"foreignKey:MenuID" json:"menu_items,omitempty"`
	Creator   Admin      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// MenuItem 菜谱详情
type MenuItem struct {
	ID       uint           `gorm:"primarykey" json:"id"`
	MenuID   uint           `gorm:"not null;index" json:"menu_id"`       // 菜谱ID
	Date     time.Time      `gorm:"not null;index" json:"date"`          // 日期
	MealType int            `gorm:"not null" json:"meal_type"`           // 餐次：1-早餐，2-午餐，3-晚餐，4-值班餐
	DishID   uint           `gorm:"not null" json:"dish_id"`             // 菜品ID
	Sort     int            `gorm:"default:0" json:"sort"`               // 排序
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Menu WeeklyMenu `gorm:"foreignKey:MenuID" json:"menu,omitempty"`
	Dish Dish       `gorm:"foreignKey:DishID" json:"dish,omitempty"`
}

// MealType 餐次类型常量
const (
	MealTypeBreakfast = 1 // 早餐
	MealTypeLunch     = 2 // 午餐
	MealTypeDinner    = 3 // 晚餐
	MealTypeDuty      = 4 // 值班餐
)

// MenuStatus 菜谱状态常量
const (
	MenuStatusDraft     = 0 // 草稿
	MenuStatusPublished = 1 // 已发布
)

// GetMealTypeName 获取餐次名称
func GetMealTypeName(mealType int) string {
	switch mealType {
	case MealTypeBreakfast:
		return "早餐"
	case MealTypeLunch:
		return "午餐"
	case MealTypeDinner:
		return "晚餐"
	case MealTypeDuty:
		return "值班餐"
	default:
		return "未知"
	}
}

// GetStatusName 获取状态名称
func GetStatusName(status int) string {
	switch status {
	case MenuStatusDraft:
		return "草稿"
	case MenuStatusPublished:
		return "已发布"
	default:
		return "未知"
	}
}