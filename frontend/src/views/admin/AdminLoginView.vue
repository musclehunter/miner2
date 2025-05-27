<template>
  <div class="min-h-screen bg-gray-100 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
        管理者ログイン
      </h2>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div class="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
        <div v-if="error" class="mb-4 p-3 bg-red-100 text-red-700 rounded">
          {{ error }}
        </div>
        
        <form class="space-y-6" @submit.prevent="handleLogin">
          <div>
            <label for="secretKey" class="block text-sm font-medium text-gray-700">
              管理者キー
            </label>
            <div class="mt-1">
              <input 
                id="secretKey" 
                name="secretKey" 
                type="password" 
                v-model="secretKey"
                required 
                class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" 
              />
            </div>
          </div>

          <div>
            <button 
              type="submit" 
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              :disabled="loading"
            >
              {{ loading ? 'ログイン中...' : 'ログイン' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import adminService from '@/services/admin'

export default {
  name: 'AdminLoginView',
  setup() {
    const secretKey = ref('')
    const error = ref('')
    const loading = ref(false)
    const router = useRouter()

    const handleLogin = async () => {
      loading.value = true
      error.value = ''
      
      try {
        // 管理者ログイン実行
        const result = await adminService.login(secretKey.value)
        
        // 認証情報をローカルストレージに保存
        localStorage.setItem('adminToken', result.token)
        
        // ダッシュボードへ遷移
        router.push({ name: 'adminDashboard' })
      } catch (err) {
        console.error('管理者ログインエラー:', err)
        error.value = '管理者キーが無効です。'
      } finally {
        loading.value = false
      }
    }

    return {
      secretKey,
      error,
      loading,
      handleLogin
    }
  }
}
</script>
