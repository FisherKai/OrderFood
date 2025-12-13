<template>
  <div class="profile-page">
    <van-nav-bar title="个人中心" />
    
    <div class="user-info">
      <van-image
        round
        width="60"
        height="60"
        :src="userStore.userInfo.avatar || 'https://via.placeholder.com/60'"
      />
      <div class="info">
        <div class="nickname">{{ userStore.userInfo.nickname || userStore.userInfo.username }}</div>
        <div class="username">@{{ userStore.userInfo.username }}</div>
      </div>
    </div>
    
    <van-cell-group inset>
      <van-cell title="我的订单" is-link to="/orders" />
      <van-cell title="预约点餐" is-link to="/reserve" />
    </van-cell-group>
    
    <div class="logout-button">
      <van-button round block type="danger" @click="handleLogout">退出登录</van-button>
    </div>
    
    <van-tabbar v-model="active" route>
      <van-tabbar-item to="/home" icon="home-o">首页</van-tabbar-item>
      <van-tabbar-item to="/dishes" icon="apps-o">菜品</van-tabbar-item>
      <van-tabbar-item to="/cart" icon="shopping-cart-o">购物车</van-tabbar-item>
      <van-tabbar-item to="/profile" icon="user-o">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { showConfirmDialog } from 'vant'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const active = ref(3)

const handleLogout = () => {
  showConfirmDialog({
    message: '确定要退出登录吗？'
  }).then(() => {
    userStore.logout()
    router.replace('/login')
  }).catch(() => {})
}
</script>

<style scoped lang="scss">
.profile-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 50px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px 16px;
  background: #fff;
  margin-bottom: 10px;
  
  .info {
    flex: 1;
    
    .nickname {
      font-size: 18px;
      font-weight: 600;
      margin-bottom: 4px;
    }
    
    .username {
      font-size: 14px;
      color: var(--text-color-2);
    }
  }
}

.logout-button {
  padding: 20px 16px;
}
</style>
