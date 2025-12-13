package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	OrderType     int            `gorm:"default:1;comment:1-普通订单 2-预约订单" json:"order_type"`
	Status        int            `gorm:"default:1;comment:1-待处理 2-制作中 3-已完成 4-已取消" json:"status"`
	TotalPrice    float64        `gorm:"type:decimal(10,2);not null" json:"total_price"`
	ReserveTime   *time.Time     `json:"reserve_time,omitempty"`
	PeopleCount   int            `gorm:"default:1" json:"people_count"`
	User          User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items         []OrderItem    `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	OrderID   uint           `gorm:"not null;index" json:"order_id"`
	DishID    uint           `gorm:"not null;index" json:"dish_id"`
	Quantity  int            `gorm:"not null" json:"quantity"`
	Price     float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Dish      Dish           `gorm:"foreignKey:DishID" json:"dish,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
