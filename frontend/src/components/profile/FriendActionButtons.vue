<script setup lang="ts">
type FriendshipStatus = 'none' | 'pending_sent' | 'pending_received' | 'accepted'

interface Props {
  status: FriendshipStatus
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false
})

const emit = defineEmits<{
  'send-request': []
  'accept-request': []
  'remove-friend': []
}>()
</script>

<template>
  <div class="friend-actions">
    <!-- No friendship: Add Friend Button -->
    <v-btn
      v-if="status === 'none'"
      @click="emit('send-request')"
      :loading="loading"
      color="primary"
      variant="flat"
      prepend-icon="mdi-account-plus"
      class="friend-btn"
    >
      Arkadaş Ekle
    </v-btn>

    <!-- Request sent: Show pending status -->
    <v-btn
      v-if="status === 'pending_sent'"
      disabled
      color="grey"
      variant="outlined"
      prepend-icon="mdi-clock-outline"
      class="friend-btn"
    >
      İstek Gönderildi
    </v-btn>

    <!-- Request received: Accept button -->
    <v-btn
      v-if="status === 'pending_received'"
      @click="emit('accept-request')"
      :loading="loading"
      color="success"
      variant="flat"
      prepend-icon="mdi-account-check"
      class="friend-btn"
    >
      İsteği Kabul Et
    </v-btn>

    <!-- Already friends: Remove friend button -->
    <v-btn
      v-if="status === 'accepted'"
      @click="emit('remove-friend')"
      :loading="loading"
      color="error"
      variant="outlined"
      prepend-icon="mdi-account-remove"
      class="friend-btn"
    >
      Arkadaşlıktan Çıkar
    </v-btn>
  </div>
</template>

<style scoped>
.friend-actions {
  width: 100%;
  max-width: 300px;
}

.friend-btn {
  width: 100%;
  font-weight: 600;
  text-transform: none;
  letter-spacing: 0.5px;
  transition: all 0.3s ease;
}

.friend-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

@media (max-width: 768px) {
  .friend-actions {
    max-width: 100%;
  }
}
</style>
