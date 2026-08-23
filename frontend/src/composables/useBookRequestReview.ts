import { ref } from 'vue'
import { bookRequestService } from '../services/api'
import type { ApproveBookRequestPayload } from '../services/api/bookRequestService'
import type { AdminBasicInfo, BookRequest, BookRequestStatus } from '../types'
import { useNotification } from './useNotification'
import { useErrorHandler } from './useErrorHandler'

/**
 * Admin tarafı kitap talebi inceleme işlemleri.
 */
export function useBookRequestReview() {
  const { success } = useNotification()
  const { handleError } = useErrorHandler()

  const requests = ref<BookRequest[]>([])
  const users = ref<AdminBasicInfo[]>([])
  const total = ref(0)
  const pendingCount = ref(0)
  const loading = ref(false)
  const actionLoading = ref(false)

  const fetchRequests = async (status: BookRequestStatus | '' = 'pending') => {
    loading.value = true
    try {
      const response = await bookRequestService.getAll(status, 0, 100)
      requests.value = response.data?.data ?? []
      total.value = response.data?.total ?? 0
    } catch (error) {
      handleError(error, 'Kitap istekleri yüklenemedi')
    } finally {
      loading.value = false
    }
  }

  const fetchPendingCount = async () => {
    try {
      const response = await bookRequestService.getPendingCount()
      pendingCount.value = response.data?.pending ?? 0
    } catch (error) {
      // Rozet sayısı kritik değil, sessizce geç
      pendingCount.value = 0
    }
  }

  const fetchUsers = async () => {
    try {
      const response = await bookRequestService.getUsers()
      users.value = Array.isArray(response.data) ? response.data : []
    } catch (error) {
      handleError(error, 'Kullanıcı listesi yüklenemedi')
    }
  }

  const approveRequest = async (requestId: number, payload: ApproveBookRequestPayload) => {
    actionLoading.value = true
    try {
      const response = await bookRequestService.approve(requestId, payload)
      success(response.data.message || 'Kitap isteği onaylandı')
      return true
    } catch (error) {
      handleError(error, 'Kitap isteği onaylanamadı')
      return false
    } finally {
      actionLoading.value = false
    }
  }

  const rejectRequest = async (requestId: number, adminNote: string) => {
    actionLoading.value = true
    try {
      const response = await bookRequestService.reject(requestId, adminNote)
      success(response.data.message || 'Kitap isteği reddedildi')
      return true
    } catch (error) {
      handleError(error, 'Kitap isteği reddedilemedi')
      return false
    } finally {
      actionLoading.value = false
    }
  }

  const refreshMetadata = async (requestId: number) => {
    actionLoading.value = true
    try {
      const response = await bookRequestService.refresh(requestId)
      success('Open Library verisi güncellendi')
      return response.data.request
    } catch (error) {
      handleError(error, 'Open Library verisi güncellenemedi')
      return null
    } finally {
      actionLoading.value = false
    }
  }

  return {
    requests,
    users,
    total,
    pendingCount,
    loading,
    actionLoading,
    fetchRequests,
    fetchPendingCount,
    fetchUsers,
    approveRequest,
    rejectRequest,
    refreshMetadata,
  }
}
