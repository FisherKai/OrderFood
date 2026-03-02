package router

import (
	"orderfood/controllers"
	"orderfood/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 跨域中间件
	r.Use(middleware.CORS())
	
	// 设置字符编码中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// API路由组
	api := r.Group("/api")
	{
		// 用户端路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// 菜品相关
		dishes := api.Group("/dishes")
		{
			dishes.GET("", controllers.GetDishes)
			dishes.GET("/:id", controllers.GetDishDetail)
			dishes.GET("/popular", controllers.GetPopularDishes)  // 热门菜品
		}

		// 分类相关
		api.GET("/categories", controllers.GetCategories)

		// 公告相关
		api.GET("/announcements/active", controllers.GetActiveAnnouncements)
		api.GET("/announcements/:id", controllers.GetAnnouncementDetail)

		// 菜谱相关
		api.GET("/menus/weekly", controllers.GetCurrentWeekMenu)
		api.GET("/menus/weekly/:date", controllers.GetWeekMenuByDate)
		api.GET("/menus/dish/:dish_id/rating", controllers.GetDishAvgRating) // 获取菜品平均评分（公开）

		// 商城相关
		api.GET("/shop/categories", controllers.GetShopCategories)
		api.GET("/shop/products", controllers.GetShopProducts)
		api.GET("/shop/products/:id", controllers.GetShopProductDetail)

		// 需要认证的用户路由
		user := api.Group("")
		user.Use(middleware.AuthMiddleware())
		{
			// 购物车
			user.POST("/cart/add", controllers.AddToCart)

			// 订单
			user.POST("/orders", controllers.CreateOrder)
			user.GET("/orders", controllers.GetUserOrders)
			user.POST("/orders/reserve", controllers.CreateReservation)

			// 评价
			user.POST("/reviews", controllers.CreateReview)
			user.GET("/reviews/:dishId", controllers.GetDishReviews)

			// 菜品点赞
			user.POST("/dishes/:id/like", controllers.LikeDish)
			user.DELETE("/dishes/:id/like", controllers.UnlikeDish)
			user.GET("/dishes/liked", controllers.GetUserLikedDishes)

			// 包间预订
			user.GET("/rooms", controllers.GetRooms)
			user.GET("/rooms/:id", controllers.GetRoomDetail)
			user.POST("/rooms/reserve", controllers.CreateRoomReservation)
			user.GET("/rooms/reservations", controllers.GetMyReservations)
			user.PUT("/rooms/reservations/:id", controllers.UpdateReservation)
			user.DELETE("/rooms/reservations/:id", controllers.CancelReservation)

			// 食堂评价
			user.POST("/canteen/reviews", controllers.CreateCanteenReview)
			user.GET("/canteen/reviews", controllers.GetCanteenReviews)
			user.GET("/canteen/stats", controllers.GetCanteenStats)

			// 值班餐
			user.GET("/duty/meals", controllers.GetDutyMealMenu)
			user.POST("/duty/orders", controllers.CreateDutyMealOrder)
			user.GET("/duty/orders", controllers.GetDutyMealOrders)
			user.GET("/duty/balance", controllers.GetSubsidyBalance)

			// 商城
			user.POST("/shop/orders", controllers.CreateShopOrder)
			user.GET("/shop/orders", controllers.GetShopOrders)

			// 菜谱评价
			user.POST("/menus/ratings", controllers.CreateMenuRating)
			user.GET("/menus/ratings/my", controllers.GetMyMenuRatings)
		}

		// 管理端路由
		admin := api.Group("/admin")
		{
			admin.POST("/login", controllers.AdminLogin)

			// 需要管理员认证的路由
			adminAuth := admin.Group("")
			adminAuth.Use(middleware.AdminAuthMiddleware())
			{
				// 仪表盘统计
				adminAuth.GET("/dashboard/stats", controllers.GetDashboardStats)
				adminAuth.GET("/dashboard/charts", controllers.GetDashboardChartData)
				adminAuth.GET("/orders/stats", controllers.GetOrderStatusStats)

				// 菜品管理
				adminAuth.POST("/dishes", controllers.CreateDish)
				adminAuth.PUT("/dishes/:id", controllers.UpdateDish)
				adminAuth.DELETE("/dishes/:id", controllers.DeleteDish)
				adminAuth.GET("/dishes/:id/images", controllers.GetDishImages)
				adminAuth.PUT("/dishes/:id/images/restore", controllers.RestoreDishImage)

				// 分类管理
				adminAuth.POST("/categories", controllers.CreateCategory)
				adminAuth.PUT("/categories/:id", controllers.UpdateCategory)
				adminAuth.DELETE("/categories/:id", controllers.DeleteCategory)

				// 图片管理
				adminAuth.POST("/upload", controllers.UploadImage)
				adminAuth.DELETE("/images/:id", controllers.SoftDeleteImage)
				adminAuth.GET("/images", controllers.GetAllImages)
				adminAuth.DELETE("/images/:id/physical", controllers.PhysicalDeleteImage)

				// 订单管理
				adminAuth.GET("/orders", controllers.GetAllOrders)
				adminAuth.PUT("/orders/:id/status", controllers.UpdateOrderStatus)

				// 评价管理
				adminAuth.GET("/reviews", controllers.GetAllReviews)

				// 公告管理
				adminAuth.POST("/announcements", controllers.CreateAnnouncement)
				adminAuth.PUT("/announcements/:id", controllers.UpdateAnnouncement)
				adminAuth.DELETE("/announcements/:id", controllers.DeleteAnnouncement)
				adminAuth.GET("/announcements", controllers.GetAllAnnouncements)
				adminAuth.GET("/announcements/:id", controllers.GetAnnouncementDetail)
				adminAuth.PUT("/announcements/:id/status", controllers.UpdateAnnouncementStatus)

				// 用户管理
				adminAuth.GET("/users", controllers.GetUsers)
				adminAuth.GET("/users/stats", controllers.GetUserStats)
				adminAuth.GET("/users/:id", controllers.GetUserDetail)
				adminAuth.PUT("/users/:id/status", controllers.UpdateUserStatus)
				adminAuth.PUT("/users/batch/status", controllers.BatchUpdateUserStatus)
				adminAuth.DELETE("/users/:id", controllers.DeleteUser)

				// 一周菜谱管理
				adminAuth.POST("/menus/weekly", controllers.CreateWeeklyMenu)
				adminAuth.PUT("/menus/weekly/:id", controllers.UpdateWeeklyMenu)
				adminAuth.DELETE("/menus/weekly/:id", controllers.DeleteWeeklyMenu)
				adminAuth.GET("/menus/weekly", controllers.GetWeeklyMenus)
				adminAuth.GET("/menus/weekly/:id", controllers.GetWeeklyMenuDetail)
				adminAuth.PUT("/menus/weekly/:id/publish", controllers.PublishWeeklyMenu)
				adminAuth.PUT("/menus/weekly/:id/cycle", controllers.SetCycleMenu)

				// 菜品点赞管理
				adminAuth.GET("/dishes/likes", controllers.GetDishLikeStats)
				adminAuth.GET("/dishes/likes/ranking", controllers.GetDishLikeRanking)

				// 包间管理
				adminAuth.GET("/rooms", controllers.AdminGetRooms)
				adminAuth.POST("/rooms", controllers.AdminCreateRoom)
				adminAuth.PUT("/rooms/:id", controllers.AdminUpdateRoom)
				adminAuth.DELETE("/rooms/:id", controllers.AdminDeleteRoom)
				adminAuth.GET("/rooms/reservations", controllers.AdminGetReservations)
				adminAuth.PUT("/rooms/reservations/:id/status", controllers.AdminUpdateReservationStatus)
				adminAuth.GET("/rooms/calendar", controllers.AdminGetRoomCalendar)

				// 食堂评价管理
				adminAuth.GET("/canteen/reviews", controllers.AdminGetCanteenReviews)
				adminAuth.POST("/canteen/reviews/:id/reply", controllers.AdminReplyCanteenReview)
				adminAuth.GET("/canteen/stats", controllers.AdminGetCanteenStats)
				adminAuth.DELETE("/canteen/reviews/:id", controllers.AdminDeleteCanteenReview)

				// 值班餐管理
				adminAuth.GET("/duty/settings", controllers.AdminGetDutyMealSettings)
				adminAuth.POST("/duty/settings", controllers.AdminCreateDutyMealSetting)
				adminAuth.PUT("/duty/settings/:id", controllers.AdminUpdateDutyMealSetting)
				adminAuth.DELETE("/duty/settings/:id", controllers.AdminDeleteDutyMealSetting)
				adminAuth.GET("/duty/orders", controllers.AdminGetDutyMealOrders)
				adminAuth.PUT("/duty/orders/:id/status", controllers.AdminUpdateDutyMealOrderStatus)
				adminAuth.GET("/duty/stats", controllers.AdminGetDutyMealStats)

				// 商城管理
				adminAuth.GET("/shop/categories", controllers.AdminGetShopCategories)
				adminAuth.POST("/shop/categories", controllers.AdminCreateShopCategory)
				adminAuth.GET("/shop/products", controllers.AdminGetShopProducts)
				adminAuth.POST("/shop/products", controllers.AdminCreateShopProduct)
				adminAuth.GET("/shop/orders", controllers.AdminGetShopOrders)
				adminAuth.PUT("/shop/orders/:id/status", controllers.AdminUpdateShopOrderStatus)
				adminAuth.GET("/shop/inventory", controllers.AdminGetInventoryRecords)
				adminAuth.POST("/shop/inventory/adjust", controllers.AdminAdjustInventory)

				// 菜谱评价管理
				adminAuth.GET("/menus/ratings", controllers.GetMenuItemRatings)
				adminAuth.GET("/menus/ratings/stats", controllers.GetMenuRatingStats)
				adminAuth.DELETE("/menus/ratings/:id", controllers.DeleteMenuRating)
			}
		}
	}

	return r
}
