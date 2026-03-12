<template>
  <div class="admin-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="searchForm" class="search-form" inline>
        <el-form-item label="图片类型">
          <el-select
            v-model="searchForm.type"
            placeholder="请选择类型"
            clearable
            style="width: 120px"
          >
            <el-option label="菜品图片" value="dish" />
            <el-option label="评价图片" value="review" />
            <el-option label="用户头像" value="avatar" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="searchForm.is_deleted"
            placeholder="请选择状态"
            clearable
            style="width: 120px"
          >
            <el-option label="正常" :value="false" />
            <el-option label="已删除" :value="true" />
          </el-select>
        </el-form-item>
        <el-form-item label="上传时间">
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
    
    <!-- 统计信息 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon primary">
            <el-icon><Picture /></el-icon>
          </div>
          <div class="stat-number">{{ imageStats.total }}</div>
          <div class="stat-label">总图片数</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon success">
            <el-icon><Check /></el-icon>
          </div>
          <div class="stat-number">{{ imageStats.normal }}</div>
          <div class="stat-label">正常图片</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon danger">
            <el-icon><Delete /></el-icon>
          </div>
          <div class="stat-number">{{ imageStats.deleted }}</div>
          <div class="stat-label">已删除</div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-card class="stat-card">
          <div class="stat-icon warning">
            <el-icon><FolderOpened /></el-icon>
          </div>
          <div class="stat-number">{{ formatFileSize(imageStats.totalSize) }}</div>
          <div class="stat-label">总大小</div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 图片列表 -->
    <el-card class="admin-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">图片管理</span>
          <div class="header-actions">
            <el-upload
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
            <el-button @click="loadImages">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button @click="handleCleanup" type="danger">
              <el-icon><Delete /></el-icon>
              清理垃圾
            </el-button>
          </div>
        </div>
      </template>
      
      <!-- 视图切换 -->
      <div class="view-controls">
        <el-radio-group v-model="viewMode" @change="handleViewModeChange">
          <el-radio-button label="grid">网格视图</el-radio-button>
          <el-radio-button label="list">列表视图</el-radio-button>
        </el-radio-group>
      </div>
      
      <!-- 网格视图 -->
      <div v-if="viewMode === 'grid'" class="image-grid" v-loading="tableLoading">
        <div
          v-for="image in images"
          :key="image.id"
          class="image-card"
          :class="{ deleted: image.is_deleted, selected: selectedImages.includes(image) }"
          @click="handleImageSelect(image)"
        >
          <div class="image-container">
            <img
              :src="getImageUrl(image.url)"
              :alt="image.name"
              @click.stop="previewImage(image)"
              @error="handleImageError"
            />
            <div class="image-overlay">
              <div class="image-actions">
                <el-button type="primary" size="small" @click.stop="previewImage(image)">
                  <el-icon><View /></el-icon>
                </el-button>
                <el-button
                  v-if="!image.is_deleted"
                  type="warning"
                  size="small"
                  @click.stop="handleSoftDelete(image)"
                >
                  <el-icon><Delete /></el-icon>
                </el-button>
                <el-button
                  v-if="image.is_deleted"
                  type="success"
                  size="small"
                  @click.stop="handleRestore(image)"
                >
                  <el-icon><RefreshLeft /></el-icon>
                </el-button>
                <el-button
                  type="danger"
                  size="small"
                  @click.stop="handlePhysicalDelete(image)"
                >
                  <el-icon><Close /></el-icon>
                </el-button>
              </div>
            </div>
          </div>
          <div class="image-info">
            <div class="image-name">{{ image.name }}</div>
            <div class="image-meta">
              <span class="image-size">{{ formatFileSize(image.size) }}</span>
              <el-tag :type="getImageTypeColor(image.type)" size="small">
                {{ getImageTypeText(image.type) }}
              </el-tag>
            </div>
            <div class="image-time">{{ formatDate(image.created_at) }}</div>
          </div>
        </div>
      </div>
      
      <!-- 列表视图 -->
      <el-table
        v-else
        :data="images"
        class="admin-table"
        v-loading="tableLoading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="预览" width="100">
          <template #default="{ row }">
            <img
              :src="getImageUrl(row.url)"
              class="table-image-preview"
              @click="previewImage(row)"
              @error="handleImageError"
            />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="文件名" min-width="200" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getImageTypeColor(row.type)" size="small">
              {{ getImageTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" width="100">
          <template #default="{ row }">
            {{ formatFileSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column label="使用情况" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.usage_count > 0" type="success" size="small">
              被引用 {{ row.usage_count }} 次
            </el-tag>
            <el-tag v-else type="info" size="small">未使用</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_deleted" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_deleted ? 'danger' : 'success'" size="small">
              {{ row.is_deleted ? '已删除' : '正常' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="上传时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="text" @click="previewImage(row)">
              <el-icon><View /></el-icon>
              预览
            </el-button>
            <el-button
              v-if="!row.is_deleted"
              type="text"
              @click="handleSoftDelete(row)"
            >
              <el-icon><Delete /></el-icon>
              删除
            </el-button>
            <el-button
              v-if="row.is_deleted"
              type="text"
              @click="handleRestore(row)"
              class="success"
            >
              <el-icon><RefreshLeft /></el-icon>
              恢复
            </el-button>
            <el-button
              type="text"
              @click="handlePhysicalDelete(row)"
              class="danger"
            >
              <el-icon><Close /></el-icon>
              彻底删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :page-sizes="[20, 50, 100, 200]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
    
    <!-- 图片预览对话框 -->
    <el-dialog v-model="previewVisible" title="图片预览" width="800px">
      <div v-if="currentPreviewImage" class="image-preview-container">
        <img :src="getImageUrl(currentPreviewImage.url)" class="preview-image" />
        <div class="preview-info">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="文件名">{{ currentPreviewImage.name }}</el-descriptions-item>
            <el-descriptions-item label="文件大小">{{ formatFileSize(currentPreviewImage.size) }}</el-descriptions-item>
            <el-descriptions-item label="图片类型">{{ getImageTypeText(currentPreviewImage.type) }}</el-descriptions-item>
            <el-descriptions-item label="上传时间">{{ formatDate(currentPreviewImage.created_at) }}</el-descriptions-item>
            <el-descriptions-item label="使用次数">{{ currentPreviewImage.usage_count || 0 }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="currentPreviewImage.is_deleted ? 'danger' : 'success'">
                {{ currentPreviewImage.is_deleted ? '已删除' : '正常' }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
    </el-dialog>
    
    <!-- 批量操作 -->
    <div class="batch-actions" v-if="selectedImages.length > 0">
      <el-card>
        <div class="batch-info">
          <span>已选择 {{ selectedImages.length }} 张图片</span>
          <div class="actions">
            <el-button type="warning" @click="handleBatchSoftDelete">
              批量删除
            </el-button>
            <el-button type="success" @click="handleBatchRestore">
              批量恢复
            </el-button>
            <el-button type="danger" @click="handleBatchPhysicalDelete">
              批量彻底删除
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
  Search, Refresh, Upload, Delete, Picture, Check, FolderOpened,
  View, Close, RefreshLeft
} from '@element-plus/icons-vue'
import { uploadAPI } from '@/api/upload'
import { useAdminStore } from '@/stores/admin'
import dayjs from 'dayjs'

const adminStore = useAdminStore()

// 搜索表单
const searchForm = reactive({
  type: '',
  is_deleted: '',
  dateRange: []
})

// 图片列表
const images = ref([])
const tableLoading = ref(false)
const selectedImages = ref([])
const viewMode = ref('grid')

// 图片统计
const imageStats = reactive({
  total: 0,
  normal: 0,
  deleted: 0,
  totalSize: 0
})

// 分页
const pagination = reactive({
  page: 1,
  size: 20,
  total: 0
})

// 图片预览
const previewVisible = ref(false)
const currentPreviewImage = ref(null)

// 上传配置
const uploadUrl = '/api/admin/upload'
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${adminStore.token}`
}))

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 获取图片URL
const getImageUrl = (url) => {
  if (url.startsWith('http')) return url
  return `http://localhost:8080${url}`
}

// 获取图片类型颜色
const getImageTypeColor = (type) => {
  const colorMap = {
    dish: 'primary',
    review: 'success',
    avatar: 'warning'
  }
  return colorMap[type] || ''
}

// 获取图片类型文本
const getImageTypeText = (type) => {
  const textMap = {
    dish: '菜品图片',
    review: '评价图片',
    avatar: '用户头像'
  }
  return textMap[type] || '未知类型'
}

// 加载图片统计
const loadImageStats = async () => {
  try {
    // 模拟数据，实际应该调用API
    imageStats.total = 256
    imageStats.normal = 198
    imageStats.deleted = 58
    imageStats.totalSize = 52428800 // 50MB
  } catch (error) {
    console.error('加载图片统计失败:', error)
  }
}

// 加载图片列表
const loadImages = async () => {
  try {
    tableLoading.value = true
    const params = {
      page: pagination.page,
      size: pagination.size,
      ...searchForm
    }
    
    // 处理日期范围
    if (searchForm.dateRange && searchForm.dateRange.length === 2) {
      params.start_date = searchForm.dateRange[0]
      params.end_date = searchForm.dateRange[1]
    }
    delete params.dateRange
    
    // 模拟数据，实际应该调用API
    const mockImages = [
      {
        id: 1,
        name: 'dish_001.jpg',
        url: 'https://via.placeholder.com/300x200/409eff/ffffff?text=菜品图片1',
        type: 'dish',
        size: 1024000,
        usage_count: 3,
        is_deleted: false,
        created_at: new Date()
      },
      {
        id: 2,
        name: 'review_001.jpg',
        url: 'https://via.placeholder.com/300x200/67c23a/ffffff?text=评价图片1',
        type: 'review',
        size: 512000,
        usage_count: 1,
        is_deleted: false,
        created_at: new Date(Date.now() - 3600000)
      },
      {
        id: 3,
        name: 'avatar_001.jpg',
        url: 'https://via.placeholder.com/300x200/e6a23c/ffffff?text=头像图片1',
        type: 'avatar',
        size: 256000,
        usage_count: 0,
        is_deleted: true,
        created_at: new Date(Date.now() - 7200000)
      },
      {
        id: 4,
        name: 'dish_002.jpg',
        url: 'https://via.placeholder.com/300x200/409eff/ffffff?text=菜品图片2',
        type: 'dish',
        size: 856000,
        usage_count: 2,
        is_deleted: false,
        created_at: new Date(Date.now() - 86400000)
      },
      {
        id: 5,
        name: 'review_002.jpg',
        url: 'https://via.placeholder.com/300x200/67c23a/ffffff?text=评价图片2',
        type: 'review',
        size: 678000,
        usage_count: 0,
        is_deleted: false,
        created_at: new Date(Date.now() - 172800000)
      }
    ]
    
    images.value = mockImages
    pagination.total = 256
  } catch (error) {
    ElMessage.error('加载图片列表失败')
  } finally {
    tableLoading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadImages()
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
  loadImages()
}

// 分页变化
const handleSizeChange = (size) => {
  pagination.size = size
  pagination.page = 1
  loadImages()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  loadImages()
}

// 视图模式变化
const handleViewModeChange = (mode) => {
  viewMode.value = mode
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedImages.value = selection
}

// 网格视图图片选择
const handleImageSelect = (image) => {
  const index = selectedImages.value.findIndex(item => item.id === image.id)
  if (index > -1) {
    selectedImages.value.splice(index, 1)
  } else {
    selectedImages.value.push(image)
  }
}

// 预览图片
const previewImage = (image) => {
  currentPreviewImage.value = image
  previewVisible.value = true
}

// 图片加载错误处理
const handleImageError = (event) => {
  // 避免无限循环，只设置一次占位图
  if (!event.target.dataset.errorHandled) {
    event.target.dataset.errorHandled = 'true'
    // 使用 data URL 作为占位图，避免再次网络请求
    event.target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjQiIGhlaWdodD0iNjQiIHZpZXdCb3g9IjAgMCA2NCA2NCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHJlY3Qgd2lkdGg9IjY0IiBoZWlnaHQ9IjY0IiBmaWxsPSIjRjVGNUY1Ii8+CjxwYXRoIGQ9Ik0yMC41IDI2TDMyIDM3LjVMNDMuNSAyNiIgc3Ryb2tlPSIjQ0NDIiBzdHJva2Utd2lkdGg9IjIiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCIvPgo8Y2lyY2xlIGN4PSIyNiIgY3k9IjI2IiByPSIzIiBmaWxsPSIjQ0NDIi8+Cjx0ZXh0IHg9IjMyIiB5PSI0OCIgZm9udC1mYW1pbHk9IkFyaWFsIiBmb250LXNpemU9IjEwIiBmaWxsPSIjOTk5IiB0ZXh0LWFuY2hvcj0ibWlkZGxlIj7lm77niYfliqDovb3lpLHotKU8L3RleHQ+Cjwvc3ZnPgo='
    event.target.alt = '图片加载失败'
  }
}

// 上传前验证
const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt10M = file.size / 1024 / 1024 < 10
  
  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt10M) {
    ElMessage.error('图片大小不能超过 10MB!')
    return false
  }
  return true
}

// 上传成功
const handleUploadSuccess = (response) => {
  if (response.url) {
    ElMessage.success('图片上传成功')
    loadImages()
    loadImageStats()
  } else {
    ElMessage.error('图片上传失败')
  }
}

// 软删除图片
const handleSoftDelete = async (image) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除图片"${image.name}"吗？删除后可以恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await uploadAPI.softDeleteImage(image.id)
    ElMessage.success('删除成功')
    image.is_deleted = true
    loadImageStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 恢复图片
const handleRestore = async (image) => {
  try {
    // 这里应该调用恢复图片的API
    ElMessage.success('恢复成功')
    image.is_deleted = false
    loadImageStats()
  } catch (error) {
    ElMessage.error('恢复失败')
  }
}

// 物理删除图片
const handlePhysicalDelete = async (image) => {
  if (image.usage_count > 0) {
    ElMessage.warning('该图片正在被使用，不能彻底删除')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要彻底删除图片"${image.name}"吗？此操作不可恢复！`,
      '确认彻底删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'error'
      }
    )
    
    await uploadAPI.physicalDeleteImage(image.id)
    ElMessage.success('彻底删除成功')
    loadImages()
    loadImageStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('彻底删除失败')
    }
  }
}

// 批量软删除
const handleBatchSoftDelete = async () => {
  const validImages = selectedImages.value.filter(img => !img.is_deleted)
  if (validImages.length === 0) {
    ElMessage.warning('没有可删除的图片')
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${validImages.length} 张图片吗？`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 批量删除API调用
    ElMessage.success('批量删除成功')
    validImages.forEach(img => {
      img.is_deleted = true
    })
    selectedImages.value = []
    loadImageStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

// 批量恢复
const handleBatchRestore = async () => {
  const validImages = selectedImages.value.filter(img => img.is_deleted)
  if (validImages.length === 0) {
    ElMessage.warning('没有可恢复的图片')
    return
  }
  
  try {
    // 批量恢复API调用
    ElMessage.success('批量恢复成功')
    validImages.forEach(img => {
      img.is_deleted = false
    })
    selectedImages.value = []
    loadImageStats()
  } catch (error) {
    ElMessage.error('批量恢复失败')
  }
}

// 批量物理删除
const handleBatchPhysicalDelete = async () => {
  const usedImages = selectedImages.value.filter(img => img.usage_count > 0)
  if (usedImages.length > 0) {
    ElMessage.warning(`有 ${usedImages.length} 张图片正在被使用，不能彻底删除`)
    return
  }
  
  try {
    await ElMessageBox.confirm(
      `确定要彻底删除选中的 ${selectedImages.value.length} 张图片吗？此操作不可恢复！`,
      '确认批量彻底删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'error'
      }
    )
    
    // 批量物理删除API调用
    ElMessage.success('批量彻底删除成功')
    selectedImages.value = []
    loadImages()
    loadImageStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量彻底删除失败')
    }
  }
}

// 清理垃圾图片
const handleCleanup = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清理所有未使用的已删除图片吗？此操作不可恢复！',
      '确认清理垃圾',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 清理垃圾图片API调用
    ElMessage.success('垃圾清理完成')
    loadImages()
    loadImageStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('垃圾清理失败')
    }
  }
}

onMounted(() => {
  loadImageStats()
  loadImages()
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

.view-controls {
  margin-bottom: 16px;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
  
  .image-card {
    border: 1px solid #f0f0f0;
    border-radius: 10px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.2s ease;
    
    &:hover {
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
      transform: translateY(-2px);
    }
    
    &.selected {
      border-color: #1677ff;
      box-shadow: 0 0 0 2px rgba(22, 119, 255, 0.15);
    }
    
    &.deleted {
      opacity: 0.6;
      
      .image-container img {
        filter: grayscale(100%);
      }
    }
    
    .image-container {
      position: relative;
      height: 150px;
      overflow: hidden;
      
      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.3s;
      }
      
      .image-overlay {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.45);
        display: flex;
        align-items: center;
        justify-content: center;
        opacity: 0;
        transition: opacity 0.2s;
        
        .image-actions {
          display: flex;
          gap: 8px;
        }
      }
      
      &:hover .image-overlay {
        opacity: 1;
      }
    }
    
    .image-info {
      padding: 12px;
      
      .image-name {
        font-size: 13px;
        font-weight: 500;
        color: #1a1a2e;
        margin-bottom: 4px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
      
      .image-meta {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 4px;
        
        .image-size {
          font-size: 12px;
          color: #8c8c8c;
        }
      }
      
      .image-time {
        font-size: 12px;
        color: #bfbfbf;
      }
    }
  }
}

.table-image-preview {
  width: 56px;
  height: 56px;
  object-fit: cover;
  border-radius: 6px;
  cursor: pointer;
  transition: transform 0.2s;
  
  &:hover {
    transform: scale(1.1);
  }
}

.image-preview-container {
  text-align: center;
  
  .preview-image {
    max-width: 100%;
    max-height: 400px;
    margin-bottom: 20px;
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

.success {
  color: #52c41a !important;
}

.danger {
  color: #ff4d4f !important;
  
  &:hover {
    color: #ff7875 !important;
  }
}

@media (max-width: 768px) {
  .image-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 12px;
  }
}
</style>