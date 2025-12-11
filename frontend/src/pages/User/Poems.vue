<script setup lang="ts">
import axios from "axios";
import {onMounted, ref, watch, onUnmounted} from "vue";
import PoemCard from "../../components/frontend/poem/PoemCard.vue";
import GoToButtons from "../../components/frontend/poems/GoToButtons.vue";
import {useAppStore} from "../../store/app";
import EmptyStateCard from "../../components/shared/EmptyStateCard.vue";
import LoadingStateCard from "../../components/shared/LoadingStateCard.vue";

const store     = useAppStore()
const poems     = ref([])
const page      = ref(1)
const total     = ref(1)
const last      = ref(1)
const search    = ref("")
const loading   = ref(false)
const isSearch  = ref(false)

const x         = ref(0)
const minValue  = 0;
const maxValue  = 100;
const increment = 1;
let direction   = 1;

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
  isSearch.value = false
  poems.value = []
  loading.value = true
  try {
    const res = await axios.get(`/get-poems?page=${page.value}`)
    poems.value = res.data.poems || []
    total.value = res.data.meta?.total || 0
    last.value = res.data.meta?.last_page || 1
  } catch (error) {
    console.error('Failed to fetch poems:', error)
    poems.value = []
    total.value = 0
    last.value = 1
  } finally {
    loading.value = false
  }
}

const searchPoem = async () => {
  if (!search.value || search.value.trim() === "") {
    page.value = 1
    await getAllPoems()
    return
  }

  isSearch.value = true
  poems.value = []
  loading.value = true
  try {
    const res = await axios.get(`/get-search-poems?page=${page.value}&search=${search.value}`)
    poems.value = res.data.poems || []
    total.value = res.data.meta?.total || 0
    last.value = res.data.meta?.last_page || 1
  } catch (error) {
    console.error('Failed to search poems:', error)
    poems.value = []
    total.value = 0
    last.value = 1
  } finally {
    page.value = 1
    loading.value = false
  }
}
const clearSearch = async () => {
  search.value = ""
  page.value = 1
  isSearch.value = false
  await getAllPoems()
}
// Watch page changes (for pagination)
watch(page, async (newPage, oldPage) => {
  // Skip if this is the initial page value
  if (oldPage === undefined) {
    return
  }

  if (isSearch.value) {
    await searchPoem()
  } else {
    await getAllPoems()
  }
})

const openMessagesModal = ref(false)
const mihrimahCards = ref([])
const loadingCards = ref(false)

const fetchMihrimahCards = async () => {
  try {
    loadingCards.value = true
    const response = await axios.get('/get-mihrimah-cards')
    mihrimahCards.value = response.data || []
  } catch (error) {
    console.error('Mihrimah kartları yüklenemedi:', error)
    mihrimahCards.value = []
  } finally {
    loadingCards.value = false
  }
}

const openMihrimahDialog = async () => {
  openMessagesModal.value = true
  await fetchMihrimahCards()
}

onMounted(async () => {
  try {
    const [likedRes, bookmarkedRes] = await Promise.all([
      axios.get(`/get-liked-poems-id/${store.getUserId}`),
      axios.get(`/get-bookmark-id/${store.getUserId}`)
    ])
    store.setLikedPoemIds(likedRes.data)
    store.setBookmarkedPoemIds(bookmarkedRes.data)
  } catch (error) {
    console.error('Failed to fetch poem IDs:', error)
  }

  await getAllPoems()
})

onUnmounted(() => {
  clearInterval(interval);
});

</script>

<template>
  <div class="poems-container page">
    <!-- MİHRİMAH Easter Egg -->
    <div
      v-if="store.getUser.id == 2"
      @click="openMihrimahDialog"
      class="comfortaa text-h5 py-4 text-white rounded-xl elevation-5 omn my-3 text-center font-bold mihrimah-card"
      :style="{ backgroundColor: `hsl(${x}, 20%, 60%)` }"
    >
      MİHRİMAH
    </div>

    <!-- Search Section -->
    <div class="search-section">
      <v-text-field
        class="search-field"
        v-model="search"
        label="Şiir Ara"
        placeholder="Şiir başlığı veya içeriği..."
        outlined
        clearable
        rounded
        bg-color="grey-darken-2"
        variant="solo-inverted"
        @keyup.enter="searchPoem"
        @click:clear="clearSearch"
      >
        <template v-slot:prepend-inner>
          <v-icon icon="mdi-magnify" size="20" />
        </template>
      </v-text-field>
    </div>

    <!-- Poems Grid -->
    <div class="poems-grid">
      <div v-if="!loading && poems.length > 0">
        <div class="poems-grid-container">
          <div v-for="poem in poems" :key="poem.id" class="poem-grid-item">
            <PoemCard :poem="poem" />
          </div>
        </div>
      </div>
              <!-- Pagination -->
      <v-pagination
        v-if="last > 1"
        class="pagination-control"
        v-model="page"
        :length="last"
        next-icon="mdi-arrow-right"
        prev-icon="mdi-arrow-left"
        rounded="circle"
        active-color="grey-lighten-1"
      ></v-pagination>

      <!-- Loading State -->
      <LoadingStateCard v-if="loading" text="Şiirler yükleniyor" />

      <!-- Empty State -->
      <EmptyStateCard
        v-if="!loading && poems.length === 0"
        title="Şiir Bulunamadı"
        :description="isSearch ? 'Aradığınız kriterlere uygun şiir bulunamadı.' : 'Henüz hiç şiir eklenmemiş.'"
        icon="mdi-file-document-outline"
        :button-text="isSearch ? 'Aramayı Temizle' : ''"
        :show-button="isSearch"
        @click="clearSearch"
      />
    </div>

    <GoToButtons/>

    <!-- Mihrimah Cards Dialog -->
    <v-dialog
      v-model="openMessagesModal"
      max-width="600"
      scrollable
    >
      <v-card
        class="mihrimah-dialog"
      >
        <v-card-title class="d-flex align-center justify-between bg-purple-darken-3 text-white">
          <span class="comfortaa text-h5 d-flex justify-center items-center">MİHRİMAH</span>
          <v-icon size="large" icon="mdi-flower-tulip" />
        </v-card-title>

        <v-card-text class="pa-6">
          <!-- Loading State -->
          <div v-if="loadingCards" class="text-center py-8">
            <v-progress-circular indeterminate color="purple"></v-progress-circular>
          </div>

          <!-- Cards -->
          <div v-else-if="mihrimahCards.length > 0" class="mihrimah-cards-list">
            <div
              v-for="card in mihrimahCards"
              :key="card.id"
              class="mihrimah-card-item mb-4 pa-4 rounded-lg elevation-2"
            >
              <div v-html="card.title" class="card-title text-h6 font-weight-bold mb-3"></div>
              <div v-html="card.content" class="card-content"></div>
            </div>
          </div>

          <!-- Empty State -->
          <div v-else class="text-center py-8">
            <v-icon size="64" color="grey" class="mb-3">mdi-card-text-outline</v-icon>
            <p class="text-grey">Henüz kart eklenmemiş</p>
          </div>
        </v-card-text>

        <v-card-actions class="justify-end pa-4">
          <v-btn
            color="purple-darken-3"
            variant="flat"
            @click="openMessagesModal = false"
          >
            Kapat
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Comfortaa:wght@300..700&display=swap');

/* Container */
.poems-container {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px 20px;
  display: flex;
  flex-direction: column;
}

/* MİHRİMAH Card */
.mihrimah-card {
  transition: all 0.3s ease;
}

.comfortaa {
  font-family: "Comfortaa", sans-serif;
  font-optical-sizing: auto;
  font-weight: 900;
}

.omn:hover {
  cursor: pointer;
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
}

.search-section {
  margin-bottom: 32px;
  display: flex;
  justify-content: center;
}

.search-field {
  max-width: 500px;
  width: 100%;
}

/* Poems Grid */
.poems-grid {
  width: 100%;
}

.poems-grid-container {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 24px;
  margin-bottom: 32px;
}

.poem-grid-item {
  width: 100%;
  max-width: 380px;
}

/* Pagination */
.pagination-control {
  margin-top: 48px;
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
  .poems-container {
    padding: 24px 12px;
  }

  .search-section {
    margin-bottom: 24px;
  }

  .poems-grid-container {
    gap: 16px;
  }
}

/* Mihrimah Dialog Styles */
.mihrimah-dialog {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
}

.mihrimah-cards-list {
  max-height: 500px;
  overflow-y: auto;
}

.mihrimah-card-item {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.3s ease;
}

.mihrimah-card-item:hover {
  background: rgba(255, 255, 255, 0.08);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.mihrimah-card-item :deep(p) {
  margin-bottom: 0.5rem;
}

.mihrimah-card-item :deep(strong) {
  font-weight: bold;
}

.mihrimah-card-item :deep(em) {
  font-style: italic;
}

.mihrimah-card-item :deep(a) {
  color: #ba68c8;
  text-decoration: underline;
}

.mihrimah-card-item :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  margin: 0.5rem 0;
}
</style>
