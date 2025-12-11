import { apiClient } from './client'
import type { Book } from '../../types'

export const bookService = {
  // Get all books
  getAll: () =>
    apiClient.get<Book[]>('/get-books'),

  // Get book by slug
  getBySlug: (slug: string) =>
    apiClient.get<Book>(`/get-book/${slug}`),

  // Get read book IDs for a user
  getReadBookIds: (userId: number) =>
    apiClient.get<number[]>(`/get-reads-books-ids/${userId}`),

  // Add book to read list
  addToReadList: (userId: number, book: Book) =>
    apiClient.post(`/add-book-to-reads/${userId}`, book),

  // Remove book from read list
  removeFromReadList: (userId: number, book: Book) =>
    apiClient.post(`/delete-book-from-reads/${userId}`, book),
}
