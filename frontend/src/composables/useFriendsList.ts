import { ref, computed, onMounted } from 'vue'
import { friendService } from '../services/api'
import type { Friend, FriendRequest } from '../types'

/**
 * Composable for managing friends lists and requests
 */
export function useFriendsList() {
  const friends = ref<Friend[]>([])
  const friendRequests = ref<FriendRequest[]>([])
  const sentRequests = ref<FriendRequest[]>([])

  const loadingFriends = ref(false)
  const loadingRequests = ref(false)
  const loadingSentRequests = ref(false)

  const fetchFriends = async () => {
    loadingFriends.value = true
    try {
      const response = await friendService.getFriends()
      friends.value = response.data || []
    } catch (error) {
      console.error('Failed to fetch friends:', error)
      friends.value = []
    } finally {
      loadingFriends.value = false
    }
  }

  const fetchFriendRequests = async () => {
    loadingRequests.value = true
    try {
      const response = await friendService.getRequests()
      friendRequests.value = response.data || []
    } catch (error) {
      console.error('Failed to fetch friend requests:', error)
      friendRequests.value = []
    } finally {
      loadingRequests.value = false
    }
  }

  const fetchSentRequests = async () => {
    loadingSentRequests.value = true
    try {
      const response = await friendService.getSentRequests()
      sentRequests.value = response.data || []
    } catch (error) {
      console.error('Failed to fetch sent requests:', error)
      sentRequests.value = []
    } finally {
      loadingSentRequests.value = false
    }
  }

  const fetchAll = async () => {
    await Promise.all([
      fetchFriends(),
      fetchFriendRequests(),
      fetchSentRequests()
    ])
  }

  onMounted(() => {
    fetchAll()
  })

  return {
    friends,
    friendRequests,
    sentRequests,
    loadingFriends,
    loadingRequests,
    loadingSentRequests,
    fetchFriends,
    fetchFriendRequests,
    fetchSentRequests,
    fetchAll
  }
}
