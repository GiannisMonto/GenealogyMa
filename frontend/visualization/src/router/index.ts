import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/pages/Home.vue')
  },
  {
    path: '/genealogy',
    name: 'GenealogyTree',
    component: () => import('@/pages/GenealogyTree.vue')
  },
  {
    path: '/person/:id',
    name: 'PersonDetail',
    component: () => import('@/pages/PersonDetail.vue'),
    props: true
  },
  {
    path: '/statistics',
    name: 'Statistics',
    component: () => import('@/pages/Statistics.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
