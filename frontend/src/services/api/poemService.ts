import { apiClient } from './client'
import type { Poem, PaginatedPoemsResponse, PoemResponse } from '../../types'

export const poemService = {
  // Get all poems with pagination
  getAll: (page: number = 1) =>
    apiClient.get<PaginatedPoemsResponse>(`/get-poems`, { params: { page } }),

  // Get poems with search
  search: (search: string, page: number = 1) =>
    apiClient.get<PaginatedPoemsResponse>(`/get-search-poems`, { params: { search, page } }),

  // Get single poem by slug
  getBySlug: (slug: string) =>
    apiClient.get<PoemResponse>(`/get-poem/${slug}`),

  // Get popular poems
  getPopular: (page: number = 1) =>
    apiClient.get<PaginatedPoemsResponse>(`/get-popular-poems`, { params: { page } }),

  // Get latest poems
  getLatest: (page: number = 1) =>
    apiClient.get<PaginatedPoemsResponse>(`/get-latest-poems`, { params: { page } }),

  // Get liked poems IDs for a user
  getLikedPoemIds: (userId: number) =>
    apiClient.get<number[]>(`/get-liked-poems-id/${userId}`),

  // Get liked poems paginated
  getLikedPoemsPaginated: (userId: number, limit: number, offset: number, search?: string) =>
    apiClient.get<PaginatedPoemsResponse>(`/get-liked-poems-paginated/${userId}`, {
      params: { limit, offset, search }
    }),

  // Like a poem
  like: (userId: number, poem: Poem) =>
    apiClient.post(`/add-poem-to-liked/${userId}`, poem),

  // Unlike a poem
  unlike: (userId: number, poem: Poem) =>
    apiClient.post(`/undo-poem-to-liked/${userId}`, poem),

  // Get bookmarked poem IDs for a user
  getBookmarkedPoemIds: (userId: number) =>
    apiClient.get<number[]>(`/get-bookmark-id/${userId}`),

  // Get bookmarked poems paginated
  getBookmarkedPoemsPaginated: (userId: number, limit: number, offset: number, search?: string) =>
    apiClient.get<PaginatedPoemsResponse>(`/get-bookmarks-paginated/${userId}`, {
      params: { limit, offset, search }
    }),

  // Bookmark a poem
  bookmark: (userId: number, poem: Poem) =>
    apiClient.post(`/add-bookmark/${userId}`, poem),

  // Remove bookmark
  removeBookmark: (userId: number, poem: Poem) =>
    apiClient.post(`/undo-bookmark/${userId}`, poem),
}
