package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"orderfood/database"
	"orderfood/models"
)

// GetUsers 获取用户列表（管理员）
func GetUsers(c *gin.Context) {
	var users []models.User
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	username := c.Query("username")
	phone := c.Query("phone")
	email := c.Query("email")
	status := c.Query("status")

	// 构建查询
	query := database.DB.Model(&models.User{})

	// 添加搜索条件
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetUserStats 获取用户统计信息（管理员）
func GetUserStats(c *gin.Context) {
	var stats struct {
		Total     int64 `json:"total"`
		Active    int64 `json:"active"`
		Disabled  int64 `json:"disabled"`
		NewToday  int64 `json:"new_today"`
	}

	// 总用户数
	database.DB.Model(&models.User{}).Count(&stats.Total)

	// 正常用户数
	database.DB.Model(&models.User{}).Where("status = ?", 1).Count(&stats.Active)

	// 禁用用户数
	database.DB.Model(&models.User{}).Where("status = ?", 0).Count(&stats.Disabled)

	// 今日新增用户数
	database.DB.Model(&models.User{}).
		Where("DATE(created_at) = CURDATE()").
		Count(&stats.NewToday)

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// UpdateUserStatus 更新用户状态（管理员）
func UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var updateData struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新状态
	if err := database.DB.Model(&user).Update("status", updateData.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}

// BatchUpdateUserStatus 批量更新用户状态（管理员）
func BatchUpdateUserStatus(c *gin.Context) {
	var requestData struct {
		UserIDs []uint `json:"user_ids" binding:"required"`
		Status  int    `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 批量更新
	if err := database.DB.Model(&models.User{}).
		Where("id IN ?", requestData.UserIDs).
		Update("status", requestData.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "批量更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "批量更新成功"})
}

// GetUserDetail 获取用户详情（管理员）
func GetUserDetail(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 获取用户订单统计
	var orderStats struct {
		Total     int64 `json:"total"`
		Completed int64 `json:"completed"`
	}

	database.DB.Model(&models.Order{}).Where("user_id = ?", id).Count(&orderStats.Total)
	database.DB.Model(&models.Order{}).Where("user_id = ? AND status = ?", id, 3).Count(&orderStats.Completed)

	// 获取用户评价统计
	var reviewCount int64
	database.DB.Model(&models.Review{}).Where("user_id = ?", id).Count(&reviewCount)

	c.JSON(http.StatusOK, gin.H{
		"data": user,
		"stats": gin.H{
			"orders":  orderStats,
			"reviews": reviewCount,
		},
	})
}

// DeleteUser 删除用户（管理员）
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 检查用户是否有未完成的订单
	var pendingOrderCount int64
	database.DB.Model(&models.Order{}).
		Where("user_id = ? AND status IN ?", id, []int{1, 2}).
		Count(&pendingOrderCount)

	if pendingOrderCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该用户还有未完成的订单，无法删除"})
		return
	}

	// 软删除用户（如果模型支持）或物理删除
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}