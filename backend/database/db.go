package database

import (
	"fmt"
	"orderfood/config"
	"orderfood/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	cfg := config.AppConfig.Database

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&collation=utf8mb4_unicode_ci",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Category{},
		&models.Dish{},
		&models.DishImage{},
		&models.Order{},
		&models.OrderItem{},
		&models.Review{},
		&models.ReviewImage{},
		&models.Announcement{},
		&models.WeeklyMenu{},
		&models.MenuItem{},
		&models.DishLike{},
		&models.Room{},
		&models.RoomImage{},
		&models.RoomReservation{},
		&models.CanteenReview{},
		&models.CanteenReviewImage{},
		&models.DutyMealSetting{},
		&models.DutyMealOrder{},
		&models.DutyMealOrderItem{},
		&models.ShopCategory{},
		&models.ShopProduct{},
		&models.ShopProductImage{},
		&models.ShopOrder{},
		&models.ShopOrderItem{},
		&models.InventoryRecord{},
	)
}

// SeedReviewData 插入评价测试数据
func SeedReviewData() error {
	// 检查是否已有评价数据
	var count int64
	DB.Model(&models.Review{}).Count(&count)
	if count > 0 {
		return nil // 已有数据，跳过
	}

	// 获取已完成的订单
	var orders []models.Order
	if err := DB.Where("status = ?", 3).Preload("Items").Find(&orders).Error; err != nil {
		return err
	}

	if len(orders) == 0 {
		return nil // 没有已完成订单
	}

	// 为每个已完成订单的菜品创建评价
	reviewContents := []string{
		"非常好吃，强烈推荐！味道正宗，分量足够。",
		"口感不错，下次还会再点。服务态度也很好。",
		"味道很棒，食材新鲜，值得推荐给朋友们！",
		"做得很用心，能感受到厨师的用心。",
		"性价比很高，味道也很不错，会再来的。",
	}

	for _, order := range orders {
		for i, item := range order.Items {
			review := models.Review{
				UserID:  order.UserID,
				DishID:  item.DishID,
				OrderID: order.ID,
				Rating:  4 + (i % 2), // 4或5星
				Content: reviewContents[i%len(reviewContents)],
			}
			if err := DB.Create(&review).Error; err != nil {
				continue // 忽略重复插入错误
			}
		}
	}

	return nil
}
