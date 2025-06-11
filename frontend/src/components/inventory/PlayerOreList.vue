<template>
  <div class="ore-list-container">
    <h3>所持鉱石</h3>
    <div v-if="isLoading" class="loading-indicator">読み込み中...</div>
    <div v-else-if="error" class="error-message">{{ error }}</div>
    <table v-else-if="ores.length > 0" class="ore-table">
      <thead>
        <tr>
          <th>鉱石名</th>
          <th>数量</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="ore in ores" :key="ore.ore_id">
          <td>{{ ore.name }}</td>
          <td>{{ ore.quantity }}</td>
        </tr>
      </tbody>
    </table>
    <div v-else class="no-data-message">所持している鉱石はありません。</div>
  </div>
</template>

<script>
import { computed } from 'vue';
import { useStore } from 'vuex';

export default {
  name: 'PlayerOreList',
  setup() {
    const store = useStore();

    const ores = computed(() => store.getters['inventory/getPlayerOres']);
    const isLoading = computed(() => store.getters['inventory/isLoadingInventory']);
    const error = computed(() => store.getters['inventory/getInventoryError']);

    return {
      ores,
      isLoading,
      error,
    };
  },
};
</script>

<style scoped>
.ore-list-container {
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

.ore-table {
  width: 100%;
  border-collapse: collapse;
}

.ore-table th, .ore-table td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}

.ore-table th {
  background-color: #f2f2f2;
}
</style>
