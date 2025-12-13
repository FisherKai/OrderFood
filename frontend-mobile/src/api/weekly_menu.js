import request from './request'

export const weeklyMenuAPI = {
  // 获取当前周菜谱
  getCurrentWeekMenu() {
    return request.get('/menus/weekly')
  },

  // 根据日期获取菜谱
  getWeekMenuByDate(date) {
    return request.get(`/menus/weekly/${date}`)
  }
}