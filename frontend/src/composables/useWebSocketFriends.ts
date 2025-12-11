import { onMounted, onUnmounted } from 'vue'
import { wsService } from '../services/websocket'

/**
 * Composable for managing WebSocket events related to friends
 */
export function useWebSocketFriends(callbacks: {
  onFriendRequestReceived?: () => void
  onFriendRequestUpdate?: () => void
  onFriendRequestAccepted?: () => void
  onFriendRemoved?: () => void
}) {
  const handleFriendRequestReceived = () => {
    console.log('[WebSocket] Friend request received')
    callbacks.onFriendRequestReceived?.()
  }

  const handleFriendRequestUpdate = () => {
    console.log('[WebSocket] Friend request update')
    callbacks.onFriendRequestUpdate?.()
  }

  const handleFriendRequestAccepted = () => {
    console.log('[WebSocket] Friend request accepted')
    callbacks.onFriendRequestAccepted?.()
  }

  const handleFriendRemoved = () => {
    console.log('[WebSocket] Friend removed')
    callbacks.onFriendRemoved?.()
  }

  onMounted(() => {
    wsService.on('friend_request_received', handleFriendRequestReceived)
    wsService.on('friend_request_update', handleFriendRequestUpdate)
    wsService.on('friend_request_accepted', handleFriendRequestAccepted)
    wsService.on('friend_removed', handleFriendRemoved)
  })

  onUnmounted(() => {
    wsService.off('friend_request_received', handleFriendRequestReceived)
    wsService.off('friend_request_update', handleFriendRequestUpdate)
    wsService.off('friend_request_accepted', handleFriendRequestAccepted)
    wsService.off('friend_removed', handleFriendRemoved)
  })
}
