import { createApp } from 'vue'
import router from './router'
import AppLayout from './AppLayout.vue'
import 'animate.css'
import AOS from 'aos'
import 'aos/dist/aos.css'

// 初始化AOS
AOS.init({
  duration: 800,
  easing: 'ease-in-out',
  once: true
})

const app = createApp(AppLayout)
app.use(router)
app.mount('#app') 