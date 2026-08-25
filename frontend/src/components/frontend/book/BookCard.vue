<script setup lang="ts">
//props ile book objesini alma
import {computed} from 'vue'
import {useAppStore} from "../../../store/app";
// import Book from "../../../pages/User/Book.vue";
//
// interface Book {
//   id: number
//   name: string
//   author: string
//   slug: string
//   image: string
//   page: number
//   is_deleted: boolean
//   created_at: string
//   comment: Array<Comment>|null
// }
// interface Comment {
//   id: number
//   user_id: number
//   book_id: number
//   comment: string
//   created_at: string
//   updated_at: string
// }

const store = useAppStore()
const props = defineProps<{
  book: Object
}>()
const emit = defineEmits(['addBookToReads', 'removeBookFromReads']);

const handleAddBook = () => {
  emit('addBookToReads', props.book);
};
const handleRemoveBook = () => {
  emit("removeBookFromReads", props.book);
};

// burada gelen book ID'si store'da var mı kontrol eden fonksiyon vardır
var isRead = computed(() => {
  return store.readBookIds.has(props.book.id)
})

</script>

<template>
    <v-card class="book-card position-relative" rounded="lg">
      <router-link class="text-decoration-none card-link" :to="{name: 'BookDetail', params: { slug: book.slug } }">
        <!-- Read Status Badge -->
        <div class="read-status-badge" :class="isRead ? 'read-badge' : 'unread-badge'">
          <v-icon size="14" class="mr-1">
            {{ isRead ? 'mdi-check-circle' : 'mdi-circle-outline' }}
          </v-icon>
          <span class="badge-text">{{ isRead ? 'Okundu' : 'Okunacak' }}</span>
        </div>

        <!-- Book Cover -->
        <!-- Yükseklik CSS'ten geliyor: telefonda 2 sütuna geçtiğimiz için
             sabit 280px kapak orantısız kalıyordu (bkz. .book-cover-wrapper). -->
        <div class="book-cover-wrapper">
          <v-img
            :src="book.image"
            class="book-cover"
            height="100%"
            cover
          >
            <div class="cover-overlay"></div>
          </v-img>
        </div>

        <!-- Book Info -->
        <div class="book-info pa-3 pa-sm-4">
          <h3 class="book-title text-white mb-2">{{ book.name }}</h3>
          <p class="book-author text-grey-lighten-1 mb-1">{{ book.author }}</p>
          <div class="book-meta text-grey-darken-1">
            <v-icon size="14" class="mr-1">mdi-file-document-outline</v-icon>
            <span class="meta-text">{{ book.page }} sayfa</span>
          </div>
        </div>
      </router-link>

      <!-- Action Button -->
      <!-- Etiket telefonda kısalıyor: 2 sütunda sütun ~165px ve
           "Okunanlara Ekle" sığmıyordu (bkz. .action-btn @600px). -->
      <v-card-actions class="pa-3 pa-sm-4 pt-0">
        <v-btn
          v-if="!isRead"
          block
          variant="outlined"
          color="grey-lighten-1"
          class="action-btn add-btn"
          @click.stop="handleAddBook()"
        >
          <v-icon left size="18" class="mr-2 btn-icon">mdi-plus-circle-outline</v-icon>
          <span class="btn-label-long">Okunanlara Ekle</span>
          <span class="btn-label-short">Ekle</span>
        </v-btn>

        <v-btn
          v-else
          block
          variant="outlined"
          color="green-lighten-1"
          class="action-btn remove-btn"
          @click.stop="handleRemoveBook()"
        >
          <v-icon left size="18" class="mr-2 btn-icon">mdi-check-circle</v-icon>
          <span class="btn-label-long">Okunanlardan Çıkar</span>
          <span class="btn-label-short">Çıkar</span>
        </v-btn>
      </v-card-actions>
    </v-card>
</template>

<style scoped>
/* Buton etiketi: geniş ekranda uzun, telefonda kısa hâli görünür */
.btn-label-short { display: none; }

/* Book Card Container */
.book-card {
  width: 100%;
  max-width: 320px;
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
  transition: all 0.3s ease;
  overflow: hidden;
}

.book-card:hover {
  transform: translateY(-8px);
  border-color: #616161;
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.4);
}

/* Card Link */
.card-link {
  display: block;
  color: inherit;
}

/* Read Status Badge */
.read-status-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  z-index: 10;
  display: flex;
  align-items: center;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  letter-spacing: 0.5px;
  backdrop-filter: blur(8px);
  transition: all 0.3s ease;
}

.read-badge {
  background: rgba(76, 175, 80, 0.2);
  border: 1px solid rgba(76, 175, 80, 0.4);
  color: #81c784;
}

.unread-badge {
  background: rgba(158, 158, 158, 0.2);
  border: 1px solid rgba(158, 158, 158, 0.4);
  color: #bdbdbd;
}

.badge-text {
  font-size: 11px;
  text-transform: uppercase;
}

/* Book Cover */
.book-cover-wrapper {
  position: relative;
  overflow: hidden;
  border-radius: 8px 8px 0 0;
  height: 280px;
}

.book-cover {
  transition: transform 0.4s ease;
}

.book-card:hover .book-cover {
  transform: scale(1.05);
}

.cover-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 50%;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.8) 0%, rgba(0, 0, 0, 0) 100%);
  pointer-events: none;
}

/* Book Info */
.book-info {
  background: transparent;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.book-title {
  font-size: 18px;
  font-weight: 600;
  line-height: 1.4;
  letter-spacing: 0.3px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  transition: color 0.3s ease;
}

.book-card:hover .book-title {
  color: #e0e0e0 !important;
}

.book-author {
  font-size: 14px;
  font-weight: 400;
  font-style: italic;
  letter-spacing: 0.2px;
}

.book-meta {
  display: flex;
  align-items: center;
  font-size: 12px;
  margin-top: 8px;
}

.meta-text {
  font-size: 12px;
  opacity: 0.7;
}

/* Action Button */
.action-btn {
  text-transform: none;
  font-weight: 500;
  letter-spacing: 0.5px;
  transition: all 0.3s ease;
  border-width: 1px;
}

.add-btn:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: #fff;
}

.remove-btn:hover {
  background: rgba(76, 175, 80, 0.1);
  border-color: #66bb6a;
}

/* Responsive Typography */
@media (max-width: 600px) {
  /* Telefonda satır başına 2 kart var; sütun ~165px genişliğinde.
     Kapak ve tipografi buna göre küçülüyor. */
  .book-cover-wrapper {
    height: 190px;
  }

  .book-title {
    font-size: 14px;
    line-height: 1.3;
  }

  .book-author {
    font-size: 12px;
  }

  .book-meta {
    font-size: 11px;
    margin-top: 4px;
  }

  .read-status-badge {
    top: 6px;
    right: 6px;
    padding: 3px 8px;
  }

  .badge-text {
    font-size: 9px;
  }

  /* "Okunanlara Ekle" ~165px sütuna sığmıyordu */
  .btn-label-long { display: none; }
  .btn-label-short { display: inline; }

  .action-btn {
    font-size: 12px;
    min-height: 32px;
  }

  .btn-icon {
    margin-right: 4px !important;
  }
}
</style>
