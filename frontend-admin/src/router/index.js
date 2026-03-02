import { createRouter, createWebHistory } from 'vue-router'
import { useAdminStore } from '@/stores/admin'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/layout/Layout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '仪表盘', icon: 'Odometer' }
      },
      {
        path: 'dishes',
        name: 'Dishes',
        component: () => import('@/views/Dishes.vue'),
        meta: { title: '菜品管理', icon: 'Food' }
      },
      {
        path: 'categories',
        name: 'Categories',
        component: () => import('@/views/Categories.vue'),
        meta: { title: '分类管理', icon: 'Menu' }
      },
      {
        path: 'weekly-menus',
        name: 'WeeklyMenus',
        component: () => import('@/views/WeeklyMenus.vue'),
        meta: { title: '菜谱管理', icon: 'Calendar' }
      },
      {
        path: 'menu-ratings',
        name: 'MenuRatings',
        component: () => import('@/views/MenuRatings.vue'),
        meta: { title: '菜谱评价', icon: 'ChatDotRound' }
      },
      {
        path: 'orders',
        name: 'Orders',
        component: () => import('@/views/Orders.vue'),
        meta: { title: '订单管理', icon: 'Document' }
      },
      {
        path: 'reviews',
        name: 'Reviews',
        component: () => import('@/views/Reviews.vue'),
        meta: { title: '评价管理', icon: 'Star' }
      },
      {
        path: 'announcements',
        name: 'Announcements',
        component: () => import('@/views/Announcements.vue'),
        meta: { title: '公告管理', icon: 'Bell' }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/Users.vue'),
        meta: { title: '用户管理', icon: 'User' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory('/admin/'),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const adminStore = useAdminStore()
  
  if (to.meta.requiresAuth && !adminStore.token) {
    next('/login')
  } else if (to.path === '/login' && adminStore.token) {
    next('/')
  } else {
    next()
  }
})

export default router