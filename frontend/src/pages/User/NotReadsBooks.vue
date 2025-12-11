<script setup lang="ts">
import { useAppStore } from "../../store/app";
import { useRouter } from "vue-router";
import { useLazyLoad } from '@/composables/useLazyLoad'
import BookCard from "../../components/frontend/book/BookCard.vue";
import PageHeader from "../../components/common/PageHeader.vue";
import SearchSection from "../../components/common/SearchSection.vue";
import LoadingState from "../../components/common/LoadingState.vue";
import ErrorState from "../../components/common/ErrorState.vue";
import EmptyState from "../../components/common/EmptyState.vue";
import LazyLoadingIndicator from "../../components/common/LazyLoadingIndicator.vue";
import EndOfListMessage from "../../components/common/EndOfListMessage.vue";
import axios from "axios";
import { ref } from "vue";

const router = useRouter()
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

const { items: notReadsBooks, loading, hasMore, total, error, refresh } = useLazyLoad<BookType>(
  `/get-not-reads-books-paginated/${store.getUserId}`,
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
    handleBookAddedToReads(a)
  } catch (e) {
    console.error('Failed to add book to reads:', e)
    store.removeReadBookId(a.id)
  }
}

const removeBookFromReads = async (a) => {
  store.removeReadBookId(a.id)
  try {
    await axios.post(`/delete-book-from-reads/${store.getUserId}`, a)
  } catch (e) {
    console.error('Failed to remove book from reads:', e)
    store.addReadBookId(a.id)
  }
}

const handleBookAddedToReads = (book: any) => {
  const index = notReadsBooks.value.findIndex(b => b.id === book.id)
  if (index !== -1) {
    notReadsBooks.value.splice(index, 1)
    if (total.value > 0) {
      total.value--
    }
  }
}
</script>

<template>
  <div class="page-container">
    <!-- Initial Loading -->
    <LoadingState
      v-if="loading && notReadsBooks.length === 0"
      message="Okunacak kitaplar yükleniyor..."
      color="green-lighten-1"
    />

    <!-- Error State -->
    <ErrorState v-else-if="error" :error="error" />

    <!-- Empty State -->
    <EmptyState
      v-else-if="!loading && notReadsBooks.length === 0"
      icon="mdi-check-circle-outline"
      icon-color="green-lighten-1"
      title="Tüm Kitapları Okudunuz!"
      message="Harika! Tüm kitapları bitirdiniz.<br>Yeni kitaplar eklendiğinde buradan görebilirsiniz."
      button-text="Tüm Kitaplara Dön"
      button-icon="mdi-book-multiple"
      button-color="green-lighten-1"
      @button-click="router.push('/books')"
    />

    <!-- Books Content -->
    <div v-else class="content-wrapper w-100">
      <!-- Page Header -->
      <PageHeader
        icon="mdi-book-open-variant"
        icon-color="#81c784"
        title="Okumadığınız Kitaplar"
        subtitle="Okuma listenizde bekleyenler"
        :total="total"
        total-label="kitap"
      />

      <!-- Search Section -->
      <SearchSection
        v-model="search"
        label="Kitap Ara"
        placeholder="Kitap adı veya yazar..."
        @search="handleSearch"
        @clear="clearSearch"
      />

      <!-- Books Grid -->
      <div class="books-grid d-flex flex-column align-center flex-sm-row flex-sm-wrap justify-sm-center ga-4 mb-12 w-100">
        <div
          v-for="book in notReadsBooks"
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
      <LazyLoadingIndicator
        v-if="loading && hasMore"
        message="Daha fazla kitap yükleniyor..."
        color="green-lighten-1"
      />

      <!-- End of List Message -->
      <EndOfListMessage
        v-if="!hasMore && notReadsBooks.length > 0"
        message="Tüm okunacak kitaplar gösteriliyor"
        icon-color="green-darken-1"
      />

      <!-- View All Books Button -->
      <div class="text-center mt-8">
        <v-btn
          color="green-lighten-1"
          variant="outlined"
          size="large"
          rounded="lg"
          @click="router.push('/books')"
          class="text-none"
        >
          <v-icon size="20" class="mr-2">mdi-book-multiple</v-icon>
          Tüm Kitapları Gör
        </v-btn>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px 20px;
  min-height: calc(100vh - 200px);
  animation: fadeIn 0.6s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.content-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.books-grid {
  max-width: 1400px;
  margin: 0 auto;
}

@media (max-width: 768px) {
  .page-container {
    padding: 24px 12px;
  }
}
</style>
