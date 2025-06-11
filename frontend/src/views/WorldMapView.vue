<template>
  <div class="min-h-screen bg-dark bg-cover bg-center">
    <!-- コンテンツタイトル -->
    <div class="p-4">
      <div class="container mx-auto">
        <h1 class="text-2xl font-bold text-light">ワールドマップ</h1>
      </div>
    </div>

    <div class="container mx-auto py-6 px-4">
      <!-- 説明テキスト -->
      <div class="mb-6 text-center">
        <h2 class="text-xl md:text-2xl font-bold text-light mb-2">町を選択して拠点を設営する</h2>
        <p class="text-light-dark">各町には特徴があり、採掘できる資源や市場の状況が異なります。自分に合った町を選んで採掘事業を始めましょう。</p>
      </div>

      <!-- メインコンテンツ -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- 左側：世界地図 -->
        <div class="lg:col-span-2 relative">
          <div class="card backdrop-blur-sm bg-dark-lighter/80 border border-dark-light p-4 mb-4">
            <h3 class="text-lg font-bold text-light mb-3 pl-2 border-l-4 border-primary">ワールドマップ</h3>
            
            <!-- 地図エリア -->
            <div class="relative world-map-container" ref="mapContainer">
              <img 
                src="/worldmap.png" 
                alt="ワールドマップ" 
                class="w-full h-auto" 
                @click="handleMapClick"
              />
              
              <!-- マップ上の町マーカー -->
              <div 
                v-for="town in towns" 
                :key="town.id"
                class="town-marker absolute cursor-pointer transition-all hover:scale-125"
                :class="{ 'active': selectedTownId === town.id }"
                :style="getMarkerStyle(town)"
                @click.stop="selectTown(town.id)"
              >
                <div class="w-6 h-6 rounded-full bg-primary flex items-center justify-center"
                     :title="town.name">
                  <span class="text-xs font-bold">{{ town.id }}</span>
                </div>
                <div class="town-name absolute whitespace-nowrap bg-dark-lighter px-2 py-1 rounded text-xs font-bold">
                  {{ town.name }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 右側：選択した町の詳細 -->
        <div class="lg:col-span-1">
          <div class="card backdrop-blur-sm bg-dark-lighter/80 border border-dark-light p-4">
            <h3 class="text-lg font-bold text-light mb-3 pl-2 border-l-4 border-primary">町の詳細情報</h3>
            
            <div v-if="selectedTown">
              <town-detail 
                :town="selectedTown"
                @establish-base="establishBase"
              />
            </div>
            
            <div v-else class="p-8 text-center">
              <p class="text-light-dark mb-4">マップ上の町を選択すると、詳細情報が表示されます</p>
              <img src="@/assets/images/towns/iron-hill.svg" alt="町の選択" class="w-20 h-20 mx-auto opacity-50" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 拠点設立確認モーダル -->
    <div v-if="showEstablishModal" class="fixed inset-0 bg-black/70 flex items-center justify-center p-4 z-50">
      <div class="card max-w-md w-full">
        <h3 class="text-xl font-bold text-light mb-4">拠点設立の確認</h3>
        <p v-if="selectedTown" class="text-light-dark mb-6">
          <span class="font-bold text-primary">{{ selectedTown.name }}</span> に拠点を設立します。よろしいですか？
        </p>
        <p class="text-sm text-light-dark mb-6">※一度設立すると変更できません</p>
        <p v-if="creationError" class="text-danger text-sm mb-4">{{ creationError }}</p>
        
        <div class="flex justify-end gap-3">
          <button 
            class="btn bg-dark-light text-light hover:bg-dark-light/80" 
            @click="showEstablishModal = false"
            :disabled="isCreatingBase"
          >
            キャンセル
          </button>
          <button 
            class="btn btn-primary" 
            @click="confirmEstablishBase"
            :disabled="isCreatingBase"
          >
            <span v-if="isCreatingBase">設立中...</span>
            <span v-else>設立する</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue';
import { useStore } from 'vuex';
import { useRouter } from 'vue-router';
import TownDetail from '@/components/world/TownDetail.vue';
import { towns, getTownById } from '@/data/towns';

export default {
  name: 'WorldMapView',
  components: {
    TownDetail
  },
  setup() {
    const router = useRouter();
    const store = useStore();
    const selectedTownId = ref(null);
    const showEstablishModal = ref(false);
    const mapContainer = ref(null);
    const mapScale = ref(1); // 地図のスケール係数（実際の画像サイズと表示サイズの比率）

    const isCreatingBase = ref(false);
    const creationError = ref(null);

    // 選択中の町の情報
    const selectedTown = computed(() => {
      if (!selectedTownId.value) return null;
      return getTownById(selectedTownId.value);
    });

    // マップがロードされた後に実行
    onMounted(() => {
      // 画像のロードが完了したらスケールを計算
      const mapImage = new Image();
      mapImage.onload = () => {
        if (mapContainer.value) {
          const containerWidth = mapContainer.value.clientWidth;
          mapScale.value = containerWidth / mapImage.width;
        }
      };
      mapImage.src = '/worldmap.png';
    });

    // マーカースタイルを計算（町の座標をもとに位置を決定）
    const getMarkerStyle = (town) => {
      const { position } = town;
      // 地図上の相対位置をCSSの位置に変換
      return {
        left: `${position.x * mapScale.value}px`,
        top: `${position.y * mapScale.value}px`,
        transform: 'translate(-50%, -50%)' // マーカーの中心が座標に来るように調整
      };
    };

    // 地図上でのクリック処理
    const handleMapClick = (event) => {
      // クリック位置を取得
      const rect = event.target.getBoundingClientRect();
      const x = event.clientX - rect.left;
      const y = event.clientY - rect.top;
      
      // 実際の座標に変換（スケールで割る）
      const actualX = x / mapScale.value;
      const actualY = y / mapScale.value;
      
      console.log(`Map clicked at coordinates: x=${actualX}, y=${actualY}`);
      
      // 近い町を探す
      findNearestTown(actualX, actualY);
    };
    
    // 指定した座標に最も近い町を見つける
    const findNearestTown = (x, y) => {
      let nearestTown = null;
      let minDistance = Number.MAX_VALUE;
      
      towns.forEach(town => {
        const dx = town.position.x - x;
        const dy = town.position.y - y;
        const distance = Math.sqrt(dx * dx + dy * dy);
        
        // 最小距離を更新
        if (distance < minDistance) {
          minDistance = distance;
          nearestTown = town;
        }
      });
      
      // 最も近い町があり、ある程度近い場合（ここでは50ピクセル以内）
      if (nearestTown && minDistance < 50) {
        selectTown(nearestTown.id);
      }
    };

    // 町を選択
    const selectTown = (townId) => {
      selectedTownId.value = townId;
    };

    // 拠点設立ダイアログ表示
    const establishBase = (townId) => {
      selectTown(townId);
      creationError.value = null;
      showEstablishModal.value = true;
    };

    // 拠点設立を確定
    const confirmEstablishBase = async () => {
      isCreatingBase.value = true;
      creationError.value = null;

      try {
        await store.dispatch('base/createBase', {
          townId: selectedTownId.value,
        });
        
        showEstablishModal.value = false;
        router.push('/base'); // 拠点画面へ遷移

      } catch (error) {
        creationError.value = store.getters['base/error'] || '拠点設立中に不明なエラーが発生しました。';
      } finally {
        isCreatingBase.value = false;
      }
    };

    return {
      towns,
      selectedTownId,
      selectedTown,
      showEstablishModal,
      mapContainer,
      getMarkerStyle,
      handleMapClick,
      selectTown,
      establishBase,
      confirmEstablishBase,
      isCreatingBase,
      creationError,
    };
  },
};
</script>

<style scoped>
.world-map-container {
  position: relative;
  max-width: 100%;
  overflow: hidden;
}

.town-marker {
  z-index: 10;
}

.town-marker.active div:first-child {
  border: 2px solid #ffffff;
  box-shadow: 0 0 8px rgba(255, 255, 255, 0.8);
}

.town-name {
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  opacity: 0;
  transition: opacity 0.2s ease, transform 0.2s ease;
  pointer-events: none;
  z-index: 20;
  white-space: nowrap;
}

.town-marker:hover .town-name,
.town-marker.active .town-name {
  opacity: 1;
  transform: translateX(-50%) translateY(-5px);
}
</style>
