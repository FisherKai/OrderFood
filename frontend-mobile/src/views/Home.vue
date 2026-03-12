<template>
  <div class="home">
    <!-- 页面标题 -->
    <van-nav-bar title="健康食堂">
      <template #right>
        <div v-if="!isLoggedIn" class="nav-login-btn" @click="$router.push('/login')">
          <van-icon name="user-o" size="16" />
          <span>登录</span>
        </div>
        <div v-else class="nav-user-info">
          <van-icon name="user-circle-o" size="18" color="#1989fa" />
          <span>{{ userStore.userInfo.nickname || userStore.userInfo.username || '用户' }}</span>
        </div>
      </template>
    </van-nav-bar>

    <!-- 未登录提示横幅 -->
    <div v-if="!isLoggedIn" class="login-banner" @click="$router.push('/login')">
      <div class="banner-left">
        <van-icon name="info-o" size="18" color="#1989fa" />
        <span>登录后即可使用点餐、下单等功能</span>
      </div>
      <div class="banner-right">
        <span>去登录</span>
        <van-icon name="arrow" size="14" />
      </div>
    </div>

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
import { useUserStore } from '@/stores/user'

const cartStore = useCartStore()
const userStore = useUserStore()
const active = ref(0)
const cartCount = computed(() => cartStore.totalCount || '')
const isLoggedIn = computed(() => !!userStore.token)
</script>

<style scoped lang="scss">
.home {
  min-height: 100vh;
  background: #f7f8fa;
  padding-bottom: 50px;
}

/* 导航栏右侧登录按钮 */
.nav-login-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 12px;
  background: #1989fa;
  color: #fff;
  border-radius: 16px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;

  &:active {
    opacity: 0.85;
  }
}

/* 导航栏右侧已登录用户信息 */
.nav-user-info {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #323233;
}

/* 未登录提示横幅 */
.login-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 12px 16px 0;
  padding: 12px 16px;
  background: linear-gradient(135deg, #e8f4fd, #f0f7ff);
  border-radius: 10px;
  cursor: pointer;
  transition: opacity 0.2s;

  &:active {
    opacity: 0.85;
  }

  .banner-left {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    color: #555;
  }

  .banner-right {
    display: flex;
    align-items: center;
    gap: 2px;
    font-size: 13px;
    color: #1989fa;
    font-weight: 500;
  }
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
