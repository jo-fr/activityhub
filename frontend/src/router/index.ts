import { createRouter, createWebHistory } from 'vue-router'
import ExploreView from '../views/ExploreView.vue'
import FeedDetailView from '../views/FeedDetailView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: ExploreView
    },
    {
      path: '/feed/:id',
      name: 'feedDetail',
      component: FeedDetailView,
      props: true
    }
  ]
})

export default router
