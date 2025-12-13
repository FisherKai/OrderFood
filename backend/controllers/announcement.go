package controllers

import (
	"net/http"
	"orderfood/database"
	"orderfood/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetActiveAnnouncements 获取当前有效的公告
func GetActiveAnnouncements(c *gin.Context) {
	var announcements []models.Announcement
	now := time.Now()

	if err := database.DB.Where("status = ? AND start_time <= ? AND end_time >= ?", 1, now, now).
		Order("sort DESC, created_at DESC").
		Limit(10).
		Find(&announcements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": announcements})
}

// GetAnnouncementDetail 获取公告详情
func GetAnnouncementDetail(c *gin.Context) {
	id := c.Param("id")
	var announcement models.Announcement

	if err := database.DB.First(&announcement, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "公告不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": announcement})
}

// CreateAnnouncement 创建公告（管理员）
func CreateAnnouncement(c *gin.Context) {
	var announcement models.Announcement
	if err := c.ShouldBindJSON(&announcement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证时间
	if announcement.EndTime.Before(announcement.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间不能早于开始时间"})
		return
	}

	if err := database.DB.Create(&announcement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建公告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "data": announcement})
}

// UpdateAnnouncement 更新公告（管理员）
func UpdateAnnouncement(c *gin.Context) {
	id := c.Param("id")
	var announcement models.Announcement

	if err := database.DB.First(&announcement, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "公告不存在"})
		return
	}

	var req struct {
		Title     string    `json:"title" binding:"required"`
		Content   string    `json:"content" binding:"required"`
		Type      int       `json:"type" binding:"required"`
		StartTime time.Time `json:"start_time" binding:"required"`
		EndTime   time.Time `json:"end_time" binding:"required"`
		Status    int       `json:"status"`
		Sort      int       `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.EndTime.Before(req.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间不能早于开始时间"})
		return
	}

	// 更新字段
	announcement.Title = req.Title
	announcement.Content = req.Content
	announcement.Type = req.Type
	announcement.StartTime = req.StartTime
	announcement.EndTime = req.EndTime
	announcement.Status = req.Status
	announcement.Sort = req.Sort

	if err := database.DB.Save(&announcement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新公告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": announcement})
}

// DeleteAnnouncement 删除公告（管理员）
func DeleteAnnouncement(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Announcement{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetAllAnnouncements 获取所有公告（管理员）
func GetAllAnnouncements(c *gin.Context) {
	var announcements []models.Announcement

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	offset := (page - 1) * pageSize

	query := database.DB.Order("sort DESC, created_at DESC")

	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// 类型筛选
	if announcementType := c.Query("type"); announcementType != "" {
		query = query.Where("type = ?", announcementType)
	}

	var total int64
	query.Model(&models.Announcement{}).Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&announcements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取公告失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": announcements,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// UpdateAnnouncementStatus 更新公告状态（管理员）
func UpdateAnnouncementStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status *int `json:"status" binding:"required,min=0,max=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查公告是否存在
	var announcement models.Announcement
	if err := database.DB.First(&announcement, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "公告不存在"})
		return
	}

	if err := database.DB.Model(&announcement).Update("status", *req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}
