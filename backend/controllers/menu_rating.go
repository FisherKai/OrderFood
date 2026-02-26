package controllers

import (
	"net/http"
	"orderfood/database"
	"orderfood/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateMenuRatingRequest 创建评价请求
type CreateMenuRatingRequest struct {
	MenuID     uint   `json:"menu_id" binding:"required"`
	MenuItemID uint   `json:"menu_item_id" binding:"required"`
	DishID     uint   `json:"dish_id" binding:"required"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	Comment    string `json:"comment"`
	MealType   int    `json:"meal_type" binding:"required,min=1,max=4"`
	RatingDate string `json:"rating_date" binding:"required"` // 格式: 2006-01-02
}

// CreateMenuRating 创建菜品评价（用户端）
func CreateMenuRating(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID类型错误"})
		return
	}

	var req CreateMenuRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 解析日期
	ratingDate, err := time.ParseInLocation("2006-01-02", req.RatingDate, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式错误"})
		return
	}

	// 检查是否已经评价过
	var existingRating models.MenuItemRating
	if err := database.DB.Where(
		"menu_id = ? AND dish_id = ? AND user_id = ? AND meal_type = ? AND DATE(rating_date) = DATE(?)",
		req.MenuID, req.DishID, userID, req.MealType, ratingDate,
	).First(&existingRating).Error; err == nil {
		// 已存在评价，更新它
		existingRating.Rating = req.Rating
		existingRating.Comment = req.Comment
		if err := database.DB.Save(&existingRating).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新评价失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "评价更新成功", "data": existingRating})
		return
	}

	// 创建新评价
	rating := models.MenuItemRating{
		MenuID:     req.MenuID,
		MenuItemID: req.MenuItemID,
		DishID:     req.DishID,
		UserID:     userID,
		Rating:     req.Rating,
		Comment:    req.Comment,
		MealType:   req.MealType,
		RatingDate: ratingDate,
	}

	if err := database.DB.Create(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评价失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评价成功", "data": rating})
}

// GetMyMenuRatings 获取用户对某菜谱的评价（用户端）
func GetMyMenuRatings(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID类型错误"})
		return
	}

	menuID := c.Query("menu_id")
	dateStr := c.Query("date") // 可选，按日期筛选

	query := database.DB.Model(&models.MenuItemRating{}).Where("user_id = ?", userID)

	if menuID != "" {
		query = query.Where("menu_id = ?", menuID)
	}

	if dateStr != "" {
		query = query.Where("DATE(rating_date) = DATE(?)", dateStr)
	}

	var ratings []models.MenuItemRating
	if err := query.Preload("Dish").Find(&ratings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评价失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ratings})
}

// GetMenuItemRatings 获取某菜谱菜品的所有评价（管理端）
func GetMenuItemRatings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	menuID := c.Query("menu_id")
	dishID := c.Query("dish_id")
	dateStr := c.Query("date")
	rating := c.Query("rating") // 按评分筛选

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.MenuItemRating{})

	if menuID != "" {
		query = query.Where("menu_id = ?", menuID)
	}

	if dishID != "" {
		query = query.Where("dish_id = ?", dishID)
	}

	if dateStr != "" {
		query = query.Where("DATE(rating_date) = DATE(?)", dateStr)
	}

	if rating != "" {
		ratingInt, _ := strconv.Atoi(rating)
		query = query.Where("rating = ?", ratingInt)
	}

	var total int64
	query.Count(&total)

	var ratings []models.MenuItemRating
	if err := query.
		Preload("Dish").
		Preload("User").
		Preload("Menu").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&ratings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评价列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ratings,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetMenuRatingStats 获取菜谱评价统计（管理端）
func GetMenuRatingStats(c *gin.Context) {
	menuID := c.Query("menu_id")

	type RatingStats struct {
		DishID      uint    `json:"dish_id"`
		DishName    string  `json:"dish_name"`
		TotalCount  int64   `json:"total_count"`
		AvgRating   float64 `json:"avg_rating"`
		Rating5     int64   `json:"rating_5"`
		Rating4     int64   `json:"rating_4"`
		Rating3     int64   `json:"rating_3"`
		Rating2     int64   `json:"rating_2"`
		Rating1     int64   `json:"rating_1"`
	}

	var stats []RatingStats

	query := database.DB.Table("menu_item_ratings").
		Select(`
			menu_item_ratings.dish_id,
			dishes.name as dish_name,
			COUNT(*) as total_count,
			AVG(rating) as avg_rating,
			SUM(CASE WHEN rating = 5 THEN 1 ELSE 0 END) as rating_5,
			SUM(CASE WHEN rating = 4 THEN 1 ELSE 0 END) as rating_4,
			SUM(CASE WHEN rating = 3 THEN 1 ELSE 0 END) as rating_3,
			SUM(CASE WHEN rating = 2 THEN 1 ELSE 0 END) as rating_2,
			SUM(CASE WHEN rating = 1 THEN 1 ELSE 0 END) as rating_1
		`).
		Joins("LEFT JOIN dishes ON dishes.id = menu_item_ratings.dish_id").
		Where("menu_item_ratings.deleted_at IS NULL")

	if menuID != "" {
		query = query.Where("menu_item_ratings.menu_id = ?", menuID)
	}

	if err := query.Group("menu_item_ratings.dish_id, dishes.name").
		Order("avg_rating DESC").
		Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}

	// 获取总体统计
	var overallStats struct {
		TotalRatings int64   `json:"total_ratings"`
		AvgRating    float64 `json:"avg_rating"`
	}

	overallQuery := database.DB.Table("menu_item_ratings").
		Select("COUNT(*) as total_ratings, AVG(rating) as avg_rating").
		Where("deleted_at IS NULL")

	if menuID != "" {
		overallQuery = overallQuery.Where("menu_id = ?", menuID)
	}

	overallQuery.Scan(&overallStats)

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"dish_stats": stats,
			"overall": gin.H{
				"total_ratings": overallStats.TotalRatings,
				"avg_rating":    overallStats.AvgRating,
			},
		},
	})
}

// DeleteMenuRating 删除评价（管理端）
func DeleteMenuRating(c *gin.Context) {
	id := c.Param("id")

	var rating models.MenuItemRating
	if err := database.DB.First(&rating, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评价不存在"})
		return
	}

	if err := database.DB.Delete(&rating).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除评价失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetDishAvgRating 获取菜品的平均评分（公开）
func GetDishAvgRating(c *gin.Context) {
	dishID := c.Param("dish_id")

	var result struct {
		AvgRating   float64 `json:"avg_rating"`
		TotalCount  int64   `json:"total_count"`
	}

	if err := database.DB.Table("menu_item_ratings").
		Select("AVG(rating) as avg_rating, COUNT(*) as total_count").
		Where("dish_id = ? AND deleted_at IS NULL", dishID).
		Scan(&result).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评分失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
