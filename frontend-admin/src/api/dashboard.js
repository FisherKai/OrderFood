import request from './request'

export const dashboardAPI = {
  // 获取仪表盘统计数据
  getStats() {
    return request.get('/admin/dashboard/stats')
  },

  // 获取仪表盘图表数据
  getChartData() {
    return request.get('/admin/dashboard/charts')
  }
}