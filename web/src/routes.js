import MainLayout from '@/layouts/MainLayout.vue'
import { createRouter, createWebHistory } from 'vue-router'
const routes = [
  {
    path: '/',
    component: MainLayout,
    children: [
      {
        path: 'connections',
        name: 'connections',
        component: () => import('@/pages/connections/Connections.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
