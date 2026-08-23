import { apiClient } from './client'
import type { AdminBasicInfo, Book, BookRequest, BookRequestStatus } from '../../types'

export interface BookRequestListResponse {
  data: BookRequest[]
  total: number
  offset: number
  limit: number
  has_more: boolean
}

export interface ApproveBookRequestPayload {
  name: string
  author?: string
  author_id?: number | null
  image?: string
  description?: string
  page?: number
  community: number
  visible_user_ids?: number[]
}

export const bookRequestService = {
  // Kullanıcı: ISBN ile yeni talep oluştur
  create: (isbn: string, note?: string) =>
    apiClient.post<{ message: string; request: BookRequest }>('/create-book-request', { isbn, note }),

  // Kullanıcı: kendi talepleri
  getMine: () =>
    apiClient.get<BookRequest[]>('/get-my-book-requests'),

  // Kullanıcı: bekleyen kendi talebini iptal et
  cancel: (requestId: number) =>
    apiClient.delete<{ message: string }>(`/cancel-book-request/${requestId}`),

  // Admin: talep listesi
  getAll: (status: BookRequestStatus | '' = '', offset = 0, limit = 20) =>
    apiClient.get<BookRequestListResponse>('/get-book-requests', {
      params: { status, offset, limit }
    }),

  // Admin: bekleyen talep sayısı (panel rozeti)
  getPendingCount: () =>
    apiClient.get<{ pending: number }>('/get-book-request-count'),

  // Admin: onayla ve kitabı oluştur
  approve: (requestId: number, payload: ApproveBookRequestPayload) =>
    apiClient.post<{ message: string; book: Book }>(`/approve-book-request/${requestId}`, payload),

  // Admin: reddet
  reject: (requestId: number, adminNote: string) =>
    apiClient.post<{ message: string }>(`/reject-book-request/${requestId}`, { admin_note: adminNote }),

  // Admin: Open Library verisini yeniden çek
  refresh: (requestId: number) =>
    apiClient.post<{ request: BookRequest }>(`/refresh-book-request/${requestId}`),

  // Admin: görünürlük seçicisi için kullanıcı listesi
  getUsers: () =>
    apiClient.get<AdminBasicInfo[]>('/get-admins-management'),
}
