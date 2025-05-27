<template>
  <div class="min-h-screen bg-gray-100">
    <AdminHeader />
    
    <div class="py-10">
      <header>
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h1 class="text-3xl font-bold leading-tight text-gray-900">
            管理ダッシュボード
          </h1>
        </div>
      </header>
      <main>
        <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
          <div class="px-4 py-8 sm:px-0">
            <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
              <!-- ユーザー管理カード -->
              <div class="bg-white overflow-hidden shadow rounded-lg">
                <div class="p-5">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 bg-indigo-500 rounded-md p-3">
                      <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                      </svg>
                    </div>
                    <div class="ml-5 w-0 flex-1">
                      <dl>
                        <dt class="text-sm font-medium text-gray-500 truncate">
                          ユーザー管理
                        </dt>
                        <dd>
                          <div class="text-lg font-medium text-gray-900">
                            {{ userStats.total }} ユーザー
                          </div>
                        </dd>
                      </dl>
                    </div>
                  </div>
                </div>
                <div class="bg-gray-50 px-5 py-3">
                  <div class="text-sm">
                    <router-link :to="{ name: 'adminUsers' }" class="font-medium text-indigo-600 hover:text-indigo-900">
                      ユーザー管理へ
                    </router-link>
                  </div>
                </div>
              </div>

              <!-- 未確認ユーザー管理カード -->
              <div class="bg-white overflow-hidden shadow rounded-lg">
                <div class="p-5">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 bg-yellow-500 rounded-md p-3">
                      <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                      </svg>
                    </div>
                    <div class="ml-5 w-0 flex-1">
                      <dl>
                        <dt class="text-sm font-medium text-gray-500 truncate">
                          未確認ユーザー
                        </dt>
                        <dd>
                          <div class="text-lg font-medium text-gray-900">
                            {{ pendingUserStats.total }} 人
                          </div>
                        </dd>
                      </dl>
                    </div>
                  </div>
                </div>
                <div class="bg-gray-50 px-5 py-3">
                  <div class="text-sm">
                    <router-link :to="{ name: 'adminPendingUsers' }" class="font-medium text-indigo-600 hover:text-indigo-900">
                      未確認ユーザー管理へ
                    </router-link>
                  </div>
                </div>
              </div>

              <!-- 町管理カード -->
              <div class="bg-white overflow-hidden shadow rounded-lg">
                <div class="p-5">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 bg-green-500 rounded-md p-3">
                      <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                      </svg>
                    </div>
                    <div class="ml-5 w-0 flex-1">
                      <dl>
                        <dt class="text-sm font-medium text-gray-500 truncate">
                          町管理
                        </dt>
                        <dd>
                          <div class="text-lg font-medium text-gray-900">
                            {{ townStats.total }} 町
                          </div>
                        </dd>
                      </dl>
                    </div>
                  </div>
                </div>
                <div class="bg-gray-50 px-5 py-3">
                  <div class="text-sm">
                    <router-link :to="{ name: 'adminTowns' }" class="font-medium text-indigo-600 hover:text-indigo-900">
                      町管理へ
                    </router-link>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import AdminHeader from '@/components/admin/AdminHeader.vue'
import adminService from '@/services/admin'

export default {
  name: 'AdminDashboardView',
  components: {
    AdminHeader
  },
  setup() {
    const userStats = ref({ total: 0 })
    const pendingUserStats = ref({ total: 0 })
    const townStats = ref({ total: 0 })
    const loading = ref(true)
    const error = ref('')

    const fetchDashboardData = async () => {
      loading.value = true
      error.value = ''
      
      try {
        // ユーザー数を取得
        const users = await adminService.getAllUsers()
        userStats.value.total = users.length
        
        // 未確認ユーザー数を取得
        const pendingUsers = await adminService.getAllPendingUsers()
        pendingUserStats.value.total = pendingUsers.length
        
        // 町の数を取得
        const towns = await adminService.getAllTowns()
        townStats.value.total = towns.length
      } catch (err) {
        console.error('ダッシュボードデータ取得エラー:', err)
        error.value = 'データの取得に失敗しました'
      } finally {
        loading.value = false
      }
    }

    // コンポーネントマウント時にデータを取得
    onMounted(() => {
      fetchDashboardData()
    })

    return {
      userStats,
      pendingUserStats,
      townStats,
      loading,
      error
    }
  }
}
</script>
