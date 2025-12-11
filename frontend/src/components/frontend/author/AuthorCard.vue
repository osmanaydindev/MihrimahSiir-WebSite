<script setup lang="ts">
import { computed } from 'vue'

interface Author {
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

const props = defineProps<{
  author: Author
}>()

const yearRange = computed(() => {
  if (!props.author.birth_year && !props.author.death_year) return null

  const birth = props.author.birth_year || '?'
  const death = props.author.death_year || ''

  return death ? `${birth} - ${death}` : `${birth} -`
})

const authorImage = computed(() => {
  return props.author.image || '/uploads/default-author.jpg'
})
</script>

<template>
  <v-card class="author-card" rounded="lg">
    <router-link
      class="text-decoration-none card-link"
      :to="{ name: 'AuthorDetail', params: { slug: author.slug } }"
    >
      <!-- Author Portrait -->
      <div class="author-portrait-wrapper">
        <v-img
          :src="authorImage"
          class="author-portrait"
          height="320px"
          cover
        >
          <div class="portrait-overlay"></div>
        </v-img>
      </div>

      <!-- Author Info -->
      <div class="author-info pa-4">
        <!-- Name -->
        <h3 class="author-name text-white mb-2">{{ author.name }}</h3>

        <!-- Year Range -->
        <p v-if="yearRange" class="author-years text-grey-lighten-1 mb-2">
          <v-icon size="14" class="mr-1">mdi-calendar-range</v-icon>
          {{ yearRange }}
        </p>

        <!-- Nationality -->
        <p v-if="author.nationality" class="author-nationality text-grey-lighten-1 mb-3">
          <v-icon size="14" class="mr-1">mdi-earth</v-icon>
          {{ author.nationality }}
        </p>

        <!-- Stats -->
        <div class="author-stats d-flex ga-3 mt-3">
          <div class="stat-item">
            <v-icon size="16" color="blue-lighten-1">mdi-book-open-page-variant</v-icon>
            <span class="stat-text ml-1">{{ author.poems?.length || 0 }} şiir</span>
          </div>
          <div class="stat-item">
            <v-icon size="16" color="green-lighten-1">mdi-book</v-icon>
            <span class="stat-text ml-1">{{ author.books?.length || 0 }} kitap</span>
          </div>
        </div>
      </div>
    </router-link>
  </v-card>
</template>

<style scoped>
.author-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.author-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.4);
  border-color: #616161;
}

.card-link {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* Portrait Section */
.author-portrait-wrapper {
  position: relative;
  overflow: hidden;
}

.author-portrait {
  transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.author-card:hover .author-portrait {
  transform: scale(1.05);
}

.portrait-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 50%;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.9), transparent);
}

/* Info Section */
.author-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.author-name {
  font-size: 20px;
  font-weight: 600;
  letter-spacing: 0.3px;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.3s ease;
}

.author-card:hover .author-name {
  color: #e3f2fd !important;
}

.author-years,
.author-nationality {
  font-size: 14px;
  display: flex;
  align-items: center;
  letter-spacing: 0.2px;
}

/* Stats */
.author-stats {
  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  margin-top: auto;
}

.stat-item {
  display: flex;
  align-items: center;
  font-size: 13px;
  color: #bdbdbd;
  transition: color 0.3s ease;
}

.author-card:hover .stat-item {
  color: #e0e0e0;
}

.stat-text {
  font-weight: 500;
  letter-spacing: 0.2px;
}

/* Responsive */
@media (max-width: 600px) {
  .author-portrait {
    height: 280px !important;
  }

  .author-name {
    font-size: 18px;
  }

  .author-years,
  .author-nationality {
    font-size: 13px;
  }

  .stat-item {
    font-size: 12px;
  }
}
</style>
