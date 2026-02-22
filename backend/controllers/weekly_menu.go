package controllers

import (
	"net/http"
	"orderfood/database"
	"orderfood/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateWeeklyMenuRequest struct {
	WeekStart time.Time `json:"week_start" binding:"required"`
	Title     string    `json:"title" binding:"required"`
	MenuItems []struct {
		Date     time.Time `json:"date" binding:"required"`
		MealType int       `json:"meal_type" binding:"required,min=1,max=4"`
		DishID   uint      `json:"dish_id" binding:"required"`
		Sort     int       `json:"sort"`
	} `json:"menu_items"`
}

// CreateWeeklyMenu 创建一周菜谱（管理员）
func CreateWeeklyMenu(c *gin.Context) {
	adminIDInterface, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未找到管理员ID"})
		return
	}
	
	adminID, ok := adminIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "管理员ID类型错误"})
		return
	}
	
	// 直接使用map解析JSON，避免结构体验证
	var rawReq map[string]interface{}
	if err := c.ShouldBindJSON(&rawReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 手动验证必需字段
	if rawReq["title"] == nil || rawReq["title"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}
	if rawReq["week_start"] == nil || rawReq["week_start"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "week_start is required"})
		return
	}
	
	// 解析时间
	weekStartStr, ok := rawReq["week_start"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "week_start must be a string"})
		return
	}
	
	weekStart, err := time.Parse(time.RFC3339, weekStartStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid week_start format: " + err.Error()})
		return
	}
	
	title := rawReq["title"].(string)
	
	// 解析menu_items（可选）
	var menuItems []models.MenuItem
	if menuItemsRaw, exists := rawReq["menu_items"]; exists && menuItemsRaw != nil {
		if items, ok := menuItemsRaw.([]interface{}); ok {
			for _, item := range items {
				if itemMap, ok := item.(map[string]interface{}); ok {
					// 解析每个菜谱项
					dateStr, ok := itemMap["date"].(string)
					if !ok {
						c.JSON(http.StatusBadRequest, gin.H{"error": "menu item date is required"})
						return
					}
					
					date, err := time.Parse(time.RFC3339, dateStr)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "invalid menu item date format"})
						return
					}
					
					mealType, ok := itemMap["meal_type"].(float64)
					if !ok {
						c.JSON(http.StatusBadRequest, gin.H{"error": "menu item meal_type is required"})
						return
					}
					
					dishID, ok := itemMap["dish_id"].(float64)
					if !ok {
						c.JSON(http.StatusBadRequest, gin.H{"error": "menu item dish_id is required"})
						return
					}
					
					sort := 0
					if sortRaw, exists := itemMap["sort"]; exists && sortRaw != nil {
						if sortFloat, ok := sortRaw.(float64); ok {
							sort = int(sortFloat)
						}
					}
					
					menuItem := models.MenuItem{
						Date:     date,
						MealType: int(mealType),
						DishID:   uint(dishID),
						Sort:     sort,
					}
					menuItems = append(menuItems, menuItem)
				}
			}
		}
	}

	// 计算周结束日期
	weekEnd := weekStart.AddDate(0, 0, 6)

	// 检查是否已存在该周的菜谱
	var existingMenu models.WeeklyMenu
	if err := database.DB.Where("week_start = ? AND week_end = ?", weekStart, weekEnd).First(&existingMenu).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该周已存在菜谱"})
		return
	}

	// 开始事务
	tx := database.DB.Begin()
	
	// 创建菜谱
	menu := models.WeeklyMenu{
		WeekStart: weekStart,
		WeekEnd:   weekEnd,
		Title:     title,
		Status:    models.MenuStatusDraft,
		CreatedBy: adminID,
	}

	if err := tx.Create(&menu).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建菜谱失败: " + err.Error()})
		return
	}

	// 创建菜谱详情
	for _, item := range menuItems {
		item.MenuID = menu.ID
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建菜谱详情失败: " + err.Error()})
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "data": menu})
}

// GetWeeklyMenus 获取菜谱列表（管理员）
func GetWeeklyMenus(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.WeeklyMenu{}).Preload("Creator")

	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusInt)
	}

	var total int64
	query.Count(&total)

	var menus []models.WeeklyMenu
	if err := query.Offset(offset).Limit(pageSize).Order("week_start DESC").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取菜谱列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": menus,
		"pagination": gin.H{
			"page":      page,
			"page_size": pageSize,
			"total":     total,
		},
	})
}

// GetWeeklyMenuDetail 获取菜谱详情
func GetWeeklyMenuDetail(c *gin.Context) {
	id := c.Param("id")
	var menu models.WeeklyMenu

	if err := database.DB.Preload("MenuItems.Dish.Images").Preload("MenuItems.Dish.Category").Preload("Creator").First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜谱不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": menu})
}

// UpdateWeeklyMenu 更新菜谱（管理员）
func UpdateWeeklyMenu(c *gin.Context) {
	id := c.Param("id")
	var menu models.WeeklyMenu

	if err := database.DB.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜谱不存在"})
		return
	}

	var req CreateWeeklyMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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

	// 更新菜谱基本信息
	menu.Title = req.Title
	menu.WeekStart = req.WeekStart
	menu.WeekEnd = req.WeekStart.AddDate(0, 0, 6)

	if err := tx.Save(&menu).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新菜谱失败"})
		return
	}

	// 删除原有菜谱详情
	if err := tx.Where("menu_id = ?", menu.ID).Delete(&models.MenuItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除原菜谱详情失败"})
		return
	}

	// 创建新的菜谱详情
	for _, item := range req.MenuItems {
		menuItem := models.MenuItem{
			MenuID:   menu.ID,
			Date:     item.Date,
			MealType: item.MealType,
			DishID:   item.DishID,
			Sort:     item.Sort,
		}

		if err := tx.Create(&menuItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建菜谱详情失败"})
			return
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": menu})
}

// PublishWeeklyMenu 发布菜谱（管理员）
func PublishWeeklyMenu(c *gin.Context) {
	id := c.Param("id")
	var menu models.WeeklyMenu

	if err := database.DB.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜谱不存在"})
		return
	}

	menu.Status = models.MenuStatusPublished
	if err := database.DB.Save(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布菜谱失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "发布成功", "data": menu})
}

// DeleteWeeklyMenu 删除菜谱（管理员）
func DeleteWeeklyMenu(c *gin.Context) {
	id := c.Param("id")
	var menu models.WeeklyMenu

	if err := database.DB.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "菜谱不存在"})
		return
	}

	// 开始事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除菜谱详情
	if err := tx.Where("menu_id = ?", menu.ID).Delete(&models.MenuItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除菜谱详情失败"})
		return
	}

	// 删除菜谱
	if err := tx.Delete(&menu).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除菜谱失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetCurrentWeekMenu 获取当前周菜谱（用户端）
func GetCurrentWeekMenu(c *gin.Context) {
	now := time.Now()
	// 格式化为日期字符串（只保留年月日）
	todayStr := now.Format("2006-01-02")

	var menu models.WeeklyMenu
	// 查找包含今天日期的已发布菜谱（今天在 week_start 和 week_end 之间）
	if err := database.DB.Where("DATE(?) >= DATE(week_start) AND DATE(?) <= DATE(week_end) AND status = ?", todayStr, todayStr, models.MenuStatusPublished).
		Preload("MenuItems.Dish.Images").
		Preload("MenuItems.Dish.Category").
		First(&menu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "本周菜谱未发布"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": menu})
}

// GetWeekMenuByDate 根据日期获取菜谱（用户端）
func GetWeekMenuByDate(c *gin.Context) {
	dateStr := c.Param("date")
	// 验证日期格式
	_, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式错误"})
		return
	}

	var menu models.WeeklyMenu
	// 查找包含该日期的已发布菜谱（日期在 week_start 和 week_end 之间）
	if err := database.DB.Where("DATE(?) >= DATE(week_start) AND DATE(?) <= DATE(week_end) AND status = ?", dateStr, dateStr, models.MenuStatusPublished).
		Preload("MenuItems.Dish.Images").
		Preload("MenuItems.Dish.Category").
		First(&menu).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "该周菜谱未发布"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": menu})
}