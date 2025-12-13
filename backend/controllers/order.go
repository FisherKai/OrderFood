package controllers

import (
	"net/http"
	"orderfood/database"
	"orderfood/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateOrderRequest struct {
	Items []struct {
		DishID   uint `json:"dish_id" binding:"required"`
		Quantity int  `json:"quantity" binding:"required,min=1"`
	} `json:"items" binding:"required,min=1"`
}

type CreateReservationRequest struct {
	Items []struct {
		DishID   uint `json:"dish_id" binding:"required"`
		Quantity int  `json:"quantity" binding:"required,min=1"`
	} `json:"items" binding:"required,min=1"`
	ReserveTime time.Time `json:"reserve_time" binding:"required"`
	PeopleCount int       `json:"people_count" binding:"required,min=1"`
}

// AddToCart 添加到购物车（暂时只返回成功，实际可以用Redis实现）
func AddToCart(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "添加成功"})
}

// CreateOrder 创建订单
func CreateOrder(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 开始事务
	tx := database.DB.Begin()

	// 创建订单
	order := models.Order{
		UserID:     userID.(uint),
		OrderType:  1,
		Status:     1,
		TotalPrice: 0,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单失败"})
		return
	}

	// 创建订单项并计算总价
	var totalPrice float64
	for _, item := range req.Items {
		var dish models.Dish
		if err := tx.First(&dish, item.DishID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "菜品不存在"})
			return
		}

		orderItem := models.OrderItem{
			OrderID:  order.ID,
			DishID:   item.DishID,
			Quantity: item.Quantity,
			Price:    dish.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单项失败"})
			return
		}

		totalPrice += dish.Price * float64(item.Quantity)
	}

	// 更新订单总价
	order.TotalPrice = totalPrice
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新订单失败"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "订单创建成功", "data": order})
}

// CreateReservation 创建预约订单
func CreateReservation(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req CreateReservationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查预约时间是否至少提前2小时
	if time.Until(req.ReserveTime) < 2*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "预约时间至少提前2小时"})
		return
	}

	tx := database.DB.Begin()

	order := models.Order{
		UserID:      userID.(uint),
		OrderType:   2,
		Status:      1,
		TotalPrice:  0,
		ReserveTime: &req.ReserveTime,
		PeopleCount: req.PeopleCount,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建预约订单失败"})
		return
	}

	var totalPrice float64
	for _, item := range req.Items {
		var dish models.Dish
		if err := tx.First(&dish, item.DishID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "菜品不存在"})
			return
		}

		orderItem := models.OrderItem{
			OrderID:  order.ID,
			DishID:   item.DishID,
			Quantity: item.Quantity,
			Price:    dish.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单项失败"})
			return
		}

		totalPrice += dish.Price * float64(item.Quantity)
	}

	order.TotalPrice = totalPrice
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新订单失败"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "预约订单创建成功", "data": order})
}

// GetUserOrders 获取用户订单列表
func GetUserOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var orders []models.Order

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	query := database.DB.Where("user_id = ?", userID).Preload("Items.Dish").Order("created_at DESC")

	var total int64
	query.Model(&models.Order{}).Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetAllOrders 获取所有订单（管理员）
func GetAllOrders(c *gin.Context) {
	var orders []models.Order

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	query := database.DB.Preload("User").Preload("Items.Dish").Order("created_at DESC")

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Model(&models.Order{}).Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// UpdateOrderStatus 更新订单状态（管理员）
func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status int `json:"status" binding:"required,min=1,max=4"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&models.Order{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}
