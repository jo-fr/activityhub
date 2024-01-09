import { createRouter, createWebHistory } from 'vue-router'
import ExploreView from '../views/ExploreView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: ExploreView
    }
  ]
})

export default router
