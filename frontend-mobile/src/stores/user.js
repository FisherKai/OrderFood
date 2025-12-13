import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login, register } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))

  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUserInfo = (info) => {
    userInfo.value = info
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  const loginUser = async (username, password) => {
    const res = await login({ username, password })
    setToken(res.token)
    setUserInfo(res.user)
    return res
  }

  const registerUser = async (data) => {
    return await register(data)
  }

  const logout = () => {
    token.value = ''
    userInfo.value = {}
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  return {
    token,
    userInfo,
    setToken,
    setUserInfo,
    loginUser,
    registerUser,
    logout
  }
})
