<script setup lang="ts">
interface Props {
  icon: string
  iconColor?: string
  number: number | string
  label: string
  clickable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  iconColor: 'grey-lighten-1',
  clickable: false
})

const emit = defineEmits<{
  'click': []
}>()
</script>

<template>
  <div
    class="stat-card"
    :class="{ 'stat-card-clickable': clickable }"
    @click="clickable && emit('click')"
  >
    <v-icon :size="32" :color="iconColor">{{ icon }}</v-icon>
    <div class="stat-info">
      <span class="stat-number">{{ number }}</span>
      <span class="stat-label">{{ label }}</span>
    </div>
  </div>
</template>

<style scoped>
.stat-card {
  display: flex;
  align-items: center;
  gap: clamp(16px, 3vw, 20px);
  padding: clamp(20px, 4vw, 24px);
  background: linear-gradient(135deg, #1e1e1e 0%, #2a2a2a 100%);
  border: 1px solid #424242;
  border-radius: clamp(12px, 2vw, 16px);
  transition: all 0.3s ease;
  width: 100%; /* Full width of grid cell */
  box-sizing: border-box;
  min-width: 0; /* Allow flex items to shrink */
}

.stat-card:hover {
  transform: translateY(-4px);
  border-color: #616161;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}

.stat-card-clickable {
  cursor: pointer;
}

.stat-card-clickable:hover {
  border-color: #757575;
}

.stat-info {
  display: flex;
  flex-direction: column;
  flex: 1; /* Take remaining space */
  min-width: 0; /* Allow text truncation */
}

.stat-number {
  font-size: clamp(24px, 4vw, 32px);
  font-weight: 700;
  color: #ffffff;
  line-height: 1;
  margin-bottom: 4px;
}

.stat-label {
  font-size: clamp(13px, 2vw, 14px);
  color: #9e9e9e;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
