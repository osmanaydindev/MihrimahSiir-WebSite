import { apiClient } from './client'
import type { Friend, FriendRequest } from '../../types'

export const friendService = {
  // Send friend request
  sendRequest: (username: string) =>
    apiClient.post<{ message: string }>('/send-friend-request', { username }),

  // Get friend requests (received)
  getRequests: () =>
    apiClient.get<FriendRequest[]>('/get-friend-requests'),

  // Get sent requests
  getSentRequests: () =>
    apiClient.get<FriendRequest[]>('/get-sent-requests'),

  // Accept friend request
  acceptRequest: (requestId: number) =>
    apiClient.put<{ message: string }>(`/accept-friend-request/${requestId}`),

  // Reject friend request
  rejectRequest: (requestId: number) =>
    apiClient.delete<{ message: string }>(`/reject-friend-request/${requestId}`),

  // Cancel friend request
  cancelRequest: (requestId: number) =>
    apiClient.delete<{ message: string }>(`/cancel-friend-request/${requestId}`),

  // Remove friend
  removeFriend: (friendshipId: number) =>
    apiClient.delete<{ message: string }>(`/remove-friend/${friendshipId}`),

  // Get friends list
  getFriends: () =>
    apiClient.get<Friend[]>('/get-friends'),
}
