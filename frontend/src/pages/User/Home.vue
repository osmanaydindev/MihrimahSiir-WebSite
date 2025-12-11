<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useAppStore } from "../../store/app";
import axios from "axios";
const store = useAppStore();

const homepageItems = ref([]);
const loading = ref(true);

onMounted(async () => {
  try {
    loading.value = true;

    const response = await axios.get('/get-homepage-items');
    homepageItems.value = response.data || [];
  } catch (error) {
    console.error('Anasayfa yüklenemedi:', error);
    homepageItems.value = [];
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div class="main d-flex flex-column justify-start" style="max-width: 1080px">

    <!-- Loading state -->
    <div v-if="loading" class="mt-8 mb-8 text-center">
      <v-progress-circular indeterminate color="primary"></v-progress-circular>
    </div>

    <!-- Homepage items from database -->
    <div v-else-if="homepageItems.length > 0" class="mt-2 mt-sm-8 mb-16">
      <div
        v-for="item in homepageItems"
        :key="item.id"
        class="homepage-item mb-12 w-full"
      >
        <div v-if="item.title || item.subtitle" class="mb-6">
            <div v-html="item.title" class="homepage-title text-h4 font-weight-bold mb-0"></div>
            <div v-html="item.subtitle" class="homepage-subtitle text-h5 mt-0"></div>
        </div>
        <div class="text-left" v-if="item.content">
          <div v-html="item.content" class="homepage-content"></div>
        </div>
      </div>
    </div>

    <!-- Fallback: No items message -->
    <div v-else class="mt-8 mb-8 text-center">
      <p class="text-h6 text-grey">Henüz anasayfa içeriği eklenmemiş.</p>
    </div>

  </div>

</template>

<style scoped>
.homepage-item {
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  padding-bottom: 3rem;
}

.homepage-item:last-child {
  border-bottom: none;
}

:deep(.homepage-content p) {
  margin-bottom: 1rem;
}

:deep(.homepage-content h1),
:deep(.homepage-content h2),
:deep(.homepage-content h3),
:deep(.homepage-content h4),
:deep(.homepage-content h5),
:deep(.homepage-content h6) {
  margin-top: 1.5rem;
  margin-bottom: 1rem;
  font-weight: bold;
}

:deep(.homepage-content ul),
:deep(.homepage-content ol) {
  margin-left: 2rem;
  margin-bottom: 1rem;
}

:deep(.homepage-content li) {
  margin-bottom: 0.5rem;
}

:deep(.homepage-content img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 1rem 0;
}

:deep(.homepage-content a) {
  color: #3b82f6;
  text-decoration: underline;
}

:deep(.homepage-content a:hover) {
  color: #2563eb;
}

:deep(.homepage-content strong) {
  font-weight: bold;
}

:deep(.homepage-content em) {
  font-style: italic;
}
</style>
