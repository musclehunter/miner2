import { createStore } from 'vuex'
import authService from '@/services/auth'
import inventory from './modules/inventory'
import base from './modules/base'

export default createStore({
  state: {
    user: null,
    isAuthenticated: false,
    loading: false,
    error: null,
    successMessage: null
  },
  getters: {
    isAuthenticated: state => state.isAuthenticated,
    currentUser: state => state.user,
    isLoading: state => state.loading,
    hasError: state => state.error !== null,
    getError: state => state.error,
    hasSuccessMessage: state => state.successMessage !== null,
    getSuccessMessage: state => state.successMessage
  },
  mutations: {
    SET_USER(state, user) {
      state.user = user
      state.isAuthenticated = !!user
    },
    SET_LOADING(state, status) {
      state.loading = status
    },
    SET_ERROR(state, error) {
      state.error = error
      // エラーが設定されたら成功メッセージをクリア
      state.successMessage = null
    },
    CLEAR_ERROR(state) {
      state.error = null
    },
    SET_SUCCESS_MESSAGE(state, message) {
      state.successMessage = message
      // 成功メッセージが設定されたらエラーをクリア
      state.error = null
    },
    CLEAR_SUCCESS_MESSAGE(state) {
      state.successMessage = null
    }
  },
  actions: {
    // ログイン処理
    async login({ commit }, { email, password }) {
      commit('SET_LOADING', true)
      commit('CLEAR_ERROR')
      try {
        // モック処理（バックエンドが準備できていない場合）
        if (process.env.VUE_APP_USE_MOCK === 'true') {
          console.log('モックモードでログイン中')
          await new Promise(resolve => setTimeout(resolve, 1000))
          const user = { id: '1', email, name: 'テストユーザー' }
          commit('SET_USER', user)
          localStorage.setItem('user', JSON.stringify(user))
          return true
        }
        
        // 実際のAPIコール
        console.log('本番モードでAPIログイン中', { email, password })
        const response = await authService.login({ email, password })
        
        // 明確にトークンとユーザー情報が含まれているか確認
        if (!response || !response.token || !response.user) {
          console.error('無効なレスポンス形式:', response)
          commit('SET_ERROR', 'ログインレスポンスが不正です')
          return false
        }
        
        console.log('ログイン成功:', response.user)
        commit('SET_USER', response.user)
        return true
      } catch (error) {
        console.error('ログインエラー発生:', error)
        commit('SET_ERROR', error.response?.data?.error || 'ログインに失敗しました')
        return false
      } finally {
        commit('SET_LOADING', false)
      }
    },

    // サインアップ処理
    async signup({ commit }, { email, password, name }) {
      commit('SET_LOADING', true)
      commit('CLEAR_ERROR')
      commit('CLEAR_SUCCESS_MESSAGE')
      try {
        // モック処理（バックエンドが準備できていない場合）
        if (process.env.VUE_APP_USE_MOCK === 'true') {
          await new Promise(resolve => setTimeout(resolve, 1000))
          // モックモードでも、メール確認プロセスをシミュレートするため、
          // ユーザーを認証済みとしてマークしない
          return true
        }
        
        // 実際のAPIコール
        const response = await authService.signup({ email, password, name: name || email.split('@')[0] })
        // サインアップ成功後、ユーザーはまだ認証されていない
        // メール確認後にのみ認証される
        console.log('サインアップ成功:', response)
        return true
      } catch (error) {
        commit('SET_ERROR', error.response?.data?.error || 'アカウント登録に失敗しました')
        return false
      } finally {
        commit('SET_LOADING', false)
      }
    },
    
    // ログアウト処理
    logout({ commit }) {
      authService.logout()
      commit('SET_USER', null)
    },
    
    // 起動時にローカルストレージからユーザー情報を復元
    async initAuth({ commit }) {
      const user = authService.getStoredUser()
      if (user) {
        commit('SET_USER', user)
        
        // トークンがあればユーザー情報を再取得（開発中はエラーになる可能性があるので条件分岐）
        if (authService.getToken() && process.env.VUE_APP_USE_MOCK !== 'true') {
          try {
            const response = await authService.getCurrentUser()
            if (response && response.user) {
              commit('SET_USER', response.user)
            }
          } catch (error) {
            console.error('ユーザー情報取得エラー:', error)
            // エラーがあってもログアウトはしない（interceptorがあるので）
          }
        }
      }
    }
  },
  modules: {
    inventory,
    base
  }
})
