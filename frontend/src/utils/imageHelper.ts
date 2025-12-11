// Helper function to get full image URL
export const getImageUrl = (imagePath: string | null | undefined): string => {
  if (!imagePath) {
    return '/uploads/default-profile.svg'
  }

  // If it's already a full URL, return it
  if (imagePath.startsWith('http://') || imagePath.startsWith('https://')) {
    return imagePath
  }

  // Get the API base URL from environment
  const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:8080'

  // Remove leading slash if present to avoid double slashes
  const path = imagePath.startsWith('/') ? imagePath.slice(1) : imagePath

  return `${baseURL}/${path}`
}
