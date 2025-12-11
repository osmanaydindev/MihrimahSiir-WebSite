import { defineStore } from 'pinia'
import type { User, Book, Poem, Comment } from '../types'

export const useAppStore = defineStore('app', {
  state: () => {
    return {
      user: {} as User,
      likedPoemIds: new Set<number>() as Set<number>,
      bookmarkedPoemIds: new Set<number>() as Set<number>,
      readBookIds: new Set<number>() as Set<number>,
      allBooks: [] as Book[],
      unreadBooks: [] as Book[],
    }
  },
  actions: {
    setUser(user: User) {
      this.user = user
    },

    // Set liked poem IDs from array of IDs
    setLikedPoemIds(poemIds: number[]) {
      this.likedPoemIds = new Set(poemIds)
    },

    // Set liked poem IDs from array of poem objects (for lazy loading pages)
    setLikedPoemIdsFromPoems(poems: Poem[]) {
      this.likedPoemIds = new Set(poems.map(p => p.id))
    },

    // Add a poem ID to liked poems
    addLikedPoemId(poemId: number) {
      this.likedPoemIds.add(poemId)
    },

    // Remove a poem ID from liked poems
    removeLikedPoemId(poemId: number) {
      this.likedPoemIds.delete(poemId)
    },

    // Check if a poem is liked
    isLikedPoem(poemId: number) {
      return this.likedPoemIds.has(poemId)
    },

    // Set bookmarked poem IDs from array of IDs
    setBookmarkedPoemIds(poemIds: number[]) {
      this.bookmarkedPoemIds = new Set(poemIds)
    },

    // Set bookmarked poem IDs from array of poem objects (for lazy loading pages)
    setBookmarkedPoemIdsFromPoems(poems: Poem[]) {
      this.bookmarkedPoemIds = new Set(poems.map(p => p.id))
    },

    // Add a poem ID to bookmarked poems
    addBookmarkedPoemId(poemId: number) {
      this.bookmarkedPoemIds.add(poemId)
    },

    // Remove a poem ID from bookmarked poems
    removeBookmarkedPoemId(poemId: number) {
      this.bookmarkedPoemIds.delete(poemId)
    },

    // Check if a poem is bookmarked
    isBookmarkedPoem(poemId: number) {
      return this.bookmarkedPoemIds.has(poemId)
    },

    setAllBooks(books: Book[]) {
      this.allBooks = books
    },

    // Set read book IDs from array of IDs
    setReadBookIds(bookIds: number[]) {
      this.readBookIds = new Set(bookIds)
    },

    // Set read book IDs from array of book objects (for lazy loading pages)
    setReadBookIdsFromBooks(books: Book[]) {
      this.readBookIds = new Set(books.map(b => b.id))
    },

    // Add a book ID to read books
    addReadBookId(bookId: number) {
      this.readBookIds.add(bookId)
    },

    // Remove a book ID from read books
    removeReadBookId(bookId: number) {
      this.readBookIds.delete(bookId)
    },

    // Check if a book is read
    isReadBook(bookId: number) {
      return this.readBookIds.has(bookId)
    },

  },
  getters: {
    getUser: (state) => state.user,
    getUserId: (state) => state.user.id,
    getLikedPoemIds: (state) => state.likedPoemIds,
    getBookmarkedPoemIds: (state) => state.bookmarkedPoemIds,
    getUserRole: (state) => state.user.role_id,
    getReadBookIds: (state) => state.readBookIds,
    getUnreadBooks: (state) => state.unreadBooks,
    getAllBooks: (state) => state.allBooks
  }
})
