<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <h1>健康食堂管理后台</h1>
        <p>欢迎管理员登录</p>
      </div>
      
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            size="large"
            clearable
          >
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            show-password
            clearable
          >
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleLogin"
            style="width: 100%"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>
      
      <div class="login-footer">
        <p>默认账号：admin / admin123</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import { useAdminStore } from '@/stores/admin'
import { authAPI } from '@/api/auth'

const router = useRouter()
const adminStore = useAdminStore()

const loginFormRef = ref()
const loading = ref(false)

const loginForm = reactive({
  username: 'admin',
  password: 'admin123'
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    loading.value = true
    
    const response = await authAPI.login(loginForm)
    
    if (response.token) {
      adminStore.login(response.token, response.admin || { username: loginForm.username })
      ElMessage.success('登录成功')
      router.push('/')
    } else {
      ElMessage.error('登录失败，请检查用户名和密码')
    }
  } catch (error) {
    if (error.response) {
      ElMessage.error(error.response.data?.error || '登录失败')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.login-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f2f5;
  position: relative;
  
  // 微妙的背景图案
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-image: 
      radial-gradient(circle at 25% 25%, rgba(22, 119, 255, 0.03) 0%, transparent 50%),
      radial-gradient(circle at 75% 75%, rgba(82, 196, 26, 0.03) 0%, transparent 50%);
  }
  
  .login-box {
    width: 380px;
    background: #fff;
    border-radius: 16px;
    padding: 40px 36px;
    border: 1px solid #f0f0f0;
    position: relative;
    z-index: 1;
    
    .login-header {
      text-align: center;
      margin-bottom: 36px;
      
      h1 {
        color: #1a1a2e;
        font-size: 22px;
        font-weight: 700;
        margin-bottom: 8px;
        letter-spacing: 0.5px;
      }
      
      p {
        color: #8c8c8c;
        font-size: 14px;
        font-weight: 400;
      }
    }
    
    .login-form {
      .el-form-item {
        margin-bottom: 22px;
      }
      
      :deep(.el-input) {
        .el-input__wrapper {
          border-radius: 10px;
          box-shadow: 0 0 0 1px #d9d9d9 inset;
          padding: 4px 12px;
          transition: all 0.2s ease;
          
          &:hover {
            box-shadow: 0 0 0 1px #1677ff inset;
          }
          
          &.is-focus {
            box-shadow: 0 0 0 1px #1677ff inset, 0 0 0 3px rgba(22, 119, 255, 0.08);
          }
        }
        
        .el-input__prefix {
          color: #bfbfbf;
        }
      }
      
      .el-button {
        border-radius: 10px;
        font-size: 15px;
        font-weight: 600;
        height: 46px;
        letter-spacing: 0.5px;
        background: #1677ff;
        border-color: #1677ff;
        
        &:hover {
          background: lighten(#1677ff, 8%);
          border-color: lighten(#1677ff, 8%);
        }
      }
    }
    
    .login-footer {
      text-align: center;
      margin-top: 24px;
      
      p {
        color: #8c8c8c;
        font-size: 12px;
        background-color: #fafafa;
        padding: 10px 14px;
        border-radius: 8px;
        border: 1px solid #f0f0f0;
      }
    }
  }
}

@media (max-width: 480px) {
  .login-container {
    padding: 20px;
    
    .login-box {
      width: 100%;
      padding: 32px 24px;
    }
  }
}
</style>