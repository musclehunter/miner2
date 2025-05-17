<template>
  <div class="card hover:border-primary hover:shadow-md hover:shadow-primary/20 transition-all duration-300 cursor-pointer"
       :class="{ 'border-primary': selected }"
       @click="$emit('select', town.id)">
    <div class="flex items-start gap-4">
      <!-- 町のアイコン -->
      <div class="w-16 h-16 rounded-lg bg-dark-light flex items-center justify-center">
        <img :src="require(`@/assets/images/towns/${town.imageUrl}`)" alt="町のアイコン" class="w-12 h-12" />
      </div>
      
      <!-- 町の情報 -->
      <div class="flex-1">
        <h3 class="text-xl font-bold text-light mb-1">{{ town.name }}</h3>
        <p class="text-sm text-light-dark mb-2">{{ town.description }}</p>
        
        <!-- リソース情報 -->
        <div class="flex flex-wrap gap-2 text-xs">
          <span class="px-2 py-1 rounded-full bg-dark-light text-light-dark" 
                v-if="town.resources.iron > 2"
                :class="{'bg-primary/20 text-primary': town.resources.iron > 3}">
            鉄: {{ getResourceLabel(town.resources.iron) }}
          </span>
          <span class="px-2 py-1 rounded-full bg-dark-light text-light-dark" 
                v-if="town.resources.copper > 2"
                :class="{'bg-secondary/20 text-secondary': town.resources.copper > 3}">
            銅: {{ getResourceLabel(town.resources.copper) }}
          </span>
          <span class="px-2 py-1 rounded-full bg-dark-light text-light-dark" 
                v-if="town.resources.silver > 2"
                :class="{'bg-light/20 text-light': town.resources.silver > 3}">
            銀: {{ getResourceLabel(town.resources.silver) }}
          </span>
          <span class="px-2 py-1 rounded-full bg-dark-light text-light-dark" 
                v-if="town.resources.gold > 2"
                :class="{'bg-yellow-400/20 text-yellow-400': town.resources.gold > 3}">
            金: {{ getResourceLabel(town.resources.gold) }}
          </span>
          <span class="px-2 py-1 rounded-full bg-dark-light text-light-dark" 
                v-if="town.resources.crystal > 2"
                :class="{'bg-purple-400/20 text-purple-400': town.resources.crystal > 3}">
            結晶: {{ getResourceLabel(town.resources.crystal) }}
          </span>
        </div>
      </div>
    </div>
    
    <!-- 難易度表示 -->
    <div class="mt-3 flex justify-between items-center">
      <span class="text-xs text-light-dark">
        人口: {{ town.population.toLocaleString() }}人
      </span>
      <span class="text-xs px-2 py-1 rounded-full" 
            :class="{
              'bg-green-500/20 text-green-400': town.difficultyLevel <= 2,
              'bg-yellow-500/20 text-yellow-400': town.difficultyLevel === 3,
              'bg-red-500/20 text-red-400': town.difficultyLevel >= 4
            }">
        {{ getDifficultyLabel(town.difficultyLevel) }}
      </span>
    </div>
  </div>
</template>

<script>
import { getResourceLabel, getDifficultyLabel } from '@/data/towns';

export default {
  name: 'TownCard',
  props: {
    town: {
      type: Object,
      required: true
    },
    selected: {
      type: Boolean,
      default: false
    }
  },
  methods: {
    getResourceLabel,
    getDifficultyLabel
  }
}
</script>
