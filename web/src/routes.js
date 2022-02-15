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
        component: () => import('@/pages/connections/ConnectionLayout.vue'),
        children: [
          {
            path: '',
            name: 'connections.list',
            component: () => import('@/pages/connections/ConnectionList.vue')
          },
          {
            path: 'new',
            name: 'new-connection',
            component: () => import('@/pages/connections/NewConnection.vue')
          }
        ]
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
