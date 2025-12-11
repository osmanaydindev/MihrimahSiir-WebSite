import { ref, Ref, onMounted, onUnmounted, computed, ComputedRef } from 'vue'
import axios from 'axios'

interface PaginatedResponse<T> {
  data: T[]
  total: number
  offset: number
  limit: number
  has_more: boolean
}

interface UseLazyLoadOptions {
  limit?: number
  threshold?: number // Scroll threshold in pixels from bottom
  queryParams?: Ref<Record<string, any>> // Additional query parameters (reactive)
}

export function useLazyLoad<T>(
  endpoint: string | ComputedRef<string>,
  options: UseLazyLoadOptions = {}
) {
  const {
    limit = 20,
    threshold = 300,
    queryParams
  } = options

  const items = ref<T[]>([]) as Ref<T[]>
  const loading = ref(false)
  const hasMore = ref(true)
  const offset = ref(0)
  const total = ref(0)
  const error = ref<string | null>(null)

  // Build URL with query parameters
  const buildUrl = () => {
    const baseEndpoint = typeof endpoint === 'string' ? endpoint : endpoint.value
    const params = new URLSearchParams()
    params.append('offset', offset.value.toString())
    params.append('limit', limit.toString())

    // Add additional query parameters if provided
    if (queryParams?.value) {
      Object.entries(queryParams.value).forEach(([key, value]) => {
        if (value !== null && value !== undefined && value !== '') {
          params.append(key, value.toString())
        }
      })
    }

    return `${baseEndpoint}?${params.toString()}`
  }

  const loadMore = async () => {
    if (loading.value || !hasMore.value) {
      return
    }

    loading.value = true
    error.value = null

    try {
      const url = buildUrl()
      const response = await axios.get<PaginatedResponse<T>>(url)

      const { data, total: totalCount, has_more } = response.data

      // Validate that data is an array
      if (!Array.isArray(data)) {
        console.error('Invalid response format - data is not an array:', response.data)
        error.value = 'Invalid response format from server'
        return
      }

      items.value = [...items.value, ...data]
      offset.value += limit
      hasMore.value = has_more
      total.value = totalCount
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Failed to load data'
      console.error('Lazy load error:', err)
    } finally {
      loading.value = false
    }
  }

  const refresh = async () => {
    items.value = []
    offset.value = 0
    hasMore.value = true
    total.value = 0
    error.value = null
    await loadMore()
  }

  const handleScroll = () => {
    const scrollHeight = document.documentElement.scrollHeight
    const scrollTop = document.documentElement.scrollTop
    const clientHeight = document.documentElement.clientHeight

    const distanceFromBottom = scrollHeight - scrollTop - clientHeight

    if (distanceFromBottom < threshold && !loading.value && hasMore.value) {
      loadMore()
    }
  }

  onMounted(() => {
    window.addEventListener('scroll', handleScroll)
    // Load initial data
    loadMore()
  })

  onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll)
  })

  return {
    items,
    loading,
    hasMore,
    total,
    error,
    loadMore,
    refresh
  }
}
