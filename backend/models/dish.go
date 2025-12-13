package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Icon      string         `gorm:"size:255" json:"icon"`
	Sort      int            `gorm:"default:0" json:"sort"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Dish struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CategoryID  uint           `gorm:"not null;index" json:"category_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Description string         `gorm:"type:text" json:"description"`
	Status      int            `gorm:"default:1;comment:1-上架 0-下架" json:"status"`
	Stock       int            `gorm:"default:0" json:"stock"`
	LikeCount   int            `gorm:"default:0" json:"like_count"`  // 点赞数
	IsDutyMeal  bool           `gorm:"default:false;comment:是否为值班餐菜品" json:"is_duty_meal"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Images      []DishImage    `gorm:"foreignKey:DishID" json:"images,omitempty"`
	Likes       []DishLike     `gorm:"foreignKey:DishID" json:"likes,omitempty"`  // 点赞记录
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type DishImage struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	DishID    uint           `gorm:"not null;index" json:"dish_id"`
	ImageURL  string         `gorm:"size:255;not null" json:"image_url"`
	IsMain    bool           `gorm:"default:false" json:"is_main"`
	Sort      int            `gorm:"default:0" json:"sort"`
	IsDeleted bool           `gorm:"default:false" json:"is_deleted"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt *time.Time     `json:"deleted_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at"`
}
