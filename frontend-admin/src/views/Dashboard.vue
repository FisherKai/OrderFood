<template>
  <div class="admin-container" :class="{ 'mobile': isMobile }">
    <!-- 统计卡片 -->
    <el-row :gutter="isMobile ? 10 : 20" class="stats-row">
      <el-col :xs="12" :sm="8" :md="6" :lg="6">
        <el-card class="stat-card" :body-style="{ padding: isMobile ? '12px' : '20px' }">
          <div class="stat-icon primary">
            <el-icon><Food /></el-icon>
          </div>
          <div class="stat-number">{{ stats.dishCount }}</div>
          <div class="stat-label">菜品总数</div>
        </el-card>
      </el-col>
      
      <el-col :xs="12" :sm="8" :md="6" :lg="6">
        <el-card class="stat-card" :body-style="{ padding: isMobile ? '12px' : '20px' }">
          <div class="stat-icon success">
            <el-icon><Document /></el-icon>
          </div>
          <div class="stat-number">{{ stats.orderCount }}</div>
          <div class="stat-label">订单总数</div>
        </el-card>
      </el-col>
      
      <el-col :xs="12" :sm="8" :md="6" :lg="6">
        <el-card class="stat-card" :body-style="{ padding: isMobile ? '12px' : '20px' }">
          <div class="stat-icon warning">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-number">{{ stats.userCount }}</div>
          <div class="stat-label">用户总数</div>
        </el-card>
      </el-col>
      
      <el-col :xs="12" :sm="8" :md="6" :lg="6">
        <el-card class="stat-card" :body-style="{ padding: isMobile ? '12px' : '20px' }">
          <div class="stat-icon danger">
            <el-icon><Star /></el-icon>
          </div>
          <div class="stat-number">{{ stats.reviewCount }}</div>
          <div class="stat-label">评价总数</div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 图表区域 -->
    <el-row :gutter="isMobile ? 10 : 20" class="charts-row">
      <el-col :xs="24" :lg="12">
        <el-card class="admin-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">订单趋势</span>
            </div>
          </template>
          <div class="chart-container" :style="{ height: isMobile ? '220px' : '300px' }">
            <v-chart
              class="chart"
              :option="orderChartOption"
              :loading="chartLoading"
              autoresize
            />
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :lg="12">
        <el-card class="admin-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">菜品分类统计</span>
            </div>
          </template>
          <div class="chart-container" :style="{ height: isMobile ? '220px' : '300px' }">
            <v-chart
              class="chart"
              :option="categoryChartOption"
              :loading="chartLoading"
              autoresize
            />
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 最近活动 -->
    <el-row :gutter="isMobile ? 10 : 20" class="activity-row">
      <el-col :xs="24" :lg="16">
        <el-card class="admin-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">最近订单</span>
              <el-button type="primary" link @click="$router.push('/orders')">
                查看全部
              </el-button>
            </div>
          </template>
          
          <!-- 桌面端表格 -->
          <el-table 
            v-if="!isMobile"
            :data="recentOrders" 
            style="width: 100%" 
            v-loading="tableLoading"
          >
            <el-table-column prop="id" label="订单号" width="100" />
            <el-table-column prop="user_name" label="用户" />
            <el-table-column prop="total_price" label="金额">
              <template #default="{ row }">
                ¥{{ row.total_price }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态">
              <template #default="{ row }">
                <el-tag :type="getOrderStatusType(row.status)" size="small">
                  {{ getOrderStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
          
          <!-- 移动端卡片列表 -->
          <div v-else class="mobile-order-list" v-loading="tableLoading">
            <div 
              v-for="order in recentOrders" 
              :key="order.id" 
              class="order-card"
              @click="$router.push('/orders')"
            >
              <div class="order-header">
                <span class="order-id">#{{ order.id }}</span>
                <el-tag :type="getOrderStatusType(order.status)" size="small">
                  {{ getOrderStatusText(order.status) }}
                </el-tag>
              </div>
              <div class="order-body">
                <div class="order-info">
                  <span class="label">用户</span>
                  <span class="value">{{ order.user_name }}</span>
                </div>
                <div class="order-info">
                  <span class="label">金额</span>
                  <span class="value price">¥{{ order.total_price }}</span>
                </div>
                <div class="order-info">
                  <span class="label">时间</span>
                  <span class="value">{{ formatDate(order.created_at) }}</span>
                </div>
              </div>
            </div>
            
            <el-empty v-if="recentOrders.length === 0" description="暂无订单" />
          </div>
        </el-card>
      </el-col>
      
      <el-col :xs="24" :lg="8">
        <el-card class="admin-card">
          <template #header>
            <div class="card-header">
              <span class="card-title">系统信息</span>
            </div>
          </template>
          <div class="system-info">
            <div class="info-item">
              <span class="label">系统版本</span>
              <span class="value">v1.0.0</span>
            </div>
            <div class="info-item">
              <span class="label">运行时间</span>
              <span class="value">{{ uptime }}</span>
            </div>
            <div class="info-item">
              <span class="label">当前时间</span>
              <span class="value">{{ currentTime }}</span>
            </div>
            <div class="info-item">
              <span class="label">在线管理员</span>
              <span class="value">{{ adminStore.adminInfo.username || '管理员' }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Food, Document, User, Star } from '@element-plus/icons-vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, PieChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components'
import { useAdminStore } from '@/stores/admin'
import { useDevice } from '@/composables/useDevice'
import { dashboardAPI } from '@/api/dashboard'
import { orderAPI } from '@/api/order'
import dayjs from 'dayjs'

// 注册 ECharts 组件
use([
  CanvasRenderer,
  LineChart,
  PieChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
])

const adminStore = useAdminStore()
const { isMobile, screenWidth } = useDevice()

const stats = reactive({
  dishCount: 0,
  orderCount: 0,
  userCount: 0,
  reviewCount: 0
})

const recentOrders = ref([])
const chartLoading = ref(false)
const tableLoading = ref(false)
const currentTime = ref('')
const uptime = ref('1天2小时30分钟')

let timeInterval = null

// 订单趋势图表配置
const orderChartOption = computed(() => ({
  title: {
    text: '最近7天订单趋势',
    textStyle: {
      fontSize: isMobile.value ? 14 : 16
    },
    left: isMobile.value ? 'center' : 'left'
  },
  tooltip: {
    trigger: 'axis'
  },
  grid: {
    left: isMobile.value ? '10%' : '3%',
    right: isMobile.value ? '5%' : '4%',
    bottom: '3%',
    containLabel: true
  },
  xAxis: {
    type: 'category',
    data: orderTrendData.value.dates,
    axisLabel: {
      fontSize: isMobile.value ? 10 : 12,
      rotate: isMobile.value ? 45 : 0
    }
  },
  yAxis: {
    type: 'value',
    axisLabel: {
      fontSize: isMobile.value ? 10 : 12
    }
  },
  series: [{
    data: orderTrendData.value.counts,
    type: 'line',
    smooth: true,
    itemStyle: {
      color: '#409eff'
    },
    areaStyle: {
      color: {
        type: 'linear',
        x: 0,
        y: 0,
        x2: 0,
        y2: 1,
        colorStops: [
          { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
          { offset: 1, color: 'rgba(64, 158, 255, 0.05)' }
        ]
      }
    }
  }]
}))

// 分类统计图表配置
const categoryChartOption = computed(() => ({
  title: {
    text: '菜品分类分布',
    left: 'center',
    textStyle: {
      fontSize: isMobile.value ? 14 : 16
    }
  },
  tooltip: {
    trigger: 'item',
    formatter: '{b}: {c} ({d}%)'
  },
  legend: {
    orient: isMobile.value ? 'horizontal' : 'vertical',
    left: isMobile.value ? 'center' : 'left',
    top: isMobile.value ? 'bottom' : 'middle',
    itemWidth: isMobile.value ? 10 : 14,
    itemHeight: isMobile.value ? 10 : 14,
    textStyle: {
      fontSize: isMobile.value ? 11 : 12
    }
  },
  series: [{
    type: 'pie',
    radius: isMobile.value ? ['30%', '55%'] : ['40%', '60%'],
    center: isMobile.value ? ['50%', '45%'] : ['50%', '50%'],
    data: categoryData.value,
    label: {
      show: !isMobile.value,
      fontSize: 12
    },
    emphasis: {
      itemStyle: {
        shadowBlur: 10,
        shadowOffsetX: 0,
        shadowColor: 'rgba(0, 0, 0, 0.5)'
      }
    }
  }]
}))

// 图表数据
const orderTrendData = ref({
  dates: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
  counts: [12, 19, 15, 27, 32, 25, 18]
})

const categoryData = ref([
  { value: 1048, name: '主食' },
  { value: 735, name: '凉菜' },
  { value: 580, name: '热菜' },
  { value: 484, name: '汤类' },
  { value: 300, name: '饮品' }
])

// 获取订单状态类型
const getOrderStatusType = (status) => {
  const statusMap = {
    1: '',
    2: 'warning',
    3: 'success',
    4: 'danger'
  }
  return statusMap[status] || ''
}

// 获取订单状态文本
const getOrderStatusText = (status) => {
  const statusMap = {
    1: '待处理',
    2: '制作中',
    3: '已完成',
    4: '已取消'
  }
  return statusMap[status] || '未知'
}

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('MM-DD HH:mm')
}

// 更新当前时间
const updateCurrentTime = () => {
  currentTime.value = dayjs().format('YYYY-MM-DD HH:mm:ss')
}

// 加载统计数据
const loadStats = async () => {
  try {
    const response = await dashboardAPI.getStats()
    const data = response.data
    
    stats.dishCount = data.dish_count
    stats.orderCount = data.order_count
    stats.userCount = data.user_count
    stats.reviewCount = data.review_count
  } catch (error) {
    console.error('加载统计数据失败:', error)
    ElMessage.error('加载统计数据失败')
    
    // 使用默认值
    stats.dishCount = 0
    stats.orderCount = 0
    stats.userCount = 0
    stats.reviewCount = 0
  }
}

// 加载最近订单
const loadRecentOrders = async () => {
  try {
    tableLoading.value = true
    
    // 调用真实API获取最近订单
    const response = await orderAPI.getAdminOrders({ page: 1, page_size: 5 })
    
    // 处理订单数据
    const processedOrders = response.data.map(order => ({
      ...order,
      user_name: order.user?.username || order.user?.nickname || '未知用户'
    }))
    
    recentOrders.value = processedOrders
  } catch (error) {
    console.error('加载最近订单失败:', error)
    ElMessage.error('加载最近订单失败')
  } finally {
    tableLoading.value = false
  }
}

// 加载图表数据
const loadChartData = async () => {
  try {
    chartLoading.value = true
    const response = await dashboardAPI.getChartData()
    const data = response.data
    
    // 更新订单趋势图表数据
    if (data.order_trend && data.order_trend.length > 0) {
      const dates = data.order_trend.map(item => {
        const date = new Date(item.date)
        return `${date.getMonth() + 1}/${date.getDate()}`
      })
      const counts = data.order_trend.map(item => item.count)
      
      orderTrendData.value = { dates, counts }
    }
    
    // 更新分类统计图表数据
    if (data.category_stats && data.category_stats.length > 0) {
      categoryData.value = data.category_stats.map(item => ({
        value: item.dish_count,
        name: item.category_name || '未分类'
      }))
    }
  } catch (error) {
    console.error('加载图表数据失败:', error)
  } finally {
    chartLoading.value = false
  }
}

onMounted(() => {
  loadStats()
  loadRecentOrders()
  loadChartData()
  updateCurrentTime()
  
  // 每秒更新时间
  timeInterval = setInterval(updateCurrentTime, 1000)
})

onUnmounted(() => {
  if (timeInterval) {
    clearInterval(timeInterval)
  }
})
</script>

<style scoped lang="scss">
.stats-row {
  margin-bottom: 16px;
}

.charts-row {
  margin-bottom: 16px;
}

.activity-row {
  margin-bottom: 16px;
}

.chart-container {
  height: 300px;
  
  .chart {
    height: 100%;
    width: 100%;
  }
}

.system-info {
  .info-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 14px 0;
    border-bottom: 1px solid #fafafa;
    
    &:last-child {
      border-bottom: none;
    }
    
    .label {
      color: #8c8c8c;
      font-size: 13px;
    }
    
    .value {
      color: #1a1a2e;
      font-size: 13px;
      font-weight: 500;
      font-family: 'SF Mono', 'Monaco', 'Menlo', monospace;
    }
  }
}

// 移动端订单卡片列表
.mobile-order-list {
  .order-card {
    background: #fafbfc;
    border-radius: 10px;
    padding: 14px;
    margin-bottom: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    border: 1px solid transparent;
    
    &:hover {
      background: #f0f7ff;
      border-color: rgba(22, 119, 255, 0.1);
    }
    
    &:active {
      transform: scale(0.99);
    }
    
    &:last-child {
      margin-bottom: 0;
    }
    
    .order-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 10px;
      
      .order-id {
        font-size: 14px;
        font-weight: 600;
        color: #1a1a2e;
      }
    }
    
    .order-body {
      .order-info {
        display: flex;
        justify-content: space-between;
        padding: 3px 0;
        font-size: 13px;
        
        .label {
          color: #8c8c8c;
        }
        
        .value {
          color: #595959;
          
          &.price {
            color: #ff4d4f;
            font-weight: 600;
          }
        }
      }
    }
  }
}

// 移动端样式
.admin-container.mobile {
  padding: 12px;
  
  .stats-row {
    margin-bottom: 10px;
    
    :deep(.el-col) {
      margin-bottom: 8px;
    }
  }
  
  .charts-row {
    margin-bottom: 10px;
    
    :deep(.el-col) {
      margin-bottom: 10px;
    }
  }
  
  .activity-row {
    margin-bottom: 10px;
    
    :deep(.el-col) {
      margin-bottom: 10px;
    }
  }
  
  .system-info {
    .info-item {
      padding: 10px 0;
      
      .label, .value {
        font-size: 12px;
      }
    }
  }
}

@media (max-width: 768px) {
  :deep(.el-col) {
    margin-bottom: 8px;
  }
  
  .chart-container {
    height: 220px;
  }
}
</style>