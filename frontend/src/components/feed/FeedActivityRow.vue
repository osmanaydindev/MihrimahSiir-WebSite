<script setup lang="ts">
import { getImageUrl } from '@/utils/imageHelper'
import type { FeedItem } from '@/types'

// Aktivite satırları (şiir beğenisi, kitap okundu) alıntı değil; bilerek
// beğen/kaydet taşımıyorlar — bir beğeniyi beğenmek anlamsız olurdu.
const props = defineProps<{ item: FeedItem }>()

const isPoem = props.item.kind === 'poem_like'

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
  <div class="activity-row">
    <router-link
      :to="{ name: 'UserProfile', params: { username: item.author.username } }"
      class="avatar-link"
    >
      <v-avatar size="28">
        <v-img v-if="item.author.profile_image" :src="getImageUrl(item.author.profile_image)" cover />
        <span v-else class="avatar-initial">{{ item.author.username.charAt(0).toUpperCase() }}</span>
      </v-avatar>
    </router-link>

    <div class="activity-body">
      <p class="activity-text">
        <span class="who">@{{ item.author.username }}</span>
        <template v-if="isPoem"> bir şiiri beğendi</template>
        <template v-else> bir kitabı okudu</template>
      </p>

      <router-link
        v-if="isPoem && item.poem_slug"
        :to="{ name: 'Poem', params: { slug: item.poem_slug } }"
        class="activity-target"
      >
        {{ item.poem_title }}<template v-if="item.poem_author"> — {{ item.poem_author }}</template>
      </router-link>

      <router-link
        v-else-if="item.book_slug"
        :to="{ name: 'BookDetail', params: { slug: item.book_slug } }"
        class="activity-target"
      >
        {{ item.book_name }}
      </router-link>
    </div>

    <span class="activity-time">{{ relativeTime(item.created_at) }}</span>
  </div>
</template>

<style scoped>
.activity-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 14px;
  border-left: 2px solid rgba(255, 255, 255, 0.07);
}

.avatar-link { flex-shrink: 0; text-decoration: none; }

.avatar-initial {
  font-weight: 700;
  font-size: 12px;
  color: #e0e0e0;
}

.activity-body { flex: 1; min-width: 0; }

.activity-text {
  font-size: 13px;
  color: #9e9e9e;
  margin: 0;
}

.who { color: #e0e0e0; font-weight: 600; }

.activity-target {
  display: inline-block;
  font-size: 13px;
  color: #90caf9;
  text-decoration: none;
  margin-top: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.activity-target:hover { text-decoration: underline; }

.activity-time {
  font-size: 11px;
  color: #757575;
  flex-shrink: 0;
}

@media (max-width: 599px) {
  .activity-row { padding: 8px 10px; }
  .activity-time { display: none; }
}
</style>
