package controllers

import (
	"net/http"
	"strconv"
	"time"

	"orderfood/database"
	"orderfood/models"

	"github.com/gin-gonic/gin"
)

// GetDutyMealMenu 获取值班餐菜单（用户端）
func GetDutyMealMenu(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	// 检查用户是否为值班人员
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	if !user.IsDutyStaff {
		c.JSON(http.StatusForbidden, gin.H{"error": "您不是值班人员，无法订购值班餐"})
		return
	}

	// 获取当前可用的值班餐设置
	var settings []models.DutyMealSetting
	if err := database.DB.Where("status = ?", 1).Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取值班餐设置失败"})
		return
	}

	// 检查是否在值班时间内
	var availableSettings []models.DutyMealSetting
	for _, setting := range settings {
		if setting.IsInDutyTime() {
			availableSettings = append(availableSettings, setting)
		}
	}

	// 获取值班餐专用菜品（这里假设有特定分类或标记）
	var dishes []models.Dish
	if err := database.DB.Where("status = ? AND is_duty_meal = ?", 1, true).
		Preload("Images").Find(&dishes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取值班餐菜品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"settings": availableSettings,
			"dishes":   dishes,
			"balance":  user.SubsidyBalance,
		},
	})
}

// CreateDutyMealOrder 创建值班餐订单（用户端）
func CreateDutyMealOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req struct {
		SettingID    uint `json:"setting_id" binding:"required"`
		DeliveryTime time.Time `json:"delivery_time" binding:"required"`
		Items        []struct {
			DishID   uint `json:"dish_id" binding:"required"`
			Quantity int  `json:"quantity" binding:"required,min=1"`
		} `json:"items" binding:"required,min=1"`
		Remark string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否为值班人员
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	if !user.IsDutyStaff {
		c.JSON(http.StatusForbidden, gin.H{"error": "您不是值班人员，无法订购值班餐"})
		return
	}

	// 检查值班餐设置
	var setting models.DutyMealSetting
	if err := database.DB.First(&setting, req.SettingID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "值班餐设置不存在"})
		return
	}

	if setting.Status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该值班餐设置已禁用"})
		return
	}

	// 验证配送时间
	if req.DeliveryTime.Before(time.Now().Add(30 * time.Minute)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "配送时间至少提前30分钟"})
		return
	}

	// 计算订单总额
	var totalAmount float64
	var orderItems []models.DutyMealOrderItem

	for _, item := range req.Items {
		var dish models.Dish
		if err := database.DB.First(&dish, item.DishID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "菜品不存在"})
			return
		}

		if dish.Status != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "菜品已下架"})
			return
		}

		itemTotal := dish.Price * float64(item.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, models.DutyMealOrderItem{
			DishID:   item.DishID,
			Quantity: item.Quantity,
			Price:    dish.Price,
		})
	}

	// 计算餐补使用
	subsidyUsed := setting.Subsidy
	if subsidyUsed > totalAmount {
		subsidyUsed = totalAmount
	}
	if subsidyUsed > user.SubsidyBalance {
		subsidyUsed = user.SubsidyBalance
	}

	actualPaid := totalAmount - subsidyUsed

	// 开始事务
	tx := database.DB.Begin()

	// 创建订单
	order := models.DutyMealOrder{
		UserID:       userID.(uint),
		SettingID:    req.SettingID,
		OrderDate:    time.Now(),
		DeliveryTime: req.DeliveryTime,
		TotalAmount:  totalAmount,
		SubsidyUsed:  subsidyUsed,
		ActualPaid:   actualPaid,
		Status:       models.DutyMealOrderStatusPending,
		Remark:       req.Remark,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单失败"})
		return
	}

	// 创建订单详情
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		if err := tx.Create(&orderItems[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单详情失败"})
			return
		}
	}

	// 扣减餐补余额
	if subsidyUsed > 0 {
		if err := tx.Model(&user).Update("subsidy_balance", user.SubsidyBalance-subsidyUsed).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "扣减餐补失败"})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": order,
		"message": "订单创建成功",
	})
}

// GetDutyMealOrders 获取值班餐订单列表（用户端）
func GetDutyMealOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize

	var total int64
	database.DB.Model(&models.DutyMealOrder{}).Where("user_id = ?", userID).Count(&total)

	var orders []models.DutyMealOrder
	if err := database.DB.Where("user_id = ?", userID).
		Preload("Setting").Preload("Items").Preload("Items.Dish").
		Offset(offset).Limit(pageSize).Order("created_at DESC").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      orders,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetSubsidyBalance 获取餐补余额（用户端）
func GetSubsidyBalance(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"balance":       user.SubsidyBalance,
			"is_duty_staff": user.IsDutyStaff,
		},
	})
}

// === 管理端接口 ===

// AdminGetDutyMealSettings 获取值班餐设置列表（管理端）
func AdminGetDutyMealSettings(c *gin.Context) {
	var settings []models.DutyMealSetting
	if err := database.DB.Order("created_at DESC").Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取设置列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": settings,
	})
}

// AdminCreateDutyMealSetting 创建值班餐设置（管理端）
func AdminCreateDutyMealSetting(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		StartTime   string  `json:"start_time" binding:"required"`
		EndTime     string  `json:"end_time" binding:"required"`
		Subsidy     float64 `json:"subsidy" binding:"required,min=0"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setting := models.DutyMealSetting{
		Name:        req.Name,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Subsidy:     req.Subsidy,
		Description: req.Description,
		Status:      1,
	}

	if err := database.DB.Create(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": setting,
		"message": "创建成功",
	})
}

// AdminUpdateDutyMealSetting 编辑值班餐设置（管理端）
func AdminUpdateDutyMealSetting(c *gin.Context) {
	settingID := c.Param("id")
	
	var setting models.DutyMealSetting
	if err := database.DB.First(&setting, settingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设置不存在"})
		return
	}

	var req struct {
		Name        string  `json:"name" binding:"required"`
		StartTime   string  `json:"start_time" binding:"required"`
		EndTime     string  `json:"end_time" binding:"required"`
		Subsidy     float64 `json:"subsidy" binding:"required,min=0"`
		Status      int     `json:"status" binding:"required"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setting.Name = req.Name
	setting.StartTime = req.StartTime
	setting.EndTime = req.EndTime
	setting.Subsidy = req.Subsidy
	setting.Status = req.Status
	setting.Description = req.Description

	if err := database.DB.Save(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": setting,
		"message": "更新成功",
	})
}

// AdminDeleteDutyMealSetting 删除值班餐设置（管理端）
func AdminDeleteDutyMealSetting(c *gin.Context) {
	settingID := c.Param("id")
	
	var setting models.DutyMealSetting
	if err := database.DB.First(&setting, settingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设置不存在"})
		return
	}

	// 检查是否有关联的订单
	var orderCount int64
	database.DB.Model(&models.DutyMealOrder{}).Where("setting_id = ?", settingID).Count(&orderCount)
	if orderCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该设置下有订单记录，无法删除"})
		return
	}

	if err := database.DB.Delete(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "删除成功",
	})
}

// AdminGetDutyMealOrders 获取值班餐订单列表（管理端）
func AdminGetDutyMealOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	settingID := c.Query("setting_id")

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.DutyMealOrder{})

	// 筛选条件
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if settingID != "" {
		query = query.Where("setting_id = ?", settingID)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var orders []models.DutyMealOrder
	if err := query.Preload("User").Preload("Setting").Preload("Items").Preload("Items.Dish").
		Offset(offset).Limit(pageSize).Order("created_at DESC").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      orders,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// AdminUpdateDutyMealOrderStatus 更新值班餐订单状态（管理端）
func AdminUpdateDutyMealOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	
	var order models.DutyMealOrder
	if err := database.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证状态值
	if req.Status < 1 || req.Status > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的状态值"})
		return
	}

	order.Status = req.Status
	if err := database.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "更新成功",
	})
}

// AdminGetDutyMealStats 获取值班餐统计（管理端）
func AdminGetDutyMealStats(c *gin.Context) {
	stats, err := getDutyMealStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}

	// 获取月度统计
	var monthlyStats []struct {
		Month       string  `json:"month"`
		OrderCount  int64   `json:"order_count"`
		TotalAmount float64 `json:"total_amount"`
		TotalSubsidy float64 `json:"total_subsidy"`
	}

	if err := database.DB.Raw(`
		SELECT 
			DATE_FORMAT(created_at, '%Y-%m') as month,
			COUNT(*) as order_count,
			SUM(total_amount) as total_amount,
			SUM(subsidy_used) as total_subsidy
		FROM duty_meal_orders 
		WHERE status != ? AND created_at >= DATE_SUB(NOW(), INTERVAL 12 MONTH)
		GROUP BY DATE_FORMAT(created_at, '%Y-%m')
		ORDER BY month DESC
	`, models.DutyMealOrderStatusCancelled).Scan(&monthlyStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取月度统计失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"overview":       stats,
			"monthly_stats":  monthlyStats,
		},
	})
}

// getDutyMealStats 获取值班餐统计
func getDutyMealStats() (*models.DutyMealStats, error) {
	var stats models.DutyMealStats
	
	// 总统计
	if err := database.DB.Model(&models.DutyMealOrder{}).
		Select("COUNT(*) as total_orders, COALESCE(SUM(total_amount), 0) as total_amount, COALESCE(SUM(subsidy_used), 0) as total_subsidy").
		Where("status != ?", models.DutyMealOrderStatusCancelled).
		Scan(&stats).Error; err != nil {
		return nil, err
	}

	// 本月统计
	if err := database.DB.Model(&models.DutyMealOrder{}).
		Select("COUNT(*) as monthly_orders, COALESCE(SUM(total_amount), 0) as monthly_amount, COALESCE(SUM(subsidy_used), 0) as monthly_subsidy").
		Where("status != ? AND YEAR(created_at) = YEAR(NOW()) AND MONTH(created_at) = MONTH(NOW())", models.DutyMealOrderStatusCancelled).
		Scan(&stats).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}