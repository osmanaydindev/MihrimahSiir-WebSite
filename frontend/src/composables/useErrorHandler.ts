import { AxiosError } from 'axios'
import type { ApiErrorResponse } from '../types'
import { useNotification } from './useNotification'

export function useErrorHandler() {
  const { error: showError } = useNotification()

  const handleError = (err: unknown, fallbackMessage: string = 'Bir hata oluştu') => {
    let message = fallbackMessage

    if (err instanceof AxiosError) {
      const apiError = err.response?.data as ApiErrorResponse
      message = apiError?.message || fallbackMessage
    } else if (err instanceof Error) {
      message = err.message || fallbackMessage
    }

    showError(message)

    // Log for debugging in development
    if (import.meta.env.DEV) {
      console.error('[Error Handler]:', err)
    }

    return message
  }

  const handleApiError = (
    err: unknown,
    customMessages?: Record<number, string>
  ) => {
    if (err instanceof AxiosError && customMessages) {
      const status = err.response?.status
      if (status && customMessages[status]) {
        showError(customMessages[status])
        return customMessages[status]
      }
    }

    return handleError(err)
  }

  return {
    handleError,
    handleApiError
  }
}
