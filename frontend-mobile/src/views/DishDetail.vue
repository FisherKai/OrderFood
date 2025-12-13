<template>
  <div class="dish-detail">
    <van-nav-bar title="菜品详情" left-arrow @click-left="$router.back()" />
    
    <div v-if="dish" class="content">
      <van-swipe :autoplay="3000" lazy-render>
        <van-swipe-item v-for="image in dish.images" :key="image.id">
          <img :src="image.image_url" class="dish-image" />
        </van-swipe-item>
      </van-swipe>
      
      <div class="dish-info">
        <h2>{{ dish.name }}</h2>
        <div class="price">¥{{ dish.price }}</div>
        <div class="description">{{ dish.description }}</div>
      </div>
      
      <div class="reviews-section">
        <h3>用户评价</h3>
        <van-empty v-if="!reviews.length" description="暂无评价" />
        <div v-else class="reviews-list">
          <div v-for="review in reviews" :key="review.id" class="review-item">
            <div class="review-header">
              <span class="username">{{ review.user?.nickname || review.user?.username }}</span>
              <van-rate v-model="review.rating" readonly size="12" />
            </div>
            <div class="review-content">{{ review.content }}</div>
            <div v-if="review.images?.length" class="review-images">
              <van-image
                v-for="img in review.images"
                :key="img.id"
                :src="img.image_url"
                width="60"
                height="60"
                fit="cover"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div class="footer">
      <van-button type="primary" block @click="addToCart">加入购物车</van-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { showToast } from 'vant'
import { getDishDetail } from '@/api/dish'
import { getDishReviews } from '@/api/review'
import { useCartStore } from '@/stores/cart'

const route = useRoute()
const cartStore = useCartStore()

const dish = ref(null)
const reviews = ref([])

const fetchData = async () => {
  try {
    const [dishRes, reviewRes] = await Promise.all([
      getDishDetail(route.params.id),
      getDishReviews(route.params.id, { page: 1, page_size: 10 })
    ])
    
    dish.value = dishRes.data
    reviews.value = reviewRes.data || []
  } catch (error) {
    console.error('获取详情失败', error)
  }
}

const addToCart = () => {
  if (dish.value) {
    cartStore.addItem(dish.value)
    showToast('已添加到购物车')
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.dish-detail {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 60px;
}

.dish-image {
  width: 100%;
  height: 300px;
  object-fit: cover;
}

.dish-info {
  background: #fff;
  padding: 16px;
  margin-bottom: 10px;
  
  h2 {
    font-size: 20px;
    margin-bottom: 8px;
  }
  
  .price {
    color: var(--primary-color);
    font-size: 24px;
    font-weight: 600;
    margin-bottom: 12px;
  }
  
  .description {
    color: var(--text-color-2);
    line-height: 1.6;
  }
}

.reviews-section {
  background: #fff;
  padding: 16px;
  
  h3 {
    font-size: 16px;
    margin-bottom: 16px;
  }
}

.review-item {
  padding: 12px 0;
  border-bottom: 1px solid var(--border-color);
  
  &:last-child {
    border-bottom: none;
  }
  
  .review-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
    
    .username {
      font-size: 14px;
      font-weight: 500;
    }
  }
  
  .review-content {
    font-size: 14px;
    line-height: 1.6;
    margin-bottom: 8px;
  }
  
  .review-images {
    display: flex;
    gap: 8px;
  }
}

.footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 10px 16px;
  background: #fff;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.1);
}
</style>
