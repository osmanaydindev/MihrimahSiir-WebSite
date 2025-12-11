import axios, { AxiosError, type AxiosInstance } from 'axios'
import type { ApiErrorResponse } from '../../types'

// Create axios instance with default config
export const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:8080',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  },
  withCredentials: true // For cookie-based auth
})

// Request interceptor
apiClient.interceptors.request.use(
  (config) => {
    // Add any custom headers here if needed
    // For example, add auth token if using token-based auth
    return config
  },
  (error: AxiosError) => {
    if (import.meta.env.DEV) {
      console.error('[API Request Error]:', error)
    }
    return Promise.reject(error)
  }
)

// Response interceptor
apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  (error: AxiosError<ApiErrorResponse>) => {
    if (import.meta.env.DEV) {
      console.error('[API Response Error]:', {
        status: error.response?.status,
        message: error.response?.data?.message || error.message,
        url: error.config?.url
      })
    }

    // Handle specific error codes globally
    if (error.response?.status === 401) {
      // Unauthorized - redirect to login if needed
      // You can emit an event or dispatch a store action here
      console.warn('Unauthorized access - consider redirecting to login')
    }

    if (error.response?.status === 403) {
      // Forbidden
      console.warn('Forbidden access')
    }

    if (error.response?.status === 404) {
      // Not found
      console.warn('Resource not found')
    }

    if (error.response?.status >= 500) {
      // Server error
      console.error('Server error occurred')
    }

    return Promise.reject(error)
  }
)
