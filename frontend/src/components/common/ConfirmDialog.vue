<script setup lang="ts">
interface Props {
  modelValue: boolean
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  confirmColor?: string
  icon?: string
  iconColor?: string
}

const props = withDefaults(defineProps<Props>(), {
  confirmText: 'Onayla',
  cancelText: 'İptal',
  confirmColor: 'error',
  icon: 'mdi-help-circle',
  iconColor: 'warning'
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'confirm': []
  'cancel': []
}>()

const handleConfirm = () => {
  emit('confirm')
  emit('update:modelValue', false)
}

const handleCancel = () => {
  emit('cancel')
  emit('update:modelValue', false)
}
</script>

<template>
  <v-dialog
    :model-value="modelValue"
    @update:model-value="emit('update:modelValue', $event)"
    max-width="400"
    persistent
  >
    <v-card class="confirm-card">
      <v-card-title class="d-flex align-center pa-4">
        <v-icon :icon="icon" size="28" :color="iconColor" class="mr-3" />
        <span>{{ title }}</span>
      </v-card-title>
      <v-divider />
      <v-card-text class="pa-4 text-body-1">
        {{ message }}
      </v-card-text>
      <v-divider />
      <v-card-actions class="pa-4">
        <v-spacer />
        <v-btn
          color="grey-darken-1"
          variant="text"
          @click="handleCancel"
        >
          {{ cancelText }}
        </v-btn>
        <v-btn
          :color="confirmColor"
          variant="flat"
          @click="handleConfirm"
        >
          {{ confirmText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<style scoped>
.confirm-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  border-radius: 12px;
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.5);
}

.confirm-card .v-card-title {
  font-size: 18px;
  font-weight: 600;
  color: #ffffff;
}

.confirm-card .v-card-text {
  color: #e0e0e0;
  line-height: 1.6;
}
</style>
