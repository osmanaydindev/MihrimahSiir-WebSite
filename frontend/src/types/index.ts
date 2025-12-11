// User Types
export interface User {
  id: number
  username: string
  email: string
  password?: string
  role_id: number
  profile_image?: string
  is_private?: boolean
  admin_liked_poems?: Poem[]
  admin_bookmark_poems?: Poem[]
  user_books_read?: Book[]
  comments?: Comment[]
}

// Poem Types
export interface Poem {
  id: number
  title: string
  content: string
  author?: string
  slug: string
  is_deleted: boolean
  created_at: string
  like_count?: number
}

export interface PoemResponse {
  poem: Poem
}

export interface PaginatedPoemsResponse {
  poems: Poem[]
  meta: PaginationMeta
}

// Book Types
export interface Book {
  id: number
  name: string
  author: string
  slug: string
  image: string
  page: number
  is_deleted: boolean
  created_at: string
  comments?: Comment[]
}

// Comment Types
export interface Comment {
  id: number
  user_id: number
  admin_id?: number
  book_id: number
  comment: string
  page?: number
  created_at: string
  updated_at: string
}

// Author Types
export interface Author {
  id: number
  name: string
  bio?: string
  birth_year?: number
  death_year?: number
  nationality?: string
  image?: string
  slug: string
  poems?: Poem[]
  books?: Book[]
}

// Friend Types
export interface Friend {
  friendship_id: number
  user_id: number
  username: string
  email: string
  role_id: number
  since_date: string
}

export interface FriendRequest {
  id: number
  user_id: number
  friend_id: number
  status: 'pending' | 'accepted' | 'rejected'
  created_at: string
  user: User
  friend?: User
}

export type FriendshipStatus = 'none' | 'pending_sent' | 'pending_received' | 'accepted'

export interface UserProfileResponse {
  user: User
  is_friend: boolean
  is_own_profile: boolean
  friendship_status: FriendshipStatus
}

// Pagination Types
export interface PaginationMeta {
  total: number
  last_page: number
  current_page?: number
  per_page?: number
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  offset: number
  limit: number
  has_more: boolean
}

// API Response Types
export interface ApiSuccessResponse<T = unknown> {
  data?: T
  message?: string
}

export interface ApiErrorResponse {
  message: string
  code?: string
  details?: unknown
}

// Stat Types
export interface Stat {
  icon: string
  iconColor: string
  number: number | string
  label: string
}

// Homepage Types
export interface HomepageItem {
  id: number
  title?: string
  subtitle?: string
  content?: string
}

// Mihrimah Card Types
export interface MihrimahCard {
  id: number
  title: string
  content: string
}

// Snackbar Types
export type SnackbarColor = 'success' | 'error' | 'warning' | 'info'

export interface SnackbarState {
  show: boolean
  message: string
  color: SnackbarColor
}

// Confirm Dialog Types
export interface ConfirmDialogState {
  show: boolean
  title: string
  message: string
  onConfirm: () => void
}

// WebSocket Types
export interface WebSocketMessage<T = unknown> {
  type: string
  payload: T
}

// Form Types
export interface LoginForm {
  username: string
  password: string
}

export interface RegisterForm {
  username: string
  email: string
  password: string
  password_confirmation: string
}

// Upload Types
export interface UploadResponse {
  profile_image: string
  message?: string
}

// Lazy Load Options
export interface LazyLoadOptions {
  limit?: number
  threshold?: number
  queryParams?: Record<string, string | number>
}
