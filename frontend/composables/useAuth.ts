interface User {
  id: number
  username: string
  email: string
  avatar?: string
}

export const useAuth = () => {
  // User state
  const user = useState<User | null>('user', () => null)
  const token = useState<string | null>('token', () => null)
  const isAuthenticated = computed(() => !!token.value)

  // API base URL
  const config = useRuntimeConfig()
  const apiBaseUrl = config.public.apiBaseUrl

  // Login function
  const login = async (username: string, password: string) => {
    try {
      const response = await $fetch<{ token: string; user: User }>('/api/v1/auth/login', {
        baseURL: apiBaseUrl,
        method: 'POST',
        body: {
          username,
          password
        }
      })

      if (response.token) {
        token.value = response.token
        user.value = response.user
        
        // Store in localStorage for persistence
        if (typeof window !== 'undefined') {
          localStorage.setItem('ktchat_token', response.token)
          localStorage.setItem('ktchat_user', JSON.stringify(response.user))
        }

        return { success: true, user: response.user }
      }
    } catch (error: unknown) {
      console.error('Login error:', error)
      const errorMessage = error && typeof error === 'object' && 'data' in error 
        ? (error.data as any)?.message || 'Login failed'
        : 'Login failed'
      return { 
        success: false, 
        error: errorMessage
      }
    }
  }

  // Register function
  const register = async (username: string, email: string, password: string) => {
    try {
      const response = await $fetch<{ message: string; user: User }>('/api/v1/auth/register', {
        baseURL: apiBaseUrl,
        method: 'POST',
        body: {
          username,
          email,
          password
        }
      })

      // Registration successful, but no token returned
      // User needs to login after registration
      return { success: true, user: response.user }
    } catch (error: unknown) {
      console.error('Registration error:', error)
      const errorMessage = error && typeof error === 'object' && 'data' in error 
        ? (error.data as any)?.message || 'Registration failed'
        : 'Registration failed'
      return { 
        success: false, 
        error: errorMessage
      }
    }
  }

  // Logout function
  const logout = () => {
    token.value = null
    user.value = null
    
    // Clear localStorage
    if (typeof window !== 'undefined') {
      localStorage.removeItem('ktchat_token')
      localStorage.removeItem('ktchat_user')
    }
  }

  // Initialize auth state from localStorage
  const initAuth = () => {
    if (typeof window !== 'undefined') {
      const storedToken = localStorage.getItem('ktchat_token')
      const storedUser = localStorage.getItem('ktchat_user')
      
      if (storedToken && storedUser) {
        token.value = storedToken
        user.value = JSON.parse(storedUser)
      }
    }
  }

  // Refresh token
  const refreshToken = async () => {
    if (!token.value) return false
    
    try {
      const response = await $fetch<{ token: string }>('/api/v1/auth/refresh', {
        baseURL: apiBaseUrl,
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token.value}`
        }
      })

      if (response.token) {
        token.value = response.token
        if (typeof window !== 'undefined') {
          localStorage.setItem('ktchat_token', response.token)
        }
        return true
      }
    } catch (error) {
      console.error('Token refresh failed:', error)
      logout()
      return false
    }
  }

  // Get auth headers for API requests
  const getAuthHeaders = () => {
    return token.value ? { 'Authorization': `Bearer ${token.value}` } : {}
  }

  return {
    user: readonly(user),
    token: readonly(token),
    isAuthenticated,
    login,
    register,
    logout,
    initAuth,
    refreshToken,
    getAuthHeaders
  }
} 