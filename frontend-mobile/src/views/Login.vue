<template>
  <div class="login-page">
    <van-nav-bar title="登录" left-arrow @click-left="$router.back()" />
    
    <div class="login-content">
      <div class="logo">
        <van-icon name="shop-o" size="64" color="#ff6034" />
        <h1>点餐系统</h1>
      </div>
      
      <van-form @submit="handleLogin">
        <van-cell-group inset>
          <van-field
            v-model="form.username"
            name="username"
            label="用户名"
            placeholder="请输入用户名"
            :rules="[{ required: true, message: '请输入用户名' }]"
          />
          <van-field
            v-model="form.password"
            type="password"
            name="password"
            label="密码"
            placeholder="请输入密码"
            :rules="[{ required: true, message: '请输入密码' }]"
          />
        </van-cell-group>
        
        <div class="button-group">
          <van-button round block type="primary" native-type="submit" :loading="loading">
            登录
          </van-button>
          <van-button round block plain type="primary" @click="$router.push('/register')">
            注册账号
          </van-button>
        </div>
      </van-form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  username: '',
  password: ''
})

const loading = ref(false)

const handleLogin = async () => {
  loading.value = true
  try {
    await userStore.loginUser(form.value.username, form.value.password)
    showToast('登录成功')
    router.replace('/home')
  } catch (error) {
    console.error('登录失败', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.login-page {
  min-height: 100vh;
  background: #f5f5f5;
}

.login-content {
  padding: 40px 20px;
  
  .logo {
    text-align: center;
    margin-bottom: 40px;
    
    h1 {
      margin-top: 16px;
      font-size: 24px;
      color: var(--text-color);
    }
  }
  
  .button-group {
    padding: 20px 16px;
    
    .van-button {
      margin-bottom: 12px;
    }
  }
}
</style>
