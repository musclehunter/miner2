import { createRouter, createWebHistory } from 'vue-router'
import TitleView from '../views/TitleView.vue'

const routes = [
  {
    path: '/',
    name: 'title',
    component: TitleView
  },
  {
    path: '/world-map',
    name: 'worldMap',
    // レイジーロード
    component: () => import(/* webpackChunkName: "world-map" */ '../views/WorldMapView.vue')
  },
  {
    path: '/base',
    name: 'base',
    component: () => import(/* webpackChunkName: "base" */ '../views/BaseView.vue')
  },
  {
    path: '/market',
    name: 'market',
    component: () => import(/* webpackChunkName: "market" */ '../views/MarketView.vue')
  },
  {
    path: '/workers',
    name: 'workers',
    component: () => import(/* webpackChunkName: "workers" */ '../views/WorkersView.vue')
  },
  {
    path: '/mail',
    name: 'mail',
    component: () => import(/* webpackChunkName: "mail" */ '../views/MailView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
