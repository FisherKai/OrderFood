package controllers

import (
	"net/http"
	"orderfood/database"
	"orderfood/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateReviewRequest struct {
	DishID  uint   `json:"dish_id" binding:"required"`
	OrderID uint   `json:"order_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Content string `json:"content"`
	Images  []string `json:"images"`
}

// CreateReview 创建评价
func CreateReview(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req CreateReviewRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查订单是否属于该用户且已完成
	var order models.Order
	if err := database.DB.Where("id = ? AND user_id = ? AND status = ?", req.OrderID, userID, 3).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单不存在或未完成"})
		return
	}

	// 检查是否已评价
	var existReview models.Review
	if err := database.DB.Where("user_id = ? AND dish_id = ? AND order_id = ?", userID, req.DishID, req.OrderID).First(&existReview).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已评价过该菜品"})
		return
	}

	tx := database.DB.Begin()

	review := models.Review{
		UserID:  userID.(uint),
		DishID:  req.DishID,
		OrderID: req.OrderID,
		Rating:  req.Rating,
		Content: req.Content,
	}

	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评价失败"})
		return
	}

	// 添加评价图片
	for _, imageURL := range req.Images {
		reviewImage := models.ReviewImage{
			ReviewID: review.ID,
			ImageURL: imageURL,
		}
		if err := tx.Create(&reviewImage).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "添加图片失败"})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "评价成功", "data": review})
}

// GetDishReviews 获取菜品评价列表
func GetDishReviews(c *gin.Context) {
	dishID := c.Param("dishId")
	var reviews []models.Review

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	query := database.DB.Where("dish_id = ?", dishID).Preload("User").Preload("Images").Order("created_at DESC")

	var total int64
	query.Model(&models.Review{}).Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评价失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": reviews,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetAllReviews 获取所有评价（管理员）
func GetAllReviews(c *gin.Context) {
	var reviews []models.Review

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	query := database.DB.Preload("User").Preload("Dish").Preload("Images").Order("created_at DESC")

	var total int64
	query.Model(&models.Review{}).Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取评价失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": reviews,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
