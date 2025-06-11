<template>
  <nav class="bg-indigo-600">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <span class="text-white font-bold text-xl">採掘人 管理パネル</span>
          </div>
          <div class="hidden md:block">
            <div class="ml-10 flex items-baseline space-x-4">
              <router-link 
                :to="{ name: 'adminDashboard' }" 
                class="text-white hover:bg-indigo-500 hover:bg-opacity-75 px-3 py-2 rounded-md text-sm font-medium" 
                :class="{ 'bg-indigo-700': isActive('adminDashboard') }"
              >
                ダッシュボード
              </router-link>

              <router-link 
                :to="{ name: 'adminUsers' }" 
                class="text-white hover:bg-indigo-500 hover:bg-opacity-75 px-3 py-2 rounded-md text-sm font-medium" 
                :class="{ 'bg-indigo-700': isActive('adminUsers') }"
              >
                ユーザー管理
              </router-link>

              <router-link 
                :to="{ name: 'adminPendingUsers' }" 
                class="text-white hover:bg-indigo-500 hover:bg-opacity-75 px-3 py-2 rounded-md text-sm font-medium" 
                :class="{ 'bg-indigo-700': isActive('adminPendingUsers') }"
              >
                未確認ユーザー
              </router-link>

              <router-link 
                :to="{ name: 'adminTowns' }" 
                class="text-white hover:bg-indigo-500 hover:bg-opacity-75 px-3 py-2 rounded-md text-sm font-medium" 
                :class="{ 'bg-indigo-700': isActive('adminTowns') }"
              >
                町管理
              </router-link>

              <router-link 
                :to="{ name: 'adminBases' }" 
                class="text-white hover:bg-indigo-500 hover:bg-opacity-75 px-3 py-2 rounded-md text-sm font-medium" 
                :class="{ 'bg-indigo-700': isActive('adminBases') }"
              >
                拠点管理
              </router-link>
            </div>
          </div>
        </div>
        <div class="hidden md:block">
          <div class="ml-4 flex items-center md:ml-6">
            <button 
              @click="handleLogout"
              class="bg-indigo-600 p-1 rounded-full text-white hover:text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-indigo-600 focus:ring-white"
            >
              <span class="sr-only">ログアウト</span>
              <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
            </button>
          </div>
        </div>
        <div class="-mr-2 flex md:hidden">
          <!-- モバイルメニューボタン -->
          <button 
            @click="mobileMenuOpen = !mobileMenuOpen"
            class="bg-indigo-600 inline-flex items-center justify-center p-2 rounded-md text-white hover:text-white hover:bg-indigo-500 hover:bg-opacity-75 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-indigo-600 focus:ring-white"
          >
            <span class="sr-only">メニューを開く</span>
            <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- モバイルメニュー -->
    <div class="md:hidden" v-show="mobileMenuOpen">
      <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3">
        <router-link 
          :to="{ name: 'adminDashboard' }" 
          class="text-white hover:bg-indigo-500 hover:bg-opacity-75 block px-3 py-2 rounded-md text-base font-medium" 
          :class="{ 'bg-indigo-700': isActive('adminDashboard') }"
        >
          ダッシュボード
        </router-link>

        <router-link 
          :to="{ name: 'adminUsers' }" 
          class="text-white hover:bg-indigo-500 hover:bg-opacity-75 block px-3 py-2 rounded-md text-base font-medium" 
          :class="{ 'bg-indigo-700': isActive('adminUsers') }"
        >
          ユーザー管理
        </router-link>

        <router-link 
          :to="{ name: 'adminPendingUsers' }" 
          class="text-white hover:bg-indigo-500 hover:bg-opacity-75 block px-3 py-2 rounded-md text-base font-medium" 
          :class="{ 'bg-indigo-700': isActive('adminPendingUsers') }"
        >
          未確認ユーザー
        </router-link>

        <router-link 
          :to="{ name: 'adminTowns' }" 
          class="text-white hover:bg-indigo-500 hover:bg-opacity-75 block px-3 py-2 rounded-md text-base font-medium" 
          :class="{ 'bg-indigo-700': isActive('adminTowns') }"
        >
          町管理
        </router-link>

        <router-link 
          :to="{ name: 'adminBases' }" 
          class="text-white hover:bg-indigo-500 hover:bg-opacity-75 block px-3 py-2 rounded-md text-base font-medium" 
          :class="{ 'bg-indigo-700': isActive('adminBases') }"
        >
          拠点管理
        </router-link>

        <button 
          @click="handleLogout"
          class="text-white hover:bg-indigo-500 hover:bg-opacity-75 block px-3 py-2 rounded-md text-base font-medium w-full text-left"
        >
          ログアウト
        </button>
      </div>
    </div>
  </nav>
</template>

<script>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'

export default {
  name: 'AdminHeader',
  setup() {
    const mobileMenuOpen = ref(false)
    const router = useRouter()
    const route = useRoute()

    const isActive = (routeName) => {
      return route.name === routeName
    }

    const handleLogout = () => {
      // 管理者トークンを削除
      localStorage.removeItem('adminToken')
      
      // ログイン画面へリダイレクト
      router.push({ name: 'adminLogin' })
    }

    return {
      mobileMenuOpen,
      isActive,
      handleLogout
    }
  }
}
</script>
