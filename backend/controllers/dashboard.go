package controllers

import (
	"net/http"
	"orderfood/database"
	"orderfood/models"
	"time"

	"github.com/gin-gonic/gin"
)

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	DishCount   int64 `json:"dish_count"`
	OrderCount  int64 `json:"order_count"`
	UserCount   int64 `json:"user_count"`
	ReviewCount int64 `json:"review_count"`
}

// GetDashboardStats 获取仪表盘统计数据（管理员）
func GetDashboardStats(c *gin.Context) {
	var stats DashboardStats

	// 获取菜品总数
	if err := database.DB.Model(&models.Dish{}).Count(&stats.DishCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取菜品统计失败"})
		return
	}

	// 获取订单总数
	if err := database.DB.Model(&models.Order{}).Count(&stats.OrderCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单统计失败"})
		return
	}

	// 获取用户总数
	if err := database.DB.Model(&models.User{}).Count(&stats.UserCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户统计失败"})
		return
	}

	// 获取评价总数
	if err := database.DB.Model(&models.Review{}).Count(&stats.ReviewCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评价统计失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}

// OrderStatusStats 订单状态统计
type OrderStatusStats struct {
	Pending    int64 `json:"pending"`    // 待处理
	Processing int64 `json:"processing"` // 制作中
	Completed  int64 `json:"completed"`  // 已完成
	Cancelled  int64 `json:"cancelled"`  // 已取消
}

// GetOrderStatusStats 获取订单状态统计（管理员）
func GetOrderStatusStats(c *gin.Context) {
	var stats OrderStatusStats

	// 获取各状态订单数量
	database.DB.Model(&models.Order{}).Where("status = ?", 1).Count(&stats.Pending)
	database.DB.Model(&models.Order{}).Where("status = ?", 2).Count(&stats.Processing)
	database.DB.Model(&models.Order{}).Where("status = ?", 3).Count(&stats.Completed)
	database.DB.Model(&models.Order{}).Where("status = ?", 4).Count(&stats.Cancelled)

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}

// GetDashboardChartData 获取仪表盘图表数据（管理员）
func GetDashboardChartData(c *gin.Context) {
	// 获取最近7天的订单统计
	var orderStats []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	// 获取最近7天的数据
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var count int64
		
		database.DB.Model(&models.Order{}).
			Where("DATE(created_at) = ?", date).
			Count(&count)
		
		orderStats = append(orderStats, struct {
			Date  string `json:"date"`
			Count int64  `json:"count"`
		}{
			Date:  date,
			Count: count,
		})
	}

	// 获取菜品分类统计
	var categoryStats []struct {
		CategoryName string `json:"category_name"`
		DishCount    int64  `json:"dish_count"`
	}

	database.DB.Table("dishes").
		Select("categories.name as category_name, COUNT(dishes.id) as dish_count").
		Joins("LEFT JOIN categories ON dishes.category_id = categories.id").
		Where("dishes.deleted_at IS NULL").
		Group("categories.id, categories.name").
		Scan(&categoryStats)

	// 获取订单状态统计
	var statusStats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	statusMap := map[int]string{
		1: "待处理",
		2: "制作中", 
		3: "已完成",
		4: "已取消",
	}

	for status, name := range statusMap {
		var count int64
		database.DB.Model(&models.Order{}).Where("status = ?", status).Count(&count)
		statusStats = append(statusStats, struct {
			Status string `json:"status"`
			Count  int64  `json:"count"`
		}{
			Status: name,
			Count:  count,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"order_trend":    orderStats,
			"category_stats": categoryStats,
			"status_stats":   statusStats,
		},
	})
}