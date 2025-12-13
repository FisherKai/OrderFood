package models

import (
	"time"
)

// CanteenReview 食堂评价表
type CanteenReview struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null" json:"user_id"`
	EnvironmentScore int      `gorm:"not null;check:environment_score >= 1 AND environment_score <= 5" json:"environment_score"` // 环境卫生评分 1-5
	ServiceScore     int      `gorm:"not null;check:service_score >= 1 AND service_score <= 5" json:"service_score"`           // 服务态度评分 1-5
	QualityScore     int      `gorm:"not null;check:quality_score >= 1 AND quality_score <= 5" json:"quality_score"`           // 菜品质量评分 1-5
	PriceScore       int      `gorm:"not null;check:price_score >= 1 AND price_score <= 5" json:"price_score"`                 // 价格合理性评分 1-5
	OverallScore     int      `gorm:"not null;check:overall_score >= 1 AND overall_score <= 5" json:"overall_score"`           // 整体满意度评分 1-5
	Content          string   `gorm:"type:text" json:"content"`                                                                  // 评价内容
	IsAnonymous      bool     `gorm:"default:false" json:"is_anonymous"`                                                         // 是否匿名
	AdminReply       string   `gorm:"type:text" json:"admin_reply"`                                                              // 官方回复
	AdminReplyTime   *time.Time `json:"admin_reply_time"`                                                                       // 回复时间
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// 关联
	User   User                    `gorm:"foreignKey:UserID" json:"user"`
	Images []CanteenReviewImage    `gorm:"foreignKey:ReviewID" json:"images"`
}

// CanteenReviewImage 食堂评价图片表
type CanteenReviewImage struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	ReviewID uint   `gorm:"not null" json:"review_id"`
	ImageURL string `gorm:"size:500;not null" json:"image_url"`
}

// GetAverageScore 计算平均分
func (cr *CanteenReview) GetAverageScore() float64 {
	total := cr.EnvironmentScore + cr.ServiceScore + cr.QualityScore + cr.PriceScore + cr.OverallScore
	return float64(total) / 5.0
}

// CanteenReviewStats 食堂评价统计
type CanteenReviewStats struct {
	TotalReviews        int64   `json:"total_reviews"`
	AverageEnvironment  float64 `json:"average_environment"`
	AverageService      float64 `json:"average_service"`
	AverageQuality      float64 `json:"average_quality"`
	AveragePrice        float64 `json:"average_price"`
	AverageOverall      float64 `json:"average_overall"`
	OverallAverage      float64 `json:"overall_average"`
	ScoreDistribution   map[int]int64 `json:"score_distribution"` // 评分分布
}