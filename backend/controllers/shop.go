package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"orderfood/database"
	"orderfood/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetShopCategories 获取商城分类列表（用户端）
func GetShopCategories(c *gin.Context) {
	var categories []models.ShopCategory
	if err := database.DB.Order("sort DESC, created_at ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": categories,
	})
}

// GetShopProducts 获取商城商品列表（用户端）
func GetShopProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID := c.Query("category_id")
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort", "created_at") // created_at, price, sales_count

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.ShopProduct{}).Where("status = ?", 1)

	// 分类筛选
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 搜索
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// 排序
	switch sortBy {
	case "price":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	case "sales_count":
		query = query.Order("sales_count DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var products []models.ShopProduct
	if err := query.Preload("Category").Preload("Images").
		Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取商品列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      products,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetShopProductDetail 获取商品详情（用户端）
func GetShopProductDetail(c *gin.Context) {
	productID := c.Param("id")
	
	var product models.ShopProduct
	if err := database.DB.Preload("Category").Preload("Images").
		First(&product, productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "商品不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取商品详情失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": product,
	})
}

// CreateShopOrder 创建商城订单（用户端）
func CreateShopOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req struct {
		Items []struct {
			ProductID uint `json:"product_id" binding:"required"`
			Quantity  int  `json:"quantity" binding:"required,min=1"`
		} `json:"items" binding:"required,min=1"`
		DeliveryType    int    `json:"delivery_type" binding:"required"`
		DeliveryAddress string `json:"delivery_address"`
		Remark          string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证配送方式
	if req.DeliveryType == 2 && req.DeliveryAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "配送方式为配送时，配送地址不能为空"})
		return
	}

	// 计算订单总额并验证库存
	var totalAmount float64
	var orderItems []models.ShopOrderItem

	for _, item := range req.Items {
		var product models.ShopProduct
		if err := database.DB.First(&product, item.ProductID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "商品不存在"})
			return
		}

		if product.Status != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "商品已下架"})
			return
		}

		if product.Stock < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("商品 %s 库存不足", product.Name)})
			return
		}

		itemTotal := product.Price * float64(item.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, models.ShopOrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		})
	}

	// 生成订单号
	orderNo := fmt.Sprintf("SP%d%d", time.Now().Unix(), userID)

	// 开始事务
	tx := database.DB.Begin()

	// 创建订单
	order := models.ShopOrder{
		OrderNo:         orderNo,
		UserID:          userID.(uint),
		TotalAmount:     totalAmount,
		Status:          models.ShopOrderStatusPending,
		DeliveryType:    req.DeliveryType,
		DeliveryAddress: req.DeliveryAddress,
		Remark:          req.Remark,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单失败"})
		return
	}

	// 创建订单详情并扣减库存
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		if err := tx.Create(&orderItems[i]).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单详情失败"})
			return
		}

		// 扣减库存
		if err := tx.Model(&models.ShopProduct{}).Where("id = ?", orderItems[i].ProductID).
			Update("stock", gorm.Expr("stock - ?", orderItems[i].Quantity)).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "扣减库存失败"})
			return
		}

		// 记录库存变动
		record := models.InventoryRecord{
			ProductID:    orderItems[i].ProductID,
			Type:         models.InventoryTypeOut,
			Quantity:     orderItems[i].Quantity,
			Reason:       fmt.Sprintf("订单 %s 扣减库存", orderNo),
			OperatorType: "system",
		}
		if err := tx.Create(&record).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "记录库存变动失败"})
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

// GetShopOrders 获取商城订单列表（用户端）
func GetShopOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.ShopOrder{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var orders []models.ShopOrder
	if err := query.Preload("Items").Preload("Items.Product").Preload("Items.Product.Images").
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

// === 管理端接口 ===

// AdminGetShopCategories 获取商城分类列表（管理端）
func AdminGetShopCategories(c *gin.Context) {
	var categories []models.ShopCategory
	if err := database.DB.Order("sort DESC, created_at ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": categories,
	})
}

// AdminCreateShopCategory 创建商城分类（管理端）
func AdminCreateShopCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Icon string `json:"icon"`
		Sort int    `json:"sort"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.ShopCategory{
		Name: req.Name,
		Icon: req.Icon,
		Sort: req.Sort,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建分类失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": category,
		"message": "创建成功",
	})
}

// AdminGetShopProducts 获取商城商品列表（管理端）
func AdminGetShopProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID := c.Query("category_id")
	search := c.Query("search")
	status := c.Query("status")

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.ShopProduct{})

	// 筛选条件
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var products []models.ShopProduct
	if err := query.Preload("Category").Preload("Images").
		Offset(offset).Limit(pageSize).Order("created_at DESC").
		Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取商品列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      products,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// AdminCreateShopProduct 创建商城商品（管理端）
func AdminCreateShopProduct(c *gin.Context) {
	var req struct {
		CategoryID  uint    `json:"category_id" binding:"required"`
		Name        string  `json:"name" binding:"required"`
		Price       float64 `json:"price" binding:"required,min=0"`
		OriginPrice float64 `json:"origin_price"`
		Description string  `json:"description"`
		Stock       int     `json:"stock" binding:"min=0"`
		IsPromotion bool    `json:"is_promotion"`
		Images      []string `json:"images"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证分类是否存在
	var category models.ShopCategory
	if err := database.DB.First(&category, req.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类不存在"})
		return
	}

	// 开始事务
	tx := database.DB.Begin()

	// 创建商品
	product := models.ShopProduct{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Price:       req.Price,
		OriginPrice: req.OriginPrice,
		Description: req.Description,
		Stock:       req.Stock,
		IsPromotion: req.IsPromotion,
		Status:      1,
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建商品失败"})
		return
	}

	// 保存图片
	for i, imageURL := range req.Images {
		if imageURL != "" {
			productImage := models.ShopProductImage{
				ProductID: product.ID,
				ImageURL:  imageURL,
				IsMain:    i == 0, // 第一张图片设为主图
				Sort:      i,
			}
			if err := tx.Create(&productImage).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "保存商品图片失败"})
				return
			}
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": product,
		"message": "创建成功",
	})
}

// AdminGetShopOrders 获取商城订单列表（管理端）
func AdminGetShopOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	search := c.Query("search")

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.ShopOrder{})

	// 筛选条件
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		query = query.Where("order_no LIKE ?", "%"+search+"%")
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var orders []models.ShopOrder
	if err := query.Preload("User").Preload("Items").Preload("Items.Product").
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

// AdminUpdateShopOrderStatus 更新商城订单状态（管理端）
func AdminUpdateShopOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	
	var order models.ShopOrder
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

// AdminGetInventoryRecords 获取库存记录（管理端）
func AdminGetInventoryRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	productID := c.Query("product_id")
	recordType := c.Query("type")

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.InventoryRecord{})

	// 筛选条件
	if productID != "" {
		query = query.Where("product_id = ?", productID)
	}

	if recordType != "" {
		query = query.Where("type = ?", recordType)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var records []models.InventoryRecord
	if err := query.Preload("Product").
		Offset(offset).Limit(pageSize).Order("created_at DESC").
		Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取库存记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      records,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// AdminAdjustInventory 库存调整（管理端）
func AdminAdjustInventory(c *gin.Context) {
	var req struct {
		ProductID uint   `json:"product_id" binding:"required"`
		Type      int    `json:"type" binding:"required"`      // 1-入库，2-出库，3-调整
		Quantity  int    `json:"quantity" binding:"required"`
		Reason    string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证类型
	if req.Type < 1 || req.Type > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的变动类型"})
		return
	}

	// 执行库存调整
	if err := updateStock(req.ProductID, req.Type, req.Quantity, req.Reason, 0, "admin"); err != nil {
		if err == gorm.ErrInvalidData {
			c.JSON(http.StatusBadRequest, gin.H{"error": "库存不足"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "库存调整失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "库存调整成功",
	})
}

// updateStock 更新商品库存并记录
func updateStock(productID uint, changeType int, quantity int, reason string, operatorID uint, operatorType string) error {
	tx := database.DB.Begin()

	// 获取当前库存
	var product models.ShopProduct
	if err := tx.First(&product, productID).Error; err != nil {
		tx.Rollback()
		return err
	}

	beforeStock := product.Stock
	var afterStock int

	switch changeType {
	case models.InventoryTypeIn:
		afterStock = beforeStock + quantity
	case models.InventoryTypeOut:
		afterStock = beforeStock - quantity
		if afterStock < 0 {
			tx.Rollback()
			return gorm.ErrInvalidData // 库存不足
		}
	case models.InventoryTypeAdjust:
		afterStock = quantity
	}

	// 更新库存
	if err := tx.Model(&product).Update("stock", afterStock).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 记录库存变动
	record := models.InventoryRecord{
		ProductID:    productID,
		Type:         changeType,
		Quantity:     quantity,
		BeforeStock:  beforeStock,
		AfterStock:   afterStock,
		Reason:       reason,
		OperatorID:   operatorID,
		OperatorType: operatorType,
	}

	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}