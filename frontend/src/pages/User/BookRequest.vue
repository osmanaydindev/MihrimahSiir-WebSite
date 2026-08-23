<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import PageHeader from '../../components/common/PageHeader.vue'
import LoadingState from '../../components/common/LoadingState.vue'
import EmptyState from '../../components/common/EmptyState.vue'
import SnackbarNotification from '../../components/common/SnackbarNotification.vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import { useBookRequests } from '../../composables/useBookRequests'
import { useNotification } from '../../composables/useNotification'
import { useConfirmDialog } from '../../composables/useConfirmDialog'
import { getImageUrl } from '../../utils/imageHelper'
import type { BookRequest, BookRequestStatus } from '../../types'

const { requests, loading, actionLoading, fetchMyRequests, createRequest, cancelRequest } = useBookRequests()
const { snackbar } = useNotification()
const { dialog, show: showConfirm } = useConfirmDialog()

const newIsbn = ref('')
const newNote = ref('')

const pendingRequests = computed(() => requests.value.filter(r => r.status === 'pending'))

const statusMeta = (status: BookRequestStatus) => {
  switch (status) {
    case 'approved':
      return { text: 'Onaylandı', color: 'success', icon: 'mdi-check-circle' }
    case 'rejected':
      return { text: 'Reddedildi', color: 'error', icon: 'mdi-close-circle' }
    default:
      return { text: 'İnceleniyor', color: 'warning', icon: 'mdi-clock-outline' }
  }
}

const formatDate = (value: string) => {
  if (!value) return ''
  return new Date(value).toLocaleDateString('tr-TR', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const submitRequest = async () => {
  const ok = await createRequest(newIsbn.value, newNote.value)
  if (ok) {
    newIsbn.value = ''
    newNote.value = ''
    await fetchMyRequests()
  }
}

const confirmCancel = (request: BookRequest) => {
  showConfirm(
    'İsteği iptal et',
    `"${request.fetched_title || request.isbn}" isteğini iptal etmek istediğine emin misin?`,
    async () => {
      const ok = await cancelRequest(request.id)
      if (ok) await fetchMyRequests()
    }
  )
}

onMounted(fetchMyRequests)
</script>

<template>
  <v-container class="book-request-page py-8">
    <PageHeader
      icon="mdi-book-plus-outline"
      title="Kitap İste"
      subtitle="Eklenmesini istediğin kitabın ISBN veya EAN numarasını gir, yönetici değerlendirsin."
      :total="pendingRequests.length"
      total-label="bekleyen istek"
    />

    <!-- YENİ İSTEK -->
    <v-card class="request-form mb-8 pa-4">
      <v-text-field
        v-model="newIsbn"
        label="ISBN / EAN"
        placeholder="Örn: 978-975-363-802-9"
        prepend-inner-icon="mdi-barcode-scan"
        variant="outlined"
        hint="Kitabın arka kapağındaki 13 haneli numara (10 haneli eski ISBN de olur)."
        persistent-hint
        :disabled="actionLoading"
        class="mb-4"
        @keyup.enter="submitRequest"
      />
      <v-textarea
        v-model="newNote"
        label="Not (opsiyonel)"
        placeholder="Bu kitabı neden istediğini yazabilirsin"
        variant="outlined"
        rows="2"
        auto-grow
        :disabled="actionLoading"
        class="mb-2"
      />
      <div class="d-flex justify-end">
        <v-btn
          color="primary"
          variant="flat"
          prepend-icon="mdi-send"
          :loading="actionLoading"
          style="text-transform: none"
          @click="submitRequest"
        >
          İstek Gönder
        </v-btn>
      </div>
    </v-card>

    <!-- İSTEKLERİM -->
    <h2 class="section-title text-h6 text-grey-lighten-1 mb-4">İsteklerim</h2>

    <LoadingState v-if="loading" message="İstekler yükleniyor..." />

    <EmptyState
      v-else-if="requests.length === 0"
      icon="mdi-book-search-outline"
      icon-color="grey-lighten-1"
      title="Henüz kitap isteğin yok"
      message="Yukarıdan bir ISBN girerek ilk isteğini oluşturabilirsin."
    />

    <v-card
      v-for="request in requests"
      v-else
      :key="request.id"
      class="request-card mb-4"
    >
      <div class="d-flex pa-4 ga-4">
        <v-img
          v-if="request.fetched_cover_url"
          :src="getImageUrl(request.fetched_cover_url)"
          width="70"
          max-width="70"
          height="100"
          cover
          class="flex-grow-0 rounded"
        />
        <v-icon v-else size="48" color="grey-darken-1" class="flex-grow-0 mt-4">mdi-book-outline</v-icon>

        <div class="flex-grow-1">
          <div class="d-flex align-center justify-space-between flex-wrap ga-2">
            <span class="text-subtitle-1 font-weight-medium text-white">
              {{ request.fetched_title || 'Bilgi bulunamadı' }}
            </span>
            <v-chip :color="statusMeta(request.status).color" size="small" variant="flat">
              <v-icon start size="16">{{ statusMeta(request.status).icon }}</v-icon>
              {{ statusMeta(request.status).text }}
            </v-chip>
          </div>

          <div class="text-body-2 text-grey-lighten-1 mt-1">
            <span v-if="request.fetched_authors">{{ request.fetched_authors }} · </span>
            <span>ISBN: {{ request.isbn }}</span>
            <span v-if="request.fetched_pages"> · {{ request.fetched_pages }} sayfa</span>
          </div>

          <div class="text-caption text-grey-darken-1 mt-1">
            {{ formatDate(request.created_at) }}
          </div>

          <v-alert
            v-if="!request.metadata_found && request.status === 'pending'"
            type="info"
            variant="tonal"
            density="compact"
            class="mt-3 text-body-2"
          >
            Bu ISBN için otomatik bilgi bulunamadı. Yönetici kitap bilgilerini elle girecek.
          </v-alert>

          <v-alert
            v-if="request.status === 'rejected' && request.admin_note"
            type="error"
            variant="tonal"
            density="compact"
            class="mt-3 text-body-2"
          >
            Gerekçe: {{ request.admin_note }}
          </v-alert>

          <div class="mt-3 d-flex ga-2">
            <v-btn
              v-if="request.status === 'approved' && request.created_book?.slug"
              color="success"
              variant="tonal"
              size="small"
              prepend-icon="mdi-book-open-variant"
              style="text-transform: none"
              :to="{ name: 'BookDetail', params: { slug: request.created_book.slug } }"
            >
              Kitaba Git
            </v-btn>
            <v-btn
              v-if="request.status === 'pending'"
              color="warning"
              variant="tonal"
              size="small"
              prepend-icon="mdi-cancel"
              style="text-transform: none"
              :disabled="actionLoading"
              @click="confirmCancel(request)"
            >
              İsteği İptal Et
            </v-btn>
          </div>
        </div>
      </div>
    </v-card>

    <SnackbarNotification
      v-model="snackbar.show"
      :message="snackbar.message"
      :color="snackbar.color"
    />
    <ConfirmDialog
      v-model="dialog.show"
      :title="dialog.title"
      :message="dialog.message"
      @confirm="dialog.onConfirm"
    />
  </v-container>
</template>

<style scoped>
.book-request-page {
  width: 100%;
  max-width: 900px;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  margin-right: auto;
  margin-left: auto;
}

.request-form {
  width: 100%;
  background: linear-gradient(135deg, #1e1e1e 0%, #2a2a2a 100%);
  border: 2px solid #424242;
  border-radius: 8px;
}

.section-title {
  width: 100%;
}

.request-card {
  width: 100%;
  background: linear-gradient(135deg, #1e1e1e 0%, #2a2a2a 100%);
  border: 2px solid #424242;
  border-radius: 8px;
}

@media (max-width: 600px) {
  .book-request-page {
    padding-right: 12px !important;
    padding-left: 12px !important;
  }
}
</style>
