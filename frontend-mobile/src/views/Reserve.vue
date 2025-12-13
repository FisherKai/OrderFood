<template>
  <div class="reserve-page">
    <van-nav-bar title="预约点餐" left-arrow @click-left="$router.back()" />
    
    <van-form @submit="handleSubmit">
      <van-cell-group inset>
        <van-field
          v-model="reserveTime"
          is-link
          readonly
          label="预约时间"
          placeholder="选择预约时间"
          @click="showTimePicker = true"
          :rules="[{ required: true, message: '请选择预约时间' }]"
        />
        <van-field
          v-model="peopleCount"
          type="digit"
          label="用餐人数"
          placeholder="请输入用餐人数"
          :rules="[{ required: true, message: '请输入用餐人数' }]"
        />
      </van-cell-group>
      
      <div class="selected-dishes">
        <h3>已选菜品</h3>
        <van-empty v-if="!cartStore.items.length" description="请先添加菜品" />
        <div v-else class="dishes-list">
          <van-card
            v-for="item in cartStore.items"
            :key="item.id"
            :price="item.price"
            :title="item.name"
            :thumb="item.image"
            :num="item.quantity"
          />
        </div>
      </div>
      
      <div class="total-price">
        合计: <span class="price">¥{{ cartStore.totalPrice }}</span>
      </div>
      
      <div class="submit-button">
        <van-button round block type="primary" native-type="submit">
          提交预约
        </van-button>
      </div>
    </van-form>
    
    <van-popup v-model:show="showTimePicker" position="bottom">
      <van-datetime-picker
        v-model="currentDate"
        type="datetime"
        :min-date="minDate"
        @confirm="confirmTime"
        @cancel="showTimePicker = false"
      />
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { useCartStore } from '@/stores/cart'
import { createReservation } from '@/api/order'

const router = useRouter()
const cartStore = useCartStore()

const reserveTime = ref('')
const peopleCount = ref('')
const showTimePicker = ref(false)
const currentDate = ref(new Date())
const minDate = computed(() => {
  const date = new Date()
  date.setHours(date.getHours() + 2)
  return date
})

const confirmTime = (value) => {
  const date = new Date(value)
  reserveTime.value = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
  showTimePicker.value = false
}

const handleSubmit = async () => {
  if (!cartStore.items.length) {
    showToast('请先添加菜品')
    return
  }
  
  try {
    const items = cartStore.items.map(item => ({
      dish_id: item.id,
      quantity: item.quantity
    }))
    
    await createReservation({
      items,
      reserve_time: reserveTime.value,
      people_count: parseInt(peopleCount.value)
    })
    
    showToast('预约成功')
    cartStore.clearCart()
    router.push('/orders')
  } catch (error) {
    console.error('预约失败', error)
  }
}
</script>

<style scoped lang="scss">
.reserve-page {
  min-height: 100vh;
  background: #f5f5f5;
  padding-bottom: 20px;
}

.selected-dishes {
  margin: 20px 16px;
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  
  h3 {
    font-size: 16px;
    margin-bottom: 12px;
  }
}

.total-price {
  text-align: right;
  padding: 16px;
  font-size: 18px;
  
  .price {
    color: var(--primary-color);
    font-weight: 600;
  }
}

.submit-button {
  padding: 0 16px;
}
</style>
