<template>
  <div class="min-h-screen bg-dark">
    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </div>
</template>

<script>
export default {
  name: 'App',
  created() {
    // 起動時に認証状態を復元
    this.$store.dispatch('initAuth')
  }
}
</script>

<style>
/* グローバルトランジション */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
