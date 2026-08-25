<script setup lang="ts">
import { ref, computed } from 'vue'
import { useLazyLoad } from '@/composables/useLazyLoad'
import FeedQuoteCard from '@/components/feed/FeedQuoteCard.vue'
import FeedActivityRow from '@/components/feed/FeedActivityRow.vue'
import type { FeedItem } from '@/types'

// İki sekme aynı bileşeni kullanıyor; sadece uç nokta değişiyor.
const tab = ref<'feed' | 'saved'>('feed')
const endpoint = computed(() => (tab.value === 'feed' ? '/feed' : '/saved-comments'))

const { items, loading, hasMore, error, refresh } = useLazyLoad<FeedItem>(endpoint)

function switchTab(next: 'feed' | 'saved') {
  if (tab.value === next) return
  tab.value = next
  refresh()
}
</script>

<template>
  <div class="feed-page">
    <PageHeader
      icon="mdi-timeline-text-outline"
      icon-color="#7e57c2"
      title="Akış"
      subtitle="Arkadaşlarının ve herkese açık profillerin paylaşımları"
    />

    <div class="feed-tabs">
      <button class="feed-tab" :class="{ active: tab === 'feed' }" @click="switchTab('feed')">
        Akış
      </button>
      <button class="feed-tab" :class="{ active: tab === 'saved' }" @click="switchTab('saved')">
        Kaydettiklerim
      </button>
    </div>

    <!-- İlk yükleme -->
    <LoadingState v-if="loading && items.length === 0" />

    <ErrorState v-else-if="error" :message="error" @retry="refresh" />

    <EmptyState
      v-else-if="!loading && items.length === 0"
      :icon="tab === 'feed' ? 'mdi-account-multiple-outline' : 'mdi-bookmark-outline'"
      icon-color="grey-darken-1"
      :title="tab === 'feed' ? 'Akışta henüz bir şey yok' : 'Henüz alıntı kaydetmedin'"
      :message="tab === 'feed'
        ? 'Arkadaş ekleyerek ya da profili herkese açık kullanıcıları takip ederek akışını doldurabilirsin.'
        : 'Beğendiğin alıntıları kaydedersen burada birikir.'"
    />

    <div v-else class="feed-list">
      <template v-for="item in items" :key="`${item.kind}-${item.id}-${item.created_at}`">
        <FeedQuoteCard v-if="item.kind === 'comment'" :item="item" />
        <FeedActivityRow v-else :item="item" />
      </template>

      <LazyLoadingIndicator v-if="loading && hasMore" />
      <EndOfListMessage v-if="!hasMore && items.length > 0" message="Hepsi bu kadar" />
    </div>
  </div>
</template>

<style scoped>
.feed-page {
  max-width: 680px;
  margin: 0 auto;
  width: 100%;
  padding: 40px 20px;
}

.feed-tabs {
  display: flex;
  gap: 4px;
  margin-bottom: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.feed-tab {
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  color: #9e9e9e;
  font-size: 14px;
  font-weight: 600;
  padding: 8px 12px;
  cursor: pointer;
  transition: color 0.15s, border-color 0.15s;
}

.feed-tab:hover { color: #e0e0e0; }

.feed-tab.active {
  color: #ffffff;
  border-bottom-color: #7e57c2;
}

.feed-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 48px;
}

@media (max-width: 768px) {
  .feed-page { padding: 16px 4px; }
  .feed-list { gap: 8px; }
}
</style>
