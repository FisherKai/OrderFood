package models

import (
	"time"
)

// Room 包间表
type Room struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RoomNumber  string    `gorm:"size:50;not null;unique" json:"room_number"`  // 包间编号
	RoomName    string    `gorm:"size:100;not null" json:"room_name"`          // 包间名称
	Capacity    int       `gorm:"not null" json:"capacity"`                    // 容纳人数
	Description string    `gorm:"type:text" json:"description"`                // 设施描述
	Status      int       `gorm:"default:1" json:"status"`                     // 状态：1-可用，2-维修，3-清洁中
	HourlyPrice float64   `gorm:"type:decimal(10,2);not null" json:"hourly_price"` // 小时价格
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	Images       []RoomImage       `gorm:"foreignKey:RoomID" json:"images"`
	Reservations []RoomReservation `gorm:"foreignKey:RoomID" json:"reservations"`
}

// RoomImage 包间图片表
type RoomImage struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	RoomID   uint   `gorm:"not null" json:"room_id"`
	ImageURL string `gorm:"size:500;not null" json:"image_url"`
	IsMain   bool   `gorm:"default:false" json:"is_main"` // 是否主图
	Sort     int    `gorm:"default:0" json:"sort"`        // 排序
}

// RoomReservation 包间预订表
type RoomReservation struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	RoomID       uint      `gorm:"not null" json:"room_id"`
	ReserveDate  time.Time `gorm:"not null" json:"reserve_date"`    // 预订日期
	StartTime    time.Time `gorm:"not null" json:"start_time"`      // 开始时间
	EndTime      time.Time `gorm:"not null" json:"end_time"`        // 结束时间
	GuestCount   int       `gorm:"not null" json:"guest_count"`     // 用餐人数
	Status       int       `gorm:"default:1" json:"status"`         // 状态：1-待确认，2-已确认，3-已完成，4-已取消
	TotalPrice   float64   `gorm:"type:decimal(10,2)" json:"total_price"` // 总价格
	Remark       string    `gorm:"type:text" json:"remark"`         // 备注
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user"`
	Room Room `gorm:"foreignKey:RoomID" json:"room"`
}

// 包间状态常量
const (
	RoomStatusAvailable = 1 // 可用
	RoomStatusMaintain  = 2 // 维修
	RoomStatusCleaning  = 3 // 清洁中
)

// 预订状态常量
const (
	ReservationStatusPending   = 1 // 待确认
	ReservationStatusConfirmed = 2 // 已确认
	ReservationStatusCompleted = 3 // 已完成
	ReservationStatusCancelled = 4 // 已取消
)

// GetStatusText 获取包间状态文本
func (r *Room) GetStatusText() string {
	switch r.Status {
	case RoomStatusAvailable:
		return "可用"
	case RoomStatusMaintain:
		return "维修"
	case RoomStatusCleaning:
		return "清洁中"
	default:
		return "未知"
	}
}

// GetReservationStatusText 获取预订状态文本
func (rr *RoomReservation) GetReservationStatusText() string {
	switch rr.Status {
	case ReservationStatusPending:
		return "待确认"
	case ReservationStatusConfirmed:
		return "已确认"
	case ReservationStatusCompleted:
		return "已完成"
	case ReservationStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

// CalculateDuration 计算预订时长（小时）
func (rr *RoomReservation) CalculateDuration() float64 {
	duration := rr.EndTime.Sub(rr.StartTime)
	return duration.Hours()
}