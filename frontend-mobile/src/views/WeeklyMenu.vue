<template>
  <div class="weekly-menu-page">
    <!-- 导航栏 -->
    <van-nav-bar
      title="本周菜谱"
      left-arrow
      @click-left="$router.back()"
    />
    
    <!-- 周选择器 -->
    <div class="week-selector">
      <van-button 
        icon="arrow-left" 
        size="small" 
        @click="previousWeek"
        :disabled="loading"
      />
      <div class="week-info">
        <div class="week-title">{{ weekTitle }}</div>
        <div class="week-range">{{ weekRange }}</div>
        <div v-if="isCycleMenu" class="cycle-badge">
          <van-tag type="warning" size="small">循环菜谱</van-tag>
        </div>
      </div>
      <van-button 
        icon="arrow" 
        size="small" 
        @click="nextWeek"
        :disabled="loading"
      />
    </div>
    
    <!-- 菜谱内容 -->
    <div class="menu-content" v-loading="loading">
      <div v-if="!menuData" class="empty-state">
        <van-empty 
          image="search" 
          description="本周菜谱暂未发布"
        />
      </div>
      
      <div v-else class="menu-calendar">
        <!-- 日期标签 -->
        <div class="date-tabs">
          <div 
            class="date-tab"
            :class="{ active: selectedDate === day.date }"
            v-for="day in weekDays"
            :key="day.date"
            @click="selectDate(day.date)"
          >
            <div class="day-name">{{ day.dayName }}</div>
            <div class="date-num">{{ day.dateNum }}</div>
          </div>
        </div>
        
        <!-- 选中日期的菜谱 -->
        <div class="daily-menu">
          <div 
            class="meal-section"
            v-for="meal in mealTypes"
            :key="meal.value"
          >
            <div class="meal-header">
              <van-icon :name="meal.icon" />
              <span class="meal-name">{{ meal.label }}</span>
              <span class="dish-count">({{ getMealDishes(meal.value).length }}道菜)</span>
            </div>
            
            <div class="dish-list">
              <div 
                class="dish-card"
                v-for="item in getMealItems(meal.value)"
                :key="item.id"
                @click="viewDishDetail(item)"
              >
                <div class="dish-image">
                  <van-image
                    :src="getDishImage(item.dish)"
                    fit="cover"
                    lazy-load
                  >
                    <template #error>
                      <div class="image-error">
                        <van-icon name="photo-o" />
                      </div>
                    </template>
                  </van-image>
                </div>
                <div class="dish-info">
                  <div class="dish-name">{{ item.dish.name }}</div>
                  <div class="dish-meta">
                    <span class="dish-price">¥{{ item.dish.price }}</span>
                    <span class="dish-rating" v-if="getRating(item)">
                      <van-icon name="star" color="#ffd21e" />
                      {{ getRating(item) }}分
                    </span>
                  </div>
                  <div class="dish-actions">
                    <van-button 
                      size="small" 
                      type="primary" 
                      @click.stop="addToCart(item.dish)"
                    >
                      加入购物车
                    </van-button>
                    <van-button 
                      size="small" 
                      plain
                      type="warning"
                      @click.stop="openRatingPopup(item, meal.value)"
                    >
                      {{ hasRated(item, meal.value) ? '已评价' : '评价' }}
                    </van-button>
                  </div>
                </div>
              </div>
              
              <div v-if="getMealDishes(meal.value).length === 0" class="no-dishes">
                <van-icon name="info-o" />
                <span>暂无菜品</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 菜品详情弹窗 -->
    <van-popup 
      v-model:show="dishDetailVisible" 
      position="bottom" 
      :style="{ height: '70%' }"
      round
    >
      <div class="dish-detail-popup" v-if="selectedItem">
        <div class="popup-header">
          <div class="dish-title">{{ selectedItem.dish.name }}</div>
          <van-icon name="cross" @click="dishDetailVisible = false" />
        </div>
        
        <div class="dish-images">
          <van-swipe :autoplay="3000" indicator-color="white">
            <van-swipe-item v-for="(image, index) in selectedItem.dish.images" :key="index">
              <van-image
                :src="getImageUrl(image.image_url)"
                fit="cover"
                width="100%"
                height="200px"
              />
            </van-swipe-item>
          </van-swipe>
        </div>
        
        <div class="dish-content">
          <div class="price-section">
            <span class="current-price">¥{{ selectedItem.dish.price }}</span>
            <span class="like-count">
              <van-icon name="good-job-o" />
              {{ selectedItem.dish.like_count || 0 }}
            </span>
          </div>
          
          <div class="description">
            {{ selectedItem.dish.description || '暂无描述' }}
          </div>
          
          <div class="category-info">
            <van-tag type="primary">{{ selectedItem.dish.category?.name }}</van-tag>
          </div>
        </div>
        
        <div class="popup-footer">
          <van-button 
            type="primary" 
            block 
            @click="addToCart(selectedItem.dish)"
          >
            加入购物车
          </van-button>
        </div>
      </div>
    </van-popup>
    
    <!-- 评价弹窗 -->
    <van-popup 
      v-model:show="ratingPopupVisible" 
      position="bottom" 
      :style="{ height: '50%' }"
      round
    >
      <div class="rating-popup" v-if="ratingItem">
        <div class="popup-header">
          <div class="dish-title">评价: {{ ratingItem.dish.name }}</div>
          <van-icon name="cross" @click="ratingPopupVisible = false" />
        </div>
        
        <div class="rating-content">
          <div class="rating-section">
            <div class="rating-label">您的评分</div>
            <van-rate 
              v-model="ratingForm.rating" 
              :size="28"
              color="#ffd21e"
              void-icon="star-o"
              void-color="#eee"
              allow-half
            />
            <div class="rating-text">{{ getRatingText(ratingForm.rating) }}</div>
          </div>
          
          <div class="comment-section">
            <div class="comment-label">评价内容（选填）</div>
            <van-field
              v-model="ratingForm.comment"
              rows="3"
              autosize
              type="textarea"
              maxlength="200"
              placeholder="说说您对这道菜的看法..."
              show-word-limit
            />
          </div>
        </div>
        
        <div class="popup-footer">
          <van-button 
            type="primary" 
            block 
            :loading="ratingSubmitting"
            @click="submitRating"
          >
            提交评价
          </van-button>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { showToast, showLoadingToast, closeToast } from 'vant'
import { weeklyMenuAPI } from '@/api/weekly_menu'
import dayjs from 'dayjs'
import weekOfYear from 'dayjs/plugin/weekOfYear'
import isoWeek from 'dayjs/plugin/isoWeek'

// 启用插件
dayjs.extend(weekOfYear)
dayjs.extend(isoWeek)

// 响应式数据
const loading = ref(false)
const menuData = ref(null)
const isCycleMenu = ref(false) // 是否为循环菜谱
const currentWeekStart = ref(dayjs().isoWeekday(1)) // 使用 ISO 周一
const selectedDate = ref('')
const dishDetailVisible = ref(false)
const selectedItem = ref(null) // 改为 menu_item
const myRatings = ref([]) // 用户的评价列表

// 评价相关
const ratingPopupVisible = ref(false)
const ratingItem = ref(null)
const ratingMealType = ref(1)
const ratingSubmitting = ref(false)
const ratingForm = reactive({
  rating: 5,
  comment: ''
})

// 餐次配置
const mealTypes = [
  { value: 1, label: '早餐', icon: 'sun-o' },
  { value: 2, label: '午餐', icon: 'sun' },
  { value: 3, label: '晚餐', icon: 'moon-o' },
  { value: 4, label: '值班餐', icon: 'clock-o' }
]

// 计算属性
const weekTitle = computed(() => {
  const year = currentWeekStart.value.year()
  const week = currentWeekStart.value.isoWeek()
  return `${year}年第${week}周`
})

const weekRange = computed(() => {
  const start = currentWeekStart.value.format('MM月DD日')
  const end = currentWeekStart.value.add(6, 'day').format('MM月DD日')
  return `${start} - ${end}`
})

const weekDays = computed(() => {
  const days = []
  for (let i = 0; i < 7; i++) {
    const date = currentWeekStart.value.add(i, 'day')
    days.push({
      date: date.format('YYYY-MM-DD'),
      dayName: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'][i],
      dateNum: date.format('DD')
    })
  }
  return days
})

// 方法
const loadWeekMenu = async () => {
  try {
    loading.value = true
    const dateStr = currentWeekStart.value.format('YYYY-MM-DD')
    
    try {
      const response = await weeklyMenuAPI.getWeekMenuByDate(dateStr)
      menuData.value = response.data
      isCycleMenu.value = response.is_cycle || false
      
      // 默认选择今天或周一
      const today = dayjs().format('YYYY-MM-DD')
      const isInCurrentWeek = weekDays.value.some(day => day.date === today)
      selectedDate.value = isInCurrentWeek ? today : weekDays.value[0].date
      
      // 加载用户评价
      loadMyRatings()
    } catch (error) {
      if (error.response?.status === 404) {
        menuData.value = null
        isCycleMenu.value = false
      } else {
        throw error
      }
    }
  } catch (error) {
    console.error('加载菜谱失败:', error)
    showToast('加载菜谱失败')
  } finally {
    loading.value = false
  }
}

// 加载用户的评价
const loadMyRatings = async () => {
  if (!menuData.value) return
  
  try {
    const response = await weeklyMenuAPI.getMyMenuRatings({
      menu_id: menuData.value.id
    })
    myRatings.value = response.data || []
  } catch (error) {
    console.error('加载评价失败:', error)
    myRatings.value = []
  }
}

const previousWeek = () => {
  currentWeekStart.value = currentWeekStart.value.subtract(1, 'week')
  loadWeekMenu()
}

const nextWeek = () => {
  currentWeekStart.value = currentWeekStart.value.add(1, 'week')
  loadWeekMenu()
}

const selectDate = (date) => {
  selectedDate.value = date
}

// 获取某餐次的菜品列表（返回 dish 对象）
const getMealDishes = (mealType) => {
  return getMealItems(mealType).map(item => item.dish).filter(Boolean)
}

// 获取某餐次的 menu_item 列表
const getMealItems = (mealType) => {
  if (!menuData.value || !selectedDate.value) return []
  
  // 获取选中日期是周几（0-6，周一到周日）
  const selectedDayIndex = weekDays.value.findIndex(d => d.date === selectedDate.value)
  
  return menuData.value.menu_items?.filter(item => {
    // 对于循环菜谱，匹配周几
    const itemDate = dayjs(item.date)
    const itemDayOfWeek = itemDate.isoWeekday() - 1 // 0-6
    
    return itemDayOfWeek === selectedDayIndex && item.meal_type === mealType
  }) || []
}

// 获取菜品的用户评分
const getRating = (item) => {
  const rating = myRatings.value.find(r => 
    r.dish_id === item.dish_id && 
    r.meal_type === item.meal_type
  )
  return rating?.rating || 0
}

// 检查是否已评价
const hasRated = (item, mealType) => {
  return myRatings.value.some(r => 
    r.dish_id === item.dish_id && 
    r.meal_type === mealType &&
    dayjs(r.rating_date).format('YYYY-MM-DD') === selectedDate.value
  )
}

const getDishImage = (dish) => {
  if (!dish || !dish.images || dish.images.length === 0) {
    return 'https://via.placeholder.com/200x150/f0f0f0/cccccc?text=暂无图片'
  }
  
  // 优先显示主图
  const mainImage = dish.images.find(img => img.is_main)
  const imageUrl = mainImage ? mainImage.image_url : dish.images[0].image_url
  
  return getImageUrl(imageUrl)
}

const getImageUrl = (url) => {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return `http://localhost:8080${url}`
}

const viewDishDetail = (item) => {
  selectedItem.value = item
  dishDetailVisible.value = true
}

const addToCart = (dish) => {
  // 这里实现添加到购物车的逻辑
  showToast('已添加到购物车')
  dishDetailVisible.value = false
}

// 评价相关方法
const openRatingPopup = (item, mealType) => {
  ratingItem.value = item
  ratingMealType.value = mealType
  
  // 如果已有评价，加载已有评价
  const existingRating = myRatings.value.find(r => 
    r.dish_id === item.dish_id && 
    r.meal_type === mealType &&
    dayjs(r.rating_date).format('YYYY-MM-DD') === selectedDate.value
  )
  
  if (existingRating) {
    ratingForm.rating = existingRating.rating
    ratingForm.comment = existingRating.comment || ''
  } else {
    ratingForm.rating = 5
    ratingForm.comment = ''
  }
  
  ratingPopupVisible.value = true
}

const getRatingText = (rating) => {
  if (rating <= 1) return '很差'
  if (rating <= 2) return '较差'
  if (rating <= 3) return '一般'
  if (rating <= 4) return '不错'
  return '非常好'
}

const submitRating = async () => {
  if (!ratingItem.value || !menuData.value) return
  
  if (ratingForm.rating < 1) {
    showToast('请选择评分')
    return
  }
  
  try {
    ratingSubmitting.value = true
    
    await weeklyMenuAPI.createMenuRating({
      menu_id: menuData.value.id,
      menu_item_id: ratingItem.value.id,
      dish_id: ratingItem.value.dish_id,
      rating: Math.round(ratingForm.rating),
      comment: ratingForm.comment,
      meal_type: ratingMealType.value,
      rating_date: selectedDate.value
    })
    
    showToast('评价成功')
    ratingPopupVisible.value = false
    
    // 重新加载评价
    loadMyRatings()
  } catch (error) {
    console.error('提交评价失败:', error)
    showToast(error.response?.data?.error || '评价失败，请稍后重试')
  } finally {
    ratingSubmitting.value = false
  }
}

onMounted(() => {
  loadWeekMenu()
})
</script>

<style scoped lang="scss">
.weekly-menu-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.week-selector {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background-color: white;
  margin-bottom: 8px;
  
  .week-info {
    text-align: center;
    
    .week-title {
      font-size: 16px;
      font-weight: 500;
      color: #323233;
    }
    
    .week-range {
      font-size: 12px;
      color: #969799;
      margin-top: 4px;
    }
    
    .cycle-badge {
      margin-top: 6px;
    }
  }
}

.menu-content {
  padding: 0 16px 16px;
}

.empty-state {
  background-color: white;
  border-radius: 8px;
  padding: 40px 20px;
  text-align: center;
}

.menu-calendar {
  .date-tabs {
    display: flex;
    background-color: white;
    border-radius: 8px;
    padding: 8px;
    margin-bottom: 16px;
    overflow-x: auto;
    
    .date-tab {
      flex: 1;
      min-width: 60px;
      padding: 12px 8px;
      text-align: center;
      border-radius: 6px;
      cursor: pointer;
      transition: all 0.3s;
      
      &.active {
        background-color: #1989fa;
        color: white;
      }
      
      .day-name {
        font-size: 12px;
        margin-bottom: 4px;
      }
      
      .date-num {
        font-size: 16px;
        font-weight: 500;
      }
    }
  }
}

.daily-menu {
  .meal-section {
    background-color: white;
    border-radius: 8px;
    margin-bottom: 16px;
    overflow: hidden;
    
    .meal-header {
      display: flex;
      align-items: center;
      padding: 16px;
      background-color: #f7f8fa;
      border-bottom: 1px solid #ebedf0;
      
      .van-icon {
        margin-right: 8px;
        color: #1989fa;
      }
      
      .meal-name {
        font-size: 16px;
        font-weight: 500;
        color: #323233;
      }
      
      .dish-count {
        margin-left: 8px;
        font-size: 12px;
        color: #969799;
      }
    }
    
    .dish-list {
      padding: 16px;
      
      .dish-card {
        display: flex;
        padding: 12px 0;
        border-bottom: 1px solid #f0f0f0;
        cursor: pointer;
        
        &:last-child {
          border-bottom: none;
        }
        
        .dish-image {
          width: 80px;
          height: 80px;
          border-radius: 8px;
          overflow: hidden;
          margin-right: 12px;
          flex-shrink: 0;
          
          .van-image {
            width: 100%;
            height: 100%;
          }
          
          .image-error {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 100%;
            height: 100%;
            background-color: #f7f8fa;
            color: #c8c9cc;
          }
        }
        
        .dish-info {
          flex: 1;
          display: flex;
          flex-direction: column;
          justify-content: space-between;
          min-width: 0;
          
          .dish-name {
            font-size: 16px;
            font-weight: 500;
            color: #323233;
            margin-bottom: 4px;
          }
          
          .dish-meta {
            display: flex;
            align-items: center;
            gap: 12px;
            margin-bottom: 8px;
            
            .dish-price {
              font-size: 16px;
              color: #ee0a24;
              font-weight: 600;
            }
            
            .dish-rating {
              font-size: 12px;
              color: #ff9800;
              display: flex;
              align-items: center;
              gap: 2px;
            }
          }
          
          .dish-actions {
            display: flex;
            gap: 8px;
            
            .van-button {
              --van-button-small-height: 28px;
              --van-button-small-font-size: 12px;
            }
          }
        }
      }
      
      .no-dishes {
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 40px 20px;
        color: #c8c9cc;
        
        .van-icon {
          margin-right: 8px;
        }
      }
    }
  }
}

.dish-detail-popup {
  height: 100%;
  display: flex;
  flex-direction: column;
  
  .popup-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    border-bottom: 1px solid #ebedf0;
    
    .dish-title {
      font-size: 18px;
      font-weight: 500;
      color: #323233;
    }
    
    .van-icon {
      font-size: 20px;
      color: #969799;
      cursor: pointer;
    }
  }
  
  .dish-images {
    .van-swipe {
      height: 200px;
    }
  }
  
  .dish-content {
    flex: 1;
    padding: 16px;
    
    .price-section {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 16px;
      
      .current-price {
        font-size: 24px;
        color: #ee0a24;
        font-weight: 600;
      }
      
      .like-count {
        display: flex;
        align-items: center;
        color: #969799;
        font-size: 14px;
        
        .van-icon {
          margin-right: 4px;
        }
      }
    }
    
    .description {
      font-size: 14px;
      color: #646566;
      line-height: 1.6;
      margin-bottom: 16px;
    }
    
    .category-info {
      .van-tag {
        --van-tag-primary-color: #1989fa;
      }
    }
  }
  
  .popup-footer {
    padding: 16px;
    border-top: 1px solid #ebedf0;
  }
}

// 评价弹窗样式
.rating-popup {
  height: 100%;
  display: flex;
  flex-direction: column;
  
  .popup-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    border-bottom: 1px solid #ebedf0;
    
    .dish-title {
      font-size: 18px;
      font-weight: 500;
      color: #323233;
    }
    
    .van-icon {
      font-size: 20px;
      color: #969799;
      cursor: pointer;
    }
  }
  
  .rating-content {
    flex: 1;
    padding: 20px 16px;
    overflow-y: auto;
    
    .rating-section {
      text-align: center;
      margin-bottom: 24px;
      
      .rating-label {
        font-size: 14px;
        color: #646566;
        margin-bottom: 12px;
      }
      
      .rating-text {
        margin-top: 8px;
        font-size: 14px;
        color: #ff9800;
      }
    }
    
    .comment-section {
      .comment-label {
        font-size: 14px;
        color: #646566;
        margin-bottom: 8px;
      }
      
      :deep(.van-field) {
        background-color: #f7f8fa;
        border-radius: 8px;
        
        .van-field__control {
          min-height: 80px;
        }
      }
    }
  }
  
  .popup-footer {
    padding: 16px;
    border-top: 1px solid #ebedf0;
  }
}
</style>