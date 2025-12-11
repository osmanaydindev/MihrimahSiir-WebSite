<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '../../store/app'
import axios from 'axios'
import { getImageUrl } from '../../utils/imageHelper'
import { useNotification } from '../../composables/useNotification'
import { useFriendActions } from '../../composables/useFriendActions'
import { useConfirmDialog } from '../../composables/useConfirmDialog'
import LoadingState from '../../components/common/LoadingState.vue'
import ErrorState from '../../components/common/ErrorState.vue'
import ProfileAvatar from '../../components/profile/ProfileAvatar.vue'
import PrivacyToggle from '../../components/profile/PrivacyToggle.vue'
import FriendActionButtons from '../../components/profile/FriendActionButtons.vue'
import ProfileStats from '../../components/profile/ProfileStats.vue'
import ProfileItemsDialog from '../../components/profile/ProfileItemsDialog.vue'
import SnackbarNotification from '../../components/common/SnackbarNotification.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'

const route = useRoute()
const router = useRouter()
const store = useAppStore()

// Composables
const { snackbar, success, error: showError, warning } = useNotification()
const { dialog, show: showConfirmDialog } = useConfirmDialog()
const { sendFriendRequest, acceptFriendRequest, removeFriend } = useFriendActions()

const loading = ref(true)
const user = ref<any>(null)
const error = ref<string | null>(null)
const uploadingImage = ref(false)
const updatingPrivacy = ref(false)
const isFriend = ref(false)
const isOwnProfile = ref(false)
const processingFriendRequest = ref(false)
const friendshipStatus = ref<'none' | 'pending_sent' | 'pending_received' | 'accepted'>('none')
const canViewDetails = ref(false)

// Dialog states
const showLikedPoemsDialog = ref(false)
const showBookmarkedPoemsDialog = ref(false)
const showReadBooksDialog = ref(false)
const showCommentsDialog = ref(false)

const profileImage = computed(() => {
  return getImageUrl(user.value?.profile_image)
})

const profileStats = computed(() => [
  {
    icon: 'mdi-heart',
    iconColor: '#f48fb1',
    number: user.value?.liked_poems_count || 0,
    label: 'Beğenilen Şiir',
    clickable: true
  },
  {
    icon: 'mdi-bookmark',
    iconColor: '#81c784',
    number: user.value?.bookmark_poems_count || 0,
    label: 'Kaydedilen Şiir',
    clickable: isOwnProfile.value // Sadece kendi profilinde
  },
  {
    icon: 'mdi-book-open-variant',
    iconColor: '#ffb74d',
    number: user.value?.read_books_count || 0,
    label: 'Okunan Kitap',
    clickable: true
  },
  {
    icon: 'mdi-comment-text',
    iconColor: '#ba68c8',
    number: user.value?.comments_count || 0,
    label: 'Yorum',
    clickable: isOwnProfile.value // Sadece kendi profilinde
  }
])

const handleStatClick = (statIndex: number) => {
  // Check if we have permission to view details
  if (!canViewDetails.value) {
    warning('Bu içeriği görüntülemek için yetkiniz yok')
    return
  }

  // Bookmarked poems and comments are only available for own profile
  if ((statIndex === 1 || statIndex === 3) && !isOwnProfile.value) {
    warning('Bu özellik sadece kendi profilinizde kullanılabilir')
    return
  }

  switch (statIndex) {
    case 0: // Liked Poems
      showLikedPoemsDialog.value = true
      break
    case 1: // Bookmarked Poems (only own profile)
      showBookmarkedPoemsDialog.value = true
      break
    case 2: // Read Books
      showReadBooksDialog.value = true
      break
    case 3: // Comments (only own profile)
      showCommentsDialog.value = true
      break
  }
}

const fetchUserProfile = async () => {
  loading.value = true
  error.value = null

  try {
    const username = route.params.username as string
    const response = await axios.get(`/user-profile/${username}`)
    user.value = response.data.user
    isFriend.value = response.data.is_friend
    isOwnProfile.value = response.data.is_own_profile
    friendshipStatus.value = response.data.friendship_status || 'none'
    canViewDetails.value = response.data.can_view_details || false
  } catch (err: any) {
    router.replace('/not-found')
  } finally {
    loading.value = false
  }
}

const handleProfileImageUpload = async (file: File) => {
  uploadingImage.value = true

  try {
    const formData = new FormData()
    formData.append('profile_image', file)

    const response = await axios.post('/upload-profile-image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })

    // Update user profile image
    if (user.value) {
      user.value.profile_image = response.data.profile_image
    }

    // Update store user profile image
    if (store.user) {
      store.user.profile_image = response.data.profile_image
    }

    success('Profil fotoğrafı başarıyla güncellendi!')
  } catch (err: any) {
    const errorMessage = err.response?.data?.message || 'Profil fotoğrafı yüklenirken bir hata oluştu'
    showError(errorMessage)
    console.error('Failed to upload profile image:', err)
  } finally {
    uploadingImage.value = false
  }
}

const handleSendFriendRequest = async () => {
  if (!user.value) return

  processingFriendRequest.value = true
  const result = await sendFriendRequest(user.value.username)
  processingFriendRequest.value = false

  if (result) {
    friendshipStatus.value = 'pending_sent'
  }
}

const handleAcceptFriendRequest = async () => {
  if (!user.value) return

  processingFriendRequest.value = true

  try {
    const requestsResponse = await axios.get('/get-friend-requests')
    const requests = requestsResponse.data
    const request = requests.find((r: any) => r.user?.id === user.value.id)

    if (request) {
      const result = await acceptFriendRequest(request.id)
      if (result) {
        await fetchUserProfile()
      }
    }
  } catch (err: any) {
    showError(err.response?.data?.message || 'İstek bulunamadı')
  } finally {
    processingFriendRequest.value = false
  }
}

const handleRemoveFriend = async () => {
  if (!user.value) return

  showConfirmDialog(
    'Arkadaşı Çıkar',
    'Bu kişiyi arkadaşlıktan çıkarmak istediğinizden emin misiniz?',
    async () => {
      processingFriendRequest.value = true

      try {
        const friendsResponse = await axios.get('/get-friends')
        const friends = friendsResponse.data
        const friendship = friends.find((f: any) => f.user_id === user.value.id)

        if (friendship) {
          const result = await removeFriend(friendship.friendship_id)
          if (result) {
            await fetchUserProfile()
          }
        }
      } catch (err: any) {
        showError(err.response?.data?.message || 'Arkadaş bulunamadı')
      } finally {
        processingFriendRequest.value = false
      }
    }
  )
}

const handlePrivacyToggle = async (newValue: boolean) => {
  if (!user.value) return

  updatingPrivacy.value = true

  try {
    await axios.put('/update-privacy', {
      is_private: newValue
    })

    user.value.is_private = newValue

    if (store.user) {
      store.user.is_private = newValue
    }

    success(newValue ? 'Profiliniz artık gizli' : 'Profiliniz artık herkese açık')
  } catch (err: any) {
    const errorMessage = err.response?.data?.message || 'Gizlilik ayarı güncellenirken bir hata oluştu'
    showError(errorMessage)
    console.error('Failed to update privacy:', err)
  } finally {
    updatingPrivacy.value = false
  }
}

onMounted(() => {
  fetchUserProfile()
})
</script>

<template>
  <div class="user-profile-page">
    <!-- Loading State -->
    <LoadingState
      v-if="loading"
      message="Profil yükleniyor..."
      color="grey-lighten-1"
    />

    <!-- Error State -->
    <ErrorState
      v-else-if="error"
      :error="error"
      button-text="Ana Sayfaya Dön"
      button-icon="mdi-home"
      @button-click="router.push('/home')"
    />

    <!-- Profile Content -->
    <div v-else-if="user" class="profile-content">
      <!-- Profile Header -->
      <div class="profile-header">
        <ProfileAvatar
          :image-url="profileImage"
          :size="150"
          :editable="isOwnProfile"
          :loading="uploadingImage"
          @upload="handleProfileImageUpload"
        />

        <div class="profile-info">
          <h1 class="username text-center text-md-left">{{ user.username }}</h1>

          <!-- Privacy Toggle - Only for own profile -->
          <PrivacyToggle
            v-if="isOwnProfile"
            :is-private="user.is_private"
            :loading="updatingPrivacy"
            @update:is-private="handlePrivacyToggle"
            class="mt-4"
          />

          <!-- Friend Actions - Only for other profiles -->
          <FriendActionButtons
            v-if="!isOwnProfile"
            :status="friendshipStatus"
            :loading="processingFriendRequest"
            @send-request="handleSendFriendRequest"
            @accept-request="handleAcceptFriendRequest"
            @remove-friend="handleRemoveFriend"
            class="mt-4"
          />
        </div>
      </div>

      <!-- Profile Stats - Visible if: own profile, friend, or profile is public -->
      <div class="profile-stats-wrapper">
        <ProfileStats
          v-show="isOwnProfile || isFriend || !user.is_private"
          :stats="profileStats"
          @stat-click="handleStatClick"
        />
      </div>
    </div>

    <!-- Dialogs -->
    <ProfileItemsDialog
      v-if="user"
      v-model="showLikedPoemsDialog"
      title="Beğenilen Şiirler"
      :username="user.username"
      type="poems"
      empty-message="Henüz beğenilmiş şiir yok"
    />

    <ProfileItemsDialog
      v-if="user"
      v-model="showBookmarkedPoemsDialog"
      title="Kaydedilen Şiirler"
      :username="user.username"
      type="bookmarked-poems"
      empty-message="Henüz kaydedilmiş şiir yok"
    />

    <ProfileItemsDialog
      v-if="user"
      v-model="showReadBooksDialog"
      title="Okunan Kitaplar"
      :username="user.username"
      type="books"
      empty-message="Henüz okunan kitap yok"
    />

    <ProfileItemsDialog
      v-if="user"
      v-model="showCommentsDialog"
      title="Yapılan Yorumlar"
      :username="user.username"
      type="comments"
      empty-message="Henüz yapılmış yorum yok"
    />

    <!-- Snackbar -->
    <SnackbarNotification
      v-model="snackbar.show"
      :message="snackbar.message"
      :color="snackbar.color"
    />

    <!-- Confirm Dialog -->
    <ConfirmDialog
      v-model="dialog.show"
      :title="dialog.title"
      :message="dialog.message"
      @confirm="dialog.onConfirm"
    />
  </div>
</template>

<style scoped>
.user-profile-page {
  min-height: 100vh;
  padding: clamp(20px, 4vw, 40px);
  animation: fadeIn 0.6s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

/* Profile Content */
.profile-content {
  max-width: 1200px;
  margin: 0 auto;
  animation: slideUp 0.8s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Profile Header */
.profile-header {
  display: flex;
  align-items: center;
  gap: clamp(24px, 5vw, 48px);
  padding: clamp(24px, 5vw, 48px);
  background: linear-gradient(135deg, #1e1e1e 0%, #2a2a2a 100%);
  border: 1px solid #424242;
  border-radius: clamp(12px, 2vw, 16px);
  margin-bottom: clamp(24px, 5vw, 32px);
  flex-wrap: wrap;
  justify-content: center;
  text-align: center;
}

.profile-info {
  flex: 1;
  min-width: 250px;
}

.username {
  font-size: clamp(28px, 5vw, 36px);
  font-weight: 700;
  color: #ffffff;
  margin-bottom: clamp(8px, 1.5vw, 12px);
  letter-spacing: 0.5px;;
}

/* Profile Stats Wrapper */
.profile-stats-wrapper {
  min-height: 160px; /* Consistent height to prevent layout shift */
  transition: opacity 0.3s ease;
}

/* Responsive */
@media (max-width: 768px) {
  .profile-header {
    flex-direction: column;
  }

  .profile-stats-wrapper {
    min-height: 400px; /* Taller on mobile due to single column */
  }
}
</style>
