import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/antd.css'
import { createApp } from 'vue'
import App from './App.vue'
import router from './routes'

const app = createApp(App)
app.use(Antd)
app.use(router)
app.mount('#app')
