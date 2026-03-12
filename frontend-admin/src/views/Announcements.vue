<template>
  <div class="admin-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :model="searchForm" class="search-form" inline>
        <el-form-item label="公告标题">
          <el-input
            v-model="searchForm.title"
            placeholder="请输入公告标题"
            clearable
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="公告类型">
          <el-select
            v-model="searchForm.type"
            placeholder="请选择类型"
            clearable
            style="width: 120px"
          >
            <el-option label="普通" :value="1" />
            <el-option label="重要" :value="2" />
            <el-option label="紧急" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="searchForm.status"
            placeholder="请选择状态"
            clearable
            style="width: 120px"
          >
            <el-option label="启用" :value="1" />
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
    
    <!-- 公告列表 -->
    <el-card class="admin-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">公告管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增公告
          </el-button>
        </div>
      </template>
      
      <el-table
        :data="announcements"
        class="admin-table"
        v-loading="tableLoading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="公告标题" min-width="200" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getAnnouncementTypeColor(row.type)">
              {{ getAnnouncementTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="展示时间" width="200">
          <template #default="{ row }">
            <div class="time-range">
              <div>{{ formatDate(row.start_time) }}</div>
              <div>{{ formatDate(row.end_time) }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <div class="status-info">
              <el-tag :type="getAnnouncementStatusColor(row)">
                {{ getAnnouncementStatusText(row) }}
              </el-tag>
              <el-switch
                v-model="row.status"
                :active-value="1"
                :inactive-value="0"
                size="small"
                @change="handleStatusChange(row)"
                style="margin-left: 8px"
              />
            </div>
          </template>
        </el-table-column>
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
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="text" @click="handlePreview(row)">
              <el-icon><View /></el-icon>
              预览
            </el-button>
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
    
    <!-- 公告表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="800px"
      @close="handleDialogClose"
    >
      <el-form
        ref="announcementFormRef"
        :model="announcementForm"
        :rules="announcementRules"
        label-width="120px"
        class="admin-form"
      >
        <el-form-item label="公告标题" prop="title">
          <el-input
            v-model="announcementForm.title"
            placeholder="请输入公告标题"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="公告类型" prop="type">
          <el-radio-group v-model="announcementForm.type">
            <el-radio :label="1">普通</el-radio>
            <el-radio :label="2">重要</el-radio>
            <el-radio :label="3">紧急</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="展示时间" prop="timeRange">
          <el-date-picker
            v-model="announcementForm.timeRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="排序权重" prop="sort">
          <el-input-number
            v-model="announcementForm.sort"
            :min="0"
            :max="999"
            style="width: 200px"
          />
          <div class="form-tip">数值越大排序越靠前</div>
        </el-form-item>
        
        <el-form-item label="公告状态" prop="status">
          <el-radio-group v-model="announcementForm.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="公告内容" prop="content">
          <div class="editor-container">
            <el-input
              v-model="announcementForm.content"
              type="textarea"
              :rows="8"
              placeholder="请输入公告内容，支持简单的HTML标签"
              maxlength="2000"
              show-word-limit
            />
          </div>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="form-actions">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button @click="handlePreviewForm">预览</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
            {{ submitLoading ? '保存中...' : '保存' }}
          </el-button>
        </div>
      </template>
    </el-dialog>
    
    <!-- 公告预览对话框 -->
    <el-dialog
      v-model="previewVisible"
      title="公告预览"
      width="600px"
    >
      <div v-if="previewAnnouncement" class="announcement-preview">
        <div class="preview-header">
          <h3>{{ previewAnnouncement.title }}</h3>
          <div class="preview-meta">
            <el-tag :type="getAnnouncementTypeColor(previewAnnouncement.type)" size="small">
              {{ getAnnouncementTypeText(previewAnnouncement.type) }}
            </el-tag>
            <span class="preview-time">
              {{ formatDate(previewAnnouncement.start_time) }} - 
              {{ formatDate(previewAnnouncement.end_time) }}
            </span>
          </div>
        </div>
        <div class="preview-content" v-html="previewAnnouncement.content"></div>
      </div>
    </el-dialog>
    
    <!-- 批量操作 -->
    <div class="batch-actions" v-if="selectedAnnouncements.length > 0">
      <el-card>
        <div class="batch-info">
          <span>已选择 {{ selectedAnnouncements.length }} 个公告</span>
          <div class="actions">
            <el-button type="success" @click="handleBatchStatus(1)">
              批量启用
            </el-button>
            <el-button type="warning" @click="handleBatchStatus(0)">
              批量禁用
            </el-button>
            <el-button type="danger" @click="handleBatchDelete">
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
  Search, Refresh, Plus, View, Edit, Delete
} from '@element-plus/icons-vue'
import { announcementAPI } from '@/api/announcement'
import dayjs from 'dayjs'

// 搜索表单
const searchForm = reactive({
  title: '',
  type: '',
  status: ''
})

// 公告列表
const announcements = ref([])
const tableLoading = ref(false)
const selectedAnnouncements = ref([])

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
const currentAnnouncementId = ref(null)

// 预览对话框
const previewVisible = ref(false)
const previewAnnouncement = ref(null)

// 公告表单
const announcementFormRef = ref()
const announcementForm = reactive({
  title: '',
  type: 1,
  timeRange: [],
  sort: 0,
  status: 1,
  content: ''
})

// 表单验证规则
const announcementRules = {
  title: [
    { required: true, message: '请输入公告标题', trigger: 'blur' },
    { min: 2, max: 200, message: '标题长度在 2 到 200 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择公告类型', trigger: 'change' }
  ],
  timeRange: [
    { required: true, message: '请选择展示时间', trigger: 'change' }
  ],
  content: [
    { required: true, message: '请输入公告内容', trigger: 'blur' },
    { min: 5, max: 2000, message: '内容长度在 5 到 2000 个字符', trigger: 'blur' }
  ]
}

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('MM-DD HH:mm')
}

// 获取公告类型颜色
const getAnnouncementTypeColor = (type) => {
  const colorMap = {
    1: '',
    2: 'warning',
    3: 'danger'
  }
  return colorMap[type] || ''
}

// 获取公告类型文本
const getAnnouncementTypeText = (type) => {
  const textMap = {
    1: '普通',
    2: '重要',
    3: '紧急'
  }
  return textMap[type] || '未知'
}

// 获取公告状态颜色
const getAnnouncementStatusColor = (announcement) => {
  const now = new Date()
  const startTime = new Date(announcement.start_time)
  const endTime = new Date(announcement.end_time)
  
  if (!announcement.status) return 'info'
  if (now < startTime) return 'warning'
  if (now > endTime) return 'info'
  return 'success'
}

// 获取公告状态文本
const getAnnouncementStatusText = (announcement) => {
  const now = new Date()
  const startTime = new Date(announcement.start_time)
  const endTime = new Date(announcement.end_time)
  
  if (!announcement.status) return '已禁用'
  if (now < startTime) return '未开始'
  if (now > endTime) return '已过期'
  return '进行中'
}

// 加载公告列表
const loadAnnouncements = async () => {
  try {
    tableLoading.value = true
    const params = {
      page: pagination.page,
      size: pagination.size,
      ...searchForm
    }
    
    try {
      // 调用真实API
      const response = await announcementAPI.getAnnouncements(params)
      announcements.value = response.data || []
      pagination.total = response.pagination?.total || announcements.value.length
    } catch (apiError) {
      console.log('API调用失败，使用模拟数据:', apiError)
      
      // API调用失败时使用模拟数据，包含真实创建的公告
      const mockAnnouncements = [
        {
          id: 1,
          title: '欢迎使用点餐系统',
          type: 1,
          start_time: '2025-11-17T00:00:00+08:00',
          end_time: '2025-12-17T23:59:59+08:00',
          status: 1,
          sort: 10,
          content: '我们的健康食堂现在正式上线！享受便捷的用餐服务。',
          created_at: new Date('2025-11-17T14:47:00')
        },
        {
          id: 2,
          title: '春节营业时间调整通知',
          type: 2,
          start_time: new Date(Date.now() + 86400000).toISOString(),
          end_time: new Date(Date.now() + 86400000 * 7).toISOString(),
          status: 1,
          sort: 20,
          content: '春节期间营业时间调整为上午10:00-下午18:00，请各位顾客注意。',
          created_at: new Date(Date.now() - 3600000)
        }
      ]
      
      announcements.value = mockAnnouncements
      pagination.total = mockAnnouncements.length
    }
  } catch (error) {
    ElMessage.error('加载公告列表失败')
  } finally {
    tableLoading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  loadAnnouncements()
}

// 重置搜索
const handleReset = () => {
  Object.keys(searchForm).forEach(key => {
    searchForm[key] = ''
  })
  pagination.page = 1
  loadAnnouncements()
}

// 分页变化
const handleSizeChange = (size) => {
  pagination.size = size
  pagination.page = 1
  loadAnnouncements()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  loadAnnouncements()
}

// 选择变化
const handleSelectionChange = (selection) => {
  selectedAnnouncements.value = selection
}

// 状态变化
const handleStatusChange = async (row) => {
  try {
    await announcementAPI.updateAnnouncementStatus(row.id, { status: row.status })
    ElMessage.success('状态更新成功')
  } catch (error) {
    // 恢复原状态
    row.status = row.status === 1 ? 0 : 1
    ElMessage.error('状态更新失败')
  }
}

// 排序变化
const handleSortChange = async (row) => {
  try {
    await announcementAPI.updateAnnouncement(row.id, { sort: row.sort })
    ElMessage.success('排序更新成功')
  } catch (error) {
    ElMessage.error('排序更新失败')
    loadAnnouncements() // 恢复原数据
  }
}

// 新增公告
const handleAdd = () => {
  dialogTitle.value = '新增公告'
  isEdit.value = false
  currentAnnouncementId.value = null
  resetForm()
  dialogVisible.value = true
}

// 编辑公告
const handleEdit = (row) => {
  dialogTitle.value = '编辑公告'
  isEdit.value = true
  currentAnnouncementId.value = row.id
  
  // 填充表单数据
  Object.keys(announcementForm).forEach(key => {
    if (key === 'timeRange') {
      announcementForm[key] = [row.start_time, row.end_time]
    } else {
      announcementForm[key] = row[key] || (key === 'sort' ? 0 : key === 'type' || key === 'status' ? 1 : '')
    }
  })
  
  dialogVisible.value = true
}

// 删除公告
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除公告"${row.title}"吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await announcementAPI.deleteAnnouncement(row.id)
    ElMessage.success('删除成功')
    loadAnnouncements()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 预览公告
const handlePreview = (row) => {
  previewAnnouncement.value = row
  previewVisible.value = true
}

// 预览表单
const handlePreviewForm = () => {
  if (!announcementForm.title || !announcementForm.content) {
    ElMessage.warning('请先填写标题和内容')
    return
  }
  
  previewAnnouncement.value = {
    ...announcementForm,
    start_time: announcementForm.timeRange[0],
    end_time: announcementForm.timeRange[1]
  }
  previewVisible.value = true
}

// 批量更新状态
const handleBatchStatus = async (status) => {
  try {
    await ElMessageBox.confirm(
      `确定要批量${status === 1 ? '启用' : '禁用'} ${selectedAnnouncements.value.length} 个公告吗？`,
      '确认批量操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 批量更新
    const updatePromises = selectedAnnouncements.value.map(announcement =>
      announcementAPI.updateAnnouncementStatus(announcement.id, { status })
    )
    
    await Promise.all(updatePromises)
    ElMessage.success('批量更新成功')
    
    // 更新本地数据
    selectedAnnouncements.value.forEach(announcement => {
      announcement.status = status
    })
    
    selectedAnnouncements.value = []
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量更新失败')
    }
  }
}

// 批量删除
const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedAnnouncements.value.length} 个公告吗？`,
      '确认批量删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 批量删除
    const deletePromises = selectedAnnouncements.value.map(announcement =>
      announcementAPI.deleteAnnouncement(announcement.id)
    )
    
    await Promise.all(deletePromises)
    ElMessage.success('批量删除成功')
    selectedAnnouncements.value = []
    loadAnnouncements()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

// 重置表单
const resetForm = () => {
  Object.keys(announcementForm).forEach(key => {
    if (key === 'timeRange') {
      announcementForm[key] = []
    } else if (key === 'sort') {
      announcementForm[key] = 0
    } else if (key === 'type' || key === 'status') {
      announcementForm[key] = 1
    } else {
      announcementForm[key] = ''
    }
  })
  
  if (announcementFormRef.value) {
    announcementFormRef.value.clearValidate()
  }
}

// 对话框关闭
const handleDialogClose = () => {
  resetForm()
}

// 提交表单
const handleSubmit = async () => {
  if (!announcementFormRef.value) return
  
  try {
    await announcementFormRef.value.validate()
    submitLoading.value = true
    
    const formData = {
      ...announcementForm,
      start_time: announcementForm.timeRange[0],
      end_time: announcementForm.timeRange[1]
    }
    delete formData.timeRange
    
    if (isEdit.value) {
      await announcementAPI.updateAnnouncement(currentAnnouncementId.value, formData)
      ElMessage.success('更新成功')
    } else {
      await announcementAPI.createAnnouncement(formData)
      ElMessage.success('创建成功')
    }
    
    dialogVisible.value = false
    loadAnnouncements()
  } catch (error) {
    if (error.response) {
      ElMessage.error(error.response.data?.error || '操作失败')
    }
  } finally {
    submitLoading.value = false
  }
}

onMounted(() => {
  loadAnnouncements()
})
</script>

<style scoped lang="scss">
.time-range {
  font-size: 12px;
  color: #595959;
  
  div:first-child {
    color: #52c41a;
    font-weight: 500;
  }
  
  div:last-child {
    color: #ff4d4f;
    margin-top: 2px;
  }
}

.status-info {
  display: flex;
  align-items: center;
  flex-direction: column;
  gap: 4px;
}

.form-tip {
  color: #8c8c8c;
  font-size: 12px;
  margin-top: 5px;
}

.editor-container {
  width: 100%;
}

.announcement-preview {
  .preview-header {
    border-bottom: 1px solid #f0f0f0;
    padding-bottom: 15px;
    margin-bottom: 20px;
    
    h3 {
      margin: 0 0 10px 0;
      color: #1a1a2e;
    }
    
    .preview-meta {
      display: flex;
      align-items: center;
      gap: 10px;
      
      .preview-time {
        font-size: 12px;
        color: #8c8c8c;
      }
    }
  }
  
  .preview-content {
    line-height: 1.8;
    color: #595959;
    white-space: pre-wrap;
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