import { createApp } from 'vue'
import { createPinia } from 'pinia'
import naive from 'naive-ui'
import App from './App.vue'
import router from './router'
import './style.css'

const app = createApp(App)

// 状态管理
app.use(createPinia())

// 路由
app.use(router)

// UI 组件库
app.use(naive)

app.mount('#app')
