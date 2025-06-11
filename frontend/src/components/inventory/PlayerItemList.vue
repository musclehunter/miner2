<template>
  <div class="item-list-container">
    <h3>所持アイテム</h3>
    <div v-if="isLoading" class="loading-indicator">読み込み中...</div>
    <div v-else-if="error" class="error-message">{{ error }}</div>
    <table v-else-if="items.length > 0" class="item-table">
      <thead>
        <tr>
          <th>アイテム名</th>
          <th>数量</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in items" :key="item.item_id">
          <td>{{ item.name }}</td>
          <td>{{ item.quantity }}</td>
        </tr>
      </tbody>
    </table>
    <div v-else class="no-data-message">所持しているアイテムはありません。</div>
  </div>
</template>

<script>
import { computed } from 'vue';
import { useStore } from 'vuex';

export default {
  name: 'PlayerItemList',
  setup() {
    const store = useStore();

    const items = computed(() => store.getters['inventory/getPlayerItems']);
    const isLoading = computed(() => store.getters['inventory/isLoadingInventory']);
    const error = computed(() => store.getters['inventory/getInventoryError']);

    return {
      items,
      isLoading,
      error,
    };
  },
};
</script>

<style scoped>
.item-list-container {
  padding: 16px;
  border: 1px solid #ccc;
  border-radius: 8px;
  background-color: #fff;
}

h3 {
  margin-top: 0;
  border-bottom: 2px solid #333;
  padding-bottom: 8px;
}

.loading-indicator, .error-message, .no-data-message {
  padding: 16px;
  text-align: center;
  color: #666;
}

.error-message {
  color: #d9534f;
}

.item-table {
  width: 100%;
  border-collapse: collapse;
}

.item-table th, .item-table td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}

.item-table th {
  background-color: #f2f2f2;
}
</style>
