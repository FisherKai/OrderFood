import request from './request'

export const weeklyMenuAPI = {
  // 创建一周菜谱（管理员）
  createWeeklyMenu(data) {
    return request.post('/admin/menus/weekly', data)
  },

  // 获取菜谱列表（管理员）
  getWeeklyMenus(params = {}) {
    return request.get('/admin/menus/weekly', { params })
  },

  // 获取菜谱详情
  getWeeklyMenuDetail(id) {
    return request.get(`/admin/menus/weekly/${id}`)
  },

  // 更新菜谱（管理员）
  updateWeeklyMenu(id, data) {
    return request.put(`/admin/menus/weekly/${id}`, data)
  },

  // 发布菜谱（管理员）
  publishWeeklyMenu(id) {
    return request.put(`/admin/menus/weekly/${id}/publish`)
  },

  // 删除菜谱（管理员）
  deleteWeeklyMenu(id) {
    return request.delete(`/admin/menus/weekly/${id}`)
  },

  // 获取当前周菜谱（用户端）
  getCurrentWeekMenu() {
    return request.get('/menus/weekly')
  },

  // 根据日期获取菜谱（用户端）
  getWeekMenuByDate(date) {
    return request.get(`/menus/weekly/${date}`)
  }
}