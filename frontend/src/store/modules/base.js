import { createBase as createBaseApi } from '@/services/baseService';

const state = {
  base: null,
  loading: false,
  error: null,
};

const getters = {
  base: (state) => state.base,
  loading: (state) => state.loading,
  error: (state) => state.error,
};

const mutations = {
  setBase(state, base) {
    state.base = base;
  },
  setLoading(state, loading) {
    state.loading = loading;
  },
  setError(state, error) {
    state.error = error;
  },
};

const actions = {
  async createBase({ commit }, { townId }) {
    commit('setLoading', true);
    commit('setError', null);
    try {
      const response = await createBaseApi(townId);
      commit('setBase', response.data);
      return response.data;
    } catch (error) {
      commit('setError', error.response?.data?.error || '拠点設立に失敗しました。');
      throw error;
    } finally {
      commit('setLoading', false);
    }
  },
};

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions,
};
