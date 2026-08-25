<script setup lang="ts">
import { ref } from 'vue'
import { useSanitizer } from '@/composables/useSanitizer'
import { feedService } from '@/services/api/feedService'
import { getImageUrl } from '@/utils/imageHelper'
import type { FeedItem } from '@/types'

const props = defineProps<{ item: FeedItem }>()

const { sanitizeHTML } = useSanitizer()

// Alıntı gövdesi Quill ile yazılmış HTML — ham v-html asla kullanılmaz.
const safeContent = sanitizeHTML(props.item.content ?? '')

// İyimser güncelleme + hata durumunda geri alma (usePoemActions deseni).
const liked = ref(props.item.is_liked)
const likeCount = ref(props.item.like_count)
const saved = ref(props.item.is_saved)
const busy = ref(false)

async function toggleLike() {
  if (busy.value) return
  busy.value = true
  const wasLiked = liked.value
  liked.value = !wasLiked
  likeCount.value += wasLiked ? -1 : 1
  try {
    await (wasLiked ? feedService.unlike(props.item.id) : feedService.like(props.item.id))
  } catch {
    liked.value = wasLiked
    likeCount.value += wasLiked ? 1 : -1
  } finally {
    busy.value = false
  }
}

async function toggleSave() {
  if (busy.value) return
  busy.value = true
  const wasSaved = saved.value
  saved.value = !wasSaved
  try {
    await (wasSaved ? feedService.unsave(props.item.id) : feedService.save(props.item.id))
  } catch {
    saved.value = wasSaved
  } finally {
    busy.value = false
  }
}

function relativeTime(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime()
  const min = Math.floor(diff / 60000)
  if (min < 1) return 'az önce'
  if (min < 60) return `${min} dakika önce`
  const hour = Math.floor(min / 60)
  if (hour < 24) return `${hour} saat önce`
  const day = Math.floor(hour / 24)
  if (day < 30) return `${day} gün önce`
  return new Date(iso).toLocaleDateString('tr-TR', { day: 'numeric', month: 'long', year: 'numeric' })
}
</script>

<template>
  <v-card class="feed-card" rounded="lg">
    <div class="feed-head">
      <router-link :to="{ name: 'UserProfile', params: { username: item.author.username } }" class="author-link">
        <v-avatar size="34">
          <v-img v-if="item.author.profile_image" :src="getImageUrl(item.author.profile_image)" cover />
          <span v-else class="avatar-initial">{{ item.author.username.charAt(0).toUpperCase() }}</span>
        </v-avatar>
        <span class="author-name">@{{ item.author.username }}</span>
      </router-link>
      <span class="feed-time">{{ relativeTime(item.created_at) }}</span>
    </div>

    <router-link
      v-if="item.book_slug"
      :to="{ name: 'BookDetail', params: { slug: item.book_slug } }"
      class="book-line"
    >
      <v-icon size="14" class="mr-1">mdi-book-open-page-variant-outline</v-icon>
      {{ item.book_name }}<template v-if="item.page"> · s.{{ item.page }}</template>
    </router-link>

    <h3 v-if="item.title" class="feed-title">{{ item.title }}</h3>
    <div class="feed-content" v-html="safeContent"></div>

    <div class="feed-actions">
      <button class="act" :class="{ on: liked }" :disabled="busy" @click="toggleLike">
        <v-icon size="18">{{ liked ? 'mdi-heart' : 'mdi-heart-outline' }}</v-icon>
        <span class="act-count">{{ likeCount }}</span>
      </button>
      <button class="act" :class="{ on: saved }" :disabled="busy" @click="toggleSave">
        <v-icon size="18">{{ saved ? 'mdi-bookmark' : 'mdi-bookmark-outline' }}</v-icon>
        <span class="act-label">{{ saved ? 'Kaydedildi' : 'Kaydet' }}</span>
      </button>
    </div>
  </v-card>
</template>

<style scoped>
.feed-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid rgba(255, 255, 255, 0.06);
  padding: 16px;
}

.feed-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 10px;
}

.author-link {
  display: flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
  min-width: 0;
}

.avatar-initial {
  font-weight: 700;
  font-size: 14px;
  color: #e0e0e0;
}

.author-name {
  font-weight: 600;
  font-size: 14px;
  color: #ffffff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.feed-time {
  font-size: 12px;
  color: #9e9e9e;
  flex-shrink: 0;
}

.book-line {
  display: inline-flex;
  align-items: center;
  font-size: 12px;
  color: #9e9e9e;
  text-decoration: none;
  margin-bottom: 8px;
}

.book-line:hover {
  color: #bdbdbd;
}

.feed-title {
  font-size: 15px;
  font-weight: 600;
  color: #ffffff;
  margin-bottom: 6px;
  line-height: 1.35;
}

.feed-content {
  font-size: 14px;
  line-height: 1.65;
  color: #d6d6d6;
  overflow-wrap: anywhere;
}

.feed-content :deep(p) {
  margin-bottom: 6px;
}

.feed-actions {
  display: flex;
  align-items: center;
  gap: 18px;
  margin-top: 12px;
  padding-top: 10px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.act {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: none;
  border: none;
  color: #9e9e9e;
  font-size: 13px;
  cursor: pointer;
  padding: 4px 2px;
  transition: color 0.15s;
}

.act:hover:not(:disabled) { color: #e0e0e0; }
.act:disabled { opacity: 0.6; cursor: default; }
.act.on { color: #ec407a; }
.act-count { font-variant-numeric: tabular-nums; }

@media (max-width: 599px) {
  .feed-card { padding: 12px; }
  .feed-content { font-size: 13px; }
  .act-label { display: none; }
}
</style>
