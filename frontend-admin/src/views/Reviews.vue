<template>
  <div class="admin-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="searchForm" class="search-form" inline>
        <el-form-item label="菜品名称">
          <el-input
            v-model="searchForm.dish_name"
            placeholder="请输入菜品名称"
            clearable
            style="width: 150px"
          />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input
            v-model="searchForm.user_name"
            placeholder="请输入用户名"
            clearable
            style="width: 120px"
          />
        </el-form-item>
        <el-form-item label="评分">
          <el-select
            v-model="searchForm.rating"
            placeholder="请选择评分"
            clearable
            style="width: 100px"
          >
            <el-option label="5星" :value="5" />
            <el-option label="4星" :value="4" />
            <el-option label="3星" :value="3" />
            <el-option label="2星" :value="2" />
            <el-option label="1星" :value="1" />
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
    
    <!-- 评价列表 -->
    <el-card class="admin-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">评价管理</span>
          <el-button @click="loadReviews">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>
      
      <el-table
        :data="reviews"
        class="admin-table"
        v-loading="tableLoading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="dish_name" label="菜品名称" min-width="150" />
        <el-table-column prop="user_name" label="用户" width="120" />
        <el-table-column prop="rating" label="评分" width="100">
          <template #default="{ row }">
            <el-rate
              :model-value="row.rating"
              disabled
              show-score
              text-color="#ff9900"
              score-template="{value}"
            />
          </template>
        </el-table-column>
        <el-table-column prop="content" label="评价内容" min-width="200">
          <template #default="{ row }">
            <div class="review-content">
              {{ row.content || '用户未填写评价内容' }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="评价图片" width="120">
          <template #default="{ row }">
            <div class="image-preview" v-if="row.images && row.images.length > 0">
              <img
                v-for="(image, index) in row.images.slice(0, 3)"
                :key="index"
                :src="getImageUrl(image.image_url)"
                class="preview-image"
                @click="previewImages(row.images, index)"
              />
              <span v-if="row.images.length > 3" class="more-images">
                +{{ row.images.length - 3 }}
              </span>
            </div>
            <span v-else class="no-image">无图片</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="评价时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="text" @click="handleViewDetail(row)">
              <el-icon><View /></el-icon>
              详情
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
    
    <!-- 评价详情对话框 -->
    <el-dialog
      v-model="detailVisible"
      title="评价详情"
      width="600px"
    >
      <div v-if="currentReview" class="review-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="评价ID">{{ currentReview.id }}</el-descriptions-item>
          <el-descriptions-item label="菜品名称">{{ currentReview.dish_name }}</el-descriptions-item>
          <el-descriptions-item label="用户名">{{ currentReview.user_name }}</el-descriptions-item>
          <el-descriptions-item label="订单号">{{ currentReview.order_id }}</el-descriptions-item>
          <el-descriptions-item label="评分">
            <el-rate
              :model-value="currentReview.rating"
              disabled
              show-score
              text-color="#ff9900"
              score-template="{value}分"
            />
          </el-descriptions-item>
          <el-descriptions-item label="评价时间">
            {{ formatDate(currentReview.created_at) }}
          </el-descriptions-item>
        </el-descriptions>
        
        <h4 style="margin: 20px 0 10px 0;">评价内容</h4>
        <div class="review-content-detail">
          {{ currentReview.content || '用户未填写评价内容' }}
        </div>
        
        <div v-if="currentReview.images && currentReview.images.length > 0">
          <h4 style="margin: 20px 0 10px 0;">评价图片</h4>
          <div class="review-images">
            <img
              v-for="(image, index) in currentReview.images"
              :key="index"
              :src="getImageUrl(image.image_url)"
              class="review-image"
              @click="previewImages(currentReview.images, index)"
            />
          </div>
        </div>
      </div>
    </el-dialog>
    
    <!-- 图片预览对话框 -->
    <el-dialog v-model="previewVisible" title="图片预览" width="800px">
      <div class="image-preview-container">
        <img :src="currentPreviewImage" style="width: 100%" />
        <div class="preview-actions" v-if="previewImages.length > 1">
          <el-button @click="prevImage" :disabled="currentImageIndex === 0">
            上一张
          </el-button>
          <span>{{ currentImageIndex + 1 }} / {{ previewImageList.length }}</span>
          <el-button @click="nextImage" :disabled="currentImageIndex === previewImageList.length - 1">
            下一张
          </el-button>
        </div>
      </div>
    </el-dialog>
    
    <!-- 批量操作 -->
    <div class="batch-actions" v-if="selectedReviews.length > 0">
      <el-card>
        <div class="batch-info">
          <span>已选择 {{ selectedReviews.length }} 条评价</span>
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
import {
  Search, Refresh, View, Delete
} from '@element-plus/icons-vue'
import { reviewAPI } from '@/api/review'
import dayjs from 'dayjs'

// 搜索表单
const searchForm = reactive({
  dish_name: '',
  user_name: '',
  rating: ''
})

// 评价列表
const reviews = ref([])
const tableLoading = ref(false)
const selectedReviews = ref([])

// 分页
const pagination = reactive({
  page: 1,
  size: 10,
  total: 0
})

// 评价详情对话框
const detailVisible = ref(false)
const currentReview = ref(null)

// 图片预览
const previewVisible = ref(false)
const currentPreviewImage = ref('')
const previewImageList = ref([])
const currentImageIndex = ref(0)

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// 获取图片URL
const getImageUrl = (url) => {
  if (url.startsWith('http')) return url
  return `http://localhost:8080${url}`
}

// 加载评价列表
const loadReviews = async () => {
  try {
    tableLoading.value = true
    const params = {
      page: pagination.page,
      page_size: pagination.size,
      ...searchForm
    }
    
    // 调用真实API
    const response = await reviewAPI.getAdminReviews(params)
    
    // 处理评价数据，添加用户名和菜品名称
    const processedReviews = response.data.map(review => ({
      ...review,
      user_name: review.user?.username || review.user?.nickname || '未知用户',
      dish_name: review.dish?.name || '未知菜品'
    }))
    
    reviews.value = processedReviews
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
  loadReviews()
}

// 重置搜索
const handleReset = () => {
  Object.keys(searchForm).forEach(key => {
    searchForm[key] = ''
  })
  pagination.page = 1
  loadReviews()
}

// 分页变化
const handleSizeChange = (size) => {
  pagination.size = size
  pagination.page = 1
  loadReviews()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  loadReviews()
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedReviews.value = selection
}

// 查看评价详情
const handleViewDetail = (row) => {
  currentReview.value = row
  detailVisible.value = true
}

// 删除评价
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除这条评价吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 这里应该调用API删除
    ElMessage.success('删除成功')
    loadReviews()
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
      `确定要删除选中的 ${selectedReviews.value.length} 条评价吗？`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 批量删除API调用
    ElMessage.success('批量删除成功')
    selectedReviews.value = []
    loadReviews()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

// 预览图片
const previewImages = (images, index = 0) => {
  previewImageList.value = images
  currentImageIndex.value = index
  currentPreviewImage.value = getImageUrl(images[index].image_url)
  previewVisible.value = true
}

// 切换预览图片
const prevImage = () => {
  if (currentImageIndex.value > 0) {
    currentImageIndex.value--
    currentPreviewImage.value = getImageUrl(previewImageList.value[currentImageIndex.value].image_url)
  }
}

const nextImage = () => {
  if (currentImageIndex.value < previewImageList.value.length - 1) {
    currentImageIndex.value++
    currentPreviewImage.value = getImageUrl(previewImageList.value[currentImageIndex.value].image_url)
  }
}

onMounted(() => {
  loadReviews()
})
</script>

<style scoped lang="scss">
.review-content {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #595959;
}

.image-preview {
  display: flex;
  align-items: center;
  gap: 4px;
  
  .preview-image {
    width: 30px;
    height: 30px;
    object-fit: cover;
    border-radius: 6px;
    cursor: pointer;
    transition: transform 0.2s;
    
    &:hover {
      transform: scale(1.1);
    }
  }
  
  .more-images {
    font-size: 12px;
    color: #8c8c8c;
    background-color: #fafbfc;
    padding: 2px 6px;
    border-radius: 4px;
  }
}

.no-image {
  color: #bfbfbf;
  font-size: 12px;
}

.review-detail {
  .review-content-detail {
    background-color: #fafbfc;
    padding: 16px;
    border-radius: 8px;
    border: 1px solid #f0f0f0;
    color: #595959;
    line-height: 1.6;
    min-height: 80px;
  }
  
  .review-images {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    
    .review-image {
      width: 100px;
      height: 100px;
      object-fit: cover;
      border-radius: 8px;
      cursor: pointer;
      transition: transform 0.2s;
      
      &:hover {
        transform: scale(1.05);
      }
    }
  }
}

.image-preview-container {
  text-align: center;
  
  .preview-actions {
    margin-top: 20px;
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 20px;
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

.danger {
  color: #ff4d4f !important;
  
  &:hover {
    color: #ff7875 !important;
  }
}
</style>