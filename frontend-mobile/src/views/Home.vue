<template>
  <div class="home">
    <!-- 公告轮播 -->
    <div v-if="announcements.length" class="announcement-swiper">
      <swiper
        :modules="modules"
        :autoplay="{ delay: 3000, disableOnInteraction: false }"
        :pagination="{ clickable: true }"
        :loop="true"
        class="swiper"
      >
        <swiper-slide v-for="item in announcements" :key="item.id">
          <div class="announcement-item" @click="showAnnouncementDetail(item)">
            <van-icon name="volume-o" />
            <span class="announcement-text">{{ item.title }}</span>
          </div>
        </swiper-slide>
      </swiper>
    </div>

    <!-- 分类导航 -->
    <div class="categories">
      <van-grid :column-num="4" :border="false">
        <!-- 本周菜谱入口 -->
        <van-grid-item
          text="本周菜谱"
          @click="$router.push('/weekly-menu')"
        >
          <template #icon>
            <van-icon name="calendar-o" size="32" color="#1989fa" />
          </template>
        </van-grid-item>
        
        <van-grid-item
          v-for="category in categories"
          :key="category.id"
          :text="category.name"
          @click="goToDishes(category.id)"
        >
          <template #icon>
            <van-icon name="apps-o" size="32" />
          </template>
        </van-grid-item>
      </van-grid>
    </div>

    <!-- 推荐菜品 -->
    <div class="dishes-section">
      <div class="section-header">
        <h3>推荐菜品</h3>
        <van-button type="primary" size="small" @click="$router.push('/dishes')">
          查看全部
        </van-button>
      </div>
      
      <div class="dishes-list">
        <van-card
          v-for="dish in dishes"
          :key="dish.id"
          :price="dish.price"
          :title="dish.name"
          :desc="dish.description"
          :thumb="dish.images?.[0]?.image_url"
          @click="goToDishDetail(dish.id)"
        >
          <template #footer>
            <van-button size="small" type="primary" @click.stop="addToCart(dish)">
              加入购物车
            </van-button>
          </template>
        </van-card>
      </div>
    </div>

    <!-- 底部导航 -->
    <van-tabbar v-model="active" route>
      <van-tabbar-item to="/home" icon="home-o">首页</van-tabbar-item>
      <van-tabbar-item to="/dishes" icon="apps-o">菜品</van-tabbar-item>
      <van-tabbar-item to="/cart" icon="shopping-cart-o" :badge="cartCount">购物车</van-tabbar-item>
      <van-tabbar-item to="/profile" icon="user-o">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { showToast, showDialog } from 'vant'
import { Swiper, SwiperSlide } from 'swiper/vue'
import { Autoplay, Pagination } from 'swiper/modules'
import { getDishes, getCategories } from '@/api/dish'
import { getActiveAnnouncements } from '@/api/announcement'
import { useCartStore } from '@/stores/cart'
import { useRouter } from 'vue-router'

const router = useRouter()
const cartStore = useCartStore()
const modules = [Autoplay, Pagination]

const active = ref(0)
const announcements = ref([])
const categories = ref([])
const dishes = ref([])

const cartCount = computed(() => cartStore.totalCount)

const fetchData = async () => {
  try {
    const [announcementRes, categoryRes, dishRes] = await Promise.all([
      getActiveAnnouncements(),
      getCategories(),
      getDishes({ page: 1, page_size: 10 })
    ])
    
    announcements.value = announcementRes.data || []
    categories.value = categoryRes.data || []
    dishes.value = dishRes.data || []
  } catch (error) {
    console.error('获取数据失败', error)
  }
}

const showAnnouncementDetail = (announcement) => {
  showDialog({
    title: announcement.title,
    message: announcement.content,
    confirmButtonText: '知道了'
  })
}

const goToDishes = (categoryId) => {
  router.push({ path: '/dishes', query: { category_id: categoryId } })
}

const goToDishDetail = (id) => {
  router.push(`/dish/${id}`)
}

const addToCart = (dish) => {
  cartStore.addItem(dish)
  showToast('已添加到购物车')
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.home {
  padding-bottom: 50px;
}

.announcement-swiper {
  background: #fff;
  padding: 10px 0;
  
  .swiper {
    height: 40px;
  }
  
  .announcement-item {
    display: flex;
    align-items: center;
    padding: 0 16px;
    cursor: pointer;
    
    .van-icon {
      color: var(--primary-color);
      margin-right: 8px;
    }
    
    .announcement-text {
      flex: 1;
      font-size: 14px;
      color: var(--text-color);
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}

.categories {
  margin-top: 10px;
  background: #fff;
}

.dishes-section {
  margin-top: 10px;
  background: #fff;
  padding: 16px;
  
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    
    h3 {
      font-size: 18px;
      font-weight: 600;
    }
  }
  
  .dishes-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
}
</style>
