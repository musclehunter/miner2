<template>
  <div class="min-h-screen bg-gray-100">
    <AdminHeader />
    
    <div class="py-10">
      <header>
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h1 class="text-3xl font-bold leading-tight text-gray-900">
            未確認ユーザー管理
          </h1>
        </div>
      </header>
      <main>
        <div class="max-w-7xl mx-auto sm:px-6 lg:px-8">
          <div class="px-4 py-8 sm:px-0">
            <!-- 読み込み中表示 -->
            <div v-if="loading" class="flex justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span>読み込み中...</span>
            </div>
            
            <!-- エラー表示 -->
            <div v-else-if="error" class="bg-red-50 p-4 rounded-md">
              <div class="flex">
                <div class="flex-shrink-0">
                  <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                  </svg>
                </div>
                <div class="ml-3">
                  <h3 class="text-sm font-medium text-red-800">
                    データの取得に失敗しました
                  </h3>
                  <div class="mt-2 text-sm text-red-700">
                    <p>{{ error }}</p>
                  </div>
                  <div class="mt-4">
                    <button
                      @click="fetchPendingUsers"
                      class="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-red-700 bg-red-100 hover:bg-red-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                    >
                      再試行
                    </button>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- 未確認ユーザーがいない場合 -->
            <div v-else-if="pendingUsers.length === 0" class="bg-white overflow-hidden shadow sm:rounded-md p-6 text-center">
              <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
              </svg>
              <h3 class="mt-2 text-sm font-medium text-gray-900">未確認ユーザーはいません</h3>
              <p class="mt-1 text-sm text-gray-500">現在、メール確認待ちのユーザーはいません。</p>
            </div>
            
            <!-- 未確認ユーザーリスト -->
            <div v-else class="bg-white shadow overflow-hidden sm:rounded-md">
              <ul class="divide-y divide-gray-200">
                <li v-for="user in pendingUsers" :key="user.token">
                  <div class="px-4 py-4 sm:px-6">
                    <div class="flex items-center justify-between">
                      <div class="flex items-center">
                        <div class="flex-shrink-0 bg-yellow-100 p-2 rounded-full">
                          <svg class="h-6 w-6 text-yellow-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                          </svg>
                        </div>
                        <div class="ml-3">
                          <p class="text-sm font-medium text-indigo-600 truncate">
                            {{ user.email }}
                          </p>
                          <p class="text-sm text-gray-500">
                            {{ user.name }}
                          </p>
                        </div>
                      </div>
                      <div class="flex">
                        <button
                          @click="resendVerification(user.email)"
                          class="ml-2 inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded text-indigo-700 bg-indigo-100 hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                          :disabled="resending === user.email || deleting === user.token"
                        >
                          {{ resending === user.email ? '送信中...' : '確認メール再送信' }}
                        </button>
                        <button
                          @click="confirmDelete(user)"
                          class="ml-2 inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded text-red-700 bg-red-100 hover:bg-red-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                          :disabled="deleting === user.token || resending === user.email"
                        >
                          {{ deleting === user.token ? '削除中...' : '削除' }}
                        </button>
                      </div>
                    </div>
                    <div class="mt-2 sm:flex sm:justify-between">
                      <div class="sm:flex">
                        <p class="flex items-center text-sm text-gray-500">
                          <svg class="flex-shrink-0 mr-1.5 h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd" d="M18 8a6 6 0 01-7.743 5.743L10 14l-1 1-1 1H6v-1l1-1 1-1-.257-.257A6 6 0 1118 8zm-6-4a1 1 0 10-2 0v1a1 1 0 102 0V4z" clip-rule="evenodd" />
                          </svg>
                          確認トークン: {{ user.token }}
                        </p>
                      </div>
                      <div class="mt-2">
                        <a 
                          :href="getVerificationUrl(user.token)" 
                          target="_blank"
                          class="inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded text-green-700 bg-green-100 hover:bg-green-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                        >
                          <svg class="flex-shrink-0 mr-1.5 h-4 w-4 text-green-600" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                          </svg>
                          確認リンクを開く（開発用）
                        </a>
                      </div>
                    </div>
                  </div>
                </li>
              </ul>
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
  name: 'AdminPendingUsersView',
  components: {
    AdminHeader
  },
  setup() {
    const pendingUsers = ref([])
    const loading = ref(true)
    const error = ref('')
    const resending = ref('')
    const deleting = ref('')
    const showDeleteModal = ref(false)
    const deleteTarget = ref(null)
    
    // API URL
    const apiBaseUrl = process.env.VUE_APP_API_URL || 'http://localhost:8090/api'

    // 未確認ユーザー一覧取得
    const fetchPendingUsers = async () => {
      loading.value = true
      error.value = ''
      
      try {
        pendingUsers.value = await adminService.getAllPendingUsers()
      } catch (err) {
        console.error('未確認ユーザー一覧取得エラー:', err)
        error.value = '未確認ユーザー情報の取得に失敗しました'
      } finally {
        loading.value = false
      }
    }

    // 確認メール再送信
    const resendVerification = async (email) => {
      resending.value = email
      
      try {
        await adminService.resendVerification(email)
        alert(`${email} 宛に確認メールを再送信しました`)
      } catch (err) {
        console.error('確認メール再送信エラー:', err)
        alert('確認メールの再送信に失敗しました')
      } finally {
        resending.value = ''
      }
    }
    
    // 削除確認ダイアログを表示
    const confirmDelete = (user) => {
      deleteTarget.value = user
      const confirmed = confirm(`未確認ユーザー "${user.email}" を削除しますか？\nこの操作は元に戻せません。`)
      
      if (confirmed) {
        deletePendingUser(user.token)
      }
    }
    
    // 未確認ユーザー削除
    const deletePendingUser = async (token) => {
      deleting.value = token
      
      try {
        await adminService.deletePendingUser(token)
        // 削除成功後、リストから該当ユーザーを除外
        pendingUsers.value = pendingUsers.value.filter(user => user.token !== token)
        alert('未確認ユーザーを削除しました')
      } catch (err) {
        console.error('未確認ユーザー削除エラー:', err)
        alert('未確認ユーザーの削除に失敗しました')
      } finally {
        deleting.value = ''
        deleteTarget.value = null
      }
    }
    
    // 確認リンクのURLを生成
    const getVerificationUrl = (token) => {
      // ベースURLが「/api」で終わっている場合は、それを除去
      const baseUrl = apiBaseUrl.endsWith('/api')
        ? apiBaseUrl.substring(0, apiBaseUrl.length - 4)
        : apiBaseUrl.replace('/api', '')
      
      return `${baseUrl}/api/auth/verify-email?token=${token}`
    }

    // コンポーネントマウント時にデータを取得
    onMounted(() => {
      fetchPendingUsers()
    })

    return {
      pendingUsers,
      loading,
      error,
      resending,
      deleting,
      showDeleteModal,
      deleteTarget,
      fetchPendingUsers,
      resendVerification,
      confirmDelete,
      deletePendingUser,
      getVerificationUrl
    }
  }
}
</script>
