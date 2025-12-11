<script setup lang="ts">
import { ref } from "vue";
import { useLazyLoad } from '@/composables/useLazyLoad'
import PageHeader from "../../components/common/PageHeader.vue";
import SearchSection from "../../components/common/SearchSection.vue";
import LoadingState from "../../components/common/LoadingState.vue";
import ErrorState from "../../components/common/ErrorState.vue";
import EmptyState from "../../components/common/EmptyState.vue";
import LazyLoadingIndicator from "../../components/common/LazyLoadingIndicator.vue";
import EndOfListMessage from "../../components/common/EndOfListMessage.vue";

interface ReminderType {
  id: number
  text: string
  permission: number
  created_at: string
  updated_at: string
}

// Search state
const search = ref("")
const queryParams = ref({ search: "" })

const { items: reminders, loading, hasMore, total, error, refresh } = useLazyLoad<ReminderType>(
  '/reminders-paginated',
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

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = Math.abs(now.getTime() - date.getTime())
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays === 0) return 'Bugün'
  if (diffDays === 1) return 'Dün'
  if (diffDays < 7) return `${diffDays} gün önce`

  return date.toLocaleDateString('tr-TR', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}
</script>

<template>
  <div class="reminder-page d-flex flex-column">
    <!-- Initial Loading -->
    <LoadingState
      v-if="loading && reminders.length === 0 && !search"
      message="Sözler yükleniyor..."
      color="grey-lighten-1"
    />

    <!-- Error State -->
    <ErrorState v-else-if="error" :error="error" />

    <!-- Content Wrapper (Always show header and search) -->
    <div v-else class="content-wrapper">
      <!-- Page Header -->
      <PageHeader
        icon="mdi-format-quote-close"
        icon-color="#2797b0"
        title="Sözler"
        subtitle="Hikmet, duygu dolu sözler ve hatırlatmalar"
        :total="total"
        total-label="söz"
      />

      <!-- Search Section -->
      <SearchSection
        v-model="search"
        label="Söz Ara"
        placeholder="Söz içeriğinde ara..."
        @search="handleSearch"
        @clear="clearSearch"
      />

      <!-- Empty State - No reminders at all -->
      <EmptyState
        v-if="!loading && reminders.length === 0 && !search"
        icon="mdi-format-quote-close"
        icon-color="grey-lighten-1"
        title="Henüz Söz Eklenmemiş"
        message="Şu anda görüntülenecek söz bulunmuyor."
      />

      <!-- Empty State - No search results -->
      <EmptyState
        v-else-if="!loading && reminders.length === 0 && search"
        icon="mdi-magnify-remove-outline"
        icon-color="grey-lighten-1"
        title="Sonuç Bulunamadı"
        :message="`'${search}' için sonuç bulunamadı. Farklı bir arama yapın.`"
      />

      <!-- Quotes List -->
      <div v-else class="quotes-list mt-8 mb-16">
        <div
          v-for="(reminder, index) in reminders"
          :key="reminder.id"
          class="quote-card"
          :style="{ animationDelay: `${index * 0.1}s` }"
        >
          <!-- Quote Icon -->
          <div class="quote-icon-wrapper">
            <v-icon size="32" class="quote-icon">mdi-format-quote-open</v-icon>
          </div>

          <!-- Quote Content -->
          <div class="quote-content" v-html="reminder.text"></div>

          <!-- Quote Footer -->
          <div class="quote-footer">
            <div class="quote-date">
              <v-icon size="16" class="mr-1">mdi-calendar-outline</v-icon>
              {{ formatDate(reminder.created_at) }}
            </div>
          </div>
        </div>
      </div>

      <!-- Lazy Loading Indicator -->
      <LazyLoadingIndicator
        v-if="loading && hasMore"
        message="Daha fazla söz yükleniyor..."
        color="grey-lighten-1"
      />

      <!-- End of List Message -->
      <EndOfListMessage
        v-if="!hasMore && reminders.length > 0"
        message="Tüm sözler gösteriliyor"
        icon-color="grey-darken-1"
      />
    </div>
  </div>
</template>

<style scoped>
/* Page Container */
.reminder-page {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  padding: 40px 20px;
  animation: fadeIn 0.6s ease-out;
  min-height: calc(100vh - 200px);
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
  width: 100%;
}

/* Quotes List */
.quotes-list {
  width: 100%;
}

/* Quote Card */
.quote-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  border-radius: 16px;
  padding: 32px;
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
  animation: fadeInUp 0.6s ease-out both;
  margin-bottom: 15px;
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

.quote-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #2797b0 0%, #3a9ab7 50%, #3f51b5 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.quote-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.4);
  border-color: #616161;
}

.quote-card:hover::before {
  opacity: 0.8;
}

/* Quote Icon */
.quote-icon-wrapper {
  margin-bottom: 20px;
}

.quote-icon {
  color: #2797b0;
  opacity: 0.7;
  transition: all 0.3s ease;
}

.quote-card:hover .quote-icon {
  opacity: 1;
  transform: rotate(-5deg) scale(1.1);
}

/* Quote Content */
.quote-content {
  font-size: 17px;
  line-height: 1.8;
  color: #e0e0e0;
  margin-bottom: 20px;
  letter-spacing: 0.3px;
  font-weight: 400;
  font-style: italic;
  position: relative;
}

.quote-content :deep(p) {
  margin-bottom: 6px;
}

.quote-content :deep(strong) {
  color: #ffffff;
  font-weight: 600;
}

.quote-content :deep(em) {
  color: #bdbdbd;
}

/* Quote Footer */
.quote-footer {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.quote-date {
  display: flex;
  align-items: center;
  font-size: 13px;
  color: #9e9e9e;
  font-weight: 500;
  letter-spacing: 0.3px;
  transition: color 0.3s ease;
}

.quote-card:hover .quote-date {
  color: #bdbdbd;
}

/* Responsive Design */
@media (max-width: 768px) {
  .reminder-page {
    padding: 24px 12px;
  }

  .quote-card {
    padding: 24px;
  }

  .quote-content {
    font-size: 16px;
    line-height: 1.7;
  }
}

@media (max-width: 480px) {
  .quote-card {
    padding: 20px;
  }

  .quote-content {
    font-size: 15px;
  }
}
</style>
