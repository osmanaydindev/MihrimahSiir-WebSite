<script setup lang="ts">
interface Props {
  message?: string
  color?: string
}

const props = withDefaults(defineProps<Props>(), {
  message: 'Daha fazla yükleniyor...',
  color: 'grey-lighten-1'
})
</script>

<template>
  <div class="loading-more-container">
    <v-card class="loading-more-card">
      <v-card-text class="d-flex align-center justify-center ga-3">
        <v-progress-circular
          indeterminate
          :color="color"
          size="24"
          width="2"
        ></v-progress-circular>
        <span class="text-grey-lighten-1">{{ message }}</span>
      </v-card-text>
    </v-card>
  </div>
</template>

<style scoped>
.loading-more-container {
  display: flex;
  justify-content: center;
  margin-top: 32px;
  margin-bottom: 32px;
}

.loading-more-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  padding: 16px 32px;
  border-radius: 12px;
}
</style>
