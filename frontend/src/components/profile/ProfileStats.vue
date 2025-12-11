<script setup lang="ts">
import StatCard from '../common/StatCard.vue'

interface Stat {
  icon: string
  iconColor: string
  number: number
  label: string
  clickable?: boolean
}

interface Props {
  stats: Stat[]
}

defineProps<Props>()

const emit = defineEmits<{
  'stat-click': [index: number]
}>()

const handleStatClick = (index: number) => {
  emit('stat-click', index)
}
</script>

<template>
  <div class="profile-stats">
    <StatCard
      v-for="(stat, index) in stats"
      :key="index"
      :icon="stat.icon"
      :icon-color="stat.iconColor"
      :number="stat.number"
      :label="stat.label"
      :clickable="stat.clickable"
      @click="handleStatClick(index)"
    />
  </div>
</template>

<style scoped>
.profile-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr); /* 4 equal columns */
  gap: clamp(16px, 3vw, 24px);
}

/* Tablet: 2 columns */
@media (max-width: 900px) {
  .profile-stats {
    grid-template-columns: repeat(2, 1fr);
  }
}

/* Mobile: 1 column */
@media (max-width: 480px) {
  .profile-stats {
    grid-template-columns: 1fr;
  }
}
</style>
