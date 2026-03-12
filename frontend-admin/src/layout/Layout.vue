<template>
  <el-container class="layout-container" :class="{ 'is-mobile': isMobile }">
    <!-- 移动端遮罩层 -->
    <div 
      v-if="isMobile && sidebarVisible" 
      class="sidebar-mask" 
      @click="closeSidebar"
    ></div>
    
    <!-- 侧边栏 -->
    <el-aside 
      :width="sidebarWidth" 
      class="sidebar"
      :class="{ 'sidebar-mobile': isMobile, 'sidebar-visible': sidebarVisible }"
    >
      <div class="logo" @click="isMobile && closeSidebar()">
        <el-icon class="logo-icon"><Food /></el-icon>
        <span v-if="!isCollapse || isMobile">健康食堂</span>
      </div>
      
      <el-menu
        :default-active="$route.path"
        class="sidebar-menu"
        :collapse="!isMobile && isCollapse"
        router
        background-color="transparent"
        text-color="#595959"
        active-text-color="#1677ff"
        @select="handleMenuSelect"
      >
        <el-menu-item
          v-for="route in menuRoutes"
          :key="route.path"
          :index="route.path"
        >
          <el-icon><component :is="route.meta.icon" /></el-icon>
          <template #title>{{ route.meta.title }}</template>
        </el-menu-item>
      </el-menu>
      
      <!-- 移动端底部退出按钮 -->
      <div v-if="isMobile" class="mobile-logout">
        <el-button type="danger" plain @click="handleLogout" class="logout-btn">
          <el-icon><SwitchButton /></el-icon>
          退出登录
        </el-button>
      </div>
    </el-aside>
    
    <!-- 主内容区 -->
    <el-container class="main-container">
      <!-- 顶部导航 -->
      <el-header class="header" :class="{ 'header-mobile': isMobile }">
        <div class="header-left">
          <el-button
            :type="isMobile ? 'primary' : 'text'"
            :circle="isMobile"
            @click="toggleSidebar"
            class="menu-btn"
          >
            <el-icon :size="isMobile ? 20 : 18">
              <Operation v-if="isMobile" />
              <Expand v-else-if="isCollapse" />
              <Fold v-else />
            </el-icon>
          </el-button>
          
          <div class="page-title" v-if="isMobile">
            {{ $route.meta.title || '仪表盘' }}
          </div>
          
          <el-breadcrumb v-else separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="$route.meta.title">
              {{ $route.meta.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="header-right">
          <!-- 移动端只显示头像 -->
          <template v-if="isMobile">
            <el-avatar :size="32" class="mobile-avatar">
              <el-icon><Avatar /></el-icon>
            </el-avatar>
          </template>
          
          <!-- 桌面端显示完整下拉菜单 -->
          <el-dropdown v-else @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32">
                <el-icon><Avatar /></el-icon>
              </el-avatar>
              <span class="username">{{ adminStore.adminInfo.username || '管理员' }}</span>
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <!-- 主内容 -->
      <el-main class="main-content" :class="{ 'main-mobile': isMobile }">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import {
  Food, Expand, Fold, ArrowDown, SwitchButton, Avatar, Operation
} from '@element-plus/icons-vue'
import { useAdminStore } from '@/stores/admin'
import { useDevice } from '@/composables/useDevice'

const router = useRouter()
const route = useRoute()
const adminStore = useAdminStore()
const { isMobile, isTablet, screenWidth } = useDevice()

const isCollapse = ref(false)
const sidebarVisible = ref(false)

// 侧边栏宽度
const sidebarWidth = computed(() => {
  if (isMobile.value) {
    return '280px'
  }
  return isCollapse.value ? '64px' : '240px'
})

// 获取菜单路由
const menuRoutes = computed(() => {
  return router.getRoutes()
    .find(r => r.path === '/')
    ?.children?.filter(child => child.meta?.title) || []
})

// 切换侧边栏
const toggleSidebar = () => {
  if (isMobile.value) {
    sidebarVisible.value = !sidebarVisible.value
  } else {
    isCollapse.value = !isCollapse.value
  }
}

// 关闭侧边栏
const closeSidebar = () => {
  sidebarVisible.value = false
}

// 菜单选择处理
const handleMenuSelect = () => {
  if (isMobile.value) {
    closeSidebar()
  }
}

// 处理命令
const handleCommand = async (command) => {
  if (command === 'logout') {
    await handleLogout()
  }
}

// 退出登录
const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确认退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    adminStore.logout()
    ElMessage.success('退出登录成功')
    router.push('/login')
  } catch {
    // 用户取消
  }
}

// 监听路由变化，移动端自动关闭侧边栏
watch(route, () => {
  if (isMobile.value) {
    closeSidebar()
  }
})

// 监听设备变化
watch(isMobile, (newVal) => {
  if (!newVal) {
    sidebarVisible.value = false
  }
})
</script>

<style scoped lang="scss">
.layout-container {
  height: 100vh;
  position: relative;
  
  // 遮罩层
  .sidebar-mask {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.3);
    z-index: 998;
    animation: fadeIn 0.25s ease;
    backdrop-filter: blur(2px);
  }
  
  .sidebar {
    background-color: #fff;
    border-right: 1px solid #f0f0f0;
    transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    overflow: hidden;
    
    .logo {
      height: 56px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #1a1a2e;
      font-size: 17px;
      font-weight: 700;
      letter-spacing: 0.5px;
      border-bottom: 1px solid #f0f0f0;
      white-space: nowrap;
      
      .logo-icon {
        font-size: 22px;
        margin-right: 8px;
        color: #1677ff;
        flex-shrink: 0;
      }
      
      .el-icon {
        font-size: 22px;
      }
    }
    
    .sidebar-menu {
      border: none;
      height: calc(100vh - 56px);
      overflow-y: auto;
      padding: 8px;
      
      &::-webkit-scrollbar {
        width: 4px;
      }
      
      &::-webkit-scrollbar-thumb {
        background-color: rgba(0, 0, 0, 0.08);
        border-radius: 4px;
      }
      
      :deep(.el-menu-item) {
        border-radius: 8px;
        margin-bottom: 2px;
        height: 44px;
        line-height: 44px;
        font-size: 14px;
        transition: all 0.2s ease;
        
        .el-icon {
          font-size: 18px;
        }
        
        &:hover {
          background-color: #f5f7fa;
          color: #1677ff;
        }
        
        &.is-active {
          background-color: #e8f4ff !important;
          color: #1677ff !important;
          font-weight: 500;
          
          .el-icon {
            color: #1677ff !important;
          }
          
          span {
            color: #1677ff !important;
          }
        }
      }
    }
    
    // 移动端侧边栏样式
    &.sidebar-mobile {
      position: fixed;
      top: 0;
      left: -280px;
      height: 100vh;
      z-index: 999;
      transition: left 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      box-shadow: none;
      
      &.sidebar-visible {
        left: 0;
        box-shadow: 4px 0 16px rgba(0, 0, 0, 0.08);
      }
      
      .sidebar-menu {
        height: calc(100vh - 56px - 72px);
      }
    }
    
    .mobile-logout {
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      padding: 12px 16px;
      background-color: #fafafa;
      border-top: 1px solid #f0f0f0;
      
      .logout-btn {
        width: 100%;
        height: 40px;
        font-size: 14px;
        border-radius: 8px;
      }
    }
  }
  
  .main-container {
    flex: 1;
    overflow: hidden;
  }
  
  .header {
    background-color: #fff;
    border-bottom: 1px solid #f0f0f0;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 24px;
    height: 56px;
    
    .header-left {
      display: flex;
      align-items: center;
      
      .menu-btn {
        margin-right: 16px;
        font-size: 18px;
        
        &:not(.el-button--primary) {
          color: #8c8c8c;
          
          &:hover {
            color: #1677ff;
          }
        }
      }
      
      .page-title {
        font-size: 16px;
        font-weight: 600;
        color: #1a1a2e;
      }
      
      .el-breadcrumb {
        font-size: 14px;
        
        :deep(.el-breadcrumb__inner) {
          color: #8c8c8c;
          
          &.is-link:hover {
            color: #1677ff;
          }
        }
        
        :deep(.el-breadcrumb__item:last-child .el-breadcrumb__inner) {
          color: #1a1a2e;
          font-weight: 500;
        }
      }
    }
    
    .header-right {
      .user-info {
        display: flex;
        align-items: center;
        cursor: pointer;
        padding: 4px 12px;
        border-radius: 20px;
        transition: background-color 0.2s;
        
        .username {
          margin: 0 8px;
          color: #595959;
          font-size: 14px;
        }
        
        &:hover {
          background-color: #f5f7fa;
          
          .username {
            color: #1677ff;
          }
        }
      }
      
      .mobile-avatar {
        cursor: pointer;
      }
    }
    
    // 移动端头部样式
    &.header-mobile {
      padding: 0 12px;
      height: 52px;
      
      .header-left {
        .menu-btn {
          margin-right: 12px;
        }
      }
    }
  }
  
  .main-content {
    padding: 0;
    background-color: #f5f7fa;
    overflow-y: auto;
    
    &.main-mobile {
      padding: 0;
    }
  }
  
  // 移动端整体布局
  &.is-mobile {
    .main-container {
      width: 100%;
    }
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
</style>
