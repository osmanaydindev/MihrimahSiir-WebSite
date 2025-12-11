<script setup lang="ts">
import {onMounted, ref} from "vue";
import {useAppStore} from "../../../store/app";
import {useRouter} from "vue-router";
import axios from "axios";
import {useField,useForm} from "vee-validate";
import CommentList from "../../../components/frontend/bookDetail/commentList.vue";
import NoComment from "../../../components/frontend/bookDetail/noComment.vue";
import SnackbarNotification from "../../../components/common/SnackbarNotification.vue"
import { useNotification } from "../../../composables/useNotification";
import '../../../styles/BookDetail.css';


const store = useAppStore()
const router = useRouter()
const { snackbar, error: showError, success } = useNotification()
const props = defineProps({
  slug: {
    type: String,
    required: true
  }
})

const comments = ref([])
const book = ref({})
const showCommentForm = ref(false)
const isLoading = ref(false)
const isSubmitting = ref(false)
const formHeight = ref(0)
const isOpenEditDialog = ref(false)

const toggleCommentForm = () => {
  showCommentForm.value = !showCommentForm.value
  if (showCommentForm.value) {
    formHeight.value = document.getElementById('comment-form').offsetHeight
  }else {
    formHeight.value = 0
  }
}

const { handleSubmit, handleReset } = useForm({
  validationSchema: {
    content (value){
      if(!value) return "Bu alan boş bırakılamaz."
      return true
    },
    page (value) {
      // Page is optional now
      if (value && !(/^[0-9]+$/.test(value))) return "Sadece sayı giriniz.";
      return true
    },
  },
})
const content = useField("content")
const page = useField("page")

const submitCreateComment = handleSubmit(async (values) => {
  values.book_id = book.value.id.toString()
  values.admin_id = store.getUserId.toString()
  try {
    axios.post('/add-comment', values)
      .then(async(response) => {
        comments.value = response.data
        content.value.value = ""
        page.value.value = ""
        toggleCommentForm()
      })
      .catch((error) => {
        showError(`Admin olusturma basarisiz: ${error.response?.data?.message || error.message}`)
      })
  }catch (e) {
  }
});
const handleDeleteComment = (comment) => {
  comments.value = comments.value.filter((element) => element.id !== comment.id);
}
const submitEditComment = handleSubmit(async (values) => {
  toggleCommentForm()
  try {
    // await axios.put('/update-admin/'+selectedAdminId.value, values)
    await axios.put('/update-admin/', values)
      .then(async() => {
        // console.log("update-admin then icinde")
        // console.log(response)
        // admins.value = response.data
        // dialog.value = false

        // selectedAdminId.value = 0
      })
      .catch(() => {
        // alert(`Admin olusturma basarisiz: ${error}`)
      })
  }catch (e) {
  }
});
const getReadsBooks = async () => {
  try {
    const res = await axios.get(`/get-reads-books-ids/${store.getUserId}`)
    store.setReadBookIds(res.data)
  } catch (error) {
    console.error('Failed to fetch read book IDs:', error)
  }
}

const toggleReadStatus = async () => {
  if (!book.value || !book.value.id) return

  try {
    if (store.isReadBook(book.value.id)) {
      await axios.post(`/delete-book-from-reads/${store.getUserId}`, book.value)
      store.removeReadBookId(book.value.id)
    } else {
      await axios.post(`/add-book-to-reads/${store.getUserId}`, book.value)
      store.addReadBookId(book.value.id)
    }
  } catch (error) {
    console.error('Failed to toggle read status:', error)
  }
}

onMounted(async () => {
  await getReadsBooks()
  book.value = store.allBooks.find((item) => item.slug === props.slug)

    try {
      await axios.get(`/get-book/${props.slug}`)
        .then((res) => {
          // Check if book data exists
          if (!res.data || !res.data.id) {
            router.push({ name: 'NotFound' })
            return
          }
          book.value = res.data
          // Handle comments that might be null/undefined
          const bookComments = book?.value?.comments || []
          comments.value = bookComments.slice().sort((a, b) => a.id - b.id)
        }).catch(() => {
        // If API returns error (404, 500, etc), redirect to NotFound
        router.push({ name: 'NotFound' })
      })
    }catch (e) {
      // If exception occurs, redirect to NotFound
      router.push({ name: 'NotFound' })
    }

})
</script>

<template>
  <div class="book-detail-container">
    <div v-if="book" class="book-detail-wrapper">

      <!-- Hero Section -->
      <div class="book-hero">
        <!-- Read Status Badge -->
        <div
          @click="toggleReadStatus"
          class="read-status-badge-detail clickable"
          :class="store.isReadBook(book.id) ? 'read-badge-detail' : 'unread-badge-detail'"
        >
          <v-icon size="16" class="mr-1">
            {{ store.isReadBook(book.id) ? 'mdi-check-circle' : 'mdi-circle-outline' }}
          </v-icon>
          <span>{{ store.isReadBook(book.id) ? 'Okundu' : 'Okunacak' }}</span>
        </div>

        <div class="book-hero-content">
          <!-- Book Cover -->
          <div class="book-cover-large">
            <v-img
              :src="book.image"
              class="book-cover-img"
              height="360"
              cover
            />
          </div>

          <!-- Book Info -->
          <div class="book-info-main">
            <h1 class="book-title-large">{{ book.name }}</h1>
            <p class="book-author-large">{{ book.author_data?.name || book.author || 'Bilinmeyen Yazar' }}</p>

            <div class="book-meta-row">
              <div class="book-meta-item">
                <v-icon size="18">mdi-file-document-outline</v-icon>
                <span>{{ book.page }} sayfa</span>
              </div>
              <div class="book-meta-item">
                <v-icon size="18">mdi-calendar-outline</v-icon>
                <span>{{ book.created_at }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Comments Section -->
      <div class="comments-section">
        <div class="comments-header">
          <h2 class="comments-title">
            <v-icon size="24">mdi-comment-text-multiple-outline</v-icon>
            <span>Yorumlar <span class="comments-count">({{ comments.length }})</span></span>
          </h2>

          <v-btn
            @click="toggleCommentForm"
            :loading="isLoading"
            variant="outlined"
            color="grey-lighten-1"
            class="add-comment-btn"
          >
            <v-icon size="18" class="mr-2">{{ showCommentForm ? 'mdi-close' : 'mdi-plus-circle-outline' }}</v-icon>
            {{ showCommentForm ? 'Vazgeç' : 'Yorum Ekle' }}
          </v-btn>
        </div>

        <!-- Comment Form Dialog -->
        <v-dialog v-model="showCommentForm" max-width="600">
          <v-card class="comment-form-card">
            <v-card-title class="comment-form-title pa-5">
              <v-icon size="20" class="mr-2">mdi-comment-edit-outline</v-icon>
              Yeni Yorum
            </v-card-title>

            <v-card-text class="pa-5">
              <v-form @submit.prevent="submitCreateComment">
                <v-textarea
                  v-model="content.value.value"
                  label="Yorumunuz"
                  :error-messages="content.errorMessage.value"
                  variant="outlined"
                  rows="4"
                  class="mb-3"
                  placeholder="Kitap hakkındaki düşüncelerinizi paylaşın..."
                />

                <v-text-field
                  v-model="page.value.value"
                  :error-messages="page.errorMessage.value"
                  label="Sayfa Numarası"
                  variant="outlined"
                  type="number"
                  placeholder="Örn: 42"
                />
              </v-form>
            </v-card-text>

            <v-card-actions class="pa-5 pt-0">
              <v-spacer />
              <v-btn variant="text" @click="toggleCommentForm">
                Vazgeç
              </v-btn>
              <v-btn
                color="green-lighten-1"
                variant="flat"
                :loading="isSubmitting"
                @click="submitCreateComment"
              >
                <v-icon size="18" class="mr-2">mdi-send</v-icon>
                Gönder
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>

        <!-- Comments List -->
        <div class="comments-list-wrapper">
          <comment-list
            :comments="comments"
            :book="book"
            :form-height="formHeight"
            :show-comment-form="showCommentForm"
            @deleteComment="handleDeleteComment"
          />

          <!-- No Comments -->
          <no-comment :comments="comments" />
        </div>
      </div>

    </div>
  </div>

  <!-- Snackbar for notifications -->
  <SnackbarNotification
    v-model="snackbar.show"
    :message="snackbar.message"
    :color="snackbar.color"
  />
</template>

