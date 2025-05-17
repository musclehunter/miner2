<template>
  <div class="card backdrop-blur-sm bg-dark-lighter/80 border border-dark-light">
    <!-- 町の詳細情報 -->
    <div class="flex flex-col md:flex-row gap-6">
      <!-- 町のアイコン（大きめ） -->
      <div class="w-24 h-24 rounded-lg bg-dark-light flex items-center justify-center mx-auto md:mx-0">
        <img :src="require(`@/assets/images/towns/${town.imageUrl}`)" alt="町のアイコン" class="w-20 h-20" />
      </div>

      <!-- 町の詳細情報 -->
      <div class="flex-1">
        <h2 class="text-2xl font-bold text-light mb-2">{{ town.name }}</h2>
        <p class="text-light-dark mb-4">{{ town.description }}</p>

        <!-- 特徴リスト -->
        <h3 class="text-lg font-semibold text-light-dark mb-2">特徴</h3>
        <ul class="list-disc list-inside text-light-dark mb-4 space-y-1">
          <li v-for="(feature, index) in town.features" :key="index">{{ feature }}</li>
        </ul>

        <!-- リソース詳細 -->
        <h3 class="text-lg font-semibold text-light-dark mb-2">資源</h3>
        <div class="grid grid-cols-2 md:grid-cols-5 gap-2 mb-4">
          <div class="flex flex-col items-center">
            <div class="w-8 h-8 rounded-full flex items-center justify-center" 
                 :class="getResourceColorClass('iron', town.resources.iron)">
              <span class="text-sm font-bold">Fe</span>
            </div>
            <span class="text-xs mt-1 text-light-dark">鉄: {{ getResourceLabel(town.resources.iron) }}</span>
          </div>
          <div class="flex flex-col items-center">
            <div class="w-8 h-8 rounded-full flex items-center justify-center" 
                 :class="getResourceColorClass('copper', town.resources.copper)">
              <span class="text-sm font-bold">Cu</span>
            </div>
            <span class="text-xs mt-1 text-light-dark">銅: {{ getResourceLabel(town.resources.copper) }}</span>
          </div>
          <div class="flex flex-col items-center">
            <div class="w-8 h-8 rounded-full flex items-center justify-center" 
                 :class="getResourceColorClass('silver', town.resources.silver)">
              <span class="text-sm font-bold">Ag</span>
            </div>
            <span class="text-xs mt-1 text-light-dark">銀: {{ getResourceLabel(town.resources.silver) }}</span>
          </div>
          <div class="flex flex-col items-center">
            <div class="w-8 h-8 rounded-full flex items-center justify-center" 
                 :class="getResourceColorClass('gold', town.resources.gold)">
              <span class="text-sm font-bold">Au</span>
            </div>
            <span class="text-xs mt-1 text-light-dark">金: {{ getResourceLabel(town.resources.gold) }}</span>
          </div>
          <div class="flex flex-col items-center">
            <div class="w-8 h-8 rounded-full flex items-center justify-center" 
                 :class="getResourceColorClass('crystal', town.resources.crystal)">
              <span class="text-sm font-bold">Cr</span>
            </div>
            <span class="text-xs mt-1 text-light-dark">結晶: {{ getResourceLabel(town.resources.crystal) }}</span>
          </div>
        </div>

        <!-- 市場情報 -->
        <h3 class="text-lg font-semibold text-light-dark mb-2">市場情報</h3>
        <div class="flex flex-wrap gap-4 text-light-dark">
          <div>
            <span class="font-semibold">購買力:</span> 
            <span class="inline-flex">
              <span v-for="i in 5" :key="i" class="text-lg">
                <span v-if="i <= town.marketTraits.buyingPower" class="text-yellow-400">$</span>
                <span v-else class="text-gray-600">$</span>
              </span>
            </span>
          </div>
          <div>
            <span class="font-semibold">価格変動:</span> 
            <span class="inline-flex">
              <span v-for="i in 5" :key="i" class="text-lg">
                <span v-if="i <= town.marketTraits.priceVariation" class="text-primary">↕</span>
                <span v-else class="text-gray-600">↕</span>
              </span>
            </span>
          </div>
          <div>
            <span class="font-semibold">特別需要:</span> 
            <span class="ml-1 text-primary">
              {{ getResourceName(town.marketTraits.specialtyDemand) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 分離線 -->
    <div class="border-t border-dark-light my-4"></div>

    <!-- 基本情報と拠点設立ボタン -->
    <div class="flex flex-col md:flex-row justify-between items-center gap-4">
      <div class="text-light-dark">
        <div class="flex items-center">
          <span class="text-sm mr-4">
            <span class="font-semibold">人口:</span> {{ town.population.toLocaleString() }}人
          </span>
          <span class="text-sm px-2 py-1 rounded-full" 
                :class="{
                  'bg-green-500/20 text-green-400': town.difficultyLevel <= 2,
                  'bg-yellow-500/20 text-yellow-400': town.difficultyLevel === 3,
                  'bg-red-500/20 text-red-400': town.difficultyLevel >= 4
                }">
            {{ getDifficultyLabel(town.difficultyLevel) }}
          </span>
        </div>
      </div>
      <button class="btn btn-primary" @click="$emit('establish-base', town.id)">
        この町に拠点を設立する
      </button>
    </div>
  </div>
</template>

<script>
import { getResourceLabel, getDifficultyLabel } from '@/data/towns';

export default {
  name: 'TownDetail',
  props: {
    town: {
      type: Object,
      required: true
    }
  },
  methods: {
    getResourceLabel,
    getDifficultyLabel,
    getResourceName(resourceType) {
      const resourceNames = {
        iron: '鉄',
        copper: '銅',
        silver: '銀',
        gold: '金',
        crystal: '結晶'
      };
      return resourceNames[resourceType] || resourceType;
    },
    getResourceColorClass(type, value) {
      // リソースの量に応じた色を返す
      const baseClass = 'bg-dark-light text-light-dark';
      
      if (value <= 2) return baseClass;
      
      const colorClasses = {
        iron: 'bg-primary/20 text-primary',
        copper: 'bg-secondary/20 text-secondary',
        silver: 'bg-light/20 text-light',
        gold: 'bg-yellow-400/20 text-yellow-400',
        crystal: 'bg-purple-400/20 text-purple-400'
      };
      
      return colorClasses[type] || baseClass;
    }
  }
}
</script>
