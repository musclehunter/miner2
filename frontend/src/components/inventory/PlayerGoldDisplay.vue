<template>
  <div class="gold-display">
    <span class="gold-icon">💰</span>
    <span class="gold-amount">{{ formattedGold }} G</span>
  </div>
</template>

<script>
import { computed } from 'vue';
import { useStore } from 'vuex';

export default {
  name: 'PlayerGoldDisplay',
  setup() {
    const store = useStore();

    const gold = computed(() => store.getters['inventory/getPlayerGold']);

    const formattedGold = computed(() => {
      return gold.value.toLocaleString();
    });

    return {
      formattedGold,
    };
  },
};
</script>

<style scoped>
.gold-display {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background-color: #f9f9f9;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 1.1em;
  font-weight: bold;
  color: #333;
}

.gold-icon {
  margin-right: 8px;
  font-size: 1.4em;
}

.gold-amount {
  color: #DAA520; /* Gold color */
}
</style>
