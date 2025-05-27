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

// 管理者認証確認関数
const isAdminAuthenticated = () => {
  // ローカルストレージから管理者トークンを取得
  const adminToken = localStorage.getItem('adminToken')
  
  // 管理者トークンが存在する場合は認証済み
  return !!adminToken
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
  },
  // 管理者用ルート
  {
    path: '/admin/login',
    name: 'adminLogin',
    component: () => import(/* webpackChunkName: "admin-login" */ '../views/admin/AdminLoginView.vue')
  },
  {
    path: '/admin',
    name: 'adminDashboard',
    component: () => import(/* webpackChunkName: "admin-dashboard" */ '../views/admin/AdminDashboardView.vue'),
    meta: { requiresAdmin: true } // 管理者認証が必要
  },
  {
    path: '/admin/users',
    name: 'adminUsers',
    component: () => import(/* webpackChunkName: "admin-users" */ '../views/admin/AdminUsersView.vue'),
    meta: { requiresAdmin: true }
  },
  {
    path: '/admin/pending-users',
    name: 'adminPendingUsers',
    component: () => import(/* webpackChunkName: "admin-pending-users" */ '../views/admin/AdminPendingUsersView.vue'),
    meta: { requiresAdmin: true }
  },
  {
    path: '/admin/towns',
    name: 'adminTowns',
    component: () => import(/* webpackChunkName: "admin-towns" */ '../views/admin/AdminTownsView.vue'),
    meta: { requiresAdmin: true }
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
  const requiresAdmin = to.matched.some(record => record.meta.requiresAdmin)
  const isUserAuthenticated = isAuthenticated()
  const isAdmin = isAdminAuthenticated()
  
  console.log(`ルート遷移: ${from.path} -> ${to.path}, ユーザー認証必要: ${requiresAuth}, 管理者認証必要: ${requiresAdmin}`)
  console.log(`ユーザー認証状態: ${isUserAuthenticated}, 管理者認証状態: ${isAdmin}`)
  
  if (requiresAuth && !isUserAuthenticated) {
    // ユーザー認証が必要だが認証されていない場合、タイトル画面にリダイレクト
    console.warn('未認証ユーザーのアクセスを拒否します')
    next({ name: 'title' })
  } else if (requiresAdmin && !isAdmin) {
    // 管理者認証が必要だが認証されていない場合、管理者ログイン画面にリダイレクト
    console.warn('未認証管理者のアクセスを拒否します')
    next({ name: 'adminLogin' })
  } else {
    // それ以外は通常の遷移を許可
    next()
  }
})

export default router
