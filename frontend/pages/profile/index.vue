<template>
  <div v-if="!isInitializing" class="min-h-screen bg-gradient-to-br from-blue-50 to-purple-50 py-8">
    <div class="max-w-4xl mx-auto px-4">
      <!-- Header -->
      <div class="bg-white rounded-2xl shadow-lg p-8 mb-8">
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center gap-4">
            <div class="w-16 h-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center">
              <span class="text-white text-2xl font-bold">
                {{ user?.username?.charAt(0).toUpperCase() || 'U' }}
              </span>
            </div>
            <div>
              <h1 class="text-3xl font-bold text-gray-900">{{ user?.username || 'User' }}</h1>
              <p class="text-gray-600">{{ user?.email || 'user@example.com' }}</p>
            </div>
          </div>
          <UButton
            icon="i-heroicons-arrow-left"
            variant="ghost"
            @click="navigateTo('/')"
          >
            Back to Chat
          </UButton>
        </div>
      </div>

      <!-- Profile Content -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- User Information -->
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <h2 class="text-xl font-semibold text-gray-900 mb-6">User Information</h2>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Username</label>
              <div class="bg-gray-50 px-4 py-3 rounded-lg">
                {{ user?.username || 'Not set' }}
              </div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Email</label>
              <div class="bg-gray-50 px-4 py-3 rounded-lg">
                {{ user?.email || 'Not set' }}
              </div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">User ID</label>
              <div class="bg-gray-50 px-4 py-3 rounded-lg font-mono text-sm">
                {{ user?.id || 'Not available' }}
              </div>
            </div>
          </div>
        </div>

        <!-- Account Actions -->
        <div class="bg-white rounded-2xl shadow-lg p-8">
          <h2 class="text-xl font-semibold text-gray-900 mb-6">Account Actions</h2>
          
          <div class="space-y-4">
            <UButton
              icon="i-heroicons-key"
              variant="outline"
              class="w-full justify-start"
              @click="handleChangePassword"
            >
              Change Password
            </UButton>
            
            <UButton
              icon="i-heroicons-bell"
              variant="outline"
              class="w-full justify-start"
              @click="handleNotifications"
            >
              Notification Settings
            </UButton>
            
            <UButton
              icon="i-heroicons-shield-check"
              variant="outline"
              class="w-full justify-start"
              @click="handleSecurity"
            >
              Security Settings
            </UButton>
            
            <UButton
              icon="i-heroicons-arrow-right-on-rectangle"
              variant="outline"
              color="red"
              class="w-full justify-start"
              @click="handleLogout"
            >
              Logout
            </UButton>
          </div>
        </div>
      </div>

      <!-- Statistics -->
      <div class="bg-white rounded-2xl shadow-lg p-8 mt-8">
        <h2 class="text-xl font-semibold text-gray-900 mb-6">Chat Statistics</h2>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div class="text-center">
            <div class="text-3xl font-bold text-blue-600 mb-2">0</div>
            <div class="text-gray-600">Total Messages</div>
          </div>
          <div class="text-center">
            <div class="text-3xl font-bold text-purple-600 mb-2">0</div>
            <div class="text-gray-600">Chat Rooms</div>
          </div>
          <div class="text-center">
            <div class="text-3xl font-bold text-green-600 mb-2">0</div>
            <div class="text-gray-600">Files Shared</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineOptions({
  name: 'ProfilePage'
})

// Apply authentication middleware
definePageMeta({
  middleware: 'auth'
})

// Composables
const { user, logout, isInitializing } = useAuth()

// Handle logout
const handleLogout = async () => {
  try {
    // Logout and clear auth state
    logout()
    
    // Navigate to login page
    await navigateTo('/login')
  } catch (error) {
    console.error('Logout error:', error)
    // Force logout even if there's an error
    logout()
    await navigateTo('/login')
  }
}

// Placeholder functions for future implementation
const handleChangePassword = () => {
  // TODO: Implement password change functionality
  console.log('Change password clicked')
}

const handleNotifications = () => {
  // TODO: Implement notification settings
  console.log('Notifications clicked')
}

const handleSecurity = () => {
  // TODO: Implement security settings
  console.log('Security clicked')
}
</script> 