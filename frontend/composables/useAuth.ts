interface User {
  id: string
  username: string
  email: string
  avatar?: string
}

interface CognitoError {
  code: string
  message: string
}

interface AuthResult {
  success: boolean
  user?: User
  error?: string
  errorCode?: string
}

export const useAuth = () => {
  // User state
  const user = useState<User | null>('user', () => null)
  const token = useState<string | null>('token', () => null)
  const isInitializing = useState<boolean>('auth_initializing', () => true)
  
  // Initialize auth state from localStorage
  const initializeAuth = () => {
    if (import.meta.client) {
      console.log('Initializing auth state from localStorage...')
      isInitializing.value = true
      
      try {
        const storedToken = localStorage.getItem('ktchat_token')
        const storedUser = localStorage.getItem('ktchat_user')
        
        console.log('Stored token exists:', !!storedToken)
        console.log('Stored user exists:', !!storedUser)
        
        if (storedToken && storedUser) {
          // Basic token validation - check if it's a valid JWT format
          try {
            const parts = storedToken.split('.')
            if (parts.length === 3 && parts[1]) {
              // Decode the payload to check expiration
              const payload = JSON.parse(atob(parts[1]))
              const now = Math.floor(Date.now() / 1000)
              
              console.log('Token expiration:', payload.exp)
              console.log('Current time:', now)
              
              // If token is expired, remove it
              if (payload.exp && payload.exp < now) {
                console.log('Token expired, removing from localStorage')
                localStorage.removeItem('ktchat_token')
                localStorage.removeItem('ktchat_user')
                isInitializing.value = false
                return
              }
              
              // Set the state
              console.log('Setting auth state from localStorage')
              token.value = storedToken
              user.value = JSON.parse(storedUser)
            }
          } catch (error) {
            console.log('Token validation failed:', error)
            // If token is malformed, remove it
            localStorage.removeItem('ktchat_token')
            localStorage.removeItem('ktchat_user')
          }
        }
      } catch (error) {
        console.error('Error initializing auth:', error)
      }
      
      // Mark initialization as complete
      isInitializing.value = false
    } else {
      console.log('Not on client side, skipping auth initialization')
      isInitializing.value = false
    }
  }
  
  const isAuthenticated = computed(() => !!token.value)

  // API base URL
  const config = useRuntimeConfig()
  const apiBaseUrl = config.public.apiBaseUrl

  // Login function
  const login = async (username: string, password: string): Promise<AuthResult> => {
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
      let errorMessage = 'Login failed'
      let errorCode = 'UNKNOWN_ERROR'
      
      if (error && typeof error === 'object' && 'data' in error) {
        const errorData = error.data as { error?: CognitoError }
        if (errorData?.error) {
          errorMessage = errorData.error.message
          errorCode = errorData.error.code
        }
      }
      
      return { 
        success: false, 
        error: errorMessage,
        errorCode: errorCode
      }
    }
    
    return { success: false, error: 'Login failed' }
  }

  // Register function
  const register = async (username: string, email: string, password: string): Promise<AuthResult> => {
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
      let errorMessage = 'Registration failed'
      let errorCode = 'UNKNOWN_ERROR'
      
      if (error && typeof error === 'object' && 'data' in error) {
        const errorData = error.data as { error?: CognitoError }
        if (errorData?.error) {
          errorMessage = errorData.error.message
          errorCode = errorData.error.code
        }
      }
      
      return { 
        success: false, 
        error: errorMessage,
        errorCode: errorCode
      }
    }
  }

  // Logout function
  const logout = () => {
    console.log('Logging out user...')
    
    // Clear state
    token.value = null
    user.value = null
    
    // Clear localStorage
    if (typeof window !== 'undefined') {
      localStorage.removeItem('ktchat_token')
      localStorage.removeItem('ktchat_user')
      console.log('Cleared localStorage')
    }
    
    console.log('Logout completed')
  }

  // Initialize auth state from localStorage
  const initAuth = () => {
    initializeAuth()
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

  // Debug function to check auth state
  const debugAuthState = () => {
    if (import.meta.client) {
      console.log('Auth State Debug:')
      console.log('- Token in state:', !!token.value)
      console.log('- User in state:', !!user.value)
      console.log('- Token in localStorage:', !!localStorage.getItem('ktchat_token'))
      console.log('- User in localStorage:', !!localStorage.getItem('ktchat_user'))
      console.log('- isAuthenticated:', isAuthenticated.value)
      console.log('- Client side:', import.meta.client)
      console.log('- Token value:', token.value ? token.value.substring(0, 20) + '...' : 'null')
    }
  }

  // Initialize auth state immediately when composable is created
  if (import.meta.client) {
    // Initialize immediately
    initializeAuth()
  }

  return {
    user: readonly(user),
    token: readonly(token),
    isAuthenticated,
    isInitializing: readonly(isInitializing),
    login,
    register,
    logout,
    initAuth,
    refreshToken,
    getAuthHeaders,
    debugAuthState,
    initializeAuth
  }
} 