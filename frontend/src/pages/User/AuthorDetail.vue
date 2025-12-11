<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import PoemCard from '../../components/frontend/poem/PoemCard.vue'
import BookCard from '../../components/frontend/book/BookCard.vue'
import LoadingState from '../../components/common/LoadingState.vue'
import ErrorState from '../../components/common/ErrorState.vue'
import EmptyState from '../../components/common/EmptyState.vue'

interface Author {
  id: number
  name: string
  bio?: string
  birth_year?: number
  death_year?: number
  nationality?: string
  image?: string
  slug: string
  poems?: any[]
  books?: any[]
}

const route = useRoute()
const loading = ref(true)
const author = ref<Author | null>(null)
const error = ref<string | null>(null)
const activeTab = ref('poems')

const yearRange = computed(() => {
  if (!author.value?.birth_year && !author.value?.death_year) return null

  const birth = author.value.birth_year || '?'
  const death = author.value.death_year || ''

  return death ? `${birth} - ${death}` : `${birth} -`
})

const authorImage = computed(() => {
  return author.value?.image || '/uploads/default-author.jpg'
})
const goTabs = (activeTabStr: string) => {
  activeTab.value = activeTabStr
  document.getElementById('tabSection')?.scrollIntoView({ behavior: 'smooth' })
}
const fetchAuthor = async () => {
  loading.value = true
  error.value = null

  try {
    const response = await axios.get(`/get-author/${route.params.slug}`)
    author.value = response.data
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Yazar bulunamadı'
    console.error('Failed to fetch author:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchAuthor()
})
</script>

<template>
  <div class="author-detail-page">
    <!-- Loading State -->
    <LoadingState
      v-if="loading"
      message="Yazar bilgileri yükleniyor..."
      color="grey-lighten-1"
    />

    <!-- Error State -->
    <ErrorState
      v-else-if="error"
      :error="error"
      button-text="Yazarlara Dön"
      button-icon="mdi-arrow-left"
      @button-click="$router.push('/authors')"
    />

    <!-- Author Content -->
    <div v-else-if="author" class="author-content">
      <!-- Author Header -->
      <div class="author-header">
        <div class="header-background"></div>
        <div class="header-content-wrapper">
          <div class="header-content">
            <!-- Author Portrait -->
            <div class="portrait-container">
              <v-img
                :src="authorImage"
                class="author-portrait"
                cover
              >
                <div class="portrait-border"></div>
              </v-img>
            </div>

            <!-- Author Info -->
            <div class="author-main-info">
              <h1 class="author-name">{{ author.name }}</h1>

              <div class="author-meta">
                <div v-if="yearRange" class="meta-item" >
                  <v-icon size="18" color="grey-lighten-1">mdi-calendar-range</v-icon>
                  <span>{{ yearRange }}</span>
                </div>
                <div v-if="author.nationality" class="meta-item">
                  <v-icon size="18" color="grey-lighten-1">mdi-earth</v-icon>
                  <span>{{ author.nationality }}</span>
                </div>
              </div>

              <!-- Stats -->
              <div class="author-stats">
                <div class="stat-box text-decoration-none" @click.prevent="goTabs('poems')">
                  <v-icon size="24" color="blue-lighten-1">mdi-book-open-page-variant</v-icon>
                  <div class="stat-info">
                    <span class="stat-number">{{ author.poems?.length || 0 }}</span>
                    <span class="stat-label">Şiir</span>
                  </div>
                </div>
                <div class="stat-box text-decoration-none" @click.prevent="goTabs('books')">
                  <v-icon size="24" color="green-lighten-1">mdi-book</v-icon>
                  <div class="stat-info">
                    <span class="stat-number">{{ author.books?.length || 0 }}</span>
                    <span class="stat-label">Kitap</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Biography -->
      <div v-if="author.bio" class="biography-section">
        <v-card class="bio-card">
          <v-card-title class="bio-title">
            <v-icon left color="grey-lighten-1">mdi-text-box-outline</v-icon>
            Biyografi
          </v-card-title>
          <v-card-text class="bio-content">
            <p>{{ author.bio }}</p>
          </v-card-text>
        </v-card>
      </div>

      <!-- Works Section -->
      <div class="works-section" id="tabSection">
        <!-- Tabs -->
        <v-tabs
          v-model="activeTab"
          class="works-tabs"
          bg-color="transparent"
          color="grey-lighten-4"
        >
          <v-tab value="poems" class="tab-button">
            <v-icon left size="20">mdi-book-open-page-variant</v-icon>
            Şiirleri ({{ author.poems?.length || 0 }})
          </v-tab>
          <v-tab value="books" class="tab-button">
            <v-icon left size="20">mdi-book</v-icon>
            Kitapları ({{ author.books?.length || 0 }})
          </v-tab>
        </v-tabs>

        <!-- Tab Content -->
        <div class="tab-content">
          <!-- Poems Tab -->
          <div v-show="activeTab === 'poems'">
            <div v-if="author.poems && author.poems.length > 0" class="poems-grid">
              <div v-for="poem in author.poems" :key="poem.id" class="poem-grid-item">
                <PoemCard :poem="poem" />
              </div>
            </div>
            <EmptyState
              v-else
              icon="mdi-book-open-page-variant-outline"
              icon-color="grey-darken-1"
              title="Henüz şiir eklenmemiş"
              message="Bu yazara ait şiirler burada görünecek"
            />
          </div>

          <!-- Books Tab -->
          <div v-show="activeTab === 'books'">
            <div v-if="author.books && author.books.length > 0" class="books-grid">
              <div v-for="book in author.books" :key="book.id" class="book-grid-item">
                <BookCard :book="book" />
              </div>
            </div>
            <EmptyState
              v-else
              icon="mdi-book-outline"
              icon-color="grey-darken-1"
              title="Henüz kitap eklenmemiş"
              message="Bu yazara ait kitaplar burada görünecek"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.author-detail-page {
  min-height: 100vh;
  animation: fadeIn 0.6s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

/* Author Header */
.author-header {
  position: relative;
  margin-bottom: clamp(24px, 5vw, 48px);
  padding-bottom: clamp(16px, 3vw, 32px);
}

.header-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 100%;
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 50%, #1a1a1a 100%);
  border-bottom: 2px solid #424242;
}

.header-content-wrapper {
  position: relative;
  max-width: 1400px;
  margin: 0 auto;
  padding: clamp(20px, 5vw, 40px);
}

.header-content {
  display: flex;
  gap: clamp(20px, 5vw, 40px);
  align-items: flex-start;
  flex-wrap: wrap;
  justify-content: center;
}

.portrait-container {
  position: relative;
  flex-shrink: 0;
}

.author-portrait {
  width: clamp(140px, 30vw, 240px);
  height: clamp(140px, 30vw, 240px);
  border-radius: clamp(12px, 2vw, 16px);
  border: clamp(2px, 0.5vw, 4px) solid #424242;
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.6);
  transition: transform 0.3s ease;
}

.author-portrait:hover {
  transform: scale(1.05);
}

.portrait-border {
  position: absolute;
  inset: 0;
  border-radius: clamp(10px, 2vw, 12px);
  background: linear-gradient(135deg, rgba(39, 151, 176, 0.1), transparent);
}

.author-main-info {
  flex: 1;
  min-width: min(100%, 280px);
  max-width: 100%;
}

.author-name {
  font-size: clamp(24px, 6vw, 48px);
  font-weight: 700;
  color: #ffffff;
  margin-bottom: clamp(12px, 2vw, 16px);
  letter-spacing: 0.5px;
  line-height: 1.2;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  word-wrap: break-word;
}

.author-meta {
  display: flex;
  gap: clamp(16px, 3vw, 24px);
  margin-bottom: clamp(20px, 4vw, 32px);
  flex-wrap: wrap;
  justify-content: center;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: clamp(14px, 2vw, 16px);
  color: #bdbdbd;
  font-weight: 500;
  letter-spacing: 0.3px;
}

.author-stats {
  display: flex;
  gap: clamp(12px, 3vw, 24px);
  flex-wrap: wrap;
  justify-content: center;
}

.stat-box {
  display: flex;
  align-items: center;
  gap: clamp(8px, 2vw, 12px);
  padding: clamp(12px, 2vw, 16px) clamp(16px, 3vw, 24px);
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  border-radius: clamp(8px, 1.5vw, 12px);
  transition: all 0.3s ease;
  min-width: min-content;
}

.stat-box:hover {
  transform: translateY(-4px);
  border-color: #616161;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.3);
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-number {
  font-size: clamp(20px, 4vw, 24px);
  font-weight: 700;
  color: #ffffff;
  line-height: 1;
}

.stat-label {
  font-size: clamp(12px, 2vw, 13px);
  color: #9e9e9e;
  font-weight: 500;
  letter-spacing: 0.3px;
}

/* Biography */
.biography-section {
  max-width: 1400px;
  margin: 0 auto clamp(32px, 6vw, 48px);
  padding: 0 clamp(16px, 4vw, 20px);
}

.bio-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
}

.bio-title {
  font-size: clamp(18px, 3vw, 20px);
  font-weight: 600;
  color: #e0e0e0;
  padding: clamp(20px, 4vw, 24px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.bio-content {
  font-size: clamp(14px, 2.5vw, 16px);
  line-height: 1.8;
  color: #bdbdbd;
  padding: clamp(20px, 4vw, 24px);
  letter-spacing: 0.3px;
  word-wrap: break-word;
}

/* Works Section */
.works-section {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 clamp(16px, 4vw, 20px) clamp(32px, 6vw, 40px);
}

.works-tabs {
  margin-bottom: clamp(24px, 5vw, 32px);
  border-bottom: 1px solid #424242;
}

.tab-button {
  text-transform: none;
  font-size: clamp(14px, 2.5vw, 16px);
  font-weight: 600;
  letter-spacing: 0.3px;
}
/* Pasif tab - koyu gri */
.works-tabs .v-tab {
  color: #9d9999 !important;
  opacity: 0.7 !important;
}

/* Aktif tab - parlak beyaz */
.works-tabs .v-tab.v-tab--selected {
  color: #ffffff !important;
  opacity: 1 !important;
}

.tab-content {
  min-height: clamp(300px, 60vw, 400px);
}

.poems-grid,
.books-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(clamp(240px, 45vw, 280px), 1fr));
  gap: clamp(16px, 3vw, 24px);
}

/* Mobile-specific adjustments */
@media (max-width: 768px) {
  .header-content {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
}

/* Very small screens - force single column for grids */
@media (max-width: 480px) {
  .poems-grid,
  .books-grid {
    grid-template-columns: 1fr;
  }
}
</style>
