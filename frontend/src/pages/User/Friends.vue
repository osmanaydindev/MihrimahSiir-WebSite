<script setup lang="ts">
import { ref } from 'vue'
import LoadingState from '../../components/common/LoadingState.vue'
import EmptyState from '../../components/common/EmptyState.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import SnackbarNotification from '../../components/common/SnackbarNotification.vue'
import FriendshipCard from '../../components/friends/FriendshipCard.vue'
import { useFriendsList } from '../../composables/useFriendsList'
import { useFriendActions } from '../../composables/useFriendActions'
import { useNotification } from '../../composables/useNotification'
import { useConfirmDialog } from '../../composables/useConfirmDialog'
import { useWebSocketFriends } from '../../composables/useWebSocketFriends'
import { useSearch } from '../../composables/useSearch'

const tab = ref('friends')
const newFriendUsername = ref('')

// Use composables
const { snackbar, warning } = useNotification()
const { dialog, show: showConfirmDialog } = useConfirmDialog()
const {
  friends,
  friendRequests,
  sentRequests,
  loadingFriends,
  loadingRequests,
  loadingSentRequests,
  fetchFriends,
  fetchFriendRequests,
  fetchSentRequests
} = useFriendsList()

const {
  loading: actionLoading,
  sendFriendRequest: sendRequest,
  acceptFriendRequest,
  rejectFriendRequest,
  cancelFriendRequest,
  removeFriend: removeFriendAction
} = useFriendActions()

// Search functionality
const { searchQuery: friendSearchQuery, filteredItems: filteredFriends } = useSearch(
  friends,
  (friend) => [friend.username, friend.email]
)

const { searchQuery: requestSearchQuery, filteredItems: filteredRequests } = useSearch(
  friendRequests,
  (request) => [request.user.username, request.user.email]
)

const { searchQuery: sentRequestSearchQuery, filteredItems: filteredSentRequests } = useSearch(
  sentRequests,
  (request) => [request.friend.username, request.friend.email]
)

// WebSocket event handlers
useWebSocketFriends({
  onFriendRequestReceived: fetchFriendRequests,
  onFriendRequestUpdate: async () => {
    await Promise.all([fetchFriendRequests(), fetchSentRequests()])
  },
  onFriendRequestAccepted: async () => {
    await Promise.all([fetchSentRequests(), fetchFriends()])
  },
  onFriendRemoved: fetchFriends
})

// Actions
const sendFriendRequest = async () => {
  if (!newFriendUsername.value.trim()) {
    warning('Lütfen kullanıcı adı girin')
    return
  }

  const success = await sendRequest(newFriendUsername.value.trim())
  if (success) {
    newFriendUsername.value = ''
    await fetchSentRequests()
  }
}

const acceptRequest = async (requestId: number) => {
  const success = await acceptFriendRequest(requestId)
  if (success) {
    await Promise.all([fetchFriendRequests(), fetchFriends()])
  }
}

const rejectRequest = async (requestId: number) => {
  const success = await rejectFriendRequest(requestId)
  if (success) {
    await fetchFriendRequests()
  }
}

const cancelRequest = async (requestId: number) => {
  const success = await cancelFriendRequest(requestId)
  if (success) {
    await fetchSentRequests()
  }
}

const removeFriend = async (friendshipId: number) => {
  showConfirmDialog(
    'Arkadaşı Çıkar',
    'Bu arkadaşı çıkarmak istediğinize emin misiniz?',
    async () => {
      const success = await removeFriendAction(friendshipId)
      if (success) {
        await fetchFriends()
      }
    }
  )
}
</script>

<template>
  <div class="friends-container d-flex flex-column">
    <div class="page-header">
      <h1 class="page-title">
        <v-icon size="32" class="mr-3">mdi-account-group</v-icon>
        Arkadaş Yönet
      </h1>
    </div>

    <!-- Add Friend Search -->
    <v-card class="add-friend-search-card mb-6">
      <div class="search-wrapper">
        <v-text-field
          v-model="newFriendUsername"
          placeholder="Arkadaş eklemek için kullanıcı adı girin..."
          prepend-inner-icon="mdi-account-search"
          variant="solo"
          flat
          hide-details
          :disabled="actionLoading"
          @keyup.enter="sendFriendRequest"
          class="search-input"
        >
          <template v-slot:append>
            <v-btn
              color="primary"
              :loading="actionLoading"
              @click="sendFriendRequest"
              icon
              size="small"
            >
              <v-icon>mdi-send</v-icon>
            </v-btn>
          </template>
        </v-text-field>
      </div>
    </v-card>

    <!-- Tabs -->
    <v-tabs v-model="tab" class="mb-6" color="grey-lighten-1">
      <v-tab value="friends">
        <v-icon class="mr-2">mdi-account-multiple</v-icon>
        Arkadaşlarım ({{ friends.length }})
      </v-tab>
      <v-tab value="requests">
        <v-icon class="mr-2">mdi-account-clock</v-icon>
        Gelen İstekler
        <v-badge v-if="friendRequests.length > 0" :content="friendRequests.length" color="error" class="ml-2" />
      </v-tab>
      <v-tab value="sent">
        <v-icon class="mr-2">mdi-send</v-icon>
        Gönderilen İstekler ({{ sentRequests.length }})
      </v-tab>
    </v-tabs>

    <!-- Tab Content -->
    <v-window v-model="tab">
      <!-- Friends List -->
      <v-window-item value="friends">
        <LoadingState v-if="loadingFriends" message="Arkadaşlar yükleniyor..." />

        <EmptyState
          v-else-if="friends.length === 0"
          icon="mdi-account-off"
          icon-color="grey-lighten-1"
          title="Henüz arkadaşınız yok"
          message="Yukarıdaki arama kutusundan yeni arkadaşlar ekleyebilirsiniz"
        />

        <div v-else>
          <v-card class="friends-search-card mb-4">
            <v-text-field
              v-model="friendSearchQuery"
              placeholder="Arkadaşlarınızı ara..."
              prepend-inner-icon="mdi-magnify"
              variant="solo"
              flat
              hide-details
              clearable
              class="search-input"
            />
          </v-card>

          <EmptyState
            v-if="filteredFriends.length === 0"
            icon="mdi-account-search"
            icon-color="grey-lighten-1"
            title="Arkadaş bulunamadı"
            message="Arama kriterlerinize uygun arkadaş bulunamadı"
          />

          <div v-else class="friends-grid">
            <FriendshipCard
              v-for="friend in filteredFriends"
              :key="friend.friendship_id"
              type="friend"
              :data="friend"
              @remove="removeFriend"
            />
          </div>
        </div>
      </v-window-item>

      <!-- Friend Requests (Received) -->
      <v-window-item value="requests">
        <LoadingState v-if="loadingRequests" message="İstekler yükleniyor..." />

        <EmptyState
          v-else-if="friendRequests.length === 0"
          icon="mdi-inbox"
          icon-color="grey-lighten-1"
          title="Gelen arkadaşlık isteği yok"
          message="Gelen arkadaşlık istekleri burada görünecek"
        />

        <div v-else>
          <v-card class="friends-search-card mb-4">
            <v-text-field
              v-model="requestSearchQuery"
              placeholder="İsteklerde ara..."
              prepend-inner-icon="mdi-magnify"
              variant="solo"
              flat
              hide-details
              clearable
              class="search-input"
            />
          </v-card>

          <EmptyState
            v-if="filteredRequests.length === 0"
            icon="mdi-account-search"
            icon-color="grey-lighten-1"
            title="İstek bulunamadı"
            message="Arama kriterlerinize uygun istek bulunamadı"
          />

          <div v-else class="requests-list">
            <FriendshipCard
              v-for="request in filteredRequests"
              :key="request.id"
              type="received-request"
              :data="request"
              @accept="acceptRequest"
              @reject="rejectRequest"
            />
          </div>
        </div>
      </v-window-item>

      <!-- Sent Requests -->
      <v-window-item value="sent">
        <LoadingState v-if="loadingSentRequests" message="Gönderilen istekler yükleniyor..." />

        <EmptyState
          v-else-if="sentRequests.length === 0"
          icon="mdi-send-outline"
          icon-color="grey-lighten-1"
          title="Gönderilen arkadaşlık isteği yok"
          message="Yukarıdaki arama kutusundan yeni arkadaşlık istekleri gönderebilirsiniz"
        />

        <div v-else>
          <v-card class="friends-search-card mb-4">
            <v-text-field
              v-model="sentRequestSearchQuery"
              placeholder="Gönderilen isteklerde ara..."
              prepend-inner-icon="mdi-magnify"
              variant="solo"
              flat
              hide-details
              clearable
              class="search-input"
            />
          </v-card>

          <EmptyState
            v-if="filteredSentRequests.length === 0"
            icon="mdi-account-search"
            icon-color="grey-lighten-1"
            title="İstek bulunamadı"
            message="Arama kriterlerinize uygun gönderilen istek bulunamadı"
          />

          <div v-else class="requests-list">
            <FriendshipCard
              v-for="request in filteredSentRequests"
              :key="request.id"
              type="sent-request"
              :data="request"
              @cancel="cancelRequest"
            />
          </div>
        </div>
      </v-window-item>
    </v-window>

    <!-- Snackbar for notifications -->
    <SnackbarNotification
      v-model="snackbar.show"
      :message="snackbar.message"
      :color="snackbar.color"
    />

    <!-- Confirm dialog -->
    <ConfirmDialog
      v-model="dialog.show"
      :title="dialog.title"
      :message="dialog.message"
      @confirm="dialog.onConfirm"
    />
  </div>
</template>

<style scoped>
.friends-container {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
}

.page-header {
  margin-bottom: 32px;
}

.page-title {
  font-size: 32px;
  font-weight: 700;
  color: #ffffff;
  display: flex;
  align-items: center;
}

.friends-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.requests-list {
  max-width: 600px;
}

.add-friend-search-card {
  background: linear-gradient(135deg, #1e1e1e 0%, #2a2a2a 100%);
  border: 1px solid #424242;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
}

.friends-search-card {
  background: linear-gradient(135deg, #1e1e1e 0%, #2a2a2a 100%);
  border: 1px solid #424242;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  padding: 8px 12px;
}

.search-wrapper {
  padding: 12px 16px;
}

.search-input :deep(.v-field) {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  font-size: 15px;
}

.search-input :deep(.v-field__input) {
  padding: 8px 0;
}

.search-input :deep(.v-field__prepend-inner) {
  color: #9e9e9e;
}

.search-input :deep(input::placeholder) {
  color: #757575;
  opacity: 1;
}

@media (max-width: 768px) {
  .friends-container {
    padding: 24px 12px;
  }

  .page-title {
    font-size: 24px;
  }

  .friends-grid {
    grid-template-columns: 1fr;
  }
}
</style>
