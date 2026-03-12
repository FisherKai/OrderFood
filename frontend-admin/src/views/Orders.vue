<template>
  <div class="admin-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="searchForm" class="search-form" inline>
        <el-form-item label="订单号">
          <el-input
            v-model="searchForm.id"
            placeholder="请输入订单号"
            clearable
            style="width: 150px"
          />
        </el-form-item>
        <el-form-item label="订单类型">
          <el-select
            v-model="searchForm.order_type"
            placeholder="请选择类型"
            clearable
            style="width: 120px"
          >
            <el-option label="普通订单" :value="1" />
            <el-option label="预约订单" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            clearable
            style="width: 120px"
          >
            <el-option label="待处理" :value="1" />
            <el-option label="制作中" :value="2" />
            <el-option label="已完成" :value="3" />
            <el-option label="已取消" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item class="search-actions">
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon warning">
            <el-icon><Clock /></el-icon>
          </div>
          <div class="stat-number">{{ orderStats.pending }}</div>
          <div class="stat-label">待处理</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon primary">
            <el-icon><Loading /></el-icon>
          </div>
          <div class="stat-number">{{ orderStats.processing }}</div>
          <div class="stat-label">制作中</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon success">
            <el-icon><Check /></el-icon>
          </div>
          <div class="stat-number">{{ orderStats.completed }}</div>
          <div class="stat-label">已完成</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon danger">
            <el-icon><Close /></el-icon>
          </div>
          <div class="stat-number">{{ orderStats.cancelled }}</div>
          <div class="stat-label">已取消</div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 订单列表 -->
    <el-card class="admin-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">订单列表</span>
          <div class="header-actions">
            <el-button @click="loadOrders">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button type="primary" @click="exportOrders">
              <el-icon><Download /></el-icon>
              导出
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table
        :data="orders"
        class="admin-table"
        v-loading="tableLoading"
        @selection-change="handleSelectionChange"
        @expand-change="handleExpandChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column type="expand" width="30">
          <template #default="{ row }">
            <div class="order-detail">
              <h4>订单详情</h4>
              <el-table :data="row.items" size="small">
                <el-table-column prop="dish_name" label="菜品名称" />
                <el-table-column prop="quantity" label="数量" width="80" />
                <el-table-column prop="price" label="单价" width="100">
                  <template #default="{ row: item }">
                    ¥{{ item.price }}
                  </template>
                </el-table-column>
                <el-table-column label="小计" width="100">
                  <template #default="{ row: item }">
                    ¥{{ (item.price * item.quantity).toFixed(2) }}
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="id" label="订单号" width="100" />
        <el-table-column label="订单类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.order_type === 2 ? 'warning' : ''">
              {{ getOrderTypeText(row.order_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user_name" label="用户" width="120" />
        <el-table-column prop="total_price" label="订单金额" width="120">
          <template #default="{ row }">
            <span class="price">¥{{ row.total_price }}</span>
          </template>
        </el-table-column>
        <el-table-column label="预约时间" width="140">
          <template #default="{ row }">
            <span v-if="row.reserve_time">
              {{ formatDate(row.reserve_time) }}
            </span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="people_count" label="用餐人数" width="100">
          <template #default="{ row }">
            {{ row.people_count || 1 }}人
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getOrderStatusType(row.status)" class="status-tag">
              {{ getOrderStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="下单时间" width="140">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" :width="isMobile ? 180 : 260">
          <template #default="{ row }">
            <div class="action-buttons">
              <!-- 待处理 -> 制作中 -->
              <el-button
                v-if="row.status === 1"
                type="primary"
                size="small"
                @click="handleStatusChange(row, 2)"
              >
                <el-icon v-if="!isMobile"><Loading /></el-icon>
                {{ isMobile ? '接单' : '接单制作' }}
              </el-button>
              
              <!-- 制作中 -> 已完成 -->
              <el-button
                v-if="row.status === 2"
                type="success"
                size="small"
                @click="handleStatusChange(row, 3)"
              >
                <el-icon v-if="!isMobile"><Check /></el-icon>
                {{ isMobile ? '完成' : '完成订单' }}
              </el-button>
              
              <!-- 取消订单（待处理/制作中都可取消） -->
              <el-button
                v-if="row.status < 3"
                type="danger"
                size="small"
                plain
                @click="handleStatusChange(row, 4)"
              >
                <el-icon v-if="!isMobile"><Close /></el-icon>
                取消
              </el-button>
              
              <!-- 查看详情 -->
              <el-button
                type="info"
                size="small"
                plain
                @click="handleViewDetail(row)"
              >
                <el-icon v-if="!isMobile"><View /></el-icon>
                详情
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    
    <!-- 订单详情对话框 -->
    <el-dialog
      v-model="detailVisible"
      title="订单详情"
      width="700px"
    >
      <div v-if="currentOrder" class="order-detail-dialog">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="订单号">{{ currentOrder.id }}</el-descriptions-item>
          <el-descriptions-item label="订单类型">
            <el-tag :type="currentOrder.order_type === 2 ? 'warning' : ''">
              {{ getOrderTypeText(currentOrder.order_type) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="用户">{{ currentOrder.user_name }}</el-descriptions-item>
          <el-descriptions-item label="用餐人数">{{ currentOrder.people_count || 1 }}人</el-descriptions-item>
          <el-descriptions-item label="订单金额">
            <span class="price">¥{{ currentOrder.total_price }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="订单状态">
            <el-tag :type="getOrderStatusType(currentOrder.status)">
              {{ getOrderStatusText(currentOrder.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="预约时间" v-if="currentOrder.reserve_time">
            {{ formatDate(currentOrder.reserve_time) }}
          </el-descriptions-item>
          <el-descriptions-item label="下单时间">
            {{ formatDate(currentOrder.created_at) }}
          </el-descriptions-item>
        </el-descriptions>
        
        <h4 style="margin: 20px 0 10px 0;">订单商品</h4>
        <el-table :data="currentOrder.items" border>
          <el-table-column prop="dish_name" label="菜品名称" />
          <el-table-column prop="quantity" label="数量" width="80" />
          <el-table-column prop="price" label="单价" width="100">
            <template #default="{ row: item }">
              ¥{{ item.price }}
            </template>
          </el-table-column>
          <el-table-column label="小计" width="100">
            <template #default="{ row: item }">
              ¥{{ (item.price * item.quantity).toFixed(2) }}
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>
    
    <!-- 批量操作 -->
    <div class="batch-actions" v-if="selectedOrders.length > 0">
      <el-card>
        <div class="batch-info">
          <span>已选择 {{ selectedOrders.length }} 个订单</span>
          <div class="actions">
            <el-button type="warning" @click="handleBatchStatus(2)">
              批量接单
            </el-button>
            <el-button type="success" @click="handleBatchStatus(3)">
              批量完成
            </el-button>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search, Refresh, Download, Clock, Loading, Check, Close,
  View
} from '@element-plus/icons-vue'
import { orderAPI } from '@/api/order'
import { useDevice } from '@/composables/useDevice'
import dayjs from 'dayjs'

// 设备检测
const { isMobile } = useDevice()

// 搜索表单
const searchForm = reactive({
  id: '',
  order_type: '',
  status: '',
  dateRange: []
})

// 订单列表
const orders = ref([])
const tableLoading = ref(false)
const selectedOrders = ref([])

// 订单统计
const orderStats = reactive({
  pending: 0,
  processing: 0,
  completed: 0,
  cancelled: 0
})

// 分页
const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

// 订单详情对话框
const detailVisible = ref(false)
const currentOrder = ref(null)

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('MM-DD HH:mm')
}

// 获取订单类型文本
const getOrderTypeText = (type) => {
  const typeMap = {
    1: '普通订单',
    2: '预约订单'
  }
  return typeMap[type] || '未知'
}

// 获取订单状态类型
const getOrderStatusType = (status) => {
  const statusMap = {
    1: 'warning',
    2: '',
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

// 加载订单统计
// 加载订单统计
const loadOrderStats = async () => {
  try {
    const response = await orderAPI.getOrderStatusStats()
    const data = response.data
    
    orderStats.pending = data.pending
    orderStats.processing = data.processing
    orderStats.completed = data.completed
    orderStats.cancelled = data.cancelled
  } catch (error) {
    console.error('加载订单统计失败:', error)
    // 使用默认值
    orderStats.pending = 0
    orderStats.processing = 0
    orderStats.completed = 0
    orderStats.cancelled = 0
  }
}

// 加载订单列表
const loadOrders = async () => {
  try {
    tableLoading.value = true
    const params = {
      page: pagination.page,
      page_size: pagination.size,
      ...searchForm
    }
    
    // 处理日期范围
    if (searchForm.dateRange && searchForm.dateRange.length === 2) {
      params.start_date = searchForm.dateRange[0]
      params.end_date = searchForm.dateRange[1]
    }
    delete params.dateRange
    
    // 调用真实API
    const response = await orderAPI.getAdminOrders(params)
    
    // 处理订单数据，添加用户名和菜品名称
    const processedOrders = response.data.map(order => ({
      ...order,
      user_name: order.user?.username || order.user?.nickname || '未知用户',
      items: order.items?.map(item => ({
        ...item,
        dish_name: item.dish?.name || '未知菜品'
      })) || []
    }))
    
    orders.value = processedOrders
    pagination.total = response.pagination?.total || 0
  } catch (error) {
    console.error('加载订单列表失败:', error)
    ElMessage.error('加载订单列表失败')
  } finally {
    tableLoading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadOrders()
}

// 重置搜索
const handleReset = () => {
  Object.keys(searchForm).forEach(key => {
    if (key === 'dateRange') {
      searchForm[key] = []
    } else {
      searchForm[key] = ''
    }
  })
  pagination.page = 1
  loadOrders()
}

// 分页变化
const handleSizeChange = (size) => {
  pagination.size = size
  pagination.page = 1
  loadOrders()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  loadOrders()
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedOrders.value = selection
}

// 展开变化
const handleExpandChange = (row, expandedRows) => {
  // 可以在这里加载订单详情
}

// 更新订单状态
const handleStatusChange = async (row, newStatus) => {
  try {
    await ElMessageBox.confirm(
      `确定要将订单 ${row.id} 的状态更改为"${getOrderStatusText(newStatus)}"吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await orderAPI.updateOrderStatus(row.id, { status: newStatus })
    ElMessage.success('状态更新成功')
    
    // 更新本地数据
    row.status = newStatus
    loadOrderStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('状态更新失败')
    }
  }
}

// 查看订单详情
const handleViewDetail = (row) => {
  currentOrder.value = row
  detailVisible.value = true
}

// 批量更新状态
const handleBatchStatus = async (newStatus) => {
  const validOrders = selectedOrders.value.filter(order => {
    if (newStatus === 2) return order.status === 1
    if (newStatus === 3) return order.status === 2
    return false
  })
  
  if (validOrders.length === 0) {
    ElMessage.warning('没有可操作的订单')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要批量更新 ${validOrders.length} 个订单的状态吗？`,
      '确认批量操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 批量更新
    const updatePromises = validOrders.map(order =>
      orderAPI.updateOrderStatus(order.id, { status: newStatus })
    )
    
    await Promise.all(updatePromises)
    ElMessage.success('批量更新成功')
    
    // 更新本地数据
    validOrders.forEach(order => {
      order.status = newStatus
    })
    
    selectedOrders.value = []
    loadOrderStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量更新失败')
    }
  }
}

// 导出订单
const exportOrders = () => {
  ElMessage.info('导出功能开发中...')
}

onMounted(() => {
  loadOrderStats()
  loadOrders()
})
</script>

<style scoped lang="scss">
.stats-row {
  margin-bottom: 16px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.order-detail {
  padding: 20px;
  background-color: #fafbfc;
  border-radius: 8px;
  
  h4 {
    margin: 0 0 15px 0;
    color: #1a1a2e;
  }
}

.order-detail-dialog {
  .price {
    color: #d48806;
    font-weight: 600;
    font-size: 16px;
  }
}

.price {
  color: #d48806;
  font-weight: 600;
}

.text-muted {
  color: #8c8c8c;
}

.pagination-container {
  margin-top: 16px;
  text-align: right;
}

.batch-actions {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 1000;
  
  .el-card {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
    border-radius: 10px;
  }
  
  .batch-info {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 20px;
    
    .actions {
      margin-left: 20px;
      display: flex;
      gap: 10px;
    }
  }
}

:deep(.el-table__expand-icon) {
  .el-icon {
    font-size: 12px;
  }
}

// 操作按钮样式
.action-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  align-items: center;
  
  .el-button {
    margin: 0 !important;
    
    @media (max-width: 768px) {
      padding: 4px 8px;
      font-size: 12px;
      
      .el-icon {
        display: none;
      }
    }
  }
}

// 表格容器允许溢出
:deep(.el-card__body) {
  overflow: visible !important;
}
</style>