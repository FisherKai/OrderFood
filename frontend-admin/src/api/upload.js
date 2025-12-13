import request from './request'

export const uploadAPI = {
  // 上传图片
  uploadImage(formData) {
    return request.post('/admin/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 获取所有图片列表
  getImages(params = {}) {
    return request.get('/admin/images', { params })
  },

  // 软删除图片
  softDeleteImage(id) {
    return request.delete(`/admin/images/${id}`)
  },

  // 物理删除图片
  physicalDeleteImage(id) {
    return request.delete(`/admin/images/${id}/physical`)
  }
}