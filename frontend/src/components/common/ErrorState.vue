<script setup lang="ts">
interface Props {
  error: string
}

defineProps<Props>()
</script>

<template>
  <div class="error-container">
    <v-card class="error-card">
      <v-card-text class="text-center">
        <v-icon size="64" color="error" class="mb-4">mdi-alert-circle</v-icon>
        <p class="text-error text-h6">{{ error }}</p>
      </v-card-text>
    </v-card>
  </div>
</template>

<style scoped>
.error-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 300px);
  padding: 64px 0;
}

.error-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  padding: 48px;
  border-radius: 12px;
  text-align: center;
  max-width: 500px;
  transition: all 0.3s ease;
}
</style>
