import { apiClient } from './client'
import type { Author } from '../../types'

export const authorService = {
  // Get all authors
  getAll: () =>
    apiClient.get<Author[]>('/get-authors'),

  // Get author by slug
  getBySlug: (slug: string) =>
    apiClient.get<Author>(`/get-author/${slug}`),
}
