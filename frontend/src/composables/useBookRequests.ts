import { ref } from 'vue'
import { bookRequestService } from '../services/api'
import type { BookRequest, BookRequestPreview } from '../types'
import { useNotification } from './useNotification'
import { useErrorHandler } from './useErrorHandler'

/**
 * Kullanıcı tarafı kitap talebi işlemleri.
 * Aksiyonlar boolean döner; çağıran true ise listeyi tazeler.
 */
export function useBookRequests() {
  const { success, warning } = useNotification()
  const { handleError } = useErrorHandler()

  const requests = ref<BookRequest[]>([])
  const loading = ref(false)
  const actionLoading = ref(false)
  // Talep, önizleme onaylanmadan gönderilmiyor; preview null iken form
  // ISBN adımında, dolu iken onay adımındadır.
  const preview = ref<BookRequestPreview | null>(null)
  const previewLoading = ref(false)

  const fetchMyRequests = async () => {
    loading.value = true
    try {
      const response = await bookRequestService.getMine()
      requests.value = Array.isArray(response.data) ? response.data : []
    } catch (error) {
      handleError(error, 'Kitap istekleri yüklenemedi')
    } finally {
      loading.value = false
    }
  }

  const fetchPreview = async (isbn: string) => {
    const trimmed = isbn.trim()
    if (!trimmed) {
      warning('Lütfen bir ISBN girin')
      return false
    }

    previewLoading.value = true
    try {
      const response = await bookRequestService.preview(trimmed)
      preview.value = response.data.preview
      return true
    } catch (error) {
      preview.value = null
      handleError(error, 'Kitap bilgisi getirilemedi')
      return false
    } finally {
      previewLoading.value = false
    }
  }

  const clearPreview = () => {
    preview.value = null
  }

  // Önizlemede onaylanan ISBN ile talebi açar. Sunucu meta veriyi
  // yeniden çeker; istemciden gelen önizleme verisine güvenilmiyor.
  const createRequest = async (isbn: string, note = '') => {
    const trimmed = isbn.trim()
    if (!trimmed) {
      warning('Lütfen bir ISBN girin')
      return false
    }

    actionLoading.value = true
    try {
      const response = await bookRequestService.create(trimmed, note.trim())
      success(response.data.message || 'Kitap isteğin alındı')
      preview.value = null
      return true
    } catch (error) {
      handleError(error, 'Kitap isteği gönderilemedi')
      return false
    } finally {
      actionLoading.value = false
    }
  }

  const cancelRequest = async (requestId: number) => {
    actionLoading.value = true
    try {
      const response = await bookRequestService.cancel(requestId)
      success(response.data.message || 'İstek iptal edildi')
      return true
    } catch (error) {
      handleError(error, 'İstek iptal edilemedi')
      return false
    } finally {
      actionLoading.value = false
    }
  }

  return {
    requests,
    loading,
    actionLoading,
    preview,
    previewLoading,
    fetchMyRequests,
    fetchPreview,
    clearPreview,
    createRequest,
    cancelRequest
  }
}
