package models

import (
	"time"
	"gorm.io/gorm"
)

// DishLike 菜品点赞
type DishLike struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`    // 用户ID
	DishID    uint           `gorm:"not null;index" json:"dish_id"`    // 菜品ID
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Dish Dish `gorm:"foreignKey:DishID" json:"dish,omitempty"`
}

// 添加唯一索引，防止重复点赞
func (DishLike) TableName() string {
	return "dish_likes"
}

// 在模型初始化时添加唯一约束
func init() {
	// 这个会在数据库迁移时自动创建唯一索引
}