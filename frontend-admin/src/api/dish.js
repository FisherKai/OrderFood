import request from './request'

export const dishAPI = {
  // 获取菜品列表
  getDishes(params = {}) {
    return request.get('/dishes', { params })
  },

  // 获取菜品详情
  getDishDetail(id) {
    return request.get(`/dishes/${id}`)
  },

  // 创建菜品（管理员）
  createDish(data) {
    return request.post('/admin/dishes', data)
  },

  // 更新菜品（管理员）
  updateDish(id, data) {
    return request.put(`/admin/dishes/${id}`, data)
  },

  // 删除菜品（管理员）
  deleteDish(id) {
    return request.delete(`/admin/dishes/${id}`)
  },

  // 获取菜品图片列表（包括历史图片）
  getDishImages(id) {
    return request.get(`/admin/dishes/${id}/images`)
  },

  // 恢复历史图片
  restoreDishImage(id, data) {
    return request.put(`/admin/dishes/${id}/images/restore`, data)
  }
}