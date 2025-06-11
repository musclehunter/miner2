<template>
  <nav class="bg-dark-lighter/80 backdrop-blur-sm border-b border-dark-light">
    <div class="container mx-auto px-4 py-2">
      <div class="flex justify-between items-center">
        <!-- Logo/Title -->
        <div class="flex items-center">
          <router-link to="/" class="text-xl font-bold text-light flex items-center">
            <span class="mr-2">⛏️</span>
            <span>The Miner</span>
          </router-link>
        </div>
        
        <!-- Navigation Links -->
        <div class="flex items-center space-x-4">
          <!-- 未認証時のみ表示 -->
          <template v-if="!isAuthenticated">
            <router-link to="/" class="nav-link" active-class="nav-link-active">
              ログイン
            </router-link>
          </template>
          
          <!-- 認証済みの場合に表示 -->
          <template v-if="isAuthenticated">
            <router-link to="/world-map" class="nav-link" active-class="nav-link-active">
              <span class="hidden md:inline">ワールドマップ</span>
              <span class="md:hidden">🗺️</span>
            </router-link>
            <router-link to="/base" class="nav-link" active-class="nav-link-active">
              <span class="hidden md:inline">拠点</span>
              <span class="md:hidden">🏠</span>
            </router-link>
            <router-link to="/market" class="nav-link" active-class="nav-link-active">
              <span class="hidden md:inline">マーケット</span>
              <span class="md:hidden">🏪</span>
            </router-link>
            <router-link to="/workers" class="nav-link" active-class="nav-link-active">
              <span class="hidden md:inline">従業員</span>
              <span class="md:hidden">👷</span>
            </router-link>
            <router-link to="/mail" class="nav-link" active-class="nav-link-active">
              <span class="hidden md:inline">メール</span>
              <span class="md:hidden">📧</span>
            </router-link>
            
            <!-- ログアウトボタン -->
            <button @click="logout" class="nav-link text-red-400 hover:text-red-300">
              <span class="hidden md:inline">ログアウト</span>
              <span class="md:hidden">🚪</span>
            </button>
          </template>
          
          <!-- 管理者メニュー（管理者認証時のみ表示） -->
          <template v-if="isAdmin">
            <div 
              class="relative"
              @mouseenter="openMenu"
              @mouseleave="closeMenu"
            >
              <button class="nav-link flex items-center">
                <span class="hidden md:inline">管理者</span>
                <span class="md:hidden">👑</span>
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              <div 
                :class="['absolute right-0 mt-2 w-48 bg-dark-lighter border border-dark-light rounded shadow-lg z-10', { 'block': isAdminMenuOpen, 'hidden': !isAdminMenuOpen }]"
                @mouseenter="openMenu"
                @mouseleave="closeMenu"
              >
                <router-link to="/admin" class="block px-4 py-2 text-sm text-light hover:bg-dark-light">
                  ダッシュボード
                </router-link>
                <router-link to="/admin/users" class="block px-4 py-2 text-sm text-light hover:bg-dark-light">
                  ユーザー管理
                </router-link>
                <router-link to="/admin/pending-users" class="block px-4 py-2 text-sm text-light hover:bg-dark-light">
                  承認待ちユーザー
                </router-link>
                <router-link to="/admin/towns" class="block px-4 py-2 text-sm text-light hover:bg-dark-light">
                  町管理
                </router-link>
                <router-link to="/admin/bases" class="block px-4 py-2 text-sm text-light hover:bg-dark-light">
                  拠点管理
                </router-link>
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
import { ref, computed } from 'vue';
import { useStore } from 'vuex';
import { useRouter } from 'vue-router';

export default {
  name: 'NavMenu',
  setup() {
    const store = useStore();
    const router = useRouter();
    const isAdminMenuOpen = ref(false);
    let menuTimeout = null;

    const openMenu = () => {
      clearTimeout(menuTimeout);
      isAdminMenuOpen.value = true;
    };

    const closeMenu = () => {
      menuTimeout = setTimeout(() => {
        isAdminMenuOpen.value = false;
      }, 200);
    };
    
    // 認証状態を取得
    const isAuthenticated = computed(() => {
      // ローカルストレージからトークンとユーザー情報を取得
      const token = localStorage.getItem('token');
      const user = localStorage.getItem('user');
      
      // 両方が存在する場合のみ認証済みと判定
      return !!token && !!user;
    });
    
    // 管理者認証状態を取得
    const isAdmin = computed(() => {
      // ローカルストレージから管理者トークンを取得
      const adminToken = localStorage.getItem('adminToken');
      
      // 管理者トークンが存在する場合は認証済み
      return !!adminToken;
    });
    
    // ログアウト処理
    const logout = () => {
      // ローカルストレージからトークンとユーザー情報を削除
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      
      // ストアの認証状態をリセット
      store.commit('setUser', null);
      store.commit('setToken', null);
      
      // タイトル画面に遷移
      router.push('/');
    };
    
    return {
      isAuthenticated,
      isAdmin,
      logout,
      isAdminMenuOpen,
      openMenu,
      closeMenu,
    };
  }
}
</script>

<style scoped>
.nav-link {
  @apply text-light-dark hover:text-light transition-colors px-2 py-1 rounded;
}

.nav-link-active {
  @apply text-primary font-medium;
}
</style>
