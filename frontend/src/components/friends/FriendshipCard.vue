<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { getImageUrl } from '../../utils/imageHelper'

interface User {
  id: number
  username: string
  email?: string
  profile_image?: string
}

interface FriendData {
  friendship_id: number
  user_id: number
  username: string
  email?: string
  profile_image: string
  since_date: string
}

interface RequestData {
  id: number
  user_id: number
  friend_id: number
  status: string
  created_at: string
  user?: User
  friend?: User
}

type CardType = 'friend' | 'received-request' | 'sent-request'

interface Props {
  type: CardType
  data: FriendData | RequestData
}

const props = defineProps<Props>()
const router = useRouter()

const emit = defineEmits<{
  'remove': [id: number]
  'accept': [id: number]
  'reject': [id: number]
  'cancel': [id: number]
}>()

// Computed properties based on card type
const userData = computed(() => {
  if (props.type === 'friend') {
    return props.data as FriendData
  } else if (props.type === 'received-request') {
    return (props.data as RequestData).user!
  } else {
    return (props.data as RequestData).friend!
  }
})

const username = computed(() => userData.value.username)
const email = computed(() => userData.value.email)
const profileImage = computed(() => getImageUrl(userData.value.profile_image || ''))

const displayDate = computed(() => {
  if (props.type === 'friend') {
    return formatDate((props.data as FriendData).since_date)
  } else {
    return formatDate((props.data as RequestData).created_at)
  }
})

const dateLabel = computed(() => {
  return props.type === 'friend' ? 'Arkadaş' : ''
})

const showEmail = computed(() => props.type !== 'friend')

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('tr-TR', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const navigateToProfile = () => {
  router.push(`/profile/${username.value}`)
}

const handlePrimaryAction = (event: Event) => {
  event.stopPropagation()

  if (props.type === 'friend') {
    emit('remove', (props.data as FriendData).friendship_id)
  } else if (props.type === 'received-request') {
    emit('accept', (props.data as RequestData).id)
  } else {
    emit('cancel', (props.data as RequestData).id)
  }
}

const handleSecondaryAction = (event: Event) => {
  event.stopPropagation()

  if (props.type === 'received-request') {
    emit('reject', (props.data as RequestData).id)
  }
}

// Button configurations based on type
const primaryButton = computed(() => {
  switch (props.type) {
    case 'friend':
      return { text: 'Çıkar', icon: 'mdi-account-remove', color: 'error' }
    case 'received-request':
      return { text: 'Kabul Et', icon: 'mdi-check', color: 'success' }
    case 'sent-request':
      return { text: 'İptal Et', icon: 'mdi-cancel', color: 'warning' }
  }
})

const showSecondaryButton = computed(() => props.type === 'received-request')
const secondaryButton = computed(() => ({
  text: 'Reddet',
  icon: 'mdi-close',
  color: 'error'
}))
</script>

<template>
  <v-card class="friendship-card mb-3" @click="navigateToProfile">
    <div class="card-content">
      <div class="profile-section">
        <v-avatar size="56" class="profile-avatar">
          <v-img v-if="profileImage" :src="profileImage" cover />
          <v-icon v-else size="40" color="grey-lighten-1">mdi-account-circle</v-icon>
        </v-avatar>

        <div class="user-info">
          <div class="user-name">{{ username }}</div>
          <div v-if="showEmail" class="user-email">{{ email }}</div>
          <div class="user-date">
            <v-icon size="14" class="mr-1">
              {{ type === 'friend' ? 'mdi-calendar' : 'mdi-clock-outline' }}
            </v-icon>
            {{ dateLabel }} {{ displayDate }}
          </div>
        </div>

        <div class="action-buttons">
          <v-btn
            :color="primaryButton.color"
            :variant="type === 'received-request' ? 'flat' : 'tonal'"
            size="small"
            @click="handlePrimaryAction"
            class="action-btn"
          >
            <v-icon size="18" class="mr-1">{{ primaryButton.icon }}</v-icon>
            {{ primaryButton.text }}
          </v-btn>

          <v-btn
            v-if="showSecondaryButton"
            :color="secondaryButton.color"
            variant="tonal"
            size="small"
            @click="handleSecondaryAction"
            class="action-btn"
          >
            <v-icon size="18" class="mr-1">{{ secondaryButton.icon }}</v-icon>
            {{ secondaryButton.text }}
          </v-btn>
        </div>
      </div>
    </div>
  </v-card>
</template>

<style scoped>
.friendship-card {
  background: linear-gradient(135deg, #1e1e1e 0%, #2a2a2a 100%);
  border: 2px solid #424242;
  border-radius: 16px;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.friendship-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
  border-color: #757575;
}

.card-content {
  padding: 20px;
}

.profile-section {
  display: flex;
  align-items: center;
  gap: 16px;
  justify-content: space-between;
  flex-wrap: wrap;
}

.profile-avatar {
  border: 3px solid #424242;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.friendship-card:hover .profile-avatar {
  border-color: #757575;
  transform: scale(1.05);
}

.user-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.user-name {
  font-size: 17px;
  font-weight: 700;
  color: #ffffff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color 0.3s ease;
}

.friendship-card:hover .user-name {
  color: #e0e0e0;
}

.user-email {
  font-size: 13px;
  color: #9e9e9e;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-date {
  display: flex;
  align-items: center;
  font-size: 12px;
  color: #bdbdbd;
}

.action-buttons {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
  flex-wrap: wrap;
}

.action-btn {
  transition: all 0.3s ease;
}

.action-btn:hover {
  transform: scale(1.05);
}

@media (max-width: 768px) {
  .card-content {
    padding: 16px;
  }

  .profile-section {
    gap: 12px;
  }

  .profile-avatar {
    width: 48px !important;
    height: 48px !important;
  }

  .user-name {
    font-size: 16px;
  }

  .user-info {
    min-width: 100px;
  }

  .action-buttons {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
