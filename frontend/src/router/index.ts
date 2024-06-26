import { createRouter, createWebHistory } from 'vue-router'
import SubscriptionTable from '../views/SubscriptionTable.vue'
import NotFound from '@/views/NotFound.vue'
import LoginHome from '@/views/LoginHome.vue'
import PlaylistsTable from '@/views/PlaylistsTable.vue'
import SearchResults from '@/views/SearchResults.vue'
import VideoTable from '@/views/VideoTable.vue'

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
      path: '/playlists/:name',
      component: VideoTable
    },
    {
      path: '/:catchAll(.*)',
      component: NotFound
    },
    {
      path: '/',
      component: LoginHome
    },
    {
      path: '/search',
      component: SearchResults
    },
    {
      path: '/search?=query=:query',
      component: SearchResults,
      props: (route) => ({ query: route.query.query })
    }
  ]
})

export default router
