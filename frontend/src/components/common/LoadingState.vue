<script setup lang="ts">
interface Props {
  message?: string
  color?: string
}

const props = withDefaults(defineProps<Props>(), {
  message: 'Yükleniyor...',
  color: 'grey-lighten-1'
})
</script>

<template>
  <div class="loading-container">
    <v-card class="loading-card">
      <v-card-text class="d-flex flex-column align-center justify-center ga-4">
        <v-progress-circular
          indeterminate
          :color="color"
          size="48"
          width="3"
        ></v-progress-circular>
        <span class="text-grey-lighten-1 text-h6">{{ message }}</span>
      </v-card-text>
    </v-card>
  </div>
</template>

<style scoped>
.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 300px);
  padding: 64px 0;
}

.loading-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  padding: 48px;
  border-radius: 12px;
  text-align: center;
  max-width: 500px;
  transition: all 0.3s ease;
}
</style>
