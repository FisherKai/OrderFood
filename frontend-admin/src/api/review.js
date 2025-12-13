import request from './request'

export const reviewAPI = {
  // 获取评价列表（管理员）
  getAdminReviews(params = {}) {
    return request.get('/admin/reviews', { params })
  },

  // 删除评价（管理员）
  deleteReview(id) {
    return request.delete(`/admin/reviews/${id}`)
  },

  // 批量删除评价（管理员）
  batchDeleteReviews(ids) {
    return request.delete('/admin/reviews/batch', { data: { ids } })
  }
}