<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useAppStore } from "../../../store/app";
import { useLazyLoad } from '@/composables/useLazyLoad'
import axios from "axios";
import BookCard from "../../../components/frontend/book/BookCard.vue"

const store = useAppStore()

interface BookType {
  id: number
  name: string
  author: string
  author_id?: number
  slug: string
  image: string
  page: number
  is_deleted: boolean
  created_at: string
  community: number
  comments?: any[]
  author_data?: any
}

// Search state
const search = ref("")
const queryParams = ref({ search: "" })

const { items: books, loading, hasMore, total, error, refresh } = useLazyLoad<BookType>(
  '/get-books-paginated',
  { queryParams }
)

// Handle search
const handleSearch = async () => {
  queryParams.value = { search: search.value }
  await refresh()
}

// Clear search
const clearSearch = async () => {
  search.value = ""
  queryParams.value = { search: "" }
  await refresh()
}

const addBookToReads = async (a) => {
  store.addReadBookId(a.id)
  try {
    await axios.post(`/add-book-to-reads/${store.getUserId}`, a)
  } catch (e) {
    console.error('Failed to add book to reads:', e)
    store.removeReadBookId(a.id) // Rollback on error
  }
}

const removeBookFromReads = async (a) => {
  store.removeReadBookId(a.id)
  try {
    await axios.post(`/delete-book-from-reads/${store.getUserId}`, a)
  } catch (e) {
    console.error('Failed to remove book from reads:', e)
    store.addReadBookId(a.id) // Rollback on error
  }
}

onMounted(async () => {
  try {
    const res = await axios.get(`/get-reads-books-ids/${store.getUserId}`)
    store.setReadBookIds(res.data)
  } catch (error) {
    console.error('Failed to fetch read book IDs:', error)
  }
})

</script>

<template>
  <div class="books-page">
    <!-- Page Header -->
    <div class="page-header text-center mb-8">
      <div class="header-icon-wrapper mb-4">
        <v-icon size="48" class="header-icon">mdi-book-open-page-variant</v-icon>
      </div>
      <h1 class="page-title">Kitaplar</h1>
      <p class="page-subtitle text-grey-lighten-1">Edebiyat dünyasının kapıları</p>
      <p v-if="total > 0" class="total-count text-grey-darken-1 mt-2">
        Toplam {{ total }} kitap
      </p>
    </div>

    <!-- Search Section -->
    <div class="search-section mb-8">
      <v-text-field
        class="search-field"
        v-model="search"
        label="Kitap Ara"
        placeholder="Kitap adı veya yazar..."
        outlined
        clearable
        rounded
        bg-color="grey-darken-2"
        variant="solo-inverted"
        @keyup.enter="handleSearch"
        @click:clear="clearSearch"
      >
        <template v-slot:prepend-inner>
          <v-icon icon="mdi-magnify" size="20" />
        </template>
      </v-text-field>
    </div>

    <!-- Initial Loading -->
    <div v-if="loading && books.length === 0" class="loading-container">
      <v-card class="loading-card">
        <v-card-text class="d-flex flex-column align-center justify-center ga-4">
          <v-progress-circular
            indeterminate
            color="grey-lighten-1"
            size="48"
            width="3"
          ></v-progress-circular>
          <span class="text-grey-lighten-1 text-h6">Kitaplar yükleniyor...</span>
        </v-card-text>
      </v-card>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-container">
      <v-card class="error-card">
        <v-card-text class="text-center">
          <v-icon size="64" color="error" class="mb-4">mdi-alert-circle</v-icon>
          <p class="text-error text-h6">{{ error }}</p>
        </v-card-text>
      </v-card>
    </div>

    <!-- Empty State -->
    <div v-else-if="!loading && books.length === 0" class="empty-state">
      <v-card class="empty-card">
        <v-icon size="80" color="red-lighten-1" class="mb-4">mdi-book-open-page-variant-outline</v-icon>
        <v-card-title class="text-h5 text-white font-weight-bold mb-2">
          Henüz Kitap Bulunmuyor
        </v-card-title>
        <v-card-text class="text-body-1 text-grey-lighten-1">
          Sisteme henüz eklenmiş bir kitap bulunmamaktadır.
        </v-card-text>
      </v-card>
    </div>

    <!-- Books Grid -->
    <div v-else>
      <div class="books-grid d-flex flex-column align-center flex-sm-row flex-sm-wrap justify-sm-center ga-4 mb-12 w-100">
        <div
          v-for="book in books"
          :key="book.id"
          class="book-grid-item"
          style="min-width: 300px;"
        >
          <BookCard
            :book="book"
            @addBookToReads="addBookToReads"
            @removeBookFromReads="removeBookFromReads"
          />
        </div>
      </div>

      <!-- Lazy Loading Indicator -->
      <div v-if="loading && hasMore" class="loading-more-container">
        <v-card class="loading-more-card">
          <v-card-text class="d-flex align-center justify-center ga-3">
            <v-progress-circular
              indeterminate
              color="grey-lighten-1"
              size="24"
              width="2"
            ></v-progress-circular>
            <span class="text-grey-lighten-1">Daha fazla kitap yükleniyor...</span>
          </v-card-text>
        </v-card>
      </div>

      <!-- End of List Message -->
      <div v-if="!hasMore && books.length > 0" class="end-message text-center mt-8 mb-4">
        <v-icon size="20" color="grey-darken-1" class="mr-2">mdi-check-circle</v-icon>
        <span class="text-grey-darken-1">Tüm kitaplar gösteriliyor</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.books-page {
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px 20px;
  animation: fadeIn 0.6s ease-out;
  display: flex;
  flex-direction: column;
  justify-content: start;
  align-items: center;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

/* Page Header */
.page-header {
  animation: fadeInDown 0.6s ease-out;
  width: 100%;
}

/* Search Section */
.search-section {
  width: 100%;
  display: flex;
  justify-content: center;
}

.search-field {
  max-width: 500px;
  width: 100%;
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
    border-color: #616161;
  }
}

.header-icon {
  color: #ef5350;
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

.book-grid-item {
  animation: fadeInUp 0.6s ease-out both;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Loading States */
.loading-container,
.error-container,
.empty-state {
  display: flex;
  justify-content: center;
  padding: 64px 0;
}

.loading-card,
.error-card,
.empty-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  padding: 48px;
  border-radius: 12px;
  text-align: center;
  max-width: 500px;
  transition: all 0.3s ease;
}

.empty-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
  border-color: #616161;
}

.loading-more-container {
  display: flex;
  justify-content: center;
  margin-top: 32px;
  margin-bottom: 32px;
}

.loading-more-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  padding: 16px 32px;
  border-radius: 12px;
}

.end-message {
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.3px;
}

/* Responsive Design */
@media (max-width: 768px) {
  .books-page {
    padding: 24px 12px;
  }

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
