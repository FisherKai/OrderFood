<template>
  <div class="orders-page">
    <van-nav-bar title="我的订单" />
    
    <van-list
      v-model:loading="loading"
      :finished="finished"
      finished-text="没有更多了"
      @load="loadOrders"
    >
      <div class="orders-list">
        <div v-for="order in orders" :key="order.id" class="order-item">
          <div class="order-header">
            <span>订单号: {{ order.id }}</span>
            <van-tag :type="getStatusType(order.status)">{{ getStatusText(order.status) }}</van-tag>
          </div>
          
          <div class="order-dishes">
            <div v-for="item in order.items" :key="item.id" class="dish-item">
              <img :src="item.dish?.images?.[0]?.image_url" class="dish-image" />
              <div class="dish-info">
                <div class="dish-name">{{ item.dish?.name }}</div>
                <div class="dish-price">¥{{ item.price }} x{{ item.quantity }}</div>
              </div>
            </div>
          </div>
          
          <div class="order-footer">
            <span v-if="order.order_type === 2" class="reserve-time">
              预约时间: {{ formatDate(order.reserve_time) }}
            </span>
            <span class="total-price">合计: ¥{{ order.total_price }}</span>
          </div>
        </div>
      </div>
    </van-list>
    
    <van-empty v-if="!loading && !orders.length" description="暂无订单" />
    
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
import { getUserOrders } from '@/api/order'

const active = ref(3)
const orders = ref([])
const loading = ref(false)
const finished = ref(false)
const page = ref(1)

const loadOrders = async () => {
  try {
    const res = await getUserOrders({ page: page.value, page_size: 20 })
    
    if (page.value === 1) {
      orders.value = res.data || []
    } else {
      orders.value.push(...(res.data || []))
    }
    
    loading.value = false
    
    if (orders.value.length >= res.pagination.total) {
      finished.value = true
    } else {
      page.value++
    }
  } catch (error) {
    loading.value = false
    console.error('获取订单失败', error)
  }
}

const getStatusType = (status) => {
  const types = ['', 'warning', 'primary', 'success', 'default']
  return types[status] || 'default'
}

const getStatusText = (status) => {
  const texts = ['', '待处理', '制作中', '已完成', '已取消']
  return texts[status] || '未知'
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}月${date.getDate()}日 ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
}
</script>

<style scoped lang="scss">
.orders-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 50px;
}

.orders-list {
  padding: 10px;
}

.order-item {
  background: #fff;
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 10px;
  
  .order-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: 12px;
    border-bottom: 1px solid var(--border-color);
    margin-bottom: 12px;
    font-size: 14px;
  }
  
  .order-dishes {
    .dish-item {
      display: flex;
      gap: 12px;
      margin-bottom: 8px;
      
      .dish-image {
        width: 60px;
        height: 60px;
        border-radius: 4px;
        object-fit: cover;
      }
      
      .dish-info {
        flex: 1;
        
        .dish-name {
          font-size: 14px;
          margin-bottom: 4px;
        }
        
        .dish-price {
          font-size: 12px;
          color: var(--text-color-2);
        }
      }
    }
  }
  
  .order-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 12px;
    border-top: 1px solid var(--border-color);
    margin-top: 12px;
    
    .reserve-time {
      font-size: 12px;
      color: var(--text-color-3);
    }
    
    .total-price {
      font-size: 16px;
      font-weight: 600;
      color: var(--primary-color);
    }
  }
}
</style>
