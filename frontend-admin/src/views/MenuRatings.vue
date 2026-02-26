<template>
  <div class="admin-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="searchForm" class="search-form" inline>
        <el-form-item label="菜谱">
          <el-select
            v-model="searchForm.menu_id"
            placeholder="请选择菜谱"
            clearable
            style="width: 200px"
          >
            <el-option
              v-for="menu in menuList"
              :key="menu.id"
              :label="menu.title"
              :value="menu.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="评分">
          <el-select
            v-model="searchForm.rating"
            placeholder="请选择评分"
            clearable
            style="width: 120px"
          >
            <el-option label="5星" :value="5" />
            <el-option label="4星" :value="4" />
            <el-option label="3星" :value="3" />
            <el-option label="2星" :value="2" />
            <el-option label="1星" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker
            v-model="searchForm.date"
            type="date"
            placeholder="选择日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 160px"
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
          <div class="stat-icon primary">
            <el-icon><ChatDotRound /></el-icon>
          </div>
          <div class="stat-number">{{ overallStats.total_ratings || 0 }}</div>
          <div class="stat-label">总评价数</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon warning">
            <el-icon><Star /></el-icon>
          </div>
          <div class="stat-number">{{ (overallStats.avg_rating || 0).toFixed(1) }}</div>
          <div class="stat-label">平均评分</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon success">
            <el-icon><Trophy /></el-icon>
          </div>
          <div class="stat-number">{{ dishStats.length }}</div>
          <div class="stat-label">被评价菜品</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon danger">
            <el-icon><TrendCharts /></el-icon>
          </div>
          <div class="stat-number">{{ topRatedDish?.dish_name || '-' }}</div>
          <div class="stat-label">最受好评</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 菜品评分排行 -->
    <el-card class="admin-card" v-if="dishStats.length > 0">
      <template #header>
        <div class="card-header">
          <span class="card-title">菜品评分排行</span>
        </div>
      </template>
      
      <el-table :data="dishStats" class="admin-table" max-height="300">
        <el-table-column prop="dish_name" label="菜品名称" min-width="150" />
        <el-table-column label="平均评分" width="150">
          <template #default="{ row }">
            <div class="rating-display">
              <el-rate
                :model-value="row.avg_rating"
                disabled
                show-score
                :colors="['#99A9BF', '#F7BA2A', '#FF9900']"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="total_count" label="评价数" width="100" />
        <el-table-column label="评分分布" min-width="200">
          <template #default="{ row }">
            <div class="rating-distribution">
              <span class="dist-item">5星: {{ row.rating_5 }}</span>
              <span class="dist-item">4星: {{ row.rating_4 }}</span>
              <span class="dist-item">3星: {{ row.rating_3 }}</span>
              <span class="dist-item">2星: {{ row.rating_2 }}</span>
              <span class="dist-item">1星: {{ row.rating_1 }}</span>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 评价列表 -->
    <el-card class="admin-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">评价列表</span>
          <div class="header-actions">
            <el-button @click="loadRatings">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table
        :data="ratings"
        class="admin-table"
        v-loading="tableLoading"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="菜品" min-width="150">
          <template #default="{ row }">
            <span>{{ row.dish?.name || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="用户" width="120">
          <template #default="{ row }">
            <span>{{ row.user?.username || row.user?.nickname || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="评分" width="180">
          <template #default="{ row }">
            <el-rate :model-value="row.rating" disabled />
          </template>
        </el-table-column>
        <el-table-column label="评价内容" min-width="200">
          <template #default="{ row }">
            <span class="comment-text">{{ row.comment || '无评价内容' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="餐次" width="100">
          <template #default="{ row }">
            <el-tag :type="getMealTagType(row.meal_type)">
              {{ getMealTypeName(row.meal_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="评价日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.rating_date) }}
          </template>
        </el-table-column>
        <el-table-column label="提交时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-popconfirm
              title="确定要删除这条评价吗？"
              @confirm="handleDelete(row)"
            >
              <template #reference>
                <el-button type="danger" size="small" plain>
                  删除
                </el-button>
              </template>
            </el-popconfirm>
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Search, Refresh, ChatDotRound, Star, Trophy, TrendCharts
} from '@element-plus/icons-vue'
import { weeklyMenuAPI } from '@/api/weekly_menu'
import dayjs from 'dayjs'

// 搜索表单
const searchForm = reactive({
  menu_id: '',
  rating: '',
  date: ''
})

// 数据
const ratings = ref([])
const menuList = ref([])
const dishStats = ref([])
const overallStats = reactive({
  total_ratings: 0,
  avg_rating: 0
})
const tableLoading = ref(false)

// 分页
const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

// 计算属性
const topRatedDish = computed(() => {
  if (dishStats.value.length === 0) return null
  return dishStats.value[0]
})

// 方法
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD')
}

const formatDateTime = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const getMealTypeName = (mealType) => {
  const names = { 1: '早餐', 2: '午餐', 3: '晚餐', 4: '值班餐' }
  return names[mealType] || '未知'
}

const getMealTagType = (mealType) => {
  const types = { 1: '', 2: 'success', 3: 'warning', 4: 'info' }
  return types[mealType] || ''
}

// 加载菜谱列表
const loadMenuList = async () => {
  try {
    const response = await weeklyMenuAPI.getWeeklyMenus({ page_size: 100 })
    menuList.value = response.data || []
  } catch (error) {
    console.error('加载菜谱列表失败:', error)
  }
}

// 加载统计数据
const loadStats = async () => {
  try {
    const params = {}
    if (searchForm.menu_id) {
      params.menu_id = searchForm.menu_id
    }
    
    const response = await weeklyMenuAPI.getMenuRatingStats(params)
    const data = response.data || {}
    
    dishStats.value = data.dish_stats || []
    overallStats.total_ratings = data.overall?.total_ratings || 0
    overallStats.avg_rating = data.overall?.avg_rating || 0
  } catch (error) {
    console.error('加载统计数据失败:', error)
  }
}

// 加载评价列表
const loadRatings = async () => {
  try {
    tableLoading.value = true
    const params = {
      page: pagination.page,
      page_size: pagination.size,
      ...searchForm
    }
    
    // 移除空值
    Object.keys(params).forEach(key => {
      if (params[key] === '' || params[key] === null) {
        delete params[key]
      }
    })
    
    const response = await weeklyMenuAPI.getMenuRatings(params)
    ratings.value = response.data || []
    pagination.total = response.pagination?.total || 0
  } catch (error) {
    console.error('加载评价列表失败:', error)
    ElMessage.error('加载评价列表失败')
  } finally {
    tableLoading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadRatings()
  loadStats()
}

// 重置
const handleReset = () => {
  searchForm.menu_id = ''
  searchForm.rating = ''
  searchForm.date = ''
  pagination.page = 1
  loadRatings()
  loadStats()
}

// 分页
const handleSizeChange = (size) => {
  pagination.size = size
  pagination.page = 1
  loadRatings()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  loadRatings()
}

// 删除评价
const handleDelete = async (row) => {
  try {
    await weeklyMenuAPI.deleteMenuRating(row.id)
    ElMessage.success('删除成功')
    loadRatings()
    loadStats()
  } catch (error) {
    console.error('删除评价失败:', error)
    ElMessage.error('删除失败')
  }
}

onMounted(() => {
  loadMenuList()
  loadRatings()
  loadStats()
})
</script>

<style scoped lang="scss">
.stats-row {
  margin-bottom: 20px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.rating-display {
  :deep(.el-rate) {
    --el-rate-icon-size: 16px;
  }
}

.rating-distribution {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  
  .dist-item {
    font-size: 12px;
    color: #909399;
    padding: 2px 6px;
    background: #f5f7fa;
    border-radius: 4px;
  }
}

.comment-text {
  color: #606266;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}

// 移动端适配
@media (max-width: 768px) {
  .rating-distribution {
    .dist-item {
      font-size: 10px;
      padding: 1px 4px;
    }
  }
}
</style>
