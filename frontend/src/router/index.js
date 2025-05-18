import { createRouter, createWebHistory } from 'vue-router'
import TitleView from '../views/TitleView.vue'

// 認証済みユーザーか確認する関数
const isAuthenticated = () => {
  // ローカルストレージからトークンとユーザー情報を取得
  const token = localStorage.getItem('token')
  const user = localStorage.getItem('user')
  
  // 両方が存在する場合のみ認証済みと判定
  return !!token && !!user
}

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
    component: () => import(/* webpackChunkName: "world-map" */ '../views/WorldMapView.vue'),
    meta: { requiresAuth: true } // 認証が必要
  },
  {
    path: '/base',
    name: 'base',
    component: () => import(/* webpackChunkName: "base" */ '../views/BaseView.vue'),
    meta: { requiresAuth: true } // 認証が必要
  },
  {
    path: '/market',
    name: 'market',
    component: () => import(/* webpackChunkName: "market" */ '../views/MarketView.vue'),
    meta: { requiresAuth: true } // 認証が必要
  },
  {
    path: '/workers',
    name: 'workers',
    component: () => import(/* webpackChunkName: "workers" */ '../views/WorkersView.vue'),
    meta: { requiresAuth: true } // 認証が必要
  },
  {
    path: '/mail',
    name: 'mail',
    component: () => import(/* webpackChunkName: "mail" */ '../views/MailView.vue'),
    meta: { requiresAuth: true } // 認証が必要
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

// グローバルナビゲーションガードを追加
router.beforeEach((to, from, next) => {
  // ルートが認証を必要とする場合と認証状態をチェック
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const isUserAuthenticated = isAuthenticated()
  
  console.log(`ルート遷移: ${from.path} -> ${to.path}, 認証必要: ${requiresAuth}, 認証状態: ${isUserAuthenticated}`)
  
  if (requiresAuth && !isUserAuthenticated) {
    // 認証が必要だが認証されていない場合、タイトル画面にリダイレクト
    console.warn('未認証ユーザーのアクセスを拒否します')
    next({ name: 'title' })
  } else {
    // それ以外は通常の遷移を許可
    next()
  }
})

export default router
