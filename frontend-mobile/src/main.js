import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'

// Vant组件库
import Vant from 'vant'
import 'vant/lib/index.css'

// Swiper样式
import 'swiper/css'
import 'swiper/css/pagination'

// 全局样式
import './styles/main.scss'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(Vant)

app.mount('#app')
