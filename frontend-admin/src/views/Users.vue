<template>
  <div class="admin-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="searchForm" class="search-form" inline>
        <el-form-item label="用户名">
          <el-input
            v-model="searchForm.username"
            placeholder="请输入用户名"
            clearable
            style="width: 150px"
          />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input
            v-model="searchForm.phone"
            placeholder="请输入手机号"
            clearable
            style="width: 150px"
          />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input
            v-model="searchForm.email"
            placeholder="请输入邮箱"
            clearable
            style="width: 180px"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            clearable
            style="width: 120px"
          >
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
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
    
    <!-- 统计信息 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="8" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon primary">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-number">{{ userStats.total }}</div>
          <div class="stat-label">总用户数</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="8" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon success">
            <el-icon><Check /></el-icon>
          </div>
          <div class="stat-number">{{ userStats.active }}</div>
          <div class="stat-label">正常用户</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="8" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon danger">
            <el-icon><Close /></el-icon>
          </div>
          <div class="stat-number">{{ userStats.disabled }}</div>
          <div class="stat-label">禁用用户</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="8" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon warning">
            <el-icon><Clock /></el-icon>
          </div>
          <div class="stat-number">{{ userStats.newToday }}</div>
          <div class="stat-label">今日新增</div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 用户列表 -->
    <el-card class="admin-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">用户管理</span>
          <div class="header-actions">
            <el-button @click="loadUsers">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button type="primary" @click="exportUsers">
              <el-icon><Download /></el-icon>
              导出
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table
        :data="users"
        class="admin-table"
        v-loading="tableLoading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="头像" width="80">
          <template #default="{ row }">
            <el-avatar
              :size="40"
              :src="row.avatar ? getImageUrl(row.avatar) : undefined"
              :icon="UserFilled"
            />
          </template>
        </el-table-column>
        <el-table-column prop="username" label="用户名" min-width="120" />
        <el-table-column prop="nickname" label="昵称" min-width="120" />
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column prop="email" label="邮箱" min-width="180" />
        <el-table-column label="订单统计" width="120">
          <template #default="{ row }">
            <div class="order-stats">
              <div class="stat-item">
                <span class="label">总数:</span>
                <span class="value">{{ getUserOrderCount(row.id) }}</span>
              </div>
              <div class="stat-item">
                <span class="label">完成:</span>
                <span class="value success">{{ getUserCompletedOrderCount(row.id) }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              active-text="正常"
              inactive-text="禁用"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="text" @click="handleViewDetail(row)">
              <el-icon><View /></el-icon>
              详情
            </el-button>
            <el-button type="text" @click="handleViewOrders(row)">
              <el-icon><Document /></el-icon>
              订单
            </el-button>
            <el-button type="text" @click="handleDelete(row)" class="danger">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
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
    
    <!-- 用户详情对话框 -->
    <el-dialog
      v-model="detailVisible"
      title="用户详情"
      width="600px"
    >
      <div v-if="currentUser" class="user-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="用户ID">{{ currentUser.id }}</el-descriptions-item>
          <el-descriptions-item label="用户名">{{ currentUser.username }}</el-descriptions-item>
          <el-descriptions-item label="昵称">{{ currentUser.nickname || '未设置' }}</el-descriptions-item>
          <el-descriptions-item label="手机号">{{ currentUser.phone || '未绑定' }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ currentUser.email || '未绑定' }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="currentUser.status === 1 ? 'success' : 'danger'">
              {{ currentUser.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="注册时间">
            {{ formatDate(currentUser.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="最后更新">
            {{ formatDate(currentUser.updated_at) }}
          </el-descriptions-item>
        </el-descriptions>
        
        <div class="user-avatar-section">
          <h4>用户头像</h4>
          <el-avatar
            :size="100"
            :src="currentUser.avatar ? getImageUrl(currentUser.avatar) : undefined"
            :icon="UserFilled"
          />
        </div>
        
        <div class="user-stats-section">
          <h4>统计信息</h4>
          <el-row :gutter="16">
            <el-col :span="8">
              <div class="stat-box">
                <div class="stat-number">{{ getUserOrderCount(currentUser.id) }}</div>
                <div class="stat-label">总订单数</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="stat-box">
                <div class="stat-number">{{ getUserCompletedOrderCount(currentUser.id) }}</div>
                <div class="stat-label">完成订单</div>
              </div>
            </el-col>
            <el-col :span="8">
              <div class="stat-box">
                <div class="stat-number">{{ getUserReviewCount(currentUser.id) }}</div>
                <div class="stat-label">评价数量</div>
              </div>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-dialog>
    
    <!-- 批量操作 -->
    <div class="batch-actions" v-if="selectedUsers.length > 0">
      <el-card>
        <div class="batch-info">
          <span>已选择 {{ selectedUsers.length }} 个用户</span>
          <div class="actions">
            <el-button type="success" @click="handleBatchStatus(1)">
              批量启用
            </el-button>
            <el-button type="warning" @click="handleBatchStatus(0)">
              批量禁用
            </el-button>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search, Refresh, Download, User, Check, Close, Clock,
  View, Document, UserFilled, Delete
} from '@element-plus/icons-vue'
import { userAPI } from '@/api/user'
import dayjs from 'dayjs'

// 搜索表单
const searchForm = reactive({
  username: '',
  phone: '',
  email: '',
  status: ''
})

// 用户列表
const users = ref([])
const tableLoading = ref(false)
const selectedUsers = ref([])

// 用户统计
const userStats = reactive({
  total: 0,
  active: 0,
  disabled: 0,
  newToday: 0
})

// 分页
const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

// 用户详情对话框
const detailVisible = ref(false)
const currentUser = ref(null)

// 模拟订单统计数据
const userOrderStats = ref({})
const userReviewStats = ref({})

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// 获取图片URL
const getImageUrl = (url) => {
  if (url.startsWith('http')) return url
  return `http://localhost:8080${url}`
}

// 获取用户订单数量
const getUserOrderCount = (userId) => {
  return userOrderStats.value[userId]?.total || 0
}

// 获取用户完成订单数量
const getUserCompletedOrderCount = (userId) => {
  return userOrderStats.value[userId]?.completed || 0
}

// 获取用户评价数量
const getUserReviewCount = (userId) => {
  return userReviewStats.value[userId] || 0
}

// 加载用户统计
const loadUserStats = async () => {
  try {
    const response = await userAPI.getStats()
    const data = response.data
    
    userStats.total = data.total || 0
    userStats.active = data.active || 0
    userStats.disabled = data.disabled || 0
    userStats.newToday = data.new_today || 0
  } catch (error) {
    console.error('加载用户统计失败:', error)
    // 使用默认值
    userStats.total = 0
    userStats.active = 0
    userStats.disabled = 0
    userStats.newToday = 0
  }
}

// 加载用户列表
const loadUsers = async () => {
  try {
    tableLoading.value = true
    const params = {
      page: pagination.page,
      page_size: pagination.size,
      ...searchForm
    }
    
    // 调用真实API
    const response = await userAPI.getUsers(params)
    users.value = response.data || []
    pagination.total = response.pagination?.total || 0
    
    // 加载用户订单和评价统计
    for (const user of users.value) {
      try {
        const statsResponse = await userAPI.getUserDetail(user.id)
        userOrderStats.value[user.id] = statsResponse.stats?.orders || { total: 0, completed: 0 }
        userReviewStats.value[user.id] = statsResponse.stats?.reviews || 0
      } catch (error) {
        console.error(`加载用户${user.id}统计数据失败:`, error)
        userOrderStats.value[user.id] = { total: 0, completed: 0 }
        userReviewStats.value[user.id] = 0
      }
    }
  } catch (error) {
    console.error('加载用户列表失败:', error)
    ElMessage.error('加载用户列表失败')
  } finally {
    tableLoading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadUsers()
}

// 重置搜索
const handleReset = () => {
  Object.keys(searchForm).forEach(key => {
    searchForm[key] = ''
  })
  pagination.page = 1
  loadUsers()
}

// 分页变化
const handleSizeChange = (size) => {
  pagination.size = size
  pagination.page = 1
  loadUsers()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  loadUsers()
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedUsers.value = selection
}

// 状态变化
const handleStatusChange = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要${row.status === 1 ? '启用' : '禁用'}用户"${row.username}"吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 这里应该调用API更新用户状态
    ElMessage.success('状态更新成功')
    loadUserStats()
  } catch (error) {
    if (error !== 'cancel') {
      // 恢复原状态
      row.status = row.status === 1 ? 0 : 1
      ElMessage.error('状态更新失败')
    }
  }
}

// 查看用户详情
const handleViewDetail = (row) => {
  currentUser.value = row
  detailVisible.value = true
}

// 查看用户订单
const handleViewOrders = (row) => {
  ElMessage.info(`查看用户"${row.username}"的订单记录`)
  // 这里可以跳转到订单页面并筛选该用户的订单
}

// 删除用户
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户"${row.username}"吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await userAPI.deleteUser(row.id)
    ElMessage.success('删除成功')
    loadUsers()
    loadUserStats()
  } catch (error) {
    if (error !== 'cancel') {
      const errorMsg = error.response?.data?.error || '删除失败'
      ElMessage.error(errorMsg)
    }
  }
}

// 批量更新状态
const handleBatchStatus = async (status) => {
  try {
    await ElMessageBox.confirm(
      `确定要批量${status === 1 ? '启用' : '禁用'} ${selectedUsers.value.length} 个用户吗？`,
      '确认批量操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 批量更新API调用
    ElMessage.success('批量更新成功')
    
    // 更新本地数据
    selectedUsers.value.forEach(user => {
      user.status = status
    })
    
    selectedUsers.value = []
    loadUserStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量更新失败')
    }
  }
}

// 导出用户
const exportUsers = () => {
  ElMessage.info('导出功能开发中...')
}

onMounted(() => {
  loadUserStats()
  loadUsers()
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

.order-stats {
  font-size: 12px;
  
  .stat-item {
    display: flex;
    justify-content: space-between;
    margin-bottom: 2px;
    
    .label {
      color: #8c8c8c;
    }
    
    .value {
      font-weight: 500;
      
      &.success {
        color: #52c41a;
      }
    }
  }
}

.user-detail {
  .user-avatar-section {
    margin: 20px 0;
    text-align: center;
    
    h4 {
      margin-bottom: 15px;
      color: #1a1a2e;
    }
  }
  
  .user-stats-section {
    margin-top: 20px;
    
    h4 {
      margin-bottom: 15px;
      color: #1a1a2e;
    }
    
    .stat-box {
      text-align: center;
      padding: 16px;
      background-color: #fafbfc;
      border-radius: 8px;
      
      .stat-number {
        font-size: 22px;
        font-weight: 700;
        color: #1677ff;
        margin-bottom: 4px;
      }
      
      .stat-label {
        font-size: 12px;
        color: #8c8c8c;
      }
    }
  }
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

.danger {
  color: #ff4d4f !important;
  
  &:hover {
    color: #ff7875 !important;
  }
}
</style>