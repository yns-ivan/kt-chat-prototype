export default defineNuxtRouteMiddleware(async (to) => {
  const { isAuthenticated, isInitializing, debugAuthState, initializeAuth } = useAuth()
  
  // Initialize auth state if we're on the client side
  if (import.meta.client) {
    initializeAuth()
    
    // Wait a bit for initialization to complete
    if (isInitializing.value) {
      await new Promise(resolve => setTimeout(resolve, 100))
    }
  }
  
  // Debug authentication state in development
  if (import.meta.dev) {
    debugAuthState()
  }
  
  // If still initializing, allow the route to render (will show loading state)
  if (isInitializing.value) {
    console.log('Auth middleware: Still initializing, allowing route')
    return
  }
  
  // If user is not authenticated and trying to access a protected route
  if (!isAuthenticated.value && to.path !== '/login') {
    console.log(`Auth middleware: Redirecting to /login from ${to.path}`)
    return navigateTo('/login')
  }
  
  // If user is authenticated and trying to access login page, redirect to home
  if (isAuthenticated.value && to.path === '/login') {
    console.log('Auth middleware: Redirecting to / from /login')
    return navigateTo('/')
  }
  
  console.log(`Auth middleware: Allowing access to ${to.path}`)
}) 