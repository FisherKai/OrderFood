<template>
  <div class="cart-page">
    <van-nav-bar title="购物车" />
    
    <div v-if="cartStore.items.length" class="cart-content">
      <van-checkbox-group v-model="checkedItems">
        <van-swipe-cell v-for="item in cartStore.items" :key="item.id">
          <van-card
            :price="item.price"
            :title="item.name"
            :thumb="item.image"
          >
            <template #num>
              <van-stepper v-model="item.quantity" @change="updateQuantity(item.id, item.quantity)" />
            </template>
            <template #footer>
              <van-checkbox :name="item.id" />
            </template>
          </van-card>
          <template #right>
            <van-button square type="danger" text="删除" @click="removeItem(item.id)" />
          </template>
        </van-swipe-cell>
      </van-checkbox-group>
    </div>
    
    <van-empty v-else description="购物车是空的" />
    
    <div v-if="cartStore.items.length" class="cart-footer">
      <van-checkbox v-model="checkAll" @change="toggleCheckAll">全选</van-checkbox>
      <div class="total-price">
        合计: <span class="price">¥{{ totalPrice }}</span>
      </div>
      <van-button type="primary" @click="checkout">结算</van-button>
    </div>
    
    <van-tabbar v-model="active" route>
      <van-tabbar-item to="/home" icon="home-o">首页</van-tabbar-item>
      <van-tabbar-item to="/cart" icon="shopping-cart-o" :badge="cartStore.totalCount">购物车</van-tabbar-item>
      <van-tabbar-item to="/orders" icon="orders-o">订单</van-tabbar-item>
      <van-tabbar-item to="/profile" icon="user-o">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { useCartStore } from '@/stores/cart'
import { createOrder } from '@/api/order'

const router = useRouter()
const cartStore = useCartStore()
const active = ref(1)

const checkedItems = ref([])
const checkAll = ref(false)

const totalPrice = computed(() => {
  return cartStore.items
    .filter(item => checkedItems.value.includes(item.id))
    .reduce((sum, item) => sum + item.price * item.quantity, 0)
    .toFixed(2)
})

const toggleCheckAll = (val) => {
  if (val) {
    checkedItems.value = cartStore.items.map(item => item.id)
  } else {
    checkedItems.value = []
  }
}

const updateQuantity = (id, quantity) => {
  cartStore.updateQuantity(id, quantity)
}

const removeItem = (id) => {
  showConfirmDialog({
    message: '确定要删除这个菜品吗？'
  }).then(() => {
    cartStore.removeItem(id)
    showToast('已删除')
  }).catch(() => {})
}

const checkout = async () => {
  if (!checkedItems.value.length) {
    showToast('请选择要结算的菜品')
    return
  }
  
  try {
    const items = cartStore.items
      .filter(item => checkedItems.value.includes(item.id))
      .map(item => ({
        dish_id: item.id,
        quantity: item.quantity
      }))
    
    await createOrder({ items })
    showToast('下单成功')
    
    // 清空已结算的商品
    checkedItems.value.forEach(id => cartStore.removeItem(id))
    checkedItems.value = []
    
    router.push('/orders')
  } catch (error) {
    console.error('下单失败', error)
  }
}
</script>

<style scoped lang="scss">
.cart-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 100px;
}

.cart-content {
  padding: 10px;
}

.cart-footer {
  position: fixed;
  bottom: 50px;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: #fff;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.1);
  
  .total-price {
    flex: 1;
    text-align: right;
    
    .price {
      color: var(--primary-color);
      font-size: 18px;
      font-weight: 600;
    }
  }
}
</style>
