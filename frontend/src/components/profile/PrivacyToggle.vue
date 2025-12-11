<script setup lang="ts">
interface Props {
  isPrivate: boolean
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  'update:isPrivate': [value: boolean]
}>()
</script>

<template>
  <div class="privacy-toggle">
    <v-switch
      :model-value="isPrivate"
      @update:model-value="emit('update:isPrivate', $event)"
      :loading="loading"
      color="primary"
      hide-details
      density="comfortable"
    >
      <template v-slot:label>
        <div class="d-flex align-center ga-2">
          <v-icon>{{ isPrivate ? 'mdi-lock' : 'mdi-lock-open' }}</v-icon>
          <span class="privacy-label">
            {{ isPrivate ? 'Profil Gizli' : 'Profil Herkese Açık' }}
          </span>
        </div>
      </template>
    </v-switch>
    <p class="privacy-description text-caption text-grey mt-1">
      {{ isPrivate ? 'Sadece arkadaşların profilini görebilir' : 'Herkes profilini görebilir' }}
    </p>
  </div>
</template>

<style scoped>
.privacy-toggle {
  width: 100%;
  max-width: 300px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid #424242;
  border-radius: 12px;
  transition: all 0.3s ease;
}

.privacy-toggle:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: #616161;
}

.privacy-label {
  font-size: 14px;
  font-weight: 500;
  color: #ffffff;
}

.privacy-description {
  margin-left: 44px;
  font-size: 12px;
  line-height: 1.4;
  text-align: left;
}

@media (max-width: 768px) {
  .privacy-toggle {
    max-width: 100%;
  }
}
</style>
