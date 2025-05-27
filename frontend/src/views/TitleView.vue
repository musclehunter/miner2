<template>
  <div class="flex flex-col items-center justify-center min-h-screen w-full p-4 bg-dark bg-[url('@/assets/images/mine-bg.svg')] bg-cover bg-center">
    <!-- ロゴ部分 -->
    <div class="mb-16 animate-float">
      <h1 class="title-text text-5xl md:text-7xl font-bold text-light drop-shadow-[0_0_10px_rgba(59,130,246,0.7)] mb-3">
        The Miner
      </h1>
      <p class="subtitle-text text-light-dark">掘って、売って、目指せ億万長者！</p>
    </div>

    <!-- ログインフォーム -->
    <div class="w-full max-w-md">
      <div class="card backdrop-blur-sm bg-dark-lighter/80 border border-dark-light">
        <div class="mb-6">
          <h2 class="text-2xl font-bold text-light mb-2">ログイン</h2>
          <p class="text-light-darker text-sm">アカウントにログインして冒険を始めよう</p>
        </div>
        
        <form @submit.prevent="handleLogin" class="space-y-4">
          <!-- エラーメッセージ -->
          <div v-if="error" class="bg-red-900/30 border border-red-500 text-red-200 px-4 py-2 rounded text-sm">
            {{ error }}
          </div>
          
          <!-- 成功メッセージ -->
          <div v-if="successMessage" class="bg-green-900/30 border border-green-500 text-green-200 px-4 py-2 rounded text-sm whitespace-pre-line">
            {{ successMessage }}
          </div>
          
          <!-- メールアドレス -->
          <div>
            <label for="email" class="block text-sm font-medium text-light-dark mb-1">メールアドレス</label>
            <input 
              type="email" 
              id="email" 
              v-model="email" 
              class="input-field"
              placeholder="メールアドレスを入力" 
              required 
            />
          </div>
          
          <!-- パスワード -->
          <div>
            <label for="password" class="block text-sm font-medium text-light-dark mb-1">パスワード</label>
            <input 
              type="password" 
              id="password" 
              v-model="password" 
              class="input-field"
              placeholder="パスワードを入力" 
              required 
            />
          </div>
          
          <!-- ログインボタン -->
          <div>
            <button 
              type="submit" 
              class="btn btn-primary w-full"
              :disabled="loading"
            >
              <span v-if="loading" class="inline-flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                ログイン中...
              </span>
              <span v-else>ログイン</span>
            </button>
          </div>
          
          <!-- アカウント作成リンク -->
          <div class="text-center">
            <p class="text-sm text-light-dark">
              アカウントをお持ちでない方は 
              <a href="#" class="text-primary hover:text-primary-light font-medium" @click.prevent="toggleMode">
                こちら
              </a>
            </p>
          </div>
        </form>
      </div>

      <!-- 登録フォーム -->
      <div v-if="showRegister" class="card mt-4 backdrop-blur-sm bg-dark-lighter/80 border border-dark-light">
        <div class="mb-6">
          <h2 class="text-2xl font-bold text-light mb-2">アカウント登録</h2>
          <p class="text-light-darker text-sm">新しいアカウントを作成</p>
        </div>

        <form @submit.prevent="handleRegister" class="space-y-4">
          <!-- メールアドレス -->
          <div>
            <label for="register-email" class="block text-sm font-medium text-light-dark mb-1">メールアドレス</label>
            <input 
              type="email" 
              id="register-email" 
              v-model="registerEmail" 
              class="input-field"
              placeholder="メールアドレスを入力" 
              required 
            />
          </div>
          
          <!-- パスワード -->
          <div>
            <label for="register-password" class="block text-sm font-medium text-light-dark mb-1">パスワード</label>
            <input 
              type="password" 
              id="register-password" 
              v-model="registerPassword" 
              class="input-field"
              placeholder="パスワードを入力" 
              required 
            />
          </div>
          
          <!-- 確認用パスワード -->
          <div>
            <label for="confirm-password" class="block text-sm font-medium text-light-dark mb-1">パスワード (確認)</label>
            <input 
              type="password" 
              id="confirm-password" 
              v-model="confirmPassword" 
              class="input-field"
              placeholder="パスワードを再入力" 
              required 
            />
          </div>
          
          <!-- 登録ボタン -->
          <div>
            <button 
              type="submit" 
              class="btn btn-secondary w-full"
              :disabled="loading"
            >
              <span v-if="loading" class="inline-flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                登録中...
              </span>
              <span v-else>アカウント作成</span>
            </button>
          </div>
          
          <!-- キャンセル -->
          <div class="text-center">
            <button 
              type="button"
              class="text-sm text-primary hover:text-primary-light font-medium"
              @click="showRegister = false"
            >
              ログイン画面に戻る
            </button>
          </div>
        </form>
      </div>

      <!-- ソーシャルログイン -->
      <div class="mt-6">
        <div class="relative">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-gray-600"></div>
          </div>
          <div class="relative flex justify-center text-sm">
            <span class="px-2 bg-dark text-gray-400">または</span>
          </div>
        </div>
        
        <div class="mt-6 grid grid-cols-1 gap-3">
          <button 
            type="button" 
            class="w-full py-2 px-4 border border-gray-600 rounded-md shadow-sm bg-dark-lighter text-sm font-medium text-light hover:bg-dark-light focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-dark focus:ring-primary flex justify-center items-center"
          >
            <svg class="h-5 w-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
              <path d="M20.283 10.356h-8.327v3.451h4.792c-.446 2.193-2.313 3.453-4.792 3.453a5.27 5.27 0 0 1-5.279-5.28 5.27 5.27 0 0 1 5.279-5.279c1.259 0 2.397.447 3.29 1.178l2.6-2.599c-1.584-1.381-3.615-2.233-5.89-2.233a8.908 8.908 0 0 0-8.934 8.934 8.907 8.907 0 0 0 8.934 8.934c4.467 0 8.529-3.249 8.529-8.934 0-.528-.081-1.097-.202-1.625z"></path>
            </svg>
            Googleでログイン
          </button>
        </div>
      </div>
    </div>
    
    <!-- フッター -->
    <div class="mt-16 text-sm text-light-darker text-center">
      <p>&copy; 2025 The Miner- 目指せ億万長者！</p>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue';
import { useStore } from 'vuex';
import { useRouter } from 'vue-router';

export default {
  name: 'TitleView',
  setup() {
    const store = useStore();
    const router = useRouter();
    
    const email = ref('');
    const password = ref('');
    const registerEmail = ref('');
    const registerPassword = ref('');
    const confirmPassword = ref('');
    const showRegister = ref(false);
    
    const loading = computed(() => store.state.loading);
    const error = computed(() => store.state.error);
    const successMessage = computed(() => store.state.successMessage);
    
    const handleLogin = async () => {
      console.log('ログインボタンがクリックされました');
      try {
        const success = await store.dispatch('login', {
          email: email.value,
          password: password.value
        });
        
        console.log('ログイン結果:', success);
        
        // successがtrueの場合のみワールドマップに遷移
        if (success === true) {
          console.log('ログイン成功、ワールドマップに遷移します');
          router.push('/world-map');
        } else {
          console.error('ログイン失敗');
          // エラーが表示されるので何もしない
        }
      } catch (error) {
        console.error('ログイン中にエラーが発生しました:', error);
      }
    };
    
    const handleRegister = async () => {
      // パスワードの一致チェック
      if (registerPassword.value !== confirmPassword.value) {
        store.commit('SET_ERROR', 'パスワードが一致しません');
        return;
      }
      
      // パスワードの長さチェック
      if (registerPassword.value.length < 6) {
        store.commit('SET_ERROR', 'パスワードは6文字以上で入力してください');
        return;
      }
      
      // メールアドレスの形式チェック
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!emailRegex.test(registerEmail.value)) {
        store.commit('SET_ERROR', '有効なメールアドレスを入力してください');
        return;
      }
      
      // 新しく追加したsignupアクションを使用
      const success = await store.dispatch('signup', {
        email: registerEmail.value,
        password: registerPassword.value,
        name: registerEmail.value.split('@')[0] // デフォルトの名前としてメールアドレスの前半分を使用
      });
      
      if (success) {
        // 登録成功時は、ワールドマップに遷移するのではなく、確認メッセージを表示
        store.commit('SET_SUCCESS_MESSAGE', '登録が完了しました。メールアドレスを確認して登録を完了してください。\n\n開発環境では、トークンはサーバーのログに表示されます。');
        
        // 入力フィールドをクリア
        registerEmail.value = '';
        registerPassword.value = '';
        confirmPassword.value = '';
        
        // ログインフォームに切り替え
        showRegister.value = false;
      }
    };
    
    const toggleMode = () => {
      showRegister.value = !showRegister.value;
      store.commit('CLEAR_ERROR');
    };
    
    return {
      email,
      password,
      registerEmail,
      registerPassword,
      confirmPassword,
      showRegister,
      loading,
      error,
      successMessage,
      handleLogin,
      handleRegister,
      toggleMode
    };
  }
}
</script>
