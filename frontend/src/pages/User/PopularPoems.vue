<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useAppStore } from '../../store/app'
import PoemCard from '../../components/frontend/poem/PoemCard.vue'
import PageHeader from '../../components/common/PageHeader.vue'
import LoadingState from '../../components/common/LoadingState.vue'
import EmptyState from '../../components/common/EmptyState.vue'

const store = useAppStore()
const poems = ref([])
const loading = ref(false)
const page = ref(1)
const total = ref(0)
const lastPage = ref(1)

const fetchPopularPoems = async () => {
  loading.value = true
  try {
    const res = await axios.get(`/get-popular-poems?page=${page.value}`)
    poems.value = res.data.poems || []
    total.value = res.data.meta?.total || 0
    lastPage.value = res.data.meta?.last_page || 1
  } catch (error) {
    console.error('Failed to fetch popular poems:', error)
    poems.value = []
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  // Load liked and bookmarked poem IDs
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

  await fetchPopularPoems()
})
</script>

<template>
  <div class="popular-poems-container d-flex flex-column">
    <PageHeader
      icon="mdi-fire"
      icon-color="#ff6b6b"
      title="Popüler Şiirler"
      subtitle="En çok beğenilen şiirler"
      :total="total"
      total-label="şiir"
    />

    <LoadingState v-if="loading && poems.length === 0" message="Popüler şiirler yükleniyor..." color="red-lighten-1" />

    <EmptyState
      v-else-if="!loading && poems.length === 0"
      icon="mdi-fire-off"
      icon-color="grey-lighten-1"
      title="Henüz Popüler Şiir Yok"
      message="Şiirleri beğenerek popüler şiirlerin oluşmasına katkıda bulunabilirsin!"
    />

    <div v-else class="poems-grid-container">
      <div v-for="poem in poems" :key="poem.id" class="poem-grid-item">
        <PoemCard :poem="poem" />
      </div>
    </div>

    <v-pagination
      v-if="lastPage > 1"
      v-model="page"
      :length="lastPage"
      @update:modelValue="fetchPopularPoems"
      class="mt-8"
      rounded="circle"
      active-color="red-lighten-1"
    />
  </div>
</template>

<style scoped>
.popular-poems-container {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px 20px;
}

.poems-grid-container {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 24px;
  margin-top: 32px;
}

.poem-grid-item {
  width: 100%;
  max-width: 380px;
}

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

@media (min-width: 1920px) {
  .poem-grid-item {
    width: calc(25% - 18px);
  }
}
</style>
