package controllers

import (
	"net/http"
	"strconv"
	"time"

	"orderfood/database"
	"orderfood/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateCanteenReview 提交食堂评价（用户端）
func CreateCanteenReview(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req struct {
		EnvironmentScore int      `json:"environment_score" binding:"required,min=1,max=5"`
		ServiceScore     int      `json:"service_score" binding:"required,min=1,max=5"`
		QualityScore     int      `json:"quality_score" binding:"required,min=1,max=5"`
		PriceScore       int      `json:"price_score" binding:"required,min=1,max=5"`
		OverallScore     int      `json:"overall_score" binding:"required,min=1,max=5"`
		Content          string   `json:"content"`
		IsAnonymous      bool     `json:"is_anonymous"`
		Images           []string `json:"images"` // 图片URL数组
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否已经评价过（限制：每个用户每天只能评价一次）
	// var existingReview models.CanteenReview
	// today := time.Now().Format("2006-01-02")
	// if err := database.DB.Where("user_id = ? AND DATE(created_at) = ?", userID, today).
	// 	First(&existingReview).Error; err == nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "今天已经评价过了"})
	// 	return
	// }

	// 创建评价
	review := models.CanteenReview{
		UserID:           userID.(uint),
		EnvironmentScore: req.EnvironmentScore,
		ServiceScore:     req.ServiceScore,
		QualityScore:     req.QualityScore,
		PriceScore:       req.PriceScore,
		OverallScore:     req.OverallScore,
		Content:          req.Content,
		IsAnonymous:      req.IsAnonymous,
	}

	// 开始事务
	tx := database.DB.Begin()

	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交评价失败"})
		return
	}

	// 保存图片
	for _, imageURL := range req.Images {
		if imageURL != "" {
			reviewImage := models.CanteenReviewImage{
				ReviewID: review.ID,
				ImageURL: imageURL,
			}
			if err := tx.Create(&reviewImage).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "保存评价图片失败"})
				return
			}
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": review,
		"message": "评价提交成功",
	})
}

// GetCanteenReviews 获取食堂评价列表（用户端）
func GetCanteenReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort", "created_at") // created_at, overall_score

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.CanteenReview{})

	// 排序
	switch sortBy {
	case "overall_score":
		query = query.Order("overall_score DESC, created_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var reviews []models.CanteenReview
	if err := query.Preload("User").Preload("Images").
		Offset(offset).Limit(pageSize).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评价列表失败"})
		return
	}

	// 处理匿名显示
	for i := range reviews {
		if reviews[i].IsAnonymous {
			reviews[i].User.Username = "匿名用户"
			reviews[i].User.Avatar = ""
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      reviews,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetCanteenStats 获取食堂评分统计（用户端）
func GetCanteenStats(c *gin.Context) {
	stats, err := getCanteenReviewStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": stats,
	})
}

// getCanteenReviewStats 获取食堂评价统计
func getCanteenReviewStats() (*models.CanteenReviewStats, error) {
	var stats models.CanteenReviewStats
	
	// 获取总评价数
	if err := database.DB.Model(&models.CanteenReview{}).Count(&stats.TotalReviews).Error; err != nil {
		return nil, err
	}

	if stats.TotalReviews == 0 {
		return &stats, nil
	}

	// 计算各维度平均分
	var result struct {
		AvgEnvironment float64
		AvgService     float64
		AvgQuality     float64
		AvgPrice       float64
		AvgOverall     float64
	}

	if err := database.DB.Model(&models.CanteenReview{}).
		Select("AVG(environment_score) as avg_environment, AVG(service_score) as avg_service, AVG(quality_score) as avg_quality, AVG(price_score) as avg_price, AVG(overall_score) as avg_overall").
		Scan(&result).Error; err != nil {
		return nil, err
	}

	stats.AverageEnvironment = result.AvgEnvironment
	stats.AverageService = result.AvgService
	stats.AverageQuality = result.AvgQuality
	stats.AveragePrice = result.AvgPrice
	stats.AverageOverall = result.AvgOverall
	stats.OverallAverage = (result.AvgEnvironment + result.AvgService + result.AvgQuality + result.AvgPrice + result.AvgOverall) / 5.0

	// 获取评分分布（基于整体满意度）
	stats.ScoreDistribution = make(map[int]int64)
	for i := 1; i <= 5; i++ {
		var count int64
		database.DB.Model(&models.CanteenReview{}).Where("overall_score = ?", i).Count(&count)
		stats.ScoreDistribution[i] = count
	}

	return &stats, nil
}

// === 管理端接口 ===

// AdminGetCanteenReviews 获取食堂评价列表（管理端）
func AdminGetCanteenReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort", "created_at")
	hasReply := c.Query("has_reply") // true, false, all

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.CanteenReview{})

	// 筛选条件
	switch hasReply {
	case "true":
		query = query.Where("admin_reply != ''")
	case "false":
		query = query.Where("admin_reply = '' OR admin_reply IS NULL")
	}

	// 排序
	switch sortBy {
	case "overall_score":
		query = query.Order("overall_score DESC, created_at DESC")
	case "environment_score":
		query = query.Order("environment_score DESC, created_at DESC")
	case "service_score":
		query = query.Order("service_score DESC, created_at DESC")
	case "quality_score":
		query = query.Order("quality_score DESC, created_at DESC")
	case "price_score":
		query = query.Order("price_score DESC, created_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var reviews []models.CanteenReview
	if err := query.Preload("User").Preload("Images").
		Offset(offset).Limit(pageSize).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评价列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      reviews,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// AdminReplyCanteenReview 回复食堂评价（管理端）
func AdminReplyCanteenReview(c *gin.Context) {
	reviewID := c.Param("id")
	
	var review models.CanteenReview
	if err := database.DB.First(&review, reviewID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "评价不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评价失败"})
		}
		return
	}

	var req struct {
		AdminReply string `json:"admin_reply" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新回复
	now := time.Now()
	review.AdminReply = req.AdminReply
	review.AdminReplyTime = &now

	if err := database.DB.Save(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "回复失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": review,
		"message": "回复成功",
	})
}

// AdminGetCanteenStats 获取食堂评价统计（管理端）
func AdminGetCanteenStats(c *gin.Context) {
	stats, err := getCanteenReviewStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}

	// 获取更详细的统计信息
	var monthlyStats []struct {
		Month        string  `json:"month"`
		ReviewCount  int64   `json:"review_count"`
		AvgScore     float64 `json:"avg_score"`
	}

	// 获取最近12个月的统计
	if err := database.DB.Raw(`
		SELECT 
			DATE_FORMAT(created_at, '%Y-%m') as month,
			COUNT(*) as review_count,
			AVG((environment_score + service_score + quality_score + price_score + overall_score) / 5.0) as avg_score
		FROM canteen_reviews 
		WHERE created_at >= DATE_SUB(NOW(), INTERVAL 12 MONTH)
		GROUP BY DATE_FORMAT(created_at, '%Y-%m')
		ORDER BY month DESC
	`).Scan(&monthlyStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取月度统计失败"})
		return
	}

	// 获取差评列表（评分低于3分的）
	var lowScoreReviews []models.CanteenReview
	database.DB.Where("overall_score < ?", 3).
		Preload("User").
		Order("created_at DESC").
		Limit(10).
		Find(&lowScoreReviews)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"overview":          stats,
			"monthly_stats":     monthlyStats,
			"low_score_reviews": lowScoreReviews,
		},
	})
}

// AdminDeleteCanteenReview 删除食堂评价（管理端）
func AdminDeleteCanteenReview(c *gin.Context) {
	reviewID := c.Param("id")
	
	var review models.CanteenReview
	if err := database.DB.First(&review, reviewID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评价不存在"})
		return
	}

	// 开始事务
	tx := database.DB.Begin()

	// 删除评价图片
	if err := tx.Where("review_id = ?", reviewID).Delete(&models.CanteenReviewImage{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除评价图片失败"})
		return
	}

	// 删除评价
	if err := tx.Delete(&review).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除评价失败"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "删除成功",
	})
}