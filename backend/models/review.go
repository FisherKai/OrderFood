package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	DishID    uint           `gorm:"not null;index" json:"dish_id"`
	OrderID   uint           `gorm:"not null;index" json:"order_id"`
	Rating    int            `gorm:"not null;comment:1-5星" json:"rating"`
	Content   string         `gorm:"type:text" json:"content"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Dish      Dish           `gorm:"foreignKey:DishID" json:"dish,omitempty"`
	Images    []ReviewImage  `gorm:"foreignKey:ReviewID" json:"images,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ReviewImage struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ReviewID  uint           `gorm:"not null;index" json:"review_id"`
	ImageURL  string         `gorm:"size:255;not null" json:"image_url"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
