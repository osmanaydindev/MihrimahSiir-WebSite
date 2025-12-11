<script setup lang="ts">
import { useLazyLoad } from '@/composables/useLazyLoad'
import AuthorCard from '../../components/frontend/author/AuthorCard.vue'

interface AuthorType {
  id: number
  name: string
  bio?: string
  birth_year?: number
  death_year?: number
  nationality?: string
  image?: string
  slug: string
  poems?: any[]
  books?: any[]
}

const { items: authors, loading, hasMore, total, error } = useLazyLoad<AuthorType>('/get-authors')
</script>

<template>
  <div class="authors-page">
    <!-- Page Header -->
    <div class="page-header text-center mb-8">
      <div class="header-icon-wrapper mb-4">
        <v-icon size="48" class="header-icon">mdi-account-group</v-icon>
      </div>
      <h1 class="page-title">Yazarlar</h1>
      <p class="page-subtitle text-grey-lighten-1">Şiir ve edebiyat dünyasının ustalar</p>
      <p v-if="total > 0" class="total-count text-grey-darken-1 mt-2">
        Toplam {{ total }} yazar
      </p>
    </div>

    <!-- Initial Loading -->
    <div v-if="loading && authors.length === 0" class="loading-container">
      <v-card class="loading-card">
        <v-card-text class="d-flex flex-column align-center justify-center ga-4">
          <v-progress-circular
            indeterminate
            color="grey-lighten-1"
            size="48"
            width="3"
          ></v-progress-circular>
          <span class="text-grey-lighten-1 text-h6">Yazarlar yükleniyor...</span>
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
    <div v-else-if="!loading && authors.length === 0" class="empty-state">
      <v-card class="empty-card">
        <v-icon size="80" color="grey-lighten-1" class="mb-4">mdi-account-group</v-icon>
        <v-card-title class="text-h5 text-white font-weight-bold mb-2">
          Henüz Yazar Eklenmemiş
        </v-card-title>
        <v-card-text class="text-body-1 text-grey-lighten-1">
          Şu anda görüntülenecek yazar bulunmuyor.
        </v-card-text>
      </v-card>
    </div>

    <!-- Authors Grid -->
    <div v-else>
<!--      <v-container fluid>-->
        <v-row class="w-full">
          <v-col
            v-for="author in authors"
            :key="author.id"
            class="pa-3"
            cols="12"
            sm="6"
            md="4"
            xl="3"
          >
            <div class="author-grid-item">
              <AuthorCard :author="author" />
            </div>
          </v-col>
        </v-row>
<!--      </v-container>-->


      <!-- <div class="authors-grid">
        <div
          v-for="author in authors"
          :key="author.id"
          class="author-grid-item"
        >
          <AuthorCard :author="author" />
        </div>
      </div> -->

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
            <span class="text-grey-lighten-1">Daha fazla yazar yükleniyor...</span>
          </v-card-text>
        </v-card>
      </div>

      <!-- End of List Message -->
      <div v-if="!hasMore && authors.length > 0" class="end-message text-center mt-8 mb-4">
        <v-icon size="20" color="grey-darken-1" class="mr-2">mdi-check-circle</v-icon>
        <span class="text-grey-darken-1">Tüm yazarlar gösteriliyor</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.authors-page {
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
  color: #2797b0;
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

/* Authors Grid */
.authors-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.author-grid-item {
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
  .authors-page {
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

  .authors-grid {
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 20px;
  }
}

@media (max-width: 480px) {
  .page-title {
    font-size: 28px;
  }

  .authors-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}
</style>
