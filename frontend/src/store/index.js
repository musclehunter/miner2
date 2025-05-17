import { createStore } from 'vuex'

export default createStore({
  state: {
    user: null,
    isAuthenticated: false,
    loading: false,
    error: null
  },
  getters: {
    isAuthenticated: state => state.isAuthenticated,
    currentUser: state => state.user,
    isLoading: state => state.loading,
    hasError: state => state.error !== null,
    getError: state => state.error
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
    },
    CLEAR_ERROR(state) {
      state.error = null
    }
  },
  actions: {
    // ログイン処理
    // ESLintの未使用変数警告を一時的に無効化（将来的にはAPI連携で使用予定）
    // eslint-disable-next-line no-unused-vars
    async login({ commit }, { email, password }) {
      commit('SET_LOADING', true)
      commit('CLEAR_ERROR')
      try {
        // 実際のAPIコールは後で実装
        // 今はモックデータで対応
        await new Promise(resolve => setTimeout(resolve, 1000))
        
        // ダミーのユーザーデータ
        const user = { id: '1', email, name: 'テストユーザー' }
        commit('SET_USER', user)
        
        // ローカルストレージに保存してセッション維持
        localStorage.setItem('user', JSON.stringify(user))
        return true
      } catch (error) {
        commit('SET_ERROR', error.message || 'ログインに失敗しました')
        return false
      } finally {
        commit('SET_LOADING', false)
      }
    },

    // サインアップ処理
    // eslint-disable-next-line no-unused-vars
    async signup({ commit }, { email, password, name }) {
      commit('SET_LOADING', true)
      commit('CLEAR_ERROR')
      try {
        // 実際のAPIコールは後で実装
        // 現在はモックデータで対応
        await new Promise(resolve => setTimeout(resolve, 1000))
        
        // 新規ユーザー情報
        const user = { 
          id: Date.now().toString(), // 一時的なID
          email, 
          name: name || email.split('@')[0] // 名前が指定されていない場合はメールアドレスの@前を使用
        }
        
        commit('SET_USER', user)
        
        // ローカルストレージに保存してセッション維持
        localStorage.setItem('user', JSON.stringify(user))
        return true
      } catch (error) {
        commit('SET_ERROR', error.message || 'アカウント登録に失敗しました')
        return false
      } finally {
        commit('SET_LOADING', false)
      }
    },
    
    // ログアウト処理
    logout({ commit }) {
      commit('SET_USER', null)
      localStorage.removeItem('user')
    },
    
    // 起動時にローカルストレージからユーザー情報を復元
    initAuth({ commit }) {
      const userStr = localStorage.getItem('user')
      if (userStr) {
        try {
          const user = JSON.parse(userStr)
          commit('SET_USER', user)
        } catch (e) {
          localStorage.removeItem('user')
        }
      }
    }
  },
  modules: {
    // 必要に応じて機能ごとにモジュール化
  }
})
