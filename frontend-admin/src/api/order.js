import request from './request'

export const orderAPI = {
  // 获取订单列表（管理员）
  getAdminOrders(params = {}) {
    return request.get('/admin/orders', { params })
  },

  // 获取订单列表（用户）
  getOrders(params = {}) {
    return request.get('/orders', { params })
  },

  // 更新订单状态（管理员）
  updateOrderStatus(id, data) {
    return request.put(`/admin/orders/${id}/status`, data)
  },

  // 创建订单（用户）
  createOrder(data) {
    return request.post('/orders', data)
  },

  // 获取订单状态统计（管理员）
  getOrderStatusStats() {
    return request.get('/admin/orders/stats')
  }
}