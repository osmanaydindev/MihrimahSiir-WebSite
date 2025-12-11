<script setup lang="ts">
import { watch, ref } from "vue";
import { useRouter } from "vue-router";
import { useAppStore } from "../../store/app";
import { useLazyLoad } from '@/composables/useLazyLoad'
import PoemCard from "../../components/frontend/poem/PoemCard.vue";
import PageHeader from "../../components/common/PageHeader.vue";
import SearchSection from "../../components/common/SearchSection.vue";
import LoadingState from "../../components/common/LoadingState.vue";
import ErrorState from "../../components/common/ErrorState.vue";
import EmptyState from "../../components/common/EmptyState.vue";
import LazyLoadingIndicator from "../../components/common/LazyLoadingIndicator.vue";
import EndOfListMessage from "../../components/common/EndOfListMessage.vue";

const router = useRouter()
const store = useAppStore()

interface PoemType {
  id: number
  title: string
  content: string
  author: string
  slug: string
  is_deleted: boolean
  created_at: string
}

// Search state
const search = ref("")
const queryParams = ref({ search: "" })

const { items: likedPoems, loading, hasMore, total, error, refresh } = useLazyLoad<PoemType>(
  `/get-liked-poems-paginated/${store.getUserId}`,
  { queryParams }
)

// Update store when lazy loaded poems change
watch(likedPoems, (newPoems) => {
  if (newPoems && newPoems.length > 0) {
    store.setLikedPoemIdsFromPoems(newPoems)
  }
}, { immediate: true })

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

// Handle poem unliked
const handlePoemUnliked = (poem: PoemType) => {
  const index = likedPoems.value.findIndex(p => p.id === poem.id)
  if (index !== -1) {
    likedPoems.value.splice(index, 1)
    if (total.value > 0) {
      total.value--
    }
  }
}
</script>

<template>
  <div class="liked-poems-container">
    <!-- Initial Loading -->
    <LoadingState
      v-if="loading && likedPoems.length === 0"
      message="Beğenilen şiirler yükleniyor..."
      color="pink-lighten-1"
    />

    <!-- Error State -->
    <ErrorState v-else-if="error" :error="error" />

    <!-- Empty State -->
    <EmptyState
      v-else-if="!loading && likedPoems.length === 0"
      icon="mdi-heart-outline"
      icon-color="pink-lighten-1"
      title="Henüz Beğendiğiniz Şiir Yok"
      message="Beğendiğiniz şiirleri buradan görebilirsiniz.<br>Şimdi keşfetmeye başlayın!"
      button-text="Şiirleri Keşfet"
      button-icon="mdi-book-multiple"
      button-color="pink-lighten-1"
      @button-click="router.push('/poems')"
    />

    <!-- Poems Content -->
    <div v-else class="content-wrapper w-100">
      <!-- Page Header -->
      <PageHeader
        icon="mdi-heart"
        icon-color="#f48fb1"
        title="Beğendiğiniz Şiirler"
        subtitle="Kalbinize dokunan dizeler"
        :total="total"
        total-label="şiir"
      />

      <!-- Search Section -->
      <SearchSection
        v-model="search"
        label="Şiir Ara"
        placeholder="Şiir başlığı veya içeriği..."
        @search="handleSearch"
        @clear="clearSearch"
      />

      <!-- Poems Grid -->
      <div class="poems-grid d-flex flex-column align-center flex-sm-row flex-sm-wrap justify-sm-center ga-4 mb-12 w-100">
        <div v-for="poem in likedPoems" :key="poem.id" class="poem-grid-item">
          <PoemCard :poem="poem" @poem-unliked="handlePoemUnliked" />
        </div>
      </div>

      <!-- Lazy Loading Indicator -->
      <LazyLoadingIndicator
        v-if="loading && hasMore"
        message="Daha fazla şiir yükleniyor..."
        color="pink-lighten-1"
      />

      <!-- End of List Message -->
      <EndOfListMessage
        v-if="!hasMore && likedPoems.length > 0"
        message="Tüm beğenilen şiirler gösteriliyor"
        icon-color="pink-darken-1"
      />

      <!-- View All Poems Button -->
      <div class="text-center mt-8">
        <v-btn
          color="pink-lighten-1"
          variant="outlined"
          size="large"
          rounded="lg"
          @click="router.push('/poems')"
          class="text-none"
        >
          <v-icon size="20" class="mr-2">mdi-book-multiple</v-icon>
          Tüm Şiirleri Gör
        </v-btn>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Page Container */
.liked-poems-container {
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

/* Content Wrapper */
.content-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
}

/* Grid System */
.poems-grid {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
}

.poem-grid-item {
  width: 100%;
  max-width: 380px;
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

/* Responsive Grid */
@media (min-width: 600px) {
  .poem-grid-item {
    width: calc(50% - 12px);
    min-width: 340px;
    max-width: 380px;
  }
}

@media (min-width: 960px) {
  .poem-grid-item {
    width: calc(33.333% - 16px);
    min-width: 340px;
    max-width: 380px;
  }
}

@media (min-width: 1280px) {
  .poem-grid-item {
    width: calc(33.333% - 16px);
    min-width: 360px;
    max-width: 380px;
  }
}

@media (min-width: 1920px) {
  .poem-grid-item {
    width: calc(25% - 16px);
    min-width: 360px;
    max-width: 380px;
  }
}

/* Responsive Design */
@media (max-width: 768px) {
  .liked-poems-container {
    padding: 24px 12px;
  }
}
</style>
