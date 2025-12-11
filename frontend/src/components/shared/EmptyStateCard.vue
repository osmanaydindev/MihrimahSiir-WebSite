<script setup lang="ts">
import { useRouter } from 'vue-router'

interface Props {
  title: string
  description: string
  icon: string
  iconColor?: string
  buttonText?: string
  buttonRoute?: string
  buttonColor?: string
  showButton?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  iconColor: 'grey-lighten-1',
  buttonColor: 'grey-lighten-1',
  showButton: true
})

const router = useRouter()

const handleButtonClick = () => {
  if (props.buttonRoute) {
    router.push(props.buttonRoute)
  }
}
</script>

<template>
  <div class="empty-state">
    <v-card class="empty-state-card">
      <v-icon
        :size="80"
        :color="iconColor"
        class="mb-4"
      >
        {{ icon }}
      </v-icon>
      <v-card-title class="text-h5 text-white font-weight-bold mb-2">
        {{ title }}
      </v-card-title>
      <v-card-text class="text-body-1 text-grey-lighten-1 mb-4">
        {{ description }}
      </v-card-text>
      <v-card-actions v-if="showButton && buttonText" class="d-flex justify-center">
        <v-btn
          :color="buttonColor"
          variant="outlined"
          @click="handleButtonClick"
        >
          <v-icon class="mr-2">mdi-refresh</v-icon>
          {{ buttonText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<style scoped>
.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 48px 0;
}

.empty-state-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  padding: 48px;
  border-radius: 12px;
  text-align: center;
  max-width: 500px;
  transition: transform 0.3s ease, box-shadow 0.3s ease, border-color 0.3s ease;
}

.empty-state-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
  border-color: #616161;
}
</style>
