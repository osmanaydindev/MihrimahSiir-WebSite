<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import type { Poem, Book } from '../../types'

interface Props {
  modelValue: boolean
  title: string
  username: string
  type: 'poems' | 'books' | 'bookmarked-poems' | 'comments'
  emptyMessage?: string
}

const props = withDefaults(defineProps<Props>(), {
  emptyMessage: 'Henüz bir şey yok'
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const router = useRouter()
const items = ref<Poem[] | Book[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const displayValue = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

// Lazy load items when dialog opens
watch(() => props.modelValue, async (newValue) => {
  if (newValue && items.value.length === 0) {
    await loadItems()
  }
})

const loadItems = async () => {
  loading.value = true
  error.value = null

  try {
    let endpoint = ''

    switch (props.type) {
      case 'poems':
        endpoint = `/user-profile/${props.username}/liked-poems`
        break
      case 'bookmarked-poems':
        endpoint = `/user-profile/${props.username}/bookmarked-poems`
        break
      case 'books':
        endpoint = `/user-profile/${props.username}/read-books`
        break
      case 'comments':
        endpoint = `/user-profile/${props.username}/comments`
        break
    }

    const response = await axios.get(endpoint)

    if (props.type === 'poems' || props.type === 'bookmarked-poems') {
      items.value = response.data.poems || []
    } else if (props.type === 'books') {
      items.value = response.data.books || []
    } else if (props.type === 'comments') {
      items.value = response.data.comments || []
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Veriler yüklenirken bir hata oluştu'
    console.error('Failed to load items:', err)
  } finally {
    loading.value = false
  }
}

const handleItemClick = (item: any) => {
  if (props.type === 'poems' || props.type === 'bookmarked-poems') {
    const poem = item as Poem
    router.push(`/poem/${poem.slug}`)
  } else if (props.type === 'books') {
    const book = item as Book
    router.push(`/book/${book.slug}`)
  } else if (props.type === 'comments') {
    // Navigate to the book page for comments
    if (item.book?.slug) {
      router.push(`/book/${item.book.slug}`)
    }
  }
  displayValue.value = false
}

const getItemTitle = (item: any) => {
  if (props.type === 'poems' || props.type === 'bookmarked-poems') {
    return (item as Poem).title
  } else if (props.type === 'books') {
    return (item as Book).name
  } else if (props.type === 'comments') {
    return item.title || 'İsimsiz Yorum'
  }
  return ''
}

const getItemSubtitle = (item: any) => {
  if (props.type === 'poems' || props.type === 'bookmarked-poems') {
    const poem = item as any
    return poem.author_data?.name || poem.author || 'Bilinmeyen'
  } else if (props.type === 'books') {
    const book = item as any
    return book.author_data?.name || book.author || 'Bilinmeyen'
  } else if (props.type === 'comments') {
    return item.book?.name || 'Kitap'
  }
  return ''
}

const getItemIcon = () => {
  if (props.type === 'poems' || props.type === 'bookmarked-poems') {
    return 'mdi-format-quote-close'
  } else if (props.type === 'books') {
    return 'mdi-book-open-variant'
  } else if (props.type === 'comments') {
    return 'mdi-comment-text'
  }
  return 'mdi-help'
}
</script>

<template>
  <v-dialog
    v-model="displayValue"
    max-width="700"
    scrollable
  >
    <v-card>
      <v-card-title class="d-flex align-center justify-space-between">
        <div class="d-flex align-center">
          <v-icon :icon="getItemIcon()" class="mr-2" />
          {{ title }}
        </div>
        <v-btn
          icon="mdi-close"
          variant="text"
          @click="displayValue = false"
        />
      </v-card-title>

      <v-divider />

      <v-card-text class="pa-0" style="max-height: 500px; overflow-y: hidden">
        <!-- Loading State -->
        <div v-if="loading" class="text-center py-12">
          <v-progress-circular
            indeterminate
            color="primary"
            size="48"
          />
          <div class="text-body-1 text-grey mt-4">Yükleniyor...</div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="text-center py-12">
          <v-icon
            icon="mdi-alert-circle"
            size="64"
            class="mb-4 text-error"
          />
          <div class="text-body-1 text-error">{{ error }}</div>
          <v-btn
            color="primary"
            variant="text"
            class="mt-4"
            @click="loadItems"
          >
            Tekrar Dene
          </v-btn>
        </div>

        <!-- Items List -->
        <v-list v-else-if="items.length > 0" lines="two">
          <v-virtual-scroll
            :items="items"
            :height="Math.min(items.length * 72, 500)"
            item-height="72"
          >
            <template v-slot:default="{ item }">
              <v-list-item
                :key="item.id"
                @click="handleItemClick(item)"
                class="list-item-hover"
              >
                <template v-slot:prepend>
                  <v-avatar :color="
                    type === 'poems' || type === 'bookmarked-poems' ? 'pink-lighten-4' :
                    type === 'books' ? 'orange-lighten-4' :
                    'purple-lighten-4'
                  ">
                    <v-icon :icon="getItemIcon()" />
                  </v-avatar>
                </template>

                <v-list-item-title>{{ getItemTitle(item) }}</v-list-item-title>
                <v-list-item-subtitle>{{ getItemSubtitle(item) }}</v-list-item-subtitle>

                <template v-slot:append>
                  <v-icon icon="mdi-chevron-right" />
                </template>
              </v-list-item>
              <v-divider v-if="item.id !== items[items.length - 1].id" />
            </template>
          </v-virtual-scroll>
        </v-list>

        <!-- Empty State -->
        <v-card-text v-else class="text-center py-8">
          <v-icon
            :icon="getItemIcon()"
            size="64"
            class="mb-4 text-grey"
          />
          <div class="text-body-1 text-grey">{{ emptyMessage }}</div>
        </v-card-text>
      </v-card-text>

      <v-divider />

      <v-card-actions>
        <v-spacer />
        <v-btn
          color="primary"
          variant="text"
          @click="displayValue = false"
        >
          Kapat
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<style scoped>
.list-item-hover {
  cursor: pointer;
  transition: background-color 0.2s;
}

.list-item-hover:hover {
  background-color: rgba(0, 0, 0, 0.04);
}
</style>
