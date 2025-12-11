<script setup lang="ts">
interface Props {
  modelValue: boolean
  message: string
  color?: 'success' | 'error' | 'warning' | 'info'
  timeout?: number
  location?: string
}

const props = withDefaults(defineProps<Props>(), {
  color: 'success',
  timeout: 4000,
  location: 'top right'
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const iconMap = {
  success: 'mdi-check-circle',
  error: 'mdi-alert-circle',
  warning: 'mdi-alert',
  info: 'mdi-information'
}
</script>

<template>
  <v-snackbar
    :model-value="modelValue"
    @update:model-value="emit('update:modelValue', $event)"
    :color="color"
    :timeout="timeout"
    :location="location"
    variant="elevated"
    class="custom-snackbar"
  >
    <div class="d-flex align-center">
      <v-icon
        :icon="iconMap[color]"
        size="24"
        class="mr-3"
      />
      <span class="text-body-1">{{ message }}</span>
    </div>
    <template v-slot:actions>
      <v-btn
        icon="mdi-close"
        size="small"
        variant="text"
        @click="emit('update:modelValue', false)"
      />
    </template>
  </v-snackbar>
</template>

<style scoped>
.custom-snackbar :deep(.v-snackbar__wrapper) {
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(10px);
}

.custom-snackbar :deep(.v-snackbar__content) {
  padding: 16px 20px;
}
</style>
