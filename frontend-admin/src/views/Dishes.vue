<template>
  <div class="admin-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="searchForm" class="search-form" inline>
        <el-form-item label="菜品名称">
          <el-input
            v-model="searchForm.name"
            placeholder="请输入菜品名称"
            clearable
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="分类">
          <el-select
            v-model="searchForm.category_id"
            placeholder="请选择分类"
            clearable
            style="width: 150px"
          >
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            clearable
            style="width: 120px"
          >
            <el-option label="上架" :value="1" />
            <el-option label="下架" :value="0" />
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
    
    <!-- 菜品列表 -->
    <el-card class="admin-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">菜品列表</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增菜品
          </el-button>
        </div>
      </template>
      
      <el-table
        :data="dishes"
        class="admin-table"
        v-loading="tableLoading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="菜品图片" width="100">
          <template #default="{ row }">
            <div class="image-preview" v-if="row.images && row.images.length > 0">
              <img
                :src="getImageUrl(row.images[0].image_url)"
                class="preview-image"
                @click="previewImage(row.images)"
              />
            </div>
            <span v-else class="no-image">暂无图片</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="菜品名称" min-width="150" />
        <el-table-column label="分类" width="100">
          <template #default="{ row }">
            {{ getCategoryName(row.category_id) }}
          </template>
        </el-table-column>
        <el-table-column prop="price" label="价格" width="100">
          <template #default="{ row }">
            ¥{{ row.price }}
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="80" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="1"
              :inactive-value="0"
              active-text="上架"
              inactive-text="下架"
              @change="handleStatusChange(row)"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="text" @click="handleEdit(row)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button type="text" @click="handleImageManage(row)">
              <el-icon><Picture /></el-icon>
              图片
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
    
    <!-- 菜品表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form
        ref="dishFormRef"
        :model="dishForm"
        :rules="dishRules"
        label-width="100px"
        class="admin-form"
      >
        <el-form-item label="菜品名称" prop="name">
          <el-input v-model="dishForm.name" placeholder="请输入菜品名称" />
        </el-form-item>
        
        <el-form-item label="菜品分类" prop="category_id">
          <el-select v-model="dishForm.category_id" placeholder="请选择分类" style="width: 100%">
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="菜品价格" prop="price">
          <el-input-number
            v-model="dishForm.price"
            :precision="2"
            :step="0.1"
            :min="0"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="库存数量" prop="stock">
          <el-input-number
            v-model="dishForm.stock"
            :min="0"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="菜品描述" prop="description">
          <el-input
            v-model="dishForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入菜品描述"
          />
        </el-form-item>
        
        <el-form-item label="菜品图片">
          <el-upload
            class="image-uploader"
            :action="uploadUrl"
            :headers="uploadHeaders"
            :show-file-list="false"
            :on-success="handleUploadSuccess"
            :before-upload="beforeUpload"
            accept="image/*"
            multiple
          >
            <el-button type="primary">
              <el-icon><Upload /></el-icon>
              上传图片
            </el-button>
          </el-upload>
          
          <div class="image-list" v-if="dishForm.images && dishForm.images.length > 0">
            <div
              v-for="(image, index) in dishForm.images"
              :key="`${image.image_url}-${index}`"
              class="image-item"
            >
              <img :src="getImageUrl(image.image_url)" class="preview-image" />
              <div class="image-actions">
                <el-button
                  type="text"
                  size="small"
                  @click="setMainImage(index)"
                  :class="{ active: image.is_main }"
                >
                  {{ image.is_main ? '主图' : '设为主图' }}
                </el-button>
                <el-button
                  type="text"
                  size="small"
                  @click="removeImage(index)"
                  class="danger"
                >
                  删除
                </el-button>
              </div>
            </div>
          </div>
        </el-form-item>
        
        <el-form-item label="菜品状态" prop="status">
          <el-radio-group v-model="dishForm.status">
            <el-radio :label="1">上架</el-radio>
            <el-radio :label="0">下架</el-radio>
          </el-radio-group>
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
    
    <!-- 图片预览对话框 -->
    <el-dialog v-model="previewVisible" title="图片预览" width="800px">
      <div class="image-preview-container">
        <img :src="currentPreviewImage" style="width: 100%" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search, Refresh, Plus, Edit, Delete, Picture, Upload
} from '@element-plus/icons-vue'
import { dishAPI } from '@/api/dish'
import { categoryAPI } from '@/api/category'
import { uploadAPI } from '@/api/upload'
import { useAdminStore } from '@/stores/admin'
import dayjs from 'dayjs'

const adminStore = useAdminStore()

// 搜索表单
const searchForm = reactive({
  name: '',
  category_id: '',
  status: ''
})

// 菜品列表
const dishes = ref([])
const categories = ref([])
const tableLoading = ref(false)
const selectedDishes = ref([])

// 分页
const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('')
const submitLoading = ref(false)
const isEdit = ref(false)
const currentDishId = ref(null)

// 图片预览
const previewVisible = ref(false)
const currentPreviewImage = ref('')

// 菜品表单
const dishFormRef = ref()
const dishForm = reactive({
  name: '',
  category_id: '',
  price: 0,
  stock: 0,
  description: '',
  status: 1,
  images: []
})

// 表单验证规则
const dishRules = {
  name: [
    { required: true, message: '请输入菜品名称', trigger: 'blur' }
  ],
  category_id: [
    { required: true, message: '请选择菜品分类', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入菜品价格', trigger: 'blur' },
    { type: 'number', min: 0, message: '价格不能小于0', trigger: 'blur' }
  ]
}

// 上传配置
const uploadUrl = '/api/admin/upload'
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${adminStore.token}`
}))

// 获取图片URL
const getImageUrl = (url) => {
  if (url.startsWith('http')) return url
  return `http://localhost:8080${url}`
}

// 获取分类名称
const getCategoryName = (categoryId) => {
  const category = categories.value.find(c => c.id === categoryId)
  return category ? category.name : '未知分类'
}

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// 加载分类列表
const loadCategories = async () => {
  try {
    const response = await categoryAPI.getCategories()
    categories.value = response.data || []
  } catch (error) {
    ElMessage.error('加载分类失败')
  }
}

// 加载菜品列表
const loadDishes = async () => {
  try {
    tableLoading.value = true
    const params = {
      page: pagination.page,
      size: pagination.size,
      ...searchForm
    }
    
    const response = await dishAPI.getDishes(params)
    dishes.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    ElMessage.error('加载菜品列表失败')
  } finally {
    tableLoading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadDishes()
}

// 重置搜索
const handleReset = () => {
  Object.keys(searchForm).forEach(key => {
    searchForm[key] = ''
  })
  pagination.page = 1
  loadDishes()
}

// 分页变化
const handleSizeChange = (size) => {
  pagination.size = size
  pagination.page = 1
  loadDishes()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  loadDishes()
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedDishes.value = selection
}

// 状态变化
const handleStatusChange = async (row) => {
  try {
    await dishAPI.updateDish(row.id, { status: row.status })
    ElMessage.success('状态更新成功')
  } catch (error) {
    // 恢复原状态
    row.status = row.status === 1 ? 0 : 1
    ElMessage.error('状态更新失败')
  }
}

// 新增菜品
const handleAdd = () => {
  dialogTitle.value = '新增菜品'
  isEdit.value = false
  currentDishId.value = null
  resetForm()
  dialogVisible.value = true
}

// 编辑菜品
const handleEdit = (row) => {
  dialogTitle.value = '编辑菜品'
  isEdit.value = true
  currentDishId.value = row.id
  
  // 填充表单数据
  Object.keys(dishForm).forEach(key => {
    if (key === 'images') {
      dishForm[key] = row.images || []
    } else {
      dishForm[key] = row[key]
    }
  })
  
  dialogVisible.value = true
}

// 删除菜品
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除菜品"${row.name}"吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await dishAPI.deleteDish(row.id)
    ElMessage.success('删除成功')
    loadDishes()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 图片管理
const handleImageManage = (row) => {
  // 跳转到图片管理页面或打开图片管理对话框
  ElMessage.info('图片管理功能开发中...')
}

// 预览图片
const previewImage = (images) => {
  if (images && images.length > 0) {
    currentPreviewImage.value = getImageUrl(images[0].image_url)
    previewVisible.value = true
  }
}

// 上传前验证
const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt5M = file.size / 1024 / 1024 < 5
  
  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB!')
    return false
  }
  return true
}

// 上传成功
const handleUploadSuccess = (response) => {
  if (response.url) {
    dishForm.images.push({
      image_url: response.url,
      is_main: dishForm.images.length === 0 // 第一张图片设为主图
    })
    ElMessage.success('图片上传成功')
  } else {
    ElMessage.error('图片上传失败')
  }
}

// 设置主图
const setMainImage = (index) => {
  dishForm.images.forEach((img, i) => {
    img.is_main = i === index
  })
}

// 移除图片
const removeImage = (index) => {
  if (index < 0 || index >= dishForm.images.length) return
  
  const wasMainImage = dishForm.images[index]?.is_main
  
  // 创建新数组而不是直接修改原数组，确保响应式更新
  const newImages = [...dishForm.images]
  newImages.splice(index, 1)
  
  // 如果删除的是主图，设置第一张为主图
  if (wasMainImage && newImages.length > 0) {
    newImages[0].is_main = true
  } else if (newImages.length > 0 && !newImages.some(img => img.is_main)) {
    newImages[0].is_main = true
  }
  
  // 重新赋值数组，触发响应式更新
  dishForm.images = newImages
  
  // 强制触发响应式更新
  nextTick(() => {
    console.log('图片删除后，当前图片数量:', dishForm.images.length)
  })
}

// 重置表单
const resetForm = () => {
  Object.keys(dishForm).forEach(key => {
    if (key === 'images') {
      dishForm[key] = []
    } else if (key === 'status') {
      dishForm[key] = 1
    } else if (key === 'price' || key === 'stock') {
      dishForm[key] = 0
    } else {
      dishForm[key] = ''
    }
  })
  
  if (dishFormRef.value) {
    dishFormRef.value.clearValidate()
  }
}

// 对话框关闭
const handleDialogClose = () => {
  resetForm()
}

// 提交表单
const handleSubmit = async () => {
  if (!dishFormRef.value) return
  
  try {
    await dishFormRef.value.validate()
    submitLoading.value = true
    
    const formData = { ...dishForm }
    
    if (isEdit.value) {
      await dishAPI.updateDish(currentDishId.value, formData)
      ElMessage.success('更新成功')
    } else {
      await dishAPI.createDish(formData)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    loadDishes()
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
  loadDishes()
})
</script>

<style scoped lang="scss">
.pagination-container {
  margin-top: 16px;
  text-align: right;
}

.no-image {
  color: #bfbfbf;
  font-size: 12px;
}

.image-uploader {
  margin-bottom: 15px;
}

.image-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  
  .image-item {
    position: relative;
    border: 1px solid #f0f0f0;
    border-radius: 8px;
    overflow: hidden;
    
    .preview-image {
      width: 80px;
      height: 80px;
      object-fit: cover;
      display: block;
    }
    
    .image-actions {
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      background: rgba(0, 0, 0, 0.6);
      display: flex;
      justify-content: space-between;
      padding: 4px;
      
      .el-button {
        font-size: 10px;
        padding: 2px 4px;
        color: #fff;
        
        &.active {
          color: #1677ff;
        }
        
        &.danger {
          color: #ff7875;
        }
      }
    }
  }
}

.image-preview-container {
  text-align: center;
}

.danger {
  color: #ff4d4f !important;
  
  &:hover {
    color: #ff7875 !important;
  }
}
</style>