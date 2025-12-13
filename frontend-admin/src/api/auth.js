import request from './request'

export const authAPI = {
  // 管理员登录
  login(data) {
    return request.post('/admin/login', data)
  }
}