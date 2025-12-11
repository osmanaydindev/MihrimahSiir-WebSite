<script setup lang="ts">
import {useRoute, useRouter} from "vue-router";
import {onMounted, ref, computed, watch} from "vue";
import {useAppStore} from "../../store/app";
import axios from "axios";
import PoemCard from "../../components/frontend/poem/PoemCard.vue";
import {usePoemActions} from "../../composables/usePoemActions";
import '../../styles/Poem.css';

const loading = ref(false)
const store = useAppStore()
const route = useRoute()
const router = useRouter()
const poem = ref({})
const randomPeoms = ref([])
const { likePoem, unlikePoem, bookmarkPoem, unbookmarkPoem } = usePoemActions()

const getPoem = async () => {
  window.scrollTo(0, 0)
  loading.value = true
  try {
    const res = await axios.get(`/get-poem/${route.params.slug}`)
    // Check if poem data exists
    if (!res.data || !res.data.poem || !res.data.poem.id) {
      router.push({ name: 'NotFound' })
      return
    }
    poem.value = res.data.poem
    randomPeoms.value = res.data.randomPoem || []
  } catch (error) {
    console.error('Failed to fetch poem:', error)
    // If API returns error (404, 500, etc), redirect to NotFound
    router.push({ name: 'NotFound' })
  } finally {
    loading.value = false
  }
}
const addLikedPoem = async (poemData: any) => {
  try {
    await likePoem(poemData)
  } catch (error) {
    // Error already logged in composable
  }
}

const addBookmark = async (poemData: any) => {
  try {
    await bookmarkPoem(poemData)
  } catch (error) {
    // Error already logged in composable
  }
}

const undoLikedPoem = async (poemData: any) => {
  try {
    await unlikePoem(poemData)
  } catch (error) {
    // Error already logged in composable
  }
}

const undoBookmark = async (poemData: any) => {
  try {
    await unbookmarkPoem(poemData)
  } catch (error) {
    // Error already logged in composable
  }
}

const isLiked = computed(() => {
  if (!poem.value || !poem.value.id) return false
  return store.likedPoemIds.has(poem.value.id)
})

const isBookmark = computed(() => {
  if (!poem.value || !poem.value.id) return false
  return store.bookmarkedPoemIds.has(poem.value.id)
})
watch(
   () => route.params.slug,
  async (newSlug, oldSlug) => {
    await getPoem()
  }
);
onMounted(async () => {
  await getPoem()
})

</script>

<template>
  <div class="poem-detail-container">
    <!-- Main Poem Card -->
    <div v-if="!loading" class="poem-main-card">
      <!-- Poem Header -->
      <div class="poem-header">
        <h1 class="poem-title-main">{{ poem.title }}</h1>
        <p class="poem-author-main" v-if="poem.author_data?.name || poem.author">{{ poem.author_data?.name || poem.author }}</p>
        <div class="poem-date-main">
          <v-icon size="14">mdi-calendar-outline</v-icon>
          <span>{{ poem.created_at }}</span>
        </div>
      </div>

      <!-- Poem Content -->
      <div class="poem-content-main" v-html="poem.content"></div>

      <!-- Poem Actions -->
      <div class="poem-actions-main">
        <v-btn
          v-if="!isLiked"
          @click="addLikedPoem(poem)"
          :icon="'mdi-heart-outline'"
          color="pink-lighten-1"
          variant="outlined"
          size="small"
          class="action-btn-main"
        >
        </v-btn>
        <v-btn
          v-else
          @click="undoLikedPoem(poem)"
          :icon="'mdi-heart'"
          color="pink-lighten-1"
          variant="flat"
          size="small"
          class="action-btn-main"
        >
        </v-btn>

        <v-btn
          v-if="!isBookmark"
          @click="addBookmark(poem)"
          :icon="'mdi-bookmark-outline'"
          color="blue-lighten-1"
          variant="outlined"
          size="small"
          class="action-btn-main"
        >
        </v-btn>
        <v-btn
          v-else
          @click="undoBookmark(poem)"
          :icon="'mdi-bookmark'"
          color="blue-lighten-1"
          variant="flat"
          size="small"
          class="action-btn-main"
        >
        </v-btn>
      </div>
    </div>

    <!-- Loading State -->
    <v-card v-else class="loading-card">
      <v-card-text class="d-flex align-center justify-center ga-3">
        <v-progress-circular
          indeterminate
          color="grey-lighten-1"
          size="20"
          width="2"
        ></v-progress-circular>
        <span class="text-grey-lighten-1">Şiir yükleniyor...</span>
      </v-card-text>
    </v-card>

    <!-- Related Poems Section -->
    <div class="related-poems-section">
      <h2 class="related-poems-title">Diğer Şiirler</h2>

      <div v-if="!loading" class="related-poems-grid">
        <div v-for="a in randomPeoms" :key="a.id" class="related-poem-item">
          <PoemCard :poem="a"/>
        </div>
      </div>

      <v-card v-else class="loading-card">
        <v-card-text class="d-flex align-center justify-center ga-3">
          <v-progress-circular
            indeterminate
            color="grey-lighten-1"
            size="20"
            width="2"
          ></v-progress-circular>
          <span class="text-grey-lighten-1">Diğer şiirler yükleniyor...</span>
        </v-card-text>
      </v-card>

      <v-btn
        @click="router.push('/poems')"
        color="grey-lighten-1"
        variant="outlined"
        class="all-poems-btn"
      >
        <v-icon size="18" class="mr-2">mdi-book-open-page-variant</v-icon>
        Tüm Şiirleri Gör
      </v-btn>
    </div>
  </div>
</template>

