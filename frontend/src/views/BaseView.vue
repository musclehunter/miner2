<template>
  <div class="min-h-screen bg-gray-100 p-8">
    <h1 class="text-3xl font-bold mb-6 text-center text-gray-800">拠点</h1>
    
    <div class="max-w-4xl mx-auto">
      <!-- Gold Display -->
      <div class="mb-6 flex justify-end">
        <PlayerGoldDisplay />
      </div>

      <!-- Inventory Section -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div>
          <PlayerOreList />
        </div>
        <div>
          <PlayerItemList />
        </div>
      </div>

      <!-- Back to Title Button -->
      <div class="flex justify-center mt-8">
        <router-link to="/" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
          タイトルに戻る
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { onMounted } from 'vue';
import { useStore } from 'vuex';
import PlayerGoldDisplay from '@/components/inventory/PlayerGoldDisplay.vue';
import PlayerOreList from '@/components/inventory/PlayerOreList.vue';
import PlayerItemList from '@/components/inventory/PlayerItemList.vue';

export default {
  name: 'BaseView',
  components: {
    PlayerGoldDisplay,
    PlayerOreList,
    PlayerItemList,
  },
  setup() {
    const store = useStore();

    onMounted(() => {
      store.dispatch('inventory/fetchInventory');
    });

    return {};
  },
};
</script>

<style scoped>
/* Scoped styles can be added here if needed */
</style>
