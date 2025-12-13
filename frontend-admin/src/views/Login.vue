<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <h1>点餐系统管理后台</h1>
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  
  .login-box {
    width: 400px;
    background: rgba(255, 255, 255, 0.95);
    border-radius: 10px;
    padding: 40px;
    box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
    backdrop-filter: blur(10px);
    
    .login-header {
      text-align: center;
      margin-bottom: 40px;
      
      h1 {
        color: #606266;
        font-size: 28px;
        font-weight: bold;
        margin-bottom: 10px;
      }
      
      p {
        color: #909399;
        font-size: 14px;
      }
    }
    
    .login-form {
      .el-form-item {
        margin-bottom: 25px;
      }
      
      :deep(.el-input) {
        .el-input__wrapper {
          border-radius: 8px;
          box-shadow: 0 0 0 1px #dcdfe6 inset;
          
          &:hover {
            box-shadow: 0 0 0 1px #c0c4cc inset;
          }
          
          &.is-focus {
            box-shadow: 0 0 0 1px #409eff inset;
          }
        }
      }
      
      .el-button {
        border-radius: 8px;
        font-size: 16px;
        height: 48px;
      }
    }
    
    .login-footer {
      text-align: center;
      margin-top: 30px;
      
      p {
        color: #909399;
        font-size: 12px;
        background-color: #f5f7fa;
        padding: 10px;
        border-radius: 4px;
        border: 1px dashed #e4e7ed;
      }
    }
  }
}

@media (max-width: 480px) {
  .login-container {
    padding: 20px;
    
    .login-box {
      width: 100%;
      padding: 30px 20px;
    }
  }
}
</style>