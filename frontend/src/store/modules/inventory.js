import { getMyInventory } from '@/services/inventoryService';

const state = {
  gold: 0,
  ores: [],
  items: [],
  isLoadingInventory: false,
  inventoryError: null,
};

const getters = {
  getPlayerGold: (state) => state.gold,
  getPlayerOres: (state) => state.ores,
  getPlayerItems: (state) => state.items,
  isLoadingInventory: (state) => state.isLoadingInventory,
  getInventoryError: (state) => state.inventoryError,
};

const mutations = {
  SET_INVENTORY_DATA(state, { inventory, ores, items }) {
    state.gold = inventory.gold || 0;
    state.ores = ores || [];
    state.items = items || [];
  },
  SET_INVENTORY_LOADING(state, status) {
    state.isLoadingInventory = status;
  },
  SET_INVENTORY_ERROR(state, error) {
    state.inventoryError = error;
  },
  CLEAR_INVENTORY_ERROR(state) {
    state.inventoryError = null;
  },
};

const actions = {
  async fetchInventory({ commit }) {
    commit('SET_INVENTORY_LOADING', true);
    commit('CLEAR_INVENTORY_ERROR');
    try {
      const data = await getMyInventory();
      // APIレスポンスの構造に合わせて調整
      // project_config.json と api.md を見ると、レスポンスは以下のような構造になっている
      // {
      //   "inventory": { "gold": 1000, ... },
      //   "ores": [ { "ore_id": "ore1", "name": "Iron Ore", "quantity": 10 }, ... ],
      //   "items": [ { "item_id": "item1", "name": "Pickaxe", "quantity": 1 }, ... ]
      // }
      commit('SET_INVENTORY_DATA', {
        inventory: data.inventory || { gold: 0 }, // inventory オブジェクトが存在しない場合に備える
        ores: data.ores || [],
        items: data.items || [],
      });
    } catch (error) {
      commit('SET_INVENTORY_ERROR', error.response?.data?.error || '在庫情報の取得に失敗しました');
      console.error('Error fetching inventory:', error);
    }
    commit('SET_INVENTORY_LOADING', false);
  },
};

export default {
  namespaced: true, // モジュールが名前空間を持つことを示す
  state,
  getters,
  mutations,
  actions,
};
