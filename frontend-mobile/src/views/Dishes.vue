<template>
  <div class="dishes-page">
    <van-nav-bar title="菜品列表" left-arrow @click-left="$router.back()" />
    
    <!-- 搜索栏 -->
    <van-search v-model="searchKeyword" placeholder="搜索菜品" @search="handleSearch" />
    
    <!-- 分类筛选 -->
    <van-tabs v-model:active="activeCategory" @change="handleCategoryChange">
      <van-tab title="全部" name="0" />
      <van-tab v-for="cat in categories" :key="cat.id" :title="cat.name" :name="String(cat.id)" />
    </van-tabs>
    
    <!-- 菜品列表 -->
    <van-list
      v-model:loading="loading"
      :finished="finished"
      finished-text="没有更多了"
      @load="loadDishes"
    >
      <div class="dishes-list">
        <van-card
          v-for="dish in dishes"
          :key="dish.id"
          :price="dish.price"
          :title="dish.name"
          :desc="dish.description"
          :thumb="dish.images?.[0]?.image_url"
          @click="$router.push(`/dish/${dish.id}`)"
        >
          <template #footer>
            <van-button size="small" type="primary" @click.stop="addToCart(dish)">
              加入购物车
            </van-button>
          </template>
        </van-card>
      </div>
    </van-list>
    
    <van-tabbar v-model="tabActive" route>
      <van-tabbar-item to="/home" icon="home-o">首页</van-tabbar-item>
      <van-tabbar-item to="/dishes" icon="apps-o">菜品</van-tabbar-item>
      <van-tabbar-item to="/cart" icon="shopping-cart-o" :badge="cartCount">购物车</van-tabbar-item>
      <van-tabbar-item to="/profile" icon="user-o">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { showToast } from 'vant'
import { getDishes, getCategories } from '@/api/dish'
import { useCartStore } from '@/stores/cart'

const route = useRoute()
const cartStore = useCartStore()

const tabActive = ref(1)
const searchKeyword = ref('')
const activeCategory = ref('0')
const categories = ref([])
const dishes = ref([])
const loading = ref(false)
const finished = ref(false)
const page = ref(1)

const cartCount = computed(() => cartStore.totalCount)

const fetchCategories = async () => {
  try {
    const res = await getCategories()
    categories.value = res.data || []
    
    if (route.query.category_id) {
      activeCategory.value = String(route.query.category_id)
    }
  } catch (error) {
    console.error('获取分类失败', error)
  }
}

const loadDishes = async () => {
  try {
    const params = {
      page: page.value,
      page_size: 20
    }
    
    if (activeCategory.value !== '0') {
      params.category_id = activeCategory.value
    }
    
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }
    
    const res = await getDishes(params)
    
    if (page.value === 1) {
      dishes.value = res.data || []
    } else {
      dishes.value.push(...(res.data || []))
    }
    
    loading.value = false
    
    if (dishes.value.length >= res.pagination.total) {
      finished.value = true
    } else {
      page.value++
    }
  } catch (error) {
    loading.value = false
    console.error('获取菜品失败', error)
  }
}

const handleCategoryChange = () => {
  page.value = 1
  finished.value = false
  dishes.value = []
  loadDishes()
}

const handleSearch = () => {
  page.value = 1
  finished.value = false
  dishes.value = []
  loadDishes()
}

const addToCart = (dish) => {
  cartStore.addItem(dish)
  showToast('已添加到购物车')
}

onMounted(() => {
  fetchCategories()
})
</script>

<style scoped lang="scss">
.dishes-page {
  padding-bottom: 50px;
  background: #f5f5f5;
  min-height: 100vh;
}

.dishes-list {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
</style>
