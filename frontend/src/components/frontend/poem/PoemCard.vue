<script setup lang="ts">
import {computed, ref, watch} from "vue";
import {useRouter, useRoute} from "vue-router";
import {useAppStore} from "../../../store/app";
import {usePoemActions} from "../../../composables/usePoemActions";

const store = useAppStore()
const router = useRouter()
const route = useRoute()
const { likePoem, unlikePoem, bookmarkPoem, unbookmarkPoem } = usePoemActions()

const props = defineProps({
  poem: {
    type: Object,
    required: true,
  }
})

const emit = defineEmits(['poem-liked', 'poem-unliked', 'poem-bookmarked', 'poem-unbookmarked'])

// Local reactive like count for optimistic updates
const localLikeCount = ref(props.poem.like_count || 0)

// Watch for poem changes to update local like count
watch(() => props.poem.like_count, (newCount) => {
  if (newCount !== undefined) {
    localLikeCount.value = newCount
  }
})

const isLiked = computed(() => {
  if (!props.poem || !props.poem.id) return false
  return store.likedPoemIds.has(props.poem.id)
})

const isBookmark = computed(() => {
  if (!props.poem || !props.poem.id) return false
  return store.bookmarkedPoemIds.has(props.poem.id)
})

const gotopoem = async (slug: string) => {
  await router.push(`/poem/${slug}`)
}

const addLikedPoem = async (poem: any) => {
  // Optimistic update: increment like count immediately
  localLikeCount.value++

  try {
    await likePoem(poem)
    emit('poem-liked', poem)
  } catch (error) {
    // Rollback on error
    localLikeCount.value--
  }
}

const addBookmark = async (poem: any) => {
  try {
    await bookmarkPoem(poem)
    emit('poem-bookmarked', poem)
  } catch (error) {
    // Error already logged in composable
  }
}

const undoLikedPoem = async (poem: any) => {
  // Optimistic update: decrement like count immediately
  localLikeCount.value--

  try {
    await unlikePoem(poem)
    emit('poem-unliked', poem)
  } catch (error) {
    // Rollback on error
    localLikeCount.value++
  }
}

const undoBookmark = async (poem: any) => {
  try {
    await unbookmarkPoem(poem)
    emit('poem-unbookmarked', poem)
  } catch (error) {
    // Error already logged in composable
  }
}
</script>

<template>
  <v-card
    class="poem-card position-relative"
    rounded="lg"
    @click="gotopoem(props.poem.slug)"
  >
    <!-- Poem Content Area -->
    <div class="poem-content-wrapper">
      <!-- Title -->
      <h3 class="poem-title text-white mb-2">{{ props.poem.title }}</h3>

      <!-- Author -->
      <p class="poem-author text-grey-lighten-1 mb-1" v-if="props.poem.author">{{ props.poem.author }}</p>

      <!-- Date -->
      <p class="poem-date text-grey-lighten-1 ma-0 mb-3">{{ props.poem.created_at }}</p>

      <!-- Poem Preview -->
      <div class="poem-preview">
        <div class="poem-text text-grey-lighten-2" v-html="props.poem.content"></div>
        <div class="poem-fade"></div>
      </div>
    </div>

    <!-- Actions -->
    <v-card-actions class="poem-actions pa-4">
      <v-btn
        variant="outlined"
        color="grey-lighten-1"
        class="read-more-btn flex-grow-1"
        size="small"
      >
        <v-icon size="16" class="mr-2">mdi-book-open-page-variant-outline</v-icon>
        Devamını oku
      </v-btn>

      <div class="d-flex ga-1 align-center px-2">
        <!-- Like Button & Count -->
        <div v-if="route.path !== '/bookmark'" class="d-flex align-center" style="gap: 2px;">
          <v-btn
            @click.stop="isLiked ? undoLikedPoem(props.poem) : addLikedPoem(props.poem)"
            :icon="isLiked ? 'mdi-heart' : 'mdi-heart-outline'"
            :color="'pink-lighten-1'"
            variant="text"
            size="small"
            density="compact"
            class="action-icon-btn pa-0"
            style="min-width: auto; width: auto;"
          ></v-btn>
          <span v-if="localLikeCount !== undefined" class="like-count-text pl-1 text-pink-lighten-1">{{ localLikeCount }}</span>
        </div>

        <!-- Bookmark Button -->
        <v-btn
          v-if="route.path !== '/liked-poems'"
          @click.stop="isBookmark ? undoBookmark(props.poem) : addBookmark(props.poem)"
          :icon="isBookmark ? 'mdi-bookmark' : 'mdi-bookmark-outline'"
          :color="isBookmark ? 'blue-lighten-1' : 'grey-lighten-1'"
          variant="text"
          size="small"
          class="action-icon-btn"
        ></v-btn>
      </div>
    </v-card-actions>
  </v-card>
</template>

<style scoped>
/* Poem Card Container */
.poem-card {
  width: 100%;
  max-width: 380px;
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  transition: all 0.3s ease;
  overflow: hidden;
  cursor: pointer;
  display: flex;
  flex-direction: column;
}

.poem-card:hover {
  transform: translateY(-8px);
  border-color: #616161;
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.4);
}

/* Content Wrapper */
.poem-content-wrapper {
  padding: 20px 20px 12px 20px;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
}

/* Title */
.poem-title {
  font-size: 18px;
  font-weight: 600;
  line-height: 1.4;
  letter-spacing: 0.3px;
  transition: color 0.3s ease;
}

.poem-card:hover .poem-title {
  color: #e0e0e0 !important;
}

/* Author */
.poem-author {
  font-size: 13px;
  font-weight: 500;
  opacity: 0.9;
  letter-spacing: 0.3px;
  font-style: italic;
}

/* Date */
.poem-date {
  font-size: 13px;
  font-weight: 400;
  opacity: 0.8;
  letter-spacing: 0.2px;
}

/* Like Count Text */
.like-count-text {
  font-size: 13px;
  font-weight: 500;
  min-width: 20px;
  text-align: center;
  line-height: 1;
  display: flex;
  align-items: center;
}

/* Poem Preview */
.poem-preview {
  position: relative;
  height: 155px;
  overflow: hidden;
  margin-bottom: 8px;
}

.poem-text {
  font-size: 14px;
  line-height: 1.6;
  letter-spacing: 0.2px;
  opacity: 0.9;
}

.poem-text :deep(p) {
  margin-bottom: 8px;
}

/* Fade Effect */
.poem-fade {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: linear-gradient(to bottom, rgba(29, 29, 29, 0) 0%, rgba(29, 29, 29, 0.8) 50%, rgba(29, 29, 29, 1) 100%);
  pointer-events: none;
}

/* Actions */
.poem-actions {
  display: flex;
  justify-content: space-between;
  flex-direction: row;
  flex-wrap: wrap;

  border-top: 1px solid rgba(255, 255, 255, 0.05);
  padding: 12px 16px !important;
  gap: 8px;
}

/* Read More Button */
.read-more-btn {
  text-transform: none;
  font-weight: 500;
  letter-spacing: 0.5px;
  transition: all 0.3s ease;
  border-width: 1px;
  font-size: 13px;
}

.read-more-btn:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: #fff;
}

/* Action Icon Buttons */
.action-icon-btn {
  transition: all 0.3s ease;
}

.action-icon-btn:hover {
  background: rgba(255, 255, 255, 0.05);
}

/* Responsive Typography */
@media (max-width: 600px) {
  .poem-title {
    font-size: 16px;
  }

  .poem-date {
    font-size: 12px;
  }

  .poem-text {
    font-size: 13px;
  }

  .poem-content-wrapper {
    padding: 16px 16px 10px 16px;
  }
}
</style>
