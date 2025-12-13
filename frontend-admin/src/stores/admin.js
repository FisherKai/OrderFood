import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAdminStore = defineStore('admin', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const adminInfo = ref(JSON.parse(localStorage.getItem('admin_info') || '{}'))

  // 登录
  const login = (tokenValue, info) => {
    token.value = tokenValue
    adminInfo.value = info
    localStorage.setItem('admin_token', tokenValue)
    localStorage.setItem('admin_info', JSON.stringify(info))
  }

  // 登出
  const logout = () => {
    token.value = ''
    adminInfo.value = {}
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_info')
  }

  // 初始化管理员信息
  const initAdmin = () => {
    const savedToken = localStorage.getItem('admin_token')
    const savedInfo = localStorage.getItem('admin_info')
    
    if (savedToken) {
      token.value = savedToken
    }
    if (savedInfo) {
      adminInfo.value = JSON.parse(savedInfo)
    }
  }

  return {
    token,
    adminInfo,
    login,
    logout,
    initAdmin
  }
})