import { apiClient } from './client'
import type { User, UserProfileResponse, UploadResponse } from '../../types'

export const userService = {
  // Get current user
  getUser: () =>
    apiClient.post<{ user: User }>('/user'),

  // Login
  login: (credentials: { username: string; password: string }) =>
    apiClient.post('/login', credentials),

  // Register
  register: (data: { username: string; email: string; password: string }) =>
    apiClient.post('/register', data),

  // Logout
  logout: () =>
    apiClient.post('/logout'),

  // Get user profile by username
  getProfile: (username: string) =>
    apiClient.get<UserProfileResponse>(`/user-profile/${username}`),

  // Upload profile image
  uploadProfileImage: (file: File) => {
    const formData = new FormData()
    formData.append('profile_image', file)
    return apiClient.post<UploadResponse>('/upload-profile-image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // Update privacy setting
  updatePrivacy: (isPrivate: boolean) =>
    apiClient.put<{ message: string; is_private: boolean }>('/update-privacy', {
      is_private: isPrivate
    }),
}
