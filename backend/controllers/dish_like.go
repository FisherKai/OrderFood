package controllers

import (
	"net/http"
	"orderfood/database"
	"orderfood/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LikeDish 点赞菜品
func LikeDish(c *gin.Context) {
	userID, _ := c.Get("user_id")
	dishID := c.Param("id")

	dishIDUint, err := strconv.ParseUint(dishID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的菜品ID"})
		return
	}

	// 检查菜品是否存在
	var dish models.Dish
	if err := database.DB.First(&dish, dishIDUint).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜品不存在"})
		return
	}

	// 检查是否已经点赞
	var existingLike models.DishLike
	if err := database.DB.Where("user_id = ? AND dish_id = ?", userID, dishIDUint).First(&existingLike).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已经点赞过了"})
		return
	}

	// 开始事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建点赞记录
	like := models.DishLike{
		UserID: userID.(uint),
		DishID: uint(dishIDUint),
	}

	if err := tx.Create(&like).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}

	// 更新菜品点赞数
	if err := tx.Model(&dish).Update("like_count", dish.LikeCount+1).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新点赞数失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
}

// UnlikeDish 取消点赞菜品
func UnlikeDish(c *gin.Context) {
	userID, _ := c.Get("user_id")
	dishID := c.Param("id")

	dishIDUint, err := strconv.ParseUint(dishID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的菜品ID"})
		return
	}

	// 检查菜品是否存在
	var dish models.Dish
	if err := database.DB.First(&dish, dishIDUint).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜品不存在"})
		return
	}

	// 检查是否已经点赞
	var existingLike models.DishLike
	if err := database.DB.Where("user_id = ? AND dish_id = ?", userID, dishIDUint).First(&existingLike).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "还未点赞"})
		return
	}

	// 开始事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除点赞记录
	if err := tx.Delete(&existingLike).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消点赞失败"})
		return
	}

	// 更新菜品点赞数
	newLikeCount := dish.LikeCount - 1
	if newLikeCount < 0 {
		newLikeCount = 0
	}
	if err := tx.Model(&dish).Update("like_count", newLikeCount).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新点赞数失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消点赞成功"})
}

// GetUserLikedDishes 获取用户点赞的菜品
func GetUserLikedDishes(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	offset := (page - 1) * pageSize

	var total int64
	database.DB.Model(&models.DishLike{}).Where("user_id = ?", userID).Count(&total)

	var likes []models.DishLike
	if err := database.DB.Where("user_id = ?", userID).
		Preload("Dish.Images").
		Preload("Dish.Category").
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&likes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取点赞菜品失败"})
		return
	}

	// 提取菜品信息
	dishes := make([]models.Dish, len(likes))
	for i, like := range likes {
		dishes[i] = like.Dish
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dishes,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetPopularDishes 获取热门菜品（基于点赞数）
func GetPopularDishes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	offset := (page - 1) * pageSize

	var total int64
	database.DB.Model(&models.Dish{}).Where("status = 1").Count(&total)

	var dishes []models.Dish
	if err := database.DB.Where("status = 1").
		Preload("Images").
		Preload("Category").
		Offset(offset).Limit(pageSize).
		Order("like_count DESC, created_at DESC").
		Find(&dishes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取热门菜品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dishes,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetDishLikeStats 获取菜品点赞统计（管理员）
func GetDishLikeStats(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	offset := (page - 1) * pageSize

	var total int64
	database.DB.Model(&models.Dish{}).Count(&total)

	var dishes []models.Dish
	if err := database.DB.Preload("Images").
		Preload("Category").
		Offset(offset).Limit(pageSize).
		Order("like_count DESC").
		Find(&dishes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取点赞统计失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dishes,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetDishLikeRanking 获取菜品点赞排行榜（管理员）
func GetDishLikeRanking(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var dishes []models.Dish
	if err := database.DB.Where("like_count > 0").
		Preload("Images").
		Preload("Category").
		Limit(limit).
		Order("like_count DESC").
		Find(&dishes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取点赞排行榜失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dishes})
}