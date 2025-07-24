<template>
  <div style="min-height: 100vh; display: flex; flex-direction: column;">
    <!-- Header -->
    <header style="background: rgba(255, 255, 255, 0.9); backdrop-filter: blur(10px); box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1); border-bottom: 1px solid rgba(0, 0, 0, 0.1); position: sticky; top: 0; z-index: 50;">
      <div style="max-width: 1280px; margin: 0 auto; padding: 0 1rem;">
        <div style="display: flex; justify-content: space-between; align-items: center; height: 64px;">
          <!-- Logo -->
          <div style="display: flex; align-items: center; gap: 12px;">
            <div style="width: 32px; height: 32px; background: linear-gradient(135deg, #3b82f6, #8b5cf6); border-radius: 8px; display: flex; align-items: center; justify-content: center;">
              <span style="color: white; font-size: 14px;">💬</span>
            </div>
            <h1 style="font-size: 20px; font-weight: bold; background: linear-gradient(135deg, #3b82f6, #8b5cf6); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text;">
              KT Chat
            </h1>
          </div>

          <!-- Navigation -->
          <nav style="display: flex; align-items: center; gap: 4px;">
            <NuxtLink 
              to="/" 
              style="color: #4b5563; padding: 8px 16px; border-radius: 8px; font-size: 14px; font-weight: 500; text-decoration: none; transition: all 0.2s;"
              active-style="color: #3b82f6; background: rgba(59, 130, 246, 0.1);"
            >
              💬 Chat
            </NuxtLink>
            <NuxtLink 
              to="/rooms" 
              style="color: #4b5563; padding: 8px 16px; border-radius: 8px; font-size: 14px; font-weight: 500; text-decoration: none; transition: all 0.2s;"
              active-style="color: #3b82f6; background: rgba(59, 130, 246, 0.1);"
            >
              📋 Rooms
            </NuxtLink>
            <NuxtLink 
              to="/profile" 
              style="color: #4b5563; padding: 8px 16px; border-radius: 8px; font-size: 14px; font-weight: 500; text-decoration: none; transition: all 0.2s;"
              active-style="color: #3b82f6; background: rgba(59, 130, 246, 0.1);"
            >
              👤 Profile
            </NuxtLink>
          </nav>

          <!-- User Menu -->
          <div style="display: flex; align-items: center; gap: 16px;">
            <button
              v-if="!isAuthenticated"
              @click="login"
              style="background: transparent; border: 1px solid #d1d5db; padding: 8px 16px; border-radius: 6px; font-size: 14px; cursor: pointer; transition: all 0.2s;"
            >
              Login
            </button>
            <div v-else style="display: flex; align-items: center; gap: 8px;">
              <div style="width: 32px; height: 32px; background: linear-gradient(135deg, #3b82f6, #8b5cf6); border-radius: 50%; display: flex; align-items: center; justify-content: center; color: white; font-size: 12px;">
                👤
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main style="flex: 1;">
      <slot />
    </main>
  </div>
</template>

<script setup>
defineOptions({
  name: 'DefaultLayout'
})

// Auth composable (we'll create this later)
const { isAuthenticated } = useAuth()

// Login function (placeholder)
const login = () => {
  // TODO: Implement login
  console.log('Login clicked')
}
</script> 