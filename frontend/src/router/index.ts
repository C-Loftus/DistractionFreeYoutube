import { createRouter, createWebHistory } from 'vue-router'
import SubscriptionTable from '../views/SubscriptionTable.vue'
import NotFound from '@/views/NotFound.vue'
import App from '@/App.vue'
import LoginHome from '@/views/LoginHome.vue'
import PlaylistsTable from '@/views/PlaylistsTable.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/subscriptions',
      name: 'subscriptions',
      component: SubscriptionTable
    },
    {
      path: '/playlists',
      component: PlaylistsTable
    },
    {
      path: '/:catchAll(.*)',
      component: NotFound
    },
    {
      path: '/',
      component: LoginHome
    }
  ]
})

export default router
