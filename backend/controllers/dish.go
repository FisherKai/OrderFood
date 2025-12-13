package controllers

import (
	"fmt"
	"net/http"
	"orderfood/database"
	"orderfood/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetDishes 获取菜品列表
func GetDishes(c *gin.Context) {
	var dishes []models.Dish
	query := database.DB.Preload("Category").Preload("Images", "is_deleted = ?", false).Where("status = ?", 1)

	// 分类筛选
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 搜索
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	var total int64
	query.Model(&models.Dish{}).Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&dishes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取菜品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dishes,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetDishDetail 获取菜品详情
func GetDishDetail(c *gin.Context) {
	id := c.Param("id")
	var dish models.Dish

	if err := database.DB.Preload("Category").Preload("Images", "is_deleted = ?", false).First(&dish, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜品不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": dish})
}

// CreateDish 创建菜品（管理员）
func CreateDish(c *gin.Context) {
	var dish models.Dish
	if err := c.ShouldBindJSON(&dish); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&dish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建菜品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "data": dish})
}

// UpdateDish 更新菜品（管理员）
func UpdateDish(c *gin.Context) {
	id := c.Param("id")
	var dish models.Dish

	if err := database.DB.First(&dish, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜品不存在"})
		return
	}

	// 解析请求数据
	var updateData struct {
		models.Dish
		Images []models.DishImage `json:"images"`
	}
	
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 开始事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新菜品基本信息
	dish.CategoryID = updateData.CategoryID
	dish.Name = updateData.Name
	dish.Price = updateData.Price
	dish.Description = updateData.Description
	dish.Status = updateData.Status
	dish.Stock = updateData.Stock

	if err := tx.Save(&dish).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新菜品失败"})
		return
	}

	// 处理图片数据
	// 删除该菜品的所有现有图片记录
	result := tx.Where("dish_id = ?", id).Delete(&models.DishImage{})
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除原有图片失败"})
		return
	}
	// 添加日志
	fmt.Printf("删除了 %d 条图片记录\n", result.RowsAffected)

	// 创建新的图片记录
	if len(updateData.Images) > 0 {
		for _, img := range updateData.Images {
			img.DishID = dish.ID
			img.IsDeleted = false
			
			if err := tx.Create(&img).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "创建图片失败"})
				return
			}
		}
		fmt.Printf("创建了 %d 条新图片记录\n", len(updateData.Images))
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": dish})
}

// DeleteDish 删除菜品（管理员）
func DeleteDish(c *gin.Context) {
	id := c.Param("id")

	// 检查是否有订单使用了该菜品
	var count int64
	database.DB.Model(&models.OrderItem{}).Where("dish_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该菜品已有订单，不能删除，只能下架"})
		return
	}

	if err := database.DB.Delete(&models.Dish{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetDishImages 获取菜品所有图片（包括历史图片）
func GetDishImages(c *gin.Context) {
	id := c.Param("id")
	var images []models.DishImage

	if err := database.DB.Where("dish_id = ?", id).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取图片失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": images})
}

// RestoreDishImage 恢复历史图片
func RestoreDishImage(c *gin.Context) {
	var req struct {
		ImageID uint `json:"image_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&models.DishImage{}).Where("id = ?", req.ImageID).Updates(map[string]interface{}{
		"is_deleted": false,
		"deleted_at": nil,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "恢复失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "恢复成功"})
}
