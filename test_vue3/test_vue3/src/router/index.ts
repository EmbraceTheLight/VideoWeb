import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Test from '../views/test.vue'
import Account from '../views/account.vue'
import Video from '../views/video.vue'
import Search from '../views/search.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      component: Home,
    },
    {
      path: '/test',
      component: Test,
    },
    {
      path: '/account',
      component: Account,
    },
    {
      path: '/video/VideoID=:id',
      name: 'videoDetail',
      component: Video,
    },
    {
      path: '/:pathMatch(.*)',
      component: Home,
    },
    {
      path: '/search/keyword=:keyword',
      component: Search,
    }
  ],
})

export default router
