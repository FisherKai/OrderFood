package models

import (
	"time"
	"gorm.io/gorm"
)

// MenuItemRating 菜谱菜品评价
type MenuItemRating struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	MenuID     uint           `gorm:"not null;index" json:"menu_id"`       // 菜谱ID
	MenuItemID uint           `gorm:"not null;index" json:"menu_item_id"`  // 菜谱项ID
	DishID     uint           `gorm:"not null;index" json:"dish_id"`       // 菜品ID
	UserID     uint           `gorm:"not null;index" json:"user_id"`       // 用户ID
	Rating     int            `gorm:"not null" json:"rating"`              // 评分 1-5
	Comment    string         `gorm:"size:500" json:"comment"`             // 评价内容（可选）
	MealType   int            `gorm:"not null" json:"meal_type"`           // 餐次
	RatingDate time.Time      `gorm:"not null;index" json:"rating_date"`   // 评价日期（对应菜谱中的日期）
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Menu     WeeklyMenu `gorm:"foreignKey:MenuID" json:"menu,omitempty"`
	MenuItem MenuItem   `gorm:"foreignKey:MenuItemID" json:"menu_item,omitempty"`
	Dish     Dish       `gorm:"foreignKey:DishID" json:"dish,omitempty"`
	User     User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (MenuItemRating) TableName() string {
	return "menu_item_ratings"
}
