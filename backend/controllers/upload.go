package controllers

import (
	"fmt"
	"net/http"
	"orderfood/config"
	"orderfood/database"
	"orderfood/models"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadImage 上传图片
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择文件"})
		return
	}

	// 检查文件大小
	if file.Size > config.AppConfig.Storage.ImageMaxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小超过限制"})
		return
	}

	// 检查文件格式
	ext := strings.ToLower(filepath.Ext(file.Filename))
	ext = strings.TrimPrefix(ext, ".")
	allowed := false
	for _, format := range config.AppConfig.Storage.AllowedFormats {
		if ext == format {
			allowed = true
			break
		}
	}
	if !allowed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支持的文件格式"})
		return
	}

	// 生成文件名
	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), generateRandomString(8), filepath.Ext(file.Filename))
	savePath := filepath.Join(config.AppConfig.Storage.LocalPath, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	imageURL := "/uploads/" + filename

	c.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
		"url":     imageURL,
	})
}

// SoftDeleteImage 软删除图片
func SoftDeleteImage(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()

	if err := database.DB.Model(&models.DishImage{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": &now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetAllImages 获取所有图片列表
func GetAllImages(c *gin.Context) {
	var images []models.DishImage

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	offset := (page - 1) * pageSize

	query := database.DB

	// 筛选未删除的图片
	if c.Query("include_deleted") != "true" {
		query = query.Where("is_deleted = ?", false)
	}

	var total int64
	query.Model(&models.DishImage{}).Count(&total)

	if err := query.Offset(offset).Limit(pageSize).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取图片失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": images,
		"pagination": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// PhysicalDeleteImage 物理删除图片
func PhysicalDeleteImage(c *gin.Context) {
	id := c.Param("id")

	// 这里应该同时删除文件系统中的文件
	// 为简化，暂时只删除数据库记录
	if err := database.DB.Unscoped().Delete(&models.DishImage{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
