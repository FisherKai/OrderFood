package main

import (
	"fmt"
	"log"
	"orderfood/config"
	"orderfood/database"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg := config.AppConfig.Database

	// 连接到 MySQL 服务器（不指定数据库）
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local&collation=utf8mb4_unicode_ci",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 删除数据库
	if err := db.Exec("DROP DATABASE IF EXISTS orderfood_db").Error; err != nil {
		log.Fatalf("删除数据库失败: %v", err)
	}

	// 创建数据库
	if err := db.Exec("CREATE DATABASE orderfood_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci").Error; err != nil {
		log.Fatalf("创建数据库失败: %v", err)
	}

	fmt.Println("数据库清理完成！")

	// 现在连接到新创建的数据库并初始化
	if err := database.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 自动迁移数据库表
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 插入默认数据
	if err := insertDefaultData(); err != nil {
		log.Fatalf("插入默认数据失败: %v", err)
	}

	fmt.Println("数据库重新初始化完成！")
}

func insertDefaultData() error {
	// 插入默认管理员
	adminSQL := `INSERT INTO admins (username, password, role) VALUES ('admin', '$2a$10$qFswEzYMgMER9ffdUixyduVZvDP61nnXzHIZ2VNM8Ea9TUCuQ28NG', 'admin')`
	if err := database.DB.Exec(adminSQL).Error; err != nil {
		return fmt.Errorf("插入管理员失败: %v", err)
	}

	// 插入默认分类
	categories := []map[string]interface{}{
		{"name": "主食", "icon": "🍚", "sort": 10},
		{"name": "凉菜", "icon": "🥗", "sort": 9},
		{"name": "热菜", "icon": "🍖", "sort": 8},
		{"name": "汤类", "icon": "🍲", "sort": 7},
		{"name": "饮品", "icon": "🥤", "sort": 6},
		{"name": "甜点", "icon": "🍰", "sort": 5},
	}

	for _, cat := range categories {
		categorySQL := `INSERT INTO categories (name, icon, sort) VALUES (?, ?, ?)`
		if err := database.DB.Exec(categorySQL, cat["name"], cat["icon"], cat["sort"]).Error; err != nil {
			return fmt.Errorf("插入分类失败: %v", err)
		}
	}

	// 插入示例公告
	announcementSQL := `INSERT INTO announcements (title, content, type, start_time, end_time, status, sort) VALUES (?, ?, ?, NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), 1, 10)`
	if err := database.DB.Exec(announcementSQL, "欢迎使用点餐系统", "感谢您使用我们的点餐系统，祝您用餐愉快！", 1).Error; err != nil {
		return fmt.Errorf("插入公告失败: %v", err)
	}

	return nil
}