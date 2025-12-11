import { ref, computed, type Ref, type ComputedRef } from 'vue'

/**
 * Generic search composable for filtering arrays of objects
 */
export function useSearch<T>(
  items: Ref<T[]>,
  searchFields: (item: T) => string[]
): {
  searchQuery: Ref<string>
  filteredItems: ComputedRef<T[]>
} {
  const searchQuery = ref('')

  const filteredItems = computed(() => {
    if (!searchQuery.value.trim()) {
      return items.value
    }

    const query = searchQuery.value.toLowerCase()
    return items.value.filter(item => {
      const fields = searchFields(item)
      return fields.some(field =>
        field.toLowerCase().includes(query)
      )
    })
  })

  return {
    searchQuery,
    filteredItems
  }
}
