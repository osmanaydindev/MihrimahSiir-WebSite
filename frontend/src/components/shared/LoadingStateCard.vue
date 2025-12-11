<script setup lang="ts">
interface Props {
  text?: string
  size?: number
  width?: number
  color?: string
}

withDefaults(defineProps<Props>(), {
  text: 'Yükleniyor',
  size: 24,
  width: 2,
  color: 'grey-lighten-1'
})
</script>

<template>
  <div class="loading-container">
    <v-card class="loading-card">
      <v-card-text class="d-flex align-center justify-center ga-3">
        <v-progress-circular
          indeterminate
          :color="color"
          :size="size"
          :width="width"
        ></v-progress-circular>
        <span :class="`text-${color}`">{{ text }}...</span>
      </v-card-text>
    </v-card>
  </div>
</template>

<style scoped>
.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 48px 0;
}

.loading-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  padding: 24px;
  border-radius: 12px;
}
</style>
