import { apiClient } from './client'
import type { FeedItem, PaginatedResponse } from '../../types'

export const feedService = {
  // Akış: arkadaşların + profili herkese açık kullanıcıların paylaşımları.
  // Zaman penceresi yok — yeni içerik bittiğinde geçmişe doğru sayfalar,
  // böylece akış hiçbir zaman boş görünmez.
  list: (offset = 0, limit = 20) =>
    apiClient.get<PaginatedResponse<FeedItem>>(`/feed?offset=${offset}&limit=${limit}`),

  saved: (offset = 0, limit = 20) =>
    apiClient.get<PaginatedResponse<FeedItem>>(`/saved-comments?offset=${offset}&limit=${limit}`),

  like: (commentId: number) =>
    apiClient.post<{ message: string }>(`/comments/${commentId}/like`),

  unlike: (commentId: number) =>
    apiClient.delete<{ message: string }>(`/comments/${commentId}/like`),

  save: (commentId: number) =>
    apiClient.post<{ message: string }>(`/comments/${commentId}/save`),

  unsave: (commentId: number) =>
    apiClient.delete<{ message: string }>(`/comments/${commentId}/save`),
}
