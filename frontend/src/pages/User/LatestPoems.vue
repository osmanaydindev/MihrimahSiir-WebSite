<script setup lang="ts">
import axios from "axios";
import {onMounted, ref, watch, onUnmounted, computed} from "vue";
// import {useAppStore} from "../store/app";
// import {useRouter} from "vue-router";
import PoemCard from "../../components/frontend/poem/PoemCard.vue";
import GoToButtons from "../../components/frontend/poems/GoToButtons.vue";

const loading = ref(false)
// const router = useRouter()
// const store = useAppStore()
const poems = ref([])

const x = ref(0)
const minValue = 0;
const maxValue = 100;
const increment = 1;
let direction = 1;
const updateValue = () => {
  if (x.value >= maxValue) {
    direction = -1;
  } else if (x.value <= minValue) {
    direction = 1;
  }
  x.value += direction * increment;
};

const interval = setInterval(updateValue, 10);
const getAllPoems = async () => {
  poems.value = []
  loading.value = true
  try {
    const res = await axios.get(`/get-latest-poems`)
    poems.value = res.data || []
  } catch (error) {
    console.error('Failed to fetch latest poems:', error)
    poems.value = []
  } finally {
    loading.value = false
  }
}
onMounted(async () => {
  await getAllPoems()
})
onUnmounted(() => {
  clearInterval(interval);
});

</script>

<template>
  <div class="latest-poems-container">
    <!-- Page Title -->
    <div v-if="!loading && poems.length > 0" class="page-title-section text-center mb-8">
      <h1 class="page-title d-flex align-center justify-center">
        En Son Şiirler
      </h1>
    </div>

    <!-- Poems Grid -->
    <div v-if="!loading && poems.length > 0" class="poems-grid">
      <div class="poems-grid-container">
        <div v-for="poem in poems" :key="poem.id" class="poem-grid-item">
          <PoemCard :poem="poem" />
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <v-card class="loading-card">
        <v-card-text class="d-flex align-center justify-center ga-3">
          <v-progress-circular
            indeterminate
            color="grey-lighten-1"
            size="24"
            width="2"
          ></v-progress-circular>
          <span class="text-grey-lighten-1">En son şiirler yükleniyor...</span>
        </v-card-text>
      </v-card>
    </div>

    <!-- Empty State -->
    <div v-if="!loading && poems.length === 0" class="empty-state">
      <v-card class="empty-state-card">
        <v-icon
          size="80"
          color="grey-lighten-1"
          class="mb-4"
        >
          mdi-file-document-outline
        </v-icon>
        <v-card-title class="text-h5 text-white font-weight-bold mb-2">
          Henüz Şiir Yok
        </v-card-title>
        <v-card-text class="text-body-1 text-grey-lighten-1">
          Henüz hiç şiir eklenmemiş.
        </v-card-text>
      </v-card>
    </div>

    <GoToButtons/>
  </div>
</template>

<style scoped>
/* Container */
.latest-poems-container {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px 20px;
  display: flex;
  flex-direction: column;
}

/* Page Title */
.page-title-section {
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

.page-title {
  font-size: 32px;
  font-weight: 700;
  color: #ffffff;
  letter-spacing: 0.5px;
  margin: 0;
}

/* Poems Grid */
.poems-grid {
  width: 100%;
  margin-bottom: 48px;
}

.poems-grid-container {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 24px;
}

.poem-grid-item {
  width: 100%;
  max-width: 380px;
}

/* Loading State */
.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 48px 0;
}

.loading-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  padding: 24px;
  border-radius: 12px;
}

/* Empty State */
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

/* Responsive Grid */
@media (min-width: 600px) {
  .poem-grid-item {
    width: calc(50% - 12px);
  }
}

@media (min-width: 960px) {
  .poem-grid-item {
    width: calc(33.333% - 16px);
  }
}

@media (min-width: 1280px) {
  .poem-grid-item {
    width: calc(33.333% - 16px);
  }
}

@media (min-width: 1920px) {
  .poem-grid-item {
    width: calc(25% - 18px);
  }
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .latest-poems-container {
    padding: 24px 12px;
  }

  .poems-grid-container {
    gap: 16px;
  }
}
</style>
