<template>
  <div class="min-h-screen bg-dark bg-[url('@/assets/images/world-map-bg.svg')] bg-cover bg-center">
    <!-- ヘッダー部分 -->
    <header class="p-4 bg-dark-lighter/80 backdrop-blur-sm border-b border-dark-light">
      <div class="container mx-auto flex justify-between items-center">
        <h1 class="text-2xl font-bold text-light">ワールドマップ</h1>
        <div class="flex items-center gap-4">
          <router-link to="/" class="text-light-dark hover:text-light transition-colors">
            <span class="hidden md:inline mr-1">タイトルに戻る</span>
            <span>🏠</span>
          </router-link>
        </div>
      </div>
    </header>

    <div class="container mx-auto py-6 px-4">
      <!-- 説明テキスト -->
      <div class="mb-6 text-center">
        <h2 class="text-xl md:text-2xl font-bold text-light mb-2">町を選択して拠点を設営する</h2>
        <p class="text-light-dark">各町には特徴があり、採掘できる資源や市場の状況が異なります。自分に合った町を選んで採掘事業を始めましょう。</p>
      </div>

      <!-- メインコンテンツ -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <!-- 左側：町の一覧 -->
        <div class="md:col-span-1 space-y-4">
          <h3 class="text-lg font-bold text-light mb-3 pl-2 border-l-4 border-primary">町の一覧</h3>
          
          <div class="space-y-3 max-h-[calc(100vh-250px)] overflow-y-auto pr-2">
            <town-card 
              v-for="town in towns"
              :key="town.id"
              :town="town"
              :selected="selectedTownId === town.id"
              @select="selectTown"
            />
          </div>
        </div>

        <!-- 右側：選択した町の詳細 -->
        <div class="md:col-span-2">
          <h3 class="text-lg font-bold text-light mb-3 pl-2 border-l-4 border-primary">町の詳細情報</h3>
          
          <div v-if="selectedTown">
            <town-detail 
              :town="selectedTown"
              @establish-base="establishBase"
            />
          </div>
          
          <div v-else class="card backdrop-blur-sm bg-dark-lighter/80 border border-dark-light p-8 text-center">
            <p class="text-light-dark mb-4">左側から町を選択すると、詳細情報が表示されます</p>
            <img src="@/assets/images/towns/iron-hill.svg" alt="町の選択" class="w-20 h-20 mx-auto opacity-50" />
          </div>
        </div>
      </div>
    </div>

    <!-- 拠点設立確認モーダル -->
    <div v-if="showEstablishModal" class="fixed inset-0 bg-black/70 flex items-center justify-center p-4 z-50">
      <div class="card max-w-md w-full">
        <h3 class="text-xl font-bold text-light mb-4">拠点設立の確認</h3>
        <p class="text-light-dark mb-6">{{ selectedTown.name }}に拠点を設立します。よろしいですか？</p>
        <p class="text-sm text-light-dark mb-6">※一度設立すると変更できません</p>
        
        <div class="flex justify-end gap-3">
          <button 
            class="btn bg-dark-light text-light hover:bg-dark-light/80" 
            @click="showEstablishModal = false"
          >
            キャンセル
          </button>
          <button 
            class="btn btn-primary" 
            @click="confirmEstablishBase"
          >
            設立する
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import TownCard from '@/components/world/TownCard.vue';
import TownDetail from '@/components/world/TownDetail.vue';
import { towns, getTownById } from '@/data/towns';

export default {
  name: 'WorldMapView',
  components: {
    TownCard,
    TownDetail
  },
  setup() {
    const router = useRouter();
    const selectedTownId = ref(null);
    const showEstablishModal = ref(false);

    // 選択中の町の情報
    const selectedTown = computed(() => {
      if (!selectedTownId.value) return null;
      return getTownById(selectedTownId.value);
    });

    // 町を選択
    const selectTown = (townId) => {
      selectedTownId.value = townId;
    };

    // 拠点設立ダイアログ表示
    const establishBase = (townId) => {
      selectedTownId.value = townId;
      showEstablishModal.value = true;
    };

    // 拠点設立を確定
    const confirmEstablishBase = () => {
      // TODO: APIと連携して拠点を設立する
      // ここではモック処理として遷移だけ行う
      showEstablishModal.value = false;
      
      // 拠点画面に遷移
      router.push('/base');
    };

    return {
      towns,
      selectedTownId,
      selectedTown,
      showEstablishModal,
      selectTown,
      establishBase,
      confirmEstablishBase
    };
  }
}
</script>
