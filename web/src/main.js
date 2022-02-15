import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/antd.css'
import { createApp } from 'vue'
import App from './App.vue'
import './assets/app.css'
import router from './routes'
import store from './store'

const app = createApp(App)
app.use(Antd)
app.use(router)
app.use(store)
app.mount('#app')
