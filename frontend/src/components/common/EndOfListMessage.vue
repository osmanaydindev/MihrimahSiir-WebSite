<script setup lang="ts">
interface Props {
  message?: string
  iconColor?: string
}

const props = withDefaults(defineProps<Props>(), {
  message: 'Tüm öğeler gösteriliyor',
  iconColor: 'grey-darken-1'
})
</script>

<template>
  <div class="end-message text-center mt-8 mb-4">
    <v-icon size="20" :color="iconColor" class="mr-2">mdi-check-circle</v-icon>
    <span class="text-grey-darken-1">{{ message }}</span>
  </div>
</template>

<style scoped>
.end-message {
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.3px;
}
</style>
