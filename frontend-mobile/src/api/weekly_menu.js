import request from './request'

export const weeklyMenuAPI = {
  // 获取当前周菜谱
  getCurrentWeekMenu() {
    return request.get('/menus/weekly')
  },

  // 根据日期获取菜谱
  getWeekMenuByDate(date) {
    return request.get(`/menus/weekly/${date}`)
  },

  // 创建菜品评价
  createMenuRating(data) {
    return request.post('/menus/ratings', data)
  },

  // 获取用户的评价
  getMyMenuRatings(params) {
    return request.get('/menus/ratings/my', { params })
  },

  // 获取菜品平均评分
  getDishAvgRating(dishId) {
    return request.get(`/menus/dish/${dishId}/rating`)
  }
}