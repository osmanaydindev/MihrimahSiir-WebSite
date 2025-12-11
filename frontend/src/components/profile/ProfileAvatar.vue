<script setup lang="ts">
import { ref, computed } from 'vue'
import SnackbarNotification from "../common/SnackbarNotification.vue"
import { useNotification } from "../../composables/useNotification"

interface Props {
  imageUrl: string
  size?: number
  editable?: boolean
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 150,
  editable: false,
  loading: false
})

const emit = defineEmits<{
  'upload': [file: File]
}>()

const { snackbar, error: showError, success } = useNotification()
const fileInput = ref<HTMLInputElement | null>(null)

const triggerFileInput = () => {
  if (fileInput.value) {
    fileInput.value.click()
  }
}

const handleFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]

  if (!file) return

  // Validate file size (max 1MB)
  if (file.size > 1024 * 1024) {
    showError('Dosya boyutu 1MB\'dan küçük olmalıdır')
    return
  }

  // Validate file type
  if (!file.type.startsWith('image/')) {
    showError('Lütfen bir resim dosyası seçin')
    return
  }

  emit('upload', file)

  // Clear the file input
  if (target) {
    target.value = ''
  }
}
</script>

<template>
  <div class="profile-image-container">
    <v-avatar :size="size" class="profile-avatar">
      <v-img :src="imageUrl" cover />
    </v-avatar>

    <!-- Upload button for editable -->
    <div v-if="editable" class="image-upload-overlay">
      <input
        ref="fileInput"
        type="file"
        accept="image/*"
        style="display: none"
        @change="handleFileChange"
      />
      <v-btn
        @click="triggerFileInput"
        :loading="loading"
        size="small"
        icon
        color="grey-darken-3"
        class="upload-btn"
      >
        <v-icon>mdi-camera</v-icon>
      </v-btn>
    </div>
  </div>

  <!-- Snackbar for notifications -->
  <SnackbarNotification
    v-model="snackbar.show"
    :message="snackbar.message"
    :color="snackbar.color"
  />
</template>

<style scoped>
.profile-image-container {
  position: relative;
  flex-shrink: 0;
}

.profile-avatar {
  border: 4px solid #424242;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}

.image-upload-overlay {
  position: absolute;
  bottom: 0;
  right: 0;
}

.upload-btn {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}
</style>
