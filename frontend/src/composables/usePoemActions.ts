import { useAppStore } from '@/store/app'
import { poemService } from '../services/api'
import { useErrorHandler } from './useErrorHandler'
import type { Poem } from '../types'

/**
 * Composable for handling poem actions (like, bookmark)
 * Centralizes duplicate logic from PoemCard and Poem components
 */
export function usePoemActions() {
  const store = useAppStore()
  const { handleError } = useErrorHandler()

  /**
   * Add a poem to liked poems
   */
  const likePoem = async (poem: Poem) => {
    store.addLikedPoemId(poem.id)
    try {
      await poemService.like(store.getUserId, poem)
      return true
    } catch (error) {
      handleError(error, 'Şiir beğenilemedi')
      // Rollback on error
      store.removeLikedPoemId(poem.id)
      return false
    }
  }

  /**
   * Remove a poem from liked poems
   */
  const unlikePoem = async (poem: Poem) => {
    store.removeLikedPoemId(poem.id)
    try {
      await poemService.unlike(store.getUserId, poem)
      return true
    } catch (error) {
      handleError(error, 'Beğeni kaldırılamadı')
      // Rollback on error
      store.addLikedPoemId(poem.id)
      return false
    }
  }

  /**
   * Add a poem to bookmarks
   */
  const bookmarkPoem = async (poem: Poem) => {
    store.addBookmarkedPoemId(poem.id)
    try {
      await poemService.bookmark(store.getUserId, poem)
      return true
    } catch (error) {
      handleError(error, 'İşaretlenemedi')
      // Rollback on error
      store.removeBookmarkedPoemId(poem.id)
      return false
    }
  }

  /**
   * Remove a poem from bookmarks
   */
  const unbookmarkPoem = async (poem: Poem) => {
    store.removeBookmarkedPoemId(poem.id)
    try {
      await poemService.removeBookmark(store.getUserId, poem)
      return true
    } catch (error) {
      handleError(error, 'İşaret kaldırılamadı')
      // Rollback on error
      store.addBookmarkedPoemId(poem.id)
      return false
    }
  }

  return {
    likePoem,
    unlikePoem,
    bookmarkPoem,
    unbookmarkPoem
  }
}
