package controllers

import (
	"net/http"
	"strconv"
	"time"

	"orderfood/database"
	"orderfood/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetRooms 获取包间列表（用户端）
func GetRooms(c *gin.Context) {
	var rooms []models.Room
	
	// 查询可用的包间
	if err := database.DB.Where("status = ?", models.RoomStatusAvailable).
		Preload("Images").Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取包间列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": rooms,
	})
}

// GetRoomDetail 获取包间详情（用户端）
func GetRoomDetail(c *gin.Context) {
	roomID := c.Param("id")
	
	var room models.Room
	if err := database.DB.Preload("Images").First(&room, roomID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "包间不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取包间详情失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": room,
	})
}

// CreateRoomReservation 创建包间预订（用户端）
func CreateRoomReservation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req struct {
		RoomID      uint      `json:"room_id" binding:"required"`
		ReserveDate time.Time `json:"reserve_date" binding:"required"`
		StartTime   time.Time `json:"start_time" binding:"required"`
		EndTime     time.Time `json:"end_time" binding:"required"`
		GuestCount  int       `json:"guest_count" binding:"required,min=1"`
		Remark      string    `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证时间
	if req.EndTime.Before(req.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间不能早于开始时间"})
		return
	}

	// 验证预订时间（至少提前2小时）
	if req.StartTime.Before(time.Now().Add(2 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "预订时间至少提前2小时"})
		return
	}

	// 检查包间是否存在且可用
	var room models.Room
	if err := database.DB.First(&room, req.RoomID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "包间不存在"})
		return
	}

	if room.Status != models.RoomStatusAvailable {
		c.JSON(http.StatusBadRequest, gin.H{"error": "包间当前不可用"})
		return
	}

	// 检查包间容量
	if req.GuestCount > room.Capacity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用餐人数超过包间容量"})
		return
	}

	// 检查时间冲突
	var conflictCount int64
	database.DB.Model(&models.RoomReservation{}).
		Where("room_id = ? AND status IN (?, ?) AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?))",
			req.RoomID, models.ReservationStatusPending, models.ReservationStatusConfirmed,
			req.StartTime, req.StartTime, req.EndTime, req.EndTime).
		Count(&conflictCount)

	if conflictCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该时间段已被预订"})
		return
	}

	// 计算价格
	duration := req.EndTime.Sub(req.StartTime).Hours()
	totalPrice := duration * room.HourlyPrice

	// 创建预订
	reservation := models.RoomReservation{
		UserID:      userID.(uint),
		RoomID:      req.RoomID,
		ReserveDate: req.ReserveDate,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		GuestCount:  req.GuestCount,
		TotalPrice:  totalPrice,
		Remark:      req.Remark,
		Status:      models.ReservationStatusPending,
	}

	if err := database.DB.Create(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建预订失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": reservation,
		"message": "预订成功",
	})
}

// GetMyReservations 获取我的包间预订（用户端）
func GetMyReservations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var reservations []models.RoomReservation
	if err := database.DB.Where("user_id = ?", userID).
		Preload("Room").Preload("Room.Images").
		Order("created_at DESC").Find(&reservations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取预订列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": reservations,
	})
}

// UpdateReservation 修改包间预订（用户端）
func UpdateReservation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	reservationID := c.Param("id")
	
	var reservation models.RoomReservation
	if err := database.DB.Where("id = ? AND user_id = ?", reservationID, userID).
		First(&reservation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "预订不存在"})
		return
	}

	// 只有待确认状态的预订可以修改
	if reservation.Status != models.ReservationStatusPending {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只有待确认的预订可以修改"})
		return
	}

	var req struct {
		StartTime  time.Time `json:"start_time" binding:"required"`
		EndTime    time.Time `json:"end_time" binding:"required"`
		GuestCount int       `json:"guest_count" binding:"required,min=1"`
		Remark     string    `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证时间
	if req.EndTime.Before(req.StartTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间不能早于开始时间"})
		return
	}

	// 验证预订时间（至少提前2小时）
	if req.StartTime.Before(time.Now().Add(2 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "预订时间至少提前2小时"})
		return
	}

	// 获取包间信息
	var room models.Room
	if err := database.DB.First(&room, reservation.RoomID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取包间信息失败"})
		return
	}

	// 检查包间容量
	if req.GuestCount > room.Capacity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用餐人数超过包间容量"})
		return
	}

	// 检查时间冲突（排除当前预订）
	var conflictCount int64
	database.DB.Model(&models.RoomReservation{}).
		Where("room_id = ? AND id != ? AND status IN (?, ?) AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?))",
			reservation.RoomID, reservation.ID, models.ReservationStatusPending, models.ReservationStatusConfirmed,
			req.StartTime, req.StartTime, req.EndTime, req.EndTime).
		Count(&conflictCount)

	if conflictCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该时间段已被预订"})
		return
	}

	// 重新计算价格
	duration := req.EndTime.Sub(req.StartTime).Hours()
	totalPrice := duration * room.HourlyPrice

	// 更新预订
	reservation.StartTime = req.StartTime
	reservation.EndTime = req.EndTime
	reservation.GuestCount = req.GuestCount
	reservation.TotalPrice = totalPrice
	reservation.Remark = req.Remark

	if err := database.DB.Save(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "修改预订失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": reservation,
		"message": "修改成功",
	})
}

// CancelReservation 取消包间预订（用户端）
func CancelReservation(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	reservationID := c.Param("id")
	
	var reservation models.RoomReservation
	if err := database.DB.Where("id = ? AND user_id = ?", reservationID, userID).
		First(&reservation).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "预订不存在"})
		return
	}

	// 只有待确认和已确认状态的预订可以取消
	if reservation.Status != models.ReservationStatusPending && 
	   reservation.Status != models.ReservationStatusConfirmed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该预订无法取消"})
		return
	}

	// 更新状态为已取消
	reservation.Status = models.ReservationStatusCancelled
	if err := database.DB.Save(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消预订失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "取消成功",
	})
}

// === 管理端接口 ===

// AdminGetRooms 获取包间列表（管理端）
func AdminGetRooms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	status := c.Query("status")

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.Room{})

	// 搜索条件
	if search != "" {
		query = query.Where("room_number LIKE ? OR room_name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var rooms []models.Room
	if err := query.Preload("Images").Offset(offset).Limit(pageSize).
		Order("created_at DESC").Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取包间列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      rooms,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// AdminCreateRoom 新增包间（管理端）
func AdminCreateRoom(c *gin.Context) {
	var req struct {
		RoomNumber  string  `json:"room_number" binding:"required"`
		RoomName    string  `json:"room_name" binding:"required"`
		Capacity    int     `json:"capacity" binding:"required,min=1"`
		Description string  `json:"description"`
		HourlyPrice float64 `json:"hourly_price" binding:"required,min=0"`
		Status      int     `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查包间编号是否已存在
	var existingRoom models.Room
	if err := database.DB.Where("room_number = ?", req.RoomNumber).First(&existingRoom).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "包间编号已存在"})
		return
	}

	if req.Status == 0 {
		req.Status = models.RoomStatusAvailable
	}

	room := models.Room{
		RoomNumber:  req.RoomNumber,
		RoomName:    req.RoomName,
		Capacity:    req.Capacity,
		Description: req.Description,
		HourlyPrice: req.HourlyPrice,
		Status:      req.Status,
	}

	if err := database.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建包间失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": room,
		"message": "创建成功",
	})
}

// AdminUpdateRoom 编辑包间（管理端）
func AdminUpdateRoom(c *gin.Context) {
	roomID := c.Param("id")
	
	var room models.Room
	if err := database.DB.First(&room, roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "包间不存在"})
		return
	}

	var req struct {
		RoomNumber  string  `json:"room_number" binding:"required"`
		RoomName    string  `json:"room_name" binding:"required"`
		Capacity    int     `json:"capacity" binding:"required,min=1"`
		Description string  `json:"description"`
		HourlyPrice float64 `json:"hourly_price" binding:"required,min=0"`
		Status      int     `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查包间编号是否已存在（排除当前包间）
	var existingRoom models.Room
	if err := database.DB.Where("room_number = ? AND id != ?", req.RoomNumber, roomID).
		First(&existingRoom).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "包间编号已存在"})
		return
	}

	// 更新包间信息
	room.RoomNumber = req.RoomNumber
	room.RoomName = req.RoomName
	room.Capacity = req.Capacity
	room.Description = req.Description
	room.HourlyPrice = req.HourlyPrice
	room.Status = req.Status

	if err := database.DB.Save(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新包间失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": room,
		"message": "更新成功",
	})
}

// AdminDeleteRoom 删除包间（管理端）
func AdminDeleteRoom(c *gin.Context) {
	roomID := c.Param("id")
	
	var room models.Room
	if err := database.DB.First(&room, roomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "包间不存在"})
		return
	}

	// 检查是否有未完成的预订
	var activeReservations int64
	database.DB.Model(&models.RoomReservation{}).
		Where("room_id = ? AND status IN (?, ?)", roomID, 
			models.ReservationStatusPending, models.ReservationStatusConfirmed).
		Count(&activeReservations)

	if activeReservations > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该包间有未完成的预订，无法删除"})
		return
	}

	// 删除包间及相关数据
	tx := database.DB.Begin()

	// 删除包间图片
	if err := tx.Where("room_id = ?", roomID).Delete(&models.RoomImage{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除包间图片失败"})
		return
	}

	// 删除包间
	if err := tx.Delete(&room).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除包间失败"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "删除成功",
	})
}

// AdminGetReservations 获取包间预订列表（管理端）
func AdminGetReservations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	roomID := c.Query("room_id")

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.RoomReservation{})

	// 筛选条件
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if roomID != "" {
		query = query.Where("room_id = ?", roomID)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取数据
	var reservations []models.RoomReservation
	if err := query.Preload("User").Preload("Room").Preload("Room.Images").
		Offset(offset).Limit(pageSize).Order("created_at DESC").
		Find(&reservations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取预订列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      reservations,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// AdminUpdateReservationStatus 更新预订状态（管理端）
func AdminUpdateReservationStatus(c *gin.Context) {
	reservationID := c.Param("id")
	
	var reservation models.RoomReservation
	if err := database.DB.First(&reservation, reservationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "预订不存在"})
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
	if req.Status < 1 || req.Status > 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的状态值"})
		return
	}

	reservation.Status = req.Status
	if err := database.DB.Save(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "更新成功",
	})
}

// AdminGetRoomCalendar 获取包间日历视图（管理端）
func AdminGetRoomCalendar(c *gin.Context) {
	roomID := c.Query("room_id")
	date := c.Query("date") // YYYY-MM-DD 格式

	if roomID == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必要参数"})
		return
	}

	// 解析日期
	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日期格式错误"})
		return
	}

	// 获取当天的预订
	startTime := targetDate
	endTime := targetDate.Add(24 * time.Hour)

	var reservations []models.RoomReservation
	if err := database.DB.Where("room_id = ? AND start_time >= ? AND start_time < ?", 
		roomID, startTime, endTime).
		Preload("User").Order("start_time ASC").Find(&reservations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取日历数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": reservations,
	})
}