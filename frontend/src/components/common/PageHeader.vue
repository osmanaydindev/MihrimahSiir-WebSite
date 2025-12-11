<script setup lang="ts">
interface Props {
  icon: string
  iconColor?: string
  title: string
  subtitle: string
  total?: number
  totalLabel?: string
}

const props = withDefaults(defineProps<Props>(), {
  iconColor: '#81c784',
  totalLabel: 'öğe'
})
</script>

<template>
  <div class="page-header text-center mb-8">
    <div class="header-icon-wrapper mb-4">
      <v-icon size="48" class="header-icon" :style="{ color: iconColor }">
        {{ icon }}
      </v-icon>
    </div>
    <h1 class="page-title">{{ title }}</h1>
    <p class="page-subtitle text-grey-lighten-1">{{ subtitle }}</p>
    <p v-if="total && total > 0" class="total-count text-grey-darken-1 mt-2">
      Toplam {{ total }} {{ totalLabel }}
    </p>
  </div>
</template>

<style scoped>
.page-header {
  width: 100%;
  animation: fadeInDown 0.6s ease-out;
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.header-icon-wrapper {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 96px;
  height: 96px;
  border-radius: 50%;
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 2px solid #424242;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    border-color: #424242;
  }
  50% {
    transform: scale(1.05);
    border-color: currentColor;
  }
}

.header-icon {
  opacity: 0.9;
}

.page-title {
  font-size: 42px;
  font-weight: 700;
  color: #ffffff;
  margin-bottom: 8px;
  letter-spacing: 0.5px;
}

.page-subtitle {
  font-size: 18px;
  font-weight: 400;
  letter-spacing: 0.3px;
}

.total-count {
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.2px;
}

@media (max-width: 768px) {
  .page-title {
    font-size: 32px;
  }

  .page-subtitle {
    font-size: 16px;
  }

  .header-icon-wrapper {
    width: 80px;
    height: 80px;
  }

  .header-icon {
    font-size: 40px !important;
  }
}

@media (max-width: 480px) {
  .page-title {
    font-size: 28px;
  }
}
</style>
