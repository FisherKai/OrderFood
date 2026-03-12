<template>
  <div class="admin-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="searchForm" class="search-form" inline>
        <el-form-item label="状态">
          <el-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            clearable
            style="width: 120px"
          >
            <el-option label="草稿" :value="0" />
            <el-option label="已发布" :value="1" />
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
    
    <!-- 菜谱列表 -->
    <el-card class="admin-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">一周菜谱管理</span>
          <div class="header-actions">
            <el-button type="primary" @click="handleCreate">
              <el-icon><Plus /></el-icon>
              新增菜谱
            </el-button>
            <el-button @click="loadMenus">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table
        :data="menus"
        class="admin-table"
        v-loading="tableLoading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="菜谱标题" min-width="200" />
        <el-table-column prop="week_start" label="周开始日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.week_start) }}
          </template>
        </el-table-column>
        <el-table-column prop="week_end" label="周结束日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.week_end) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '已发布' : '草稿' }}
            </el-tag>
            <el-tag v-if="row.is_cycle" type="warning" style="margin-left: 4px">
              循环
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="creator" label="创建人" width="120">
          <template #default="{ row }">
            {{ row.creator?.username || '未知' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button type="text" @click="handleView(row)">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button type="text" @click="handleEdit(row)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button 
              v-if="row.status === 0" 
              type="text" 
              @click="handlePublish(row)"
              class="success"
            >
              <el-icon><Check /></el-icon>
              发布
            </el-button>
            <el-button 
              v-if="row.status === 1"
              type="text" 
              @click="handleSetCycle(row)"
              :class="row.is_cycle ? 'warning' : 'cycle'"
            >
              <el-icon><RefreshRight /></el-icon>
              {{ row.is_cycle ? '取消循环' : '设为循环' }}
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
    
    <!-- 菜谱编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑菜谱' : '新增菜谱'"
      width="80%"
      :close-on-click-modal="false"
      @close="handleDialogClose"
    >
      <el-form
        ref="menuFormRef"
        :model="menuForm"
        :rules="menuRules"
        label-width="100px"
        class="menu-form"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="菜谱标题" prop="title">
              <el-input
                v-model="menuForm.title"
                placeholder="请输入菜谱标题"
                maxlength="100"
                show-word-limit
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="周开始日期" prop="week_start">
              <el-date-picker
                v-model="menuForm.week_start"
                type="date"
                placeholder="选择周开始日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DDTHH:mm:ssZ"
                :picker-options="{
                  disabledDate: (time) => time.getDay() !== 1
                }"
                @change="handleWeekStartChange"
              />
            </el-form-item>
          </el-col>
        </el-row>
        
        <!-- 菜谱详情 -->
        <el-form-item label="菜谱安排" prop="menu_items">
          <div class="menu-calendar">
            <div class="week-header">
              <div class="day-header" v-for="day in weekDays" :key="day.value">
                {{ day.label }}
                <span class="date">{{ getDateByDay(day.value) }}</span>
              </div>
            </div>
            
            <div class="meal-rows">
              <div class="meal-row" v-for="meal in mealTypes" :key="meal.value">
                <div class="meal-label">{{ meal.label }}</div>
                <div class="meal-cells">
                  <div 
                    class="meal-cell" 
                    v-for="day in weekDays" 
                    :key="`${meal.value}-${day.value}`"
                  >
                    <div class="dish-list">
                      <div 
                        class="dish-item"
                        v-for="(item, index) in getMenuItems(day.value, meal.value)"
                        :key="index"
                      >
                        <span class="dish-name">{{ item.dish_name }}</span>
                        <el-button 
                          type="text" 
                          size="small" 
                          @click="removeDish(day.value, meal.value, index)"
                          class="remove-btn"
                        >
                          <el-icon><Close /></el-icon>
                        </el-button>
                      </div>
                    </div>
                    <el-button 
                      type="text" 
                      size="small" 
                      @click="addDish(day.value, meal.value)"
                      class="add-dish-btn"
                    >
                      <el-icon><Plus /></el-icon>
                      添加菜品
                    </el-button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
            {{ isEdit ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>
    
    <!-- 选择菜品对话框 -->
    <el-dialog
      v-model="dishSelectVisible"
      title="选择菜品"
      width="60%"
    >
      <div class="dish-select-content">
        <el-input
          v-model="dishSearchKeyword"
          placeholder="搜索菜品"
          prefix-icon="Search"
          style="margin-bottom: 20px"
          @input="searchDishes"
        />
        
        <el-table
          :data="availableDishes"
          @selection-change="handleDishSelection"
          max-height="400"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="name" label="菜品名称" />
          <el-table-column prop="category.name" label="分类" width="100" />
          <el-table-column prop="price" label="价格" width="80">
            <template #default="{ row }">
              ¥{{ row.price }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="80">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'">
                {{ row.status === 1 ? '上架' : '下架' }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dishSelectVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmDishSelection">
            确定选择
          </el-button>
        </div>
      </template>
    </el-dialog>
    
    <!-- 查看菜谱对话框 -->
    <el-dialog
      v-model="viewDialogVisible"
      title="查看菜谱详情"
      width="80%"
      :close-on-click-modal="false"
    >
      <div v-if="viewMenuData" class="view-menu-content">
        <!-- 基本信息 -->
        <el-descriptions title="基本信息" :column="2" border>
          <el-descriptions-item label="菜谱标题">{{ viewMenuData.title }}</el-descriptions-item>
          <el-descriptions-item label="周开始日期">{{ formatDate(viewMenuData.week_start) }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="viewMenuData.status === 1 ? 'success' : 'info'">
              {{ viewMenuData.status === 1 ? '已发布' : '草稿' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDateTime(viewMenuData.created_at) }}</el-descriptions-item>
        </el-descriptions>
        
        <!-- 菜谱详情 -->
        <div class="view-menu-calendar" style="margin-top: 20px;">
          <h3>菜谱安排</h3>
          <div class="week-header">
            <div class="day-header" v-for="day in weekDays" :key="day.value">
              {{ day.label }}
              <span class="date">{{ getViewDateByDay(day.value) }}</span>
            </div>
          </div>
          
          <div class="meal-rows">
            <div class="meal-row" v-for="meal in mealTypes" :key="meal.value">
              <div class="meal-label">{{ meal.label }}</div>
              <div class="meal-cells">
                <div 
                  class="meal-cell view-mode" 
                  v-for="day in weekDays" 
                  :key="`${meal.value}-${day.value}`"
                >
                  <div class="dish-list">
                    <div 
                      class="dish-item view-item"
                      v-for="(item, index) in getViewMenuItems(day.value, meal.value)"
                      :key="index"
                    >
                      <span class="dish-name">{{ item.dish?.name || '未知菜品' }}</span>
                      <span class="dish-price">¥{{ item.dish?.price || 0 }}</span>
                    </div>
                    <div v-if="getViewMenuItems(day.value, meal.value).length === 0" class="no-dish">
                      暂无安排
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="viewDialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="handleEditFromView">编辑菜谱</el-button>
        </div>
      </template>
    </el-dialog>
    
    <!-- 批量操作 -->
    <div class="batch-actions" v-if="selectedMenus.length > 0">
      <el-card>
        <div class="batch-info">
          <span>已选择 {{ selectedMenus.length }} 个菜谱</span>
          <div class="actions">
            <el-button type="danger" @click="handleBatchDelete">
              <el-icon><Delete /></el-icon>
              批量删除
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
  Search, Refresh, Plus, View, Edit, Delete, Check, Close, RefreshRight
} from '@element-plus/icons-vue'
import { weeklyMenuAPI } from '@/api/weekly_menu'
import { dishAPI } from '@/api/dish'
import dayjs from 'dayjs'

// 搜索表单
const searchForm = reactive({
  status: ''
})

// 菜谱列表
const menus = ref([])
const tableLoading = ref(false)
const selectedMenus = ref([])

// 分页
const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

// 对话框
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentMenuId = ref(null)
const submitLoading = ref(false)

// 菜品选择对话框
const dishSelectVisible = ref(false)
const viewDialogVisible = ref(false)
const viewMenuData = ref(null)
const currentDayMeal = ref({ day: 0, meal: 1 })
const availableDishes = ref([])
const selectedDishes = ref([])
const dishSearchKeyword = ref('')

// 表单
const menuFormRef = ref()
const menuForm = reactive({
  title: '',
  week_start: '',
  menu_items: []
})

// 表单验证规则
const menuRules = {
  title: [
    { required: true, message: '请输入菜谱标题', trigger: 'blur' }
  ],
  week_start: [
    { required: true, message: '请选择周开始日期', trigger: 'change' }
  ]
}

// 周几和餐次配置
const weekDays = [
  { value: 1, label: '周一' },
  { value: 2, label: '周二' },
  { value: 3, label: '周三' },
  { value: 4, label: '周四' },
  { value: 5, label: '周五' },
  { value: 6, label: '周六' },
  { value: 0, label: '周日' }
]

const mealTypes = [
  { value: 1, label: '早餐' },
  { value: 2, label: '午餐' },
  { value: 3, label: '晚餐' },
  { value: 4, label: '值班餐' }
]

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('MM-DD')
}

const formatDateTime = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

// 根据周几获取具体日期
const getDateByDay = (dayOfWeek) => {
  if (!menuForm.week_start) return ''
  const startDate = dayjs(menuForm.week_start)
  const targetDate = startDate.add(dayOfWeek - 1, 'day')
  return targetDate.format('MM-DD')
}

// 获取指定日期和餐次的菜品
const getMenuItems = (dayOfWeek, mealType) => {
  if (!menuForm.week_start) return []
  
  const startDate = dayjs(menuForm.week_start)
  const targetDate = startDate.add(dayOfWeek - 1, 'day').format('YYYY-MM-DD')
  
  return menuForm.menu_items.filter(item => 
    dayjs(item.date).format('YYYY-MM-DD') === targetDate && 
    item.meal_type === mealType
  )
}

// 处理周开始日期变化
const handleWeekStartChange = (date) => {
  if (date) {
    // 确保选择的是周一
    const selectedDate = dayjs(date)
    if (selectedDate.day() !== 1) {
      ElMessage.warning('请选择周一作为开始日期')
      menuForm.week_start = ''
      return
    }
  }
}

// 添加菜品
const addDish = (dayOfWeek, mealType) => {
  currentDayMeal.value = { day: dayOfWeek, meal: mealType }
  loadAvailableDishes()
  dishSelectVisible.value = true
}

// 移除菜品
const removeDish = (dayOfWeek, mealType, index) => {
  if (!menuForm.week_start) return
  
  const startDate = dayjs(menuForm.week_start)
  const targetDate = startDate.add(dayOfWeek - 1, 'day').format('YYYY-MM-DD')
  
  const itemIndex = menuForm.menu_items.findIndex((item, i) => {
    const itemDate = dayjs(item.date).format('YYYY-MM-DD')
    return itemDate === targetDate && 
           item.meal_type === mealType && 
           i === menuForm.menu_items.filter(mi => 
             dayjs(mi.date).format('YYYY-MM-DD') === targetDate && 
             mi.meal_type === mealType
           ).indexOf(item) + index
  })
  
  if (itemIndex > -1) {
    menuForm.menu_items.splice(itemIndex, 1)
  }
}

// 加载可用菜品
const loadAvailableDishes = async () => {
  try {
    const response = await dishAPI.getDishes({ 
      page: 1, 
      page_size: 100, 
      status: 1 
    })
    availableDishes.value = response.data || []
  } catch (error) {
    ElMessage.error('加载菜品列表失败')
  }
}

// 搜索菜品
const searchDishes = () => {
  // 这里可以实现菜品搜索逻辑
  // 暂时使用前端过滤
}

// 处理菜品选择
const handleDishSelection = (selection) => {
  selectedDishes.value = selection
}

// 确认菜品选择
const confirmDishSelection = () => {
  if (selectedDishes.value.length === 0) {
    ElMessage.warning('请选择至少一个菜品')
    return
  }
  
  if (!menuForm.week_start) {
    ElMessage.warning('请先选择周开始日期')
    return
  }
  
  const startDate = dayjs(menuForm.week_start)
  const targetDate = startDate.add(currentDayMeal.value.day - 1, 'day').toISOString()
  
  selectedDishes.value.forEach(dish => {
    menuForm.menu_items.push({
      date: targetDate,
      meal_type: currentDayMeal.value.meal,
      dish_id: dish.id,
      dish_name: dish.name,
      sort: 0
    })
  })
  
  selectedDishes.value = []
  dishSelectVisible.value = false
}

// 加载菜谱列表
const loadMenus = async () => {
  try {
    tableLoading.value = true
    const params = {
      page: pagination.page,
      page_size: pagination.size,
      ...searchForm
    }
    
    const response = await weeklyMenuAPI.getWeeklyMenus(params)
    menus.value = response.data || []
    pagination.total = response.pagination?.total || 0
  } catch (error) {
    ElMessage.error('加载菜谱列表失败')
  } finally {
    tableLoading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadMenus()
}

// 重置搜索
const handleReset = () => {
  Object.keys(searchForm).forEach(key => {
    searchForm[key] = ''
  })
  pagination.page = 1
  loadMenus()
}

// 分页变化
const handleSizeChange = (size) => {
  pagination.size = size
  pagination.page = 1
  loadMenus()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  loadMenus()
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedMenus.value = selection
}

// 新增菜谱
const handleCreate = () => {
  isEdit.value = false
  currentMenuId.value = null
  resetForm()
  dialogVisible.value = true
}

// 查看菜谱
const handleView = async (row) => {
  try {
    const response = await weeklyMenuAPI.getWeeklyMenuDetail(row.id)
    viewMenuData.value = response.data
    viewDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取菜谱详情失败')
  }
}

// 编辑菜谱
const handleEdit = async (row) => {
  try {
    isEdit.value = true
    currentMenuId.value = row.id
    
    const response = await weeklyMenuAPI.getWeeklyMenuDetail(row.id)
    const menuData = response.data
    
    menuForm.title = menuData.title
    menuForm.week_start = dayjs(menuData.week_start).format('YYYY-MM-DD')
    menuForm.menu_items = menuData.menu_items?.map(item => ({
      ...item,
      dish_name: item.dish?.name || '未知菜品'
    })) || []
    
    dialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取菜谱详情失败')
  }
}

// 发布菜谱
const handlePublish = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要发布菜谱"${row.title}"吗？发布后用户端将可以查看。`,
      '确认发布',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await weeklyMenuAPI.publishWeeklyMenu(row.id)
    ElMessage.success('发布成功')
    loadMenus()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('发布失败')
    }
  }
}

// 设置/取消循环菜谱
const handleSetCycle = async (row) => {
  const isCycle = !row.is_cycle
  const action = isCycle ? '设为循环菜谱' : '取消循环菜谱'
  const tip = isCycle 
    ? `设为循环菜谱后，当某周没有专属菜谱时，将自动使用此菜谱的菜品安排。同一时间只能有一个循环菜谱。` 
    : `取消后，此菜谱将不再作为循环菜谱使用。`
  
  try {
    await ElMessageBox.confirm(
      `${tip}\n\n确定要${action}"${row.title}"吗？`,
      action,
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      }
    )
    
    await weeklyMenuAPI.setCycleMenu(row.id, isCycle)
    ElMessage.success(`${action}成功`)
    loadMenus()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`${action}失败`)
    }
  }
}

// 删除菜谱
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除菜谱"${row.title}"吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await weeklyMenuAPI.deleteWeeklyMenu(row.id)
    ElMessage.success('删除成功')
    loadMenus()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 批量删除
const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedMenus.value.length} 个菜谱吗？此操作不可恢复。`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 批量删除API调用
    for (const menu of selectedMenus.value) {
      await weeklyMenuAPI.deleteWeeklyMenu(menu.id)
    }
    
    ElMessage.success('批量删除成功')
    selectedMenus.value = []
    loadMenus()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

// 重置表单
const resetForm = () => {
  menuForm.title = ''
  menuForm.week_start = ''
  menuForm.menu_items = []
}

// 对话框关闭
const handleDialogClose = () => {
  resetForm()
}

// 提交表单
const handleSubmit = async () => {
  if (!menuFormRef.value) return
  
  try {
    await menuFormRef.value.validate()
    submitLoading.value = true
    
    const formData = {
      title: menuForm.title,
      week_start: menuForm.week_start || dayjs(menuForm.week_start).toISOString(),
      menu_items: menuForm.menu_items.map(item => ({
        date: dayjs(item.date).toISOString(),
        meal_type: item.meal_type,
        dish_id: item.dish_id,
        sort: item.sort || 0
      }))
    }
    
    if (isEdit.value) {
      await weeklyMenuAPI.updateWeeklyMenu(currentMenuId.value, formData)
      ElMessage.success('更新成功')
    } else {
      await weeklyMenuAPI.createWeeklyMenu(formData)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    loadMenus()
  } catch (error) {
    if (error.response) {
      ElMessage.error(error.response.data?.error || '操作失败')
    }
  } finally {
    submitLoading.value = false
  }
}

// 查看模式的辅助函数
const getViewDateByDay = (dayOfWeek) => {
  if (!viewMenuData.value?.week_start) return ''
  const startDate = dayjs(viewMenuData.value.week_start)
  const targetDate = startDate.add(dayOfWeek - 1, 'day')
  return targetDate.format('MM-DD')
}

const getViewMenuItems = (dayOfWeek, mealType) => {
  if (!viewMenuData.value?.week_start || !viewMenuData.value?.menu_items) return []
  
  const startDate = dayjs(viewMenuData.value.week_start)
  const targetDate = startDate.add(dayOfWeek - 1, 'day').format('YYYY-MM-DD')
  
  return viewMenuData.value.menu_items.filter(item => 
    dayjs(item.date).format('YYYY-MM-DD') === targetDate && 
    item.meal_type === mealType
  )
}

// 从查看模式进入编辑
const handleEditFromView = () => {
  viewDialogVisible.value = false
  handleEdit(viewMenuData.value)
}

onMounted(() => {
  loadMenus()
})
</script>

<style scoped lang="scss">
.menu-calendar {
  border: 1px solid #f0f0f0;
  border-radius: 10px;
  overflow: hidden;
  
  .week-header {
    display: flex;
    background-color: #fafbfc;
    border-bottom: 1px solid #f0f0f0;
    
    .day-header {
      flex: 1;
      padding: 12px 8px;
      text-align: center;
      border-right: 1px solid #f0f0f0;
      font-weight: 500;
      font-size: 13px;
      color: #1a1a2e;
      
      &:last-child {
        border-right: none;
      }
      
      .date {
        display: block;
        font-size: 11px;
        color: #8c8c8c;
        margin-top: 4px;
      }
    }
  }
  
  .meal-rows {
    .meal-row {
      display: flex;
      border-bottom: 1px solid #f0f0f0;
      min-height: 120px;
      
      &:last-child {
        border-bottom: none;
      }
      
      .meal-label {
        width: 80px;
        padding: 12px 8px;
        background-color: #fafbfc;
        border-right: 1px solid #f0f0f0;
        display: flex;
        align-items: center;
        justify-content: center;
        font-weight: 500;
        font-size: 13px;
        color: #595959;
      }
      
      .meal-cells {
        flex: 1;
        display: flex;
        
        .meal-cell {
          flex: 1;
          padding: 8px;
          border-right: 1px solid #f0f0f0;
          
          &:last-child {
            border-right: none;
          }
          
          .dish-list {
            margin-bottom: 8px;
            
            .dish-item {
              display: flex;
              align-items: center;
              justify-content: space-between;
              padding: 4px 8px;
              margin-bottom: 4px;
              background-color: #e8f4ff;
              border-radius: 6px;
              font-size: 12px;
              
              .dish-name {
                flex: 1;
                color: #1677ff;
              }
              
              .remove-btn {
                padding: 0;
                margin-left: 4px;
                color: #ff4d4f;
              }
            }
          }
          
          .add-dish-btn {
            width: 100%;
            padding: 8px;
            border: 1px dashed #d9d9d9;
            background-color: #fafbfc;
            color: #8c8c8c;
            font-size: 12px;
            border-radius: 6px;
            
            &:hover {
              border-color: #1677ff;
              color: #1677ff;
              background-color: #f0f7ff;
            }
          }
        }
      }
    }
  }
}

.dish-select-content {
  .el-table {
    .el-table__row {
      cursor: pointer;
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
    }
  }
}

.success {
  color: #52c41a !important;
  
  &:hover {
    color: #73d13d !important;
  }
}

.cycle {
  color: #d48806 !important;
  
  &:hover {
    color: #faad14 !important;
  }
}

.warning {
  color: #8c8c8c !important;
  
  &:hover {
    color: #bfbfbf !important;
  }
}

.danger {
  color: #ff4d4f !important;
  
  &:hover {
    color: #ff7875 !important;
  }
}

// 查看模式样式
.view-menu-content {
  .view-menu-calendar {
    .meal-cell.view-mode {
      background-color: #fafbfc;
      
      .dish-item.view-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 8px 12px;
        margin-bottom: 8px;
        background-color: #fff;
        border: 1px solid #f0f0f0;
        border-radius: 8px;
        
        .dish-name {
          flex: 1;
          font-size: 13px;
          color: #1a1a2e;
        }
        
        .dish-price {
          font-size: 12px;
          color: #ff4d4f;
          font-weight: 500;
        }
      }
      
      .no-dish {
        text-align: center;
        color: #bfbfbf;
        font-size: 12px;
        padding: 20px;
      }
    }
  }
}
</style>