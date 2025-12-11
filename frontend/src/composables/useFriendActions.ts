import { ref } from 'vue'
import { friendService } from '../services/api'
import { useNotification } from './useNotification'
import { useErrorHandler } from './useErrorHandler'

export function useFriendActions() {
  const { success } = useNotification()
  const { handleError } = useErrorHandler()
  const loading = ref(false)

  const sendFriendRequest = async (username: string) => {
    loading.value = true
    try {
      const response = await friendService.sendRequest(username)
      success(response.data.message || 'Arkadaşlık isteği gönderildi')
      return true
    } catch (error) {
      handleError(error, 'Arkadaşlık isteği gönderilemedi')
      return false
    } finally {
      loading.value = false
    }
  }

  const acceptFriendRequest = async (requestId: number) => {
    loading.value = true
    try {
      const response = await friendService.acceptRequest(requestId)
      success(response.data.message || 'Arkadaşlık isteği kabul edildi')
      return true
    } catch (error) {
      handleError(error, 'İstek kabul edilemedi')
      return false
    } finally {
      loading.value = false
    }
  }

  const rejectFriendRequest = async (requestId: number) => {
    loading.value = true
    try {
      const response = await friendService.rejectRequest(requestId)
      success(response.data.message || 'Arkadaşlık isteği reddedildi')
      return true
    } catch (error) {
      handleError(error, 'İstek reddedilemedi')
      return false
    } finally {
      loading.value = false
    }
  }

  const cancelFriendRequest = async (requestId: number) => {
    loading.value = true
    try {
      const response = await friendService.cancelRequest(requestId)
      success(response.data.message || 'Arkadaşlık isteği iptal edildi')
      return true
    } catch (error) {
      handleError(error, 'İstek iptal edilemedi')
      return false
    } finally {
      loading.value = false
    }
  }

  const removeFriend = async (friendshipId: number) => {
    loading.value = true
    try {
      const response = await friendService.removeFriend(friendshipId)
      success(response.data.message || 'Arkadaş çıkarıldı')
      return true
    } catch (error) {
      handleError(error, 'Arkadaş çıkarılamadı')
      return false
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    sendFriendRequest,
    acceptFriendRequest,
    rejectFriendRequest,
    cancelFriendRequest,
    removeFriend
  }
}
