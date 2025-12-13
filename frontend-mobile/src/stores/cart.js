import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useCartStore = defineStore('cart', () => {
  const items = ref(JSON.parse(localStorage.getItem('cartItems') || '[]'))

  const totalCount = computed(() => {
    return items.value.reduce((sum, item) => sum + item.quantity, 0)
  })

  const totalPrice = computed(() => {
    return items.value.reduce((sum, item) => sum + item.price * item.quantity, 0)
  })

  const addItem = (dish) => {
    const existItem = items.value.find(item => item.id === dish.id)
    if (existItem) {
      existItem.quantity++
    } else {
      items.value.push({
        id: dish.id,
        name: dish.name,
        price: dish.price,
        image: dish.images?.[0]?.image_url || '',
        quantity: 1
      })
    }
    saveToStorage()
  }

  const removeItem = (dishId) => {
    const index = items.value.findIndex(item => item.id === dishId)
    if (index > -1) {
      items.value.splice(index, 1)
      saveToStorage()
    }
  }

  const updateQuantity = (dishId, quantity) => {
    const item = items.value.find(item => item.id === dishId)
    if (item) {
      item.quantity = quantity
      if (item.quantity <= 0) {
        removeItem(dishId)
      } else {
        saveToStorage()
      }
    }
  }

  const clearCart = () => {
    items.value = []
    saveToStorage()
  }

  const saveToStorage = () => {
    localStorage.setItem('cartItems', JSON.stringify(items.value))
  }

  return {
    items,
    totalCount,
    totalPrice,
    addItem,
    removeItem,
    updateQuantity,
    clearCart
  }
})
