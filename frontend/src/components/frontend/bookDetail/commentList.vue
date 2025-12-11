<script setup>
import { ref, computed, onMounted } from 'vue';
import axios from "axios";
import {useAppStore} from "../../../store/app";
import { getImageUrl } from "../../../utils/imageHelper";
import SnackbarNotification from "../../common/SnackbarNotification.vue"
import { useNotification } from "../../../composables/useNotification";

const props = defineProps(['comments', 'showCommentForm', 'formHeight', 'book']);
const emit = defineEmits(['deleteComment']);
const store = useAppStore();
const user = computed(() => store.getUser || {});
const { snackbar, error: showError, success } = useNotification();

const deleteDialog = ref(false);
const commentToDelete = ref(null);

const colors = [
  'primary',
  'purple',
  'indigo',
  'deep-purple',
  'teal',
  'green',
  'deep-orange',
]

const getRandomColor = (username) => {
  if (!username) return colors[0]
  // Kullanıcı adına göre sabit bir renk döndür
  const index = username.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0)
  return colors[index % colors.length]
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('tr-TR', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

const openDeleteDialog = (comment) => {
  commentToDelete.value = comment;
  deleteDialog.value = true;
}

const confirmDelete = async () => {
  if (!commentToDelete.value) return;

  try {
    await axios.delete(`/delete-comment/${commentToDelete.value.id}`)
      .then(async(response) => {
        emit('deleteComment', commentToDelete.value);
        deleteDialog.value = false;
        commentToDelete.value = null;
      })
      .catch((error) => {
        showError(`Yorum silme başarısız: ${error.response?.data?.message || error.message}`)
      })
  } catch (e) {
    console.error('Yorum silme hatası:', e)
  }
}

const cancelDelete = () => {
  deleteDialog.value = false;
  commentToDelete.value = null;
}

</script>

<template>
  <v-list
    v-if="comments.length > 0"
    id="comment-list"
    style="z-index:1;"
    :style="{'top':formHeight+5+'px'}"
    class="bg-gray-darken-2 w-100 rounded-xl overflow-y-auto"
    :class="{ moved: showCommentForm }"
  >
    <v-list-item
      v-for="(comment, index) in props.comments"
      :key="index"
      class="rounded-lg overflow-y-auto position-relative"
      elevation="1"
    >




      <div class="d-flex flex-column w-100">
        <div class="d-flex flex-column justify-space-between align-start ga-3">
          <div class="d-flex flex-row align-center">
            <router-link
              :to="`/profile/${comment?.admin?.username}`"
              class="text-decoration-none"
            >
              <v-avatar
                size="40"
                class="me-4 comment-avatar"
              >
                <v-img
                  v-if="comment?.admin?.profile_image"
                  :src="getImageUrl(comment.admin.profile_image)"
                  cover
                />
                <div
                  v-else
                  class="avatar-fallback"
                  :style="{ backgroundColor: getRandomColor(comment?.admin?.username) }"
                >
                  <span class="text-h6 text-white font-weight-medium">
                    {{ comment?.admin?.username?.charAt(0).toUpperCase() }}
                  </span>
                </div>
              </v-avatar>
            </router-link>
            <div class="d-flex flex-column">
              <router-link
                :to="`/profile/${comment?.admin?.username}`"
                class="text-decoration-none"
              >
                <v-list-item-title class="text-subtitle-1 font-weight-bold mb-1 username-link">
                  {{ comment?.admin?.username }}
                </v-list-item-title>
              </router-link>
              <v-list-item-subtitle class="text-caption text-grey">
                {{ formatDate(comment.created_at) }}
              </v-list-item-subtitle>
            </div>
          </div>

        </div>

        <v-card-text class="text-body-1 pa-0 mt-3">
          {{ comment.content }}
        </v-card-text>

        <div class="mt-3 d-flex flex-row ga-2">
          <v-card-text
            size="small"
            color="gray"
            variant="flat"
            class="pa-0 ma-0"
          >
            Sayfa {{ comment.page }}
          </v-card-text>
          <v-chip
            v-if="user.id === comment.admin?.id"
            size="small"
            color="red"
            variant="flat"
            class="px-2"
            @click="openDeleteDialog(comment)"
          >
            Sil
          </v-chip>

        </div>
        <v-divider v-if="index !== comments.length - 1" class="mt-4"></v-divider>
      </div>
    </v-list-item>
  </v-list>

  <!-- Delete Confirmation Dialog -->
  <v-dialog v-model="deleteDialog" max-width="400">
    <v-card class="delete-dialog-card">
      <v-card-title class="text-h6 font-weight-bold pa-5">
        <v-icon class="mr-2" color="warning">mdi-alert-circle-outline</v-icon>
        Yorumu Sil
      </v-card-title>

      <v-card-text class="pa-5">
        Bu yorumu silmek istediğinizden emin misiniz? Bu işlem geri alınamaz.
      </v-card-text>

      <v-card-actions class="pa-5 pt-0">
        <v-spacer />
        <v-btn variant="text" @click="cancelDelete">
          İptal
        </v-btn>
        <v-btn color="red" variant="flat" @click="confirmDelete">
          Sil
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- Snackbar for notifications -->
  <SnackbarNotification
    v-model="snackbar.show"
    :message="snackbar.message"
    :color="snackbar.color"
  />
</template>


<style scoped>
#comment-list {
  transition: all 1s;
}
.v-list-item {
  transition: transform 0.2s, box-shadow 0.2s;
  min-height: auto !important;
}

.v-list-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}

.delete-dialog-card {
  background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
  border: 1px solid #424242;
}

/* Comment Avatar */
.comment-avatar {
  border: 2px solid #424242;
  transition: all 0.3s ease;
  cursor: pointer;
}

.comment-avatar:hover {
  border-color: #9e9e9e;
  transform: scale(1.05);
}

.avatar-fallback {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}

/* Username Link */
.username-link {
  color: #ffffff !important;
  transition: color 0.3s ease;
}

.username-link:hover {
  color: #9e9e9e !important;
}

/* Mobil görünüm için responsive düzenlemeler */
@media (max-width: 600px) {
  .v-list-item {
    padding: 12px !important;
  }

  .v-avatar {
    margin-right: 12px !important;
  }

  .text-subtitle-1 {
    font-size: 0.9rem !important;
  }

  .text-body-1 {
    font-size: 0.875rem !important;
  }
}
</style>
