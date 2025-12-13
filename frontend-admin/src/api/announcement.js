import request from './request'

export const announcementAPI = {
  // 获取公告列表（管理员）
  getAnnouncements(params = {}) {
    return request.get('/admin/announcements', { params })
  },

  // 获取公告详情
  getAnnouncementDetail(id) {
    return request.get(`/admin/announcements/${id}`)
  },

  // 创建公告（管理员）
  createAnnouncement(data) {
    return request.post('/admin/announcements', data)
  },

  // 更新公告（管理员）
  updateAnnouncement(id, data) {
    return request.put(`/admin/announcements/${id}`, data)
  },

  // 删除公告（管理员）
  deleteAnnouncement(id) {
    return request.delete(`/admin/announcements/${id}`)
  },

  // 更新公告状态（管理员）
  updateAnnouncementStatus(id, data) {
    return request.put(`/admin/announcements/${id}/status`, data)
  }
}