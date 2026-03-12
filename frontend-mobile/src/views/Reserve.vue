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
    
    <van-popup v-model:show="showTimePicker" position="bottom" round>
      <van-picker-group
        title="选择预约时间"
        :tabs="['选择日期', '选择时间']"
        @confirm="onPickerConfirm"
        @cancel="showTimePicker = false"
      >
        <van-date-picker
          v-model="pickerDate"
          :min-date="minDate"
          :max-date="maxDate"
        />
        <van-time-picker
          v-model="pickerTime"
          :min-hour="minHour"
          :max-hour="22"
          :min-minute="minMinute"
        />
      </van-picker-group>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { useCartStore } from '@/stores/cart'
import { createReservation } from '@/api/order'

const router = useRouter()
const cartStore = useCartStore()

const reserveTime = ref('')
const peopleCount = ref('')
const showTimePicker = ref(false)

// 日期选择器的值 [year, month, day]
const now = new Date()
const twoHoursLater = new Date(now.getTime() + 2 * 60 * 60 * 1000)

const pickerDate = ref([
  String(twoHoursLater.getFullYear()),
  String(twoHoursLater.getMonth() + 1).padStart(2, '0'),
  String(twoHoursLater.getDate()).padStart(2, '0')
])

// 时间选择器的值 [hour, minute]
const pickerTime = ref([
  String(twoHoursLater.getHours()).padStart(2, '0'),
  String(twoHoursLater.getMinutes()).padStart(2, '0')
])

// 最小日期：今天
const minDate = new Date(now.getFullYear(), now.getMonth(), now.getDate())
// 最大日期：30天后
const maxDate = new Date(now.getTime() + 30 * 24 * 60 * 60 * 1000)

// 动态限制：如果选的是今天，则限制最小小时和分钟
const isToday = computed(() => {
  const current = new Date()
  return (
    pickerDate.value[0] === String(current.getFullYear()) &&
    pickerDate.value[1] === String(current.getMonth() + 1).padStart(2, '0') &&
    pickerDate.value[2] === String(current.getDate()).padStart(2, '0')
  )
})

const minHour = computed(() => {
  if (isToday.value) {
    const twoHLater = new Date(Date.now() + 2 * 60 * 60 * 1000)
    return twoHLater.getHours()
  }
  return 6
})

const minMinute = computed(() => {
  if (isToday.value) {
    const twoHLater = new Date(Date.now() + 2 * 60 * 60 * 1000)
    const currentHour = parseInt(pickerTime.value[0])
    if (currentHour === twoHLater.getHours()) {
      return twoHLater.getMinutes()
    }
  }
  return 0
})

// 当日期变更且选的是今天时，自动修正时间不早于最小值
watch(pickerDate, () => {
  if (isToday.value) {
    const twoHLater = new Date(Date.now() + 2 * 60 * 60 * 1000)
    const currentHour = parseInt(pickerTime.value[0])
    const currentMin = parseInt(pickerTime.value[1])
    if (currentHour < twoHLater.getHours() || (currentHour === twoHLater.getHours() && currentMin < twoHLater.getMinutes())) {
      pickerTime.value = [
        String(twoHLater.getHours()).padStart(2, '0'),
        String(twoHLater.getMinutes()).padStart(2, '0')
      ]
    }
  }
})

const onPickerConfirm = () => {
  const [year, month, day] = pickerDate.value
  const [hour, minute] = pickerTime.value
  reserveTime.value = `${year}-${month}-${day} ${hour}:${minute}`
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
