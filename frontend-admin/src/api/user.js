import request from './request'

export const userAPI = {
  // 获取用户列表（管理员）
  getUsers(params = {}) {
    return request.get('/admin/users', { params })
  },

  // 获取用户统计概览（管理员）
  getStats() {
    return request.get('/admin/users/stats')
  },

  // 获取用户详情（管理员）
  getUserDetail(id) {
    return request.get(`/admin/users/${id}`)
  },

  // 获取单个用户统计（管理员）
  getUserStats(id) {
    return request.get(`/admin/users/${id}/stats`)
  },

  // 更新用户状态（管理员）
  updateUserStatus(id, data) {
    return request.put(`/admin/users/${id}/status`, data)
  },

  // 批量更新用户状态（管理员）
  batchUpdateUserStatus(data) {
    return request.put('/admin/users/batch/status', data)
  },

  // 删除用户（管理员）
  deleteUser(id) {
    return request.delete(`/admin/users/${id}`)
  }
}