import request from './request'

export const categoryAPI = {
  // 获取分类列表
  getCategories() {
    return request.get('/categories')
  },

  // 创建分类（管理员）
  createCategory(data) {
    return request.post('/admin/categories', data)
  },

  // 更新分类（管理员）
  updateCategory(id, data) {
    return request.put(`/admin/categories/${id}`, data)
  },

  // 删除分类（管理员）
  deleteCategory(id) {
    return request.delete(`/admin/categories/${id}`)
  }
}