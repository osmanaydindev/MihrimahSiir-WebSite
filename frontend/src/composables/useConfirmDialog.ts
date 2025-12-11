import { ref } from 'vue'
import type { ConfirmDialogState } from '../types'

// Global confirm dialog state (shared across all components)
const dialog = ref<ConfirmDialogState>({
  show: false,
  title: '',
  message: '',
  onConfirm: () => {}
})

export function useConfirmDialog() {
  const show = (title: string, message: string, onConfirm: () => void) => {
    dialog.value = {
      show: true,
      title,
      message,
      onConfirm
    }
  }

  const hide = () => {
    dialog.value.show = false
  }

  const confirm = async (title: string, message: string): Promise<boolean> => {
    return new Promise((resolve) => {
      show(title, message, () => {
        resolve(true)
        hide()
      })
    })
  }

  return {
    dialog,
    show,
    hide,
    confirm
  }
}
