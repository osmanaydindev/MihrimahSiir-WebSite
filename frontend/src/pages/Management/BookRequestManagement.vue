<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import axios from 'axios'
import SnackbarNotification from '../../components/common/SnackbarNotification.vue'
import { useBookRequestReview } from '../../composables/useBookRequestReview'
import { useNotification } from '../../composables/useNotification'
import { getImageUrl } from '../../utils/imageHelper'
import type { BookRequest, BookRequestStatus } from '../../types'

const {
  requests, users, loading, actionLoading,
  fetchRequests, fetchUsers, approveRequest, rejectRequest, refreshMetadata
} = useBookRequestReview()
const { snackbar, error: showError } = useNotification()

const statusFilter = ref<BookRequestStatus | ''>('pending')
const authors = ref<Array<{ id: number; name: string }>>([])

// Onay dialogu
const approveDialog = ref(false)
const rejectDialog = ref(false)
const selected = ref<BookRequest | null>(null)
const rejectNote = ref('')
const showAllUsers = ref(false)

const form = ref({
  name: '',
  author_id: null as number | null,
  page: 0,
  image: '',
  description: '',
  community: 2,
  visible_user_ids: [] as number[]
})

const statusOptions = [
  { title: 'Bekleyenler', value: 'pending' as const },
  { title: 'Onaylananlar', value: 'approved' as const },
  { title: 'Reddedilenler', value: 'rejected' as const },
  { title: 'Tümü', value: '' as const }
]

const communityOptions = [
  { title: 'Özel (Sadece Admin ve Üyeler)', value: 1 },
  { title: 'Herkese Açık', value: 2 },
  { title: 'Sadece Seçili Kullanıcılar', value: 3 }
]

// Admin ve üyeler zaten her kitabı görüyor; seçici varsayılan olarak
// sadece normal kullanıcıları (role_id 3) listeler.
const pickerUsers = computed(() =>
  showAllUsers.value ? users.value : users.value.filter(u => u.role_id === 3)
)

const requesterName = computed(() => selected.value?.user?.username || 'İsteği açan kullanıcı')

const statusMeta = (status: BookRequestStatus) => {
  switch (status) {
    case 'approved': return { text: 'Onaylandı', color: 'success' }
    case 'rejected': return { text: 'Reddedildi', color: 'error' }
    default: return { text: 'Bekliyor', color: 'warning' }
  }
}

const formatDate = (value: string) => {
  if (!value) return ''
  return new Date(value).toLocaleDateString('tr-TR', { year: 'numeric', month: 'long', day: 'numeric' })
}

const openApproveDialog = (request: BookRequest) => {
  selected.value = request
  showAllUsers.value = false
  form.value = {
    name: request.fetched_title || '',
    author_id: null,
    page: request.fetched_pages || 0,
    image: request.fetched_cover_url || '',
    description: request.fetched_description || '',
    community: 2,
    visible_user_ids: []
  }
  approveDialog.value = true
}

const openRejectDialog = (request: BookRequest) => {
  selected.value = request
  rejectNote.value = ''
  rejectDialog.value = true
}

const submitApprove = async () => {
  if (!selected.value) return
  if (!form.value.name.trim()) {
    showError('Kitap adı boş olamaz.')
    return
  }
  if (form.value.community === 3 && form.value.visible_user_ids.length === 0) {
    showError('Görünürlük "Sadece Seçili Kullanıcılar" iken en az bir kullanıcı seçmelisiniz.')
    return
  }

  const ok = await approveRequest(selected.value.id, {
    name: form.value.name.trim(),
    author_id: form.value.author_id,
    page: Number(form.value.page) || 0,
    image: form.value.image.trim(),
    description: form.value.description.trim(),
    community: form.value.community,
    visible_user_ids: form.value.community === 3 ? form.value.visible_user_ids : []
  })

  if (ok) {
    approveDialog.value = false
    await fetchRequests(statusFilter.value)
  }
}

const submitReject = async () => {
  if (!selected.value) return
  const ok = await rejectRequest(selected.value.id, rejectNote.value.trim())
  if (ok) {
    rejectDialog.value = false
    await fetchRequests(statusFilter.value)
  }
}

// Open Library verisi yanlış eşleşmişse yeniden çekmeyi dener
const refreshSelected = async () => {
  if (!selected.value) return
  const updated = await refreshMetadata(selected.value.id)
  if (updated) {
    selected.value = updated
    form.value.name = updated.fetched_title || form.value.name
    form.value.page = updated.fetched_pages || form.value.page
    form.value.image = updated.fetched_cover_url || form.value.image
    form.value.description = updated.fetched_description || form.value.description
  }
}

onMounted(async () => {
  await Promise.all([fetchRequests(statusFilter.value), fetchUsers()])
  try {
    const res = await axios.get('/get-all-authors-dropdown')
    authors.value = res.data || []
  } catch (e) {
    console.error('Yazar listesi alınamadı:', e)
  }
})
</script>

<template>
  <v-card
    title="Kitap İstekleri"
    class="book-request-management mt-8 mx-auto bg-grey-darken-2"
  >
    <template v-slot:text>
      <v-select
        v-model="statusFilter"
        :items="statusOptions"
        item-title="title"
        item-value="value"
        label="Durum"
        variant="solo"
        hide-details
        class="status-filter"
        @update:model-value="fetchRequests(statusFilter)"
      />
    </template>

    <v-card-text>
      <div v-if="loading" class="text-center py-8">
        <v-progress-circular indeterminate color="grey-lighten-1" />
      </div>

      <v-alert
        v-else-if="requests.length === 0"
        type="info"
        variant="tonal"
        class="my-4"
      >
        Bu filtrede kitap isteği yok.
      </v-alert>

      <!-- Table.vue sütunları route'a göre hard-code ettiği ve kapak
           adreslerinin başına /uploads eklediği için burada kendi
           tablomuzu kullanıyoruz. -->
      <v-table v-else class="request-review-table bg-white rounded">
        <thead>
          <tr>
            <th>Kapak</th>
            <th>Kitap</th>
            <th>ISBN</th>
            <th>İsteyen</th>
            <th>Tarih</th>
            <th>Durum</th>
            <th>İşlemler</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="request in requests" :key="request.id">
            <td>
              <v-img
                v-if="request.fetched_cover_url"
                :src="getImageUrl(request.fetched_cover_url)"
                width="36"
                height="50"
                cover
                class="rounded my-1"
              />
              <v-icon v-else color="grey">mdi-book-outline</v-icon>
            </td>
            <td>
              <div class="font-weight-medium">{{ request.fetched_title || '—' }}</div>
              <div class="text-caption text-grey-darken-1">{{ request.fetched_authors }}</div>
              <v-chip v-if="!request.metadata_found" color="warning" size="x-small" variant="tonal" class="mt-1">
                Bilgi bulunamadı
              </v-chip>
            </td>
            <td>{{ request.isbn }}</td>
            <td>{{ request.user?.username || '—' }}</td>
            <td>{{ formatDate(request.created_at) }}</td>
            <td>
              <v-chip :color="statusMeta(request.status).color" size="small" variant="flat">
                {{ statusMeta(request.status).text }}
              </v-chip>
            </td>
            <td>
              <template v-if="request.status === 'pending'">
                <v-btn
                  color="success"
                  variant="text"
                  size="small"
                  :disabled="actionLoading"
                  @click="openApproveDialog(request)"
                >
                  <v-icon>mdi-check</v-icon>
                  <v-tooltip activator="parent" location="top">İncele ve Onayla</v-tooltip>
                </v-btn>
                <v-btn
                  color="error"
                  variant="text"
                  size="small"
                  :disabled="actionLoading"
                  @click="openRejectDialog(request)"
                >
                  <v-icon>mdi-close</v-icon>
                  <v-tooltip activator="parent" location="top">Reddet</v-tooltip>
                </v-btn>
              </template>
              <span v-else class="text-caption text-grey-darken-1">—</span>
            </td>
          </tr>
        </tbody>
      </v-table>
    </v-card-text>
  </v-card>

  <!-- ONAY DIALOGU -->
  <v-dialog v-model="approveDialog" persistent max-width="820">
    <v-card>
      <v-card-title><span class="text-h5">Kitap İsteğini İncele</span></v-card-title>
      <v-card-text>
        <!-- Open Library yanlış kitap eşleştirebiliyor; kaynak veri
             burada salt okunur gösteriliyor ki admin fark edebilsin. -->
        <v-alert type="info" variant="tonal" density="compact" class="mb-4">
          <div class="text-body-2">
            <strong>ISBN:</strong> {{ selected?.isbn }}
            <span v-if="selected?.fetched_publisher"> · <strong>Yayınevi:</strong> {{ selected?.fetched_publisher }}</span>
            <span v-if="selected?.fetched_publish_date"> · <strong>Yayın:</strong> {{ selected?.fetched_publish_date }}</span>
          </div>
          <div v-if="selected?.user_note" class="text-body-2 mt-1">
            <strong>Kullanıcı notu:</strong> {{ selected?.user_note }}
          </div>
        </v-alert>

        <v-row>
          <v-col cols="12" sm="3" class="d-flex flex-column align-center">
            <v-img
              v-if="form.image"
              :src="getImageUrl(form.image)"
              width="140"
              height="200"
              cover
              class="rounded mb-2"
            />
            <v-icon v-else size="64" color="grey">mdi-book-outline</v-icon>
            <v-btn
              variant="text"
              size="small"
              prepend-icon="mdi-refresh"
              style="text-transform: none"
              :loading="actionLoading"
              @click="refreshSelected"
            >
              Yeniden Getir
            </v-btn>
          </v-col>

          <v-col cols="12" sm="9">
            <v-text-field v-model="form.name" label="Kitap Adı" variant="outlined" density="comfortable" />
            <v-autocomplete
              v-model="form.author_id"
              :items="authors"
              item-title="name"
              item-value="id"
              label="Yazar Seçin"
              variant="outlined"
              density="comfortable"
              clearable
              no-data-text="Yazar bulunamadı"
              :hint="selected?.fetched_authors ? `Open Library: ${selected.fetched_authors}` : ''"
              persistent-hint
            />
            <v-text-field v-model="form.page" label="Sayfa Sayısı" type="number" variant="outlined" density="comfortable" class="mt-4" />
            <v-text-field v-model="form.image" label="Kapak URL" variant="outlined" density="comfortable" />
            <v-textarea v-model="form.description" label="Özet" variant="outlined" rows="3" auto-grow density="comfortable" />

            <v-select
              v-model="form.community"
              :items="communityOptions"
              item-title="title"
              item-value="value"
              label="Görünürlük"
              variant="outlined"
              density="comfortable"
            />

            <template v-if="form.community === 3">
              <v-autocomplete
                v-model="form.visible_user_ids"
                :items="pickerUsers"
                item-title="username"
                item-value="id"
                label="Bu kitabı görebilecek kullanıcılar"
                variant="outlined"
                density="comfortable"
                multiple
                chips
                closable-chips
                no-data-text="Kullanıcı bulunamadı"
              />
              <div class="d-flex align-center justify-space-between">
                <v-chip color="primary" size="small" variant="tonal" prepend-icon="mdi-account-check">
                  {{ requesterName }} otomatik eklenir
                </v-chip>
                <v-switch
                  v-model="showAllUsers"
                  label="Admin/üyeleri de göster"
                  density="compact"
                  hide-details
                  color="primary"
                />
              </div>
            </template>
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" style="text-transform: none" @click="approveDialog = false">Vazgeç</v-btn>
        <v-btn color="success" variant="flat" style="text-transform: none" :loading="actionLoading" @click="submitApprove">
          Onayla ve Kitabı Oluştur
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- RED DIALOGU -->
  <v-dialog v-model="rejectDialog" persistent max-width="520">
    <v-card>
      <v-card-title><span class="text-h6">İsteği Reddet</span></v-card-title>
      <v-card-text>
        <p class="text-body-2 mb-4">
          "{{ selected?.fetched_title || selected?.isbn }}" isteği reddedilecek.
          Gerekçe kullanıcıya e-posta ile iletilir.
        </p>
        <v-textarea
          v-model="rejectNote"
          label="Red gerekçesi (opsiyonel)"
          variant="outlined"
          rows="3"
          auto-grow
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" style="text-transform: none" @click="rejectDialog = false">Vazgeç</v-btn>
        <v-btn color="error" variant="flat" style="text-transform: none" :loading="actionLoading" @click="submitReject">
          Reddet
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <SnackbarNotification
    v-model="snackbar.show"
    :message="snackbar.message"
    :color="snackbar.color"
  />
</template>

<style scoped>
.book-request-management {
  width: 100%;
  border-radius: 8px;
  overflow: hidden;
}

.status-filter {
  width: min(100%, 280px);
}

.request-review-table {
  overflow-x: auto;
}

.request-review-table :deep(.v-table__wrapper) {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.request-review-table :deep(table) {
  min-width: 860px;
}

.request-review-table :deep(th),
.request-review-table :deep(td) {
  white-space: nowrap;
}

@media (max-width: 600px) {
  .book-request-management {
    margin-top: 12px !important;
  }

  .request-review-table :deep(table) {
    min-width: 760px;
  }

  .request-review-table :deep(th),
  .request-review-table :deep(td) {
    padding: 0 8px !important;
    font-size: 0.82rem;
  }
}
</style>
