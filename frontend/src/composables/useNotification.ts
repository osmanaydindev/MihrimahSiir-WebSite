import { ref } from 'vue'
import type { SnackbarColor, SnackbarState } from '../types'

// Global snackbar state (shared across all components)
const snackbar = ref<SnackbarState>({
  show: false,
  message: '',
  color: 'success'
})

export function useNotification() {
  const show = (message: string, color: SnackbarColor = 'success') => {
    snackbar.value = {
      show: true,
      message,
      color
    }
  }

  const hide = () => {
    snackbar.value.show = false
  }

  const success = (message: string) => show(message, 'success')
  const error = (message: string) => show(message, 'error')
  const warning = (message: string) => show(message, 'warning')
  const info = (message: string) => show(message, 'info')

  return {
    snackbar,
    show,
    hide,
    success,
    error,
    warning,
    info
  }
}
