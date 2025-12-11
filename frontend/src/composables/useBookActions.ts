import { ref } from 'vue'
import { bookService } from '../services/api'
import { useAppStore } from '../store/app'
import { useErrorHandler } from './useErrorHandler'
import type { Book } from '../types'

export function useBookActions() {
  const store = useAppStore()
  const { handleError } = useErrorHandler()
  const loading = ref(false)

  const toggleReadStatus = async (book: Book) => {
    if (!book || !book.id) return false

    loading.value = true
    const isCurrentlyRead = store.isReadBook(book.id)

    // Optimistic update
    if (isCurrentlyRead) {
      store.removeReadBookId(book.id)
    } else {
      store.addReadBookId(book.id)
    }

    try {
      if (isCurrentlyRead) {
        await bookService.removeFromReadList(store.getUserId, book)
      } else {
        await bookService.addToReadList(store.getUserId, book)
      }
      return true
    } catch (error) {
      // Rollback on error
      if (isCurrentlyRead) {
        store.addReadBookId(book.id)
      } else {
        store.removeReadBookId(book.id)
      }
      handleError(error, 'Okuma durumu güncellenemedi')
      return false
    } finally {
      loading.value = false
    }
  }

  const markAsRead = async (book: Book) => {
    if (store.isReadBook(book.id)) return true
    return toggleReadStatus(book)
  }

  const markAsUnread = async (book: Book) => {
    if (!store.isReadBook(book.id)) return true
    return toggleReadStatus(book)
  }

  return {
    loading,
    toggleReadStatus,
    markAsRead,
    markAsUnread,
    isRead: (bookId: number) => store.isReadBook(bookId)
  }
}
