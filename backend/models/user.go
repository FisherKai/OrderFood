package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	Username       string         `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password       string         `gorm:"size:255;not null" json:"-"`
	Phone          string         `gorm:"size:20;uniqueIndex" json:"phone"`
	Email          string         `gorm:"size:100;uniqueIndex" json:"email"`
	Nickname       string         `gorm:"size:50" json:"nickname"`
	Avatar         string         `gorm:"size:255" json:"avatar"`
	Status         int            `gorm:"default:1;comment:1-正常 0-禁用" json:"status"`
	IsDutyStaff    bool           `gorm:"default:false;comment:是否为值班人员" json:"is_duty_staff"`
	SubsidyBalance float64        `gorm:"type:decimal(10,2);default:0;comment:餐补余额" json:"subsidy_balance"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type Admin struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Role      string         `gorm:"size:20;default:admin" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
