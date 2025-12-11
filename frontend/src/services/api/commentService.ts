import { apiClient } from './client'
import type { Comment } from '../../types'

export const commentService = {
  // Add comment to a book
  add: (data: { book_id: string; admin_id: string; content: string; page?: string }) =>
    apiClient.post<Comment[]>('/add-comment', data),

  // Update comment
  update: (commentId: number, data: { content: string; page?: string }) =>
    apiClient.put<{ message: string }>(`/update-comment/${commentId}`, data),

  // Delete comment
  delete: (commentId: number) =>
    apiClient.delete<{ message: string }>(`/delete-comment/${commentId}`),
}
