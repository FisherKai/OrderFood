<template>
  <div class="home">
    <!-- 页面标题 -->
    <van-nav-bar title="点餐系统" />

    <!-- 功能入口 -->
    <div class="function-entries">
      <div class="entry-card" @click="$router.push('/weekly-menu')">
        <div class="entry-icon">
          <van-icon name="calendar-o" size="48" color="#1989fa" />
        </div>
        <div class="entry-info">
          <div class="entry-title">本周菜谱</div>
          <div class="entry-desc">查看今日菜品，选菜下单</div>
        </div>
        <van-icon name="arrow" color="#c8c9cc" />
      </div>

      <div class="entry-card" @click="$router.push('/reserve')">
        <div class="entry-icon">
          <van-icon name="clock-o" size="48" color="#07c160" />
        </div>
        <div class="entry-info">
          <div class="entry-title">预定点餐</div>
          <div class="entry-desc">提前预约，按时用餐</div>
        </div>
        <van-icon name="arrow" color="#c8c9cc" />
      </div>
    </div>

    <!-- 底部导航 -->
    <van-tabbar v-model="active" route>
      <van-tabbar-item to="/home" icon="home-o">首页</van-tabbar-item>
      <van-tabbar-item to="/cart" icon="shopping-cart-o" :badge="cartCount">购物车</van-tabbar-item>
      <van-tabbar-item to="/orders" icon="orders-o">订单</van-tabbar-item>
      <van-tabbar-item to="/profile" icon="user-o">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useCartStore } from '@/stores/cart'

const cartStore = useCartStore()
const active = ref(0)
const cartCount = computed(() => cartStore.totalCount || '')
</script>

<style scoped lang="scss">
.home {
  min-height: 100vh;
  background: #f7f8fa;
  padding-bottom: 50px;
}

.function-entries {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;

  .entry-card {
    display: flex;
    align-items: center;
    background: #fff;
    border-radius: 12px;
    padding: 20px 16px;
    cursor: pointer;
    transition: box-shadow 0.2s;

    &:active {
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
    }

    .entry-icon {
      width: 64px;
      height: 64px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f5f7fa;
      border-radius: 12px;
      margin-right: 16px;
      flex-shrink: 0;
    }

    .entry-info {
      flex: 1;

      .entry-title {
        font-size: 18px;
        font-weight: 600;
        color: #323233;
        margin-bottom: 6px;
      }

      .entry-desc {
        font-size: 13px;
        color: #969799;
      }
    }
  }
}
</style>
