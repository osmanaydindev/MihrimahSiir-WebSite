<script setup lang="ts">
import { computed } from 'vue'
import { useDisplay } from 'vuetify'

interface Props {
  icon: string
  iconColor?: string
  title: string
  message: string
  buttonText?: string
  buttonIcon?: string
  buttonRoute?: string
  buttonColor?: string
}

const props = withDefaults(defineProps<Props>(), {
  iconColor: 'grey-lighten-1',
  buttonColor: 'grey-lighten-1'
})

const emit = defineEmits<{
  'button-click': []
}>()

const { xs, sm } = useDisplay()

// Responsive icon size
const iconSize = computed(() => {
  if (xs.value) return 60
  if (sm.value) return 80
  return 96
})

const handleButtonClick = () => {
  emit('button-click')
}
</script>

<template>
  <div class="empty-state-container">
    <v-card class="text-center empty-state-card elevation-8" rounded="xl">
      <v-icon :size="iconSize" :color="iconColor" class="mb-4 empty-icon">
        {{ icon }}
      </v-icon>
      <v-card-title class="empty-title text-white font-weight-bold mb-2">
        {{ title }}
      </v-card-title>
      <v-card-text class="empty-message text-grey-lighten-1 mb-4" v-html="message">
      </v-card-text>
      <v-card-actions v-if="buttonText" class="d-flex justify-center pa-0">
        <v-btn
          :color="buttonColor"
          variant="outlined"
          size="large"
          rounded="lg"
          @click="handleButtonClick"
          class="text-none action-btn"
        >
          <v-icon v-if="buttonIcon" size="20" class="mr-2">{{ buttonIcon }}</v-icon>
          <span class="d-none d-sm-inline">{{ buttonText }}</span>
          <span class="d-inline d-sm-none">{{ buttonText.split(' ').slice(-1)[0] }}</span>
        </v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<style scoped>
.empty-state-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 300px);
  padding: 40px 20px;
  width: 100%;
}

.empty-state-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  padding: 40px 32px;
  border-radius: 12px;
  text-align: center;
  max-width: 500px;
  width: 100%;
  margin: 0 auto;
  transition: all 0.3s ease;
  box-sizing: border-box;
}

.empty-state-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(76, 175, 80, 0.15);
  border-color: currentColor;
}

.empty-icon {
  opacity: 0.8;
  transition: opacity 0.3s ease;
}

.empty-state-card:hover .empty-icon {
  opacity: 1;
}

/* Responsive title and message */
.empty-title {
  font-size: clamp(1.25rem, 4vw, 1.5rem);
  word-wrap: break-word;
  overflow-wrap: break-word;
  hyphens: auto;
}

.empty-message {
  font-size: clamp(0.875rem, 3vw, 1rem);
  word-wrap: break-word;
  overflow-wrap: break-word;
  hyphens: auto;
}

.action-btn {
  max-width: 100%;
}

/* Mobile devices (< 600px) */
@media (max-width: 599px) {
  .empty-state-container {
    padding: 24px 16px;
    min-height: calc(100vh - 200px);
  }

  .empty-state-card {
    padding: 32px 20px;
    max-width: 100%;
  }

  .empty-title {
    font-size: 1.25rem;
  }

  .empty-message {
    font-size: 0.875rem;
  }
}

/* Extra small mobile devices (< 400px) */
@media (max-width: 399px) {
  .empty-state-container {
    padding: 16px 12px;
  }

  .empty-state-card {
    padding: 24px 16px;
  }

  .empty-title {
    font-size: 1.125rem;
  }

  .empty-message {
    font-size: 0.8125rem;
  }
}
</style>
