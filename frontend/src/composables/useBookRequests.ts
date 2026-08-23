import { ref } from 'vue'
import { bookRequestService } from '../services/api'
import type { BookRequest } from '../types'
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

  return { requests, loading, actionLoading, fetchMyRequests, createRequest, cancelRequest }
}
