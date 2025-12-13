<template>
  <div class="admin-container">
    <!-- 分类列表 -->
    <el-card class="admin-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">分类管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增分类
          </el-button>
        </div>
      </template>
      
      <el-table
        :data="categories"
        class="admin-table"
        v-loading="tableLoading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="分类图标" width="100">
          <template #default="{ row }">
            <span class="category-icon">{{ row.icon || '📂' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="分类名称" min-width="150" />
        <el-table-column prop="sort" label="排序" width="100">
          <template #default="{ row }">
            <el-input-number
              v-model="row.sort"
              :min="0"
              size="small"
              @change="handleSortChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="菜品数量" width="100">
          <template #default="{ row }">
            <el-tag type="info">{{ getDishCount(row.id) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="text" @click="handleEdit(row)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button type="text" @click="handleDelete(row)" class="danger">
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 分类表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      @close="handleDialogClose"
    >
      <el-form
        ref="categoryFormRef"
        :model="categoryForm"
        :rules="categoryRules"
        label-width="100px"
        class="admin-form"
      >
        <el-form-item label="分类名称" prop="name">
          <el-input
            v-model="categoryForm.name"
            placeholder="请输入分类名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="分类图标" prop="icon">
          <div class="icon-selector">
            <el-input
              v-model="categoryForm.icon"
              placeholder="请选择或输入图标"
              style="width: 200px"
            />
            <div class="icon-options">
              <span
                v-for="icon in iconOptions"
                :key="icon"
                class="icon-option"
                :class="{ active: categoryForm.icon === icon }"
                @click="categoryForm.icon = icon"
              >
                {{ icon }}
              </span>
            </div>
          </div>
        </el-form-item>
        
        <el-form-item label="排序权重" prop="sort">
          <el-input-number
            v-model="categoryForm.sort"
            :min="0"
            :max="999"
            style="width: 100%"
          />
          <div class="form-tip">数值越大排序越靠前</div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="form-actions">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
            {{ submitLoading ? '保存中...' : '保存' }}
          </el-button>
        </div>
      </template>
    </el-dialog>
    
    <!-- 批量操作 -->
    <div class="batch-actions" v-if="selectedCategories.length > 0">
      <el-card>
        <div class="batch-info">
          <span>已选择 {{ selectedCategories.length }} 个分类</span>
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
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Delete } from '@element-plus/icons-vue'
import { categoryAPI } from '@/api/category'
import { dishAPI } from '@/api/dish'
import dayjs from 'dayjs'

// 分类列表
const categories = ref([])
const tableLoading = ref(false)
const selectedCategories = ref([])
const dishCounts = ref({}) // 存储每个分类的菜品数量

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitLoading = ref(false)
const isEdit = ref(false)
const currentCategoryId = ref(null)

// 分类表单
const categoryFormRef = ref()
const categoryForm = reactive({
  name: '',
  icon: '',
  sort: 0
})

// 表单验证规则
const categoryRules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 2, max: 50, message: '分类名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  sort: [
    { type: 'number', min: 0, max: 999, message: '排序权重在 0 到 999 之间', trigger: 'blur' }
  ]
}

// 图标选项
const iconOptions = [
  '🍚', '🥗', '🍖', '🍲', '🥤', '🍰',
  '🍜', '🍕', '🍔', '🌮', '🥪', '🍱',
  '🍛', '🍝', '🍤', '🦀', '🐟', '🥘',
  '🍳', '🥓', '🧀', '🥯', '🍞', '🥖'
]

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// 获取菜品数量
const getDishCount = (categoryId) => {
  return dishCounts.value[categoryId] || 0
}

// 加载分类列表
const loadCategories = async () => {
  try {
    tableLoading.value = true
    const response = await categoryAPI.getCategories()
    categories.value = response.data || []
    
    // 加载每个分类的菜品数量
    loadDishCounts()
  } catch (error) {
    ElMessage.error('加载分类列表失败')
  } finally {
    tableLoading.value = false
  }
}

// 加载菜品数量
const loadDishCounts = async () => {
  try {
    for (const category of categories.value) {
      const response = await dishAPI.getDishes({ category_id: category.id })
      dishCounts.value[category.id] = response.total || 0
    }
  } catch (error) {
    console.error('加载菜品数量失败:', error)
  }
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedCategories.value = selection
}

// 排序变化
const handleSortChange = async (row) => {
  try {
    await categoryAPI.updateCategory(row.id, { sort: row.sort })
    ElMessage.success('排序更新成功')
    // 重新排序列表
    categories.value.sort((a, b) => b.sort - a.sort)
  } catch (error) {
    ElMessage.error('排序更新失败')
    loadCategories() // 恢复原数据
  }
}

// 新增分类
const handleAdd = () => {
  dialogTitle.value = '新增分类'
  isEdit.value = false
  currentCategoryId.value = null
  resetForm()
  dialogVisible.value = true
}

// 编辑分类
const handleEdit = (row) => {
  dialogTitle.value = '编辑分类'
  isEdit.value = true
  currentCategoryId.value = row.id
  
  // 填充表单数据
  Object.keys(categoryForm).forEach(key => {
    categoryForm[key] = row[key] || (key === 'sort' ? 0 : '')
  })
  
  dialogVisible.value = true
}

// 删除分类
const handleDelete = async (row) => {
  // 检查是否有菜品使用该分类
  const dishCount = getDishCount(row.id)
  if (dishCount > 0) {
    ElMessage.warning(`该分类下还有 ${dishCount} 个菜品，不能删除`)
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除分类"${row.name}"吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await categoryAPI.deleteCategory(row.id)
    ElMessage.success('删除成功')
    loadCategories()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 批量删除
const handleBatchDelete = async () => {
  // 检查选中的分类是否有菜品
  const hasDishesList = []
  selectedCategories.value.forEach(category => {
    const count = getDishCount(category.id)
    if (count > 0) {
      hasDishesList.push(`${category.name}(${count}个菜品)`)
    }
  })
  
  if (hasDishesList.length > 0) {
    ElMessage.warning(`以下分类下还有菜品，不能删除：${hasDishesList.join(', ')}`)
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedCategories.value.length} 个分类吗？`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 批量删除
    const deletePromises = selectedCategories.value.map(category =>
      categoryAPI.deleteCategory(category.id)
    )
    
    await Promise.all(deletePromises)
    ElMessage.success('批量删除成功')
    selectedCategories.value = []
    loadCategories()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

// 重置表单
const resetForm = () => {
  Object.keys(categoryForm).forEach(key => {
    categoryForm[key] = key === 'sort' ? 0 : ''
  })
  
  if (categoryFormRef.value) {
    categoryFormRef.value.clearValidate()
  }
}

// 对话框关闭
const handleDialogClose = () => {
  resetForm()
}

// 提交表单
const handleSubmit = async () => {
  if (!categoryFormRef.value) return
  
  try {
    await categoryFormRef.value.validate()
    submitLoading.value = true
    
    const formData = { ...categoryForm }
    
    if (isEdit.value) {
      await categoryAPI.updateCategory(currentCategoryId.value, formData)
      ElMessage.success('更新成功')
    } else {
      await categoryAPI.createCategory(formData)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    loadCategories()
  } catch (error) {
    if (error.response) {
      ElMessage.error(error.response.data?.error || '操作失败')
    }
  } finally {
    submitLoading.value = false
  }
}

onMounted(() => {
  loadCategories()
})
</script>

<style scoped lang="scss">
.category-icon {
  font-size: 24px;
  display: inline-block;
  width: 32px;
  height: 32px;
  line-height: 32px;
  text-align: center;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.icon-selector {
  .icon-options {
    margin-top: 10px;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    
    .icon-option {
      display: inline-block;
      width: 32px;
      height: 32px;
      line-height: 32px;
      text-align: center;
      font-size: 18px;
      background-color: #f5f7fa;
      border: 1px solid #e4e7ed;
      border-radius: 4px;
      cursor: pointer;
      transition: all 0.2s;
      
      &:hover {
        border-color: #409eff;
        background-color: #ecf5ff;
      }
      
      &.active {
        border-color: #409eff;
        background-color: #409eff;
        color: #fff;
      }
    }
  }
}

.form-tip {
  color: #909399;
  font-size: 12px;
  margin-top: 5px;
}

.batch-actions {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 1000;
  
  .el-card {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    border-radius: 8px;
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

.danger {
  color: #f56c6c !important;
  
  &:hover {
    color: #f78989 !important;
  }
}
</style>