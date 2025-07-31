<template>
  <div v-if="!isInitializing" style="height: 100vh; display: flex; overflow: hidden;">
    <!-- Sidebar - Chat Rooms -->
    <div style="width: 320px; background: rgba(255, 255, 255, 0.9); backdrop-filter: blur(10px); border-right: 1px solid rgba(0, 0, 0, 0.1); display: flex; flex-direction: column; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04); min-height: 0;">
      <!-- Header -->
      <div style="padding: 24px; border-bottom: 1px solid rgba(0, 0, 0, 0.1);">
        <div style="display: flex; align-items: center; justify-content: space-between; margin-bottom: 16px;">
          <div style="display: flex; align-items: center; gap: 12px;">
            <div style="width: 40px; height: 40px; background: linear-gradient(135deg, #3b82f6, #8b5cf6); border-radius: 12px; display: flex; align-items: center; justify-content: center;">
              <span style="color: white; font-size: 18px;">📋</span>
            </div>
            <h2 style="font-size: 18px; font-weight: bold; color: #111827;">Chat Rooms</h2>
          </div>

        </div>
        
        <!-- User Info -->
        <div style="display: flex; align-items: center; justify-content: space-between; padding: 12px; background: rgba(59, 130, 246, 0.1); border-radius: 8px;">
          <div style="display: flex; align-items: center; gap: 8px;">
            <div style="width: 32px; height: 32px; background: linear-gradient(135deg, #3b82f6, #8b5cf6); border-radius: 50%; display: flex; align-items: center; justify-content: center;">
              <span style="color: white; font-size: 12px; font-weight: bold;">
                {{ currentUser?.username?.charAt(0).toUpperCase() || 'U' }}
              </span>
            </div>
            <div style="min-width: 0;">
              <p style="font-size: 14px; font-weight: 600; color: #111827; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
                {{ currentUser?.username || 'User' }}
              </p>
              <p style="font-size: 12px; color: #6b7280; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
                {{ currentUser?.email || 'user@example.com' }}
              </p>
            </div>
          </div>
          <div style="display: flex; gap: 4px;">
            <button
              @click="navigateTo('/profile')"
              style="width: 28px; height: 28px; background: transparent; border: none; border-radius: 4px; display: flex; align-items: center; justify-content: center; cursor: pointer; font-size: 14px; color: #6b7280;"
              title="Profile"
            >
              👤
            </button>
            <button
              @click="handleLogout"
              style="width: 28px; height: 28px; background: transparent; border: none; border-radius: 4px; display: flex; align-items: center; justify-content: center; cursor: pointer; font-size: 14px; color: #6b7280;"
              title="Logout"
            >
              🚪
            </button>
          </div>
        </div>
      </div>

      <!-- Room List -->
      <div style="flex: 1; overflow-y: auto; min-height: 0; scroll-behavior: smooth;">
        <div v-if="loading" style="padding: 16px;">
          <div style="height: 48px; width: 100%; background: #f3f4f6; border-radius: 8px; margin-bottom: 8px;"></div>
          <div style="height: 48px; width: 100%; background: #f3f4f6; border-radius: 8px; margin-bottom: 8px;"></div>
          <div style="height: 48px; width: 100%; background: #f3f4f6; border-radius: 8px;"></div>
        </div>
        
                <div v-else-if="rooms.length === 0" style="padding: 32px; text-align: center; color: #6b7280;">
          <div style="width: 64px; height: 64px; background: linear-gradient(135deg, #dbeafe, #e9d5ff); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px;">
            <span style="font-size: 24px;">💬</span>
          </div>
          <h3 style="font-size: 18px; font-weight: 600; color: #111827; margin-bottom: 8px;">No chat rooms available</h3>
          <p style="margin-bottom: 16px;">Contact an administrator to create chat rooms</p>
 
        </div>

        <div v-else style="padding: 16px; display: flex; flex-direction: column; gap: 8px; min-height: 0;">
          <div
            v-for="room in rooms"
            :key="room.id"
            @click="selectRoom(room)"
            style="padding: 16px; border-radius: 12px; cursor: pointer; transition: all 0.2s; border: 1px solid transparent;"
            :style="{ 
              backgroundColor: selectedRoom?.id === room.id ? 'rgba(59, 130, 246, 0.1)' : 'transparent',
              borderColor: selectedRoom?.id === room.id ? '#3b82f6' : 'transparent',
              boxShadow: selectedRoom?.id === room.id ? '0 4px 6px -1px rgba(0, 0, 0, 0.1)' : 'none'
            }"
          >
            <div style="display: flex; align-items: center; gap: 12px;">
              <div style="width: 40px; height: 40px; background: linear-gradient(135deg, #3b82f6, #8b5cf6); border-radius: 8px; display: flex; align-items: center; justify-content: center;">
                <span style="color: white; font-size: 14px;">💬</span>
              </div>
              <div style="flex: 1; min-width: 0;">
                <p style="font-size: 14px; font-weight: 600; color: #111827; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
                  {{ room.name }}
                </p>
                <p style="font-size: 12px; color: #6b7280; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
                  {{ room.description || 'No description' }}
                </p>
              </div>
              <div
                v-if="room.unreadCount"
                style="background: #ef4444; color: white; font-size: 12px; font-weight: 600; padding: 2px 6px; border-radius: 10px; min-width: 20px; text-align: center;"
              >
                {{ room.unreadCount }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Chat Area -->
    <div style="flex: 1; display: flex; flex-direction: column; min-height: 0; overflow: hidden;">
      <!-- Chat Header -->
      <div v-if="selectedRoom" style="background: white; border-bottom: 1px solid #e5e7eb; padding: 16px;">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ selectedRoom.name }}
            </h3>
            <p class="text-sm text-gray-500 dark:text-gray-400">
              {{ selectedRoom.participants?.length || 0 }} participants
            </p>
          </div>
          <div class="flex items-center space-x-2">
            <UButton
              icon="i-heroicons-magnifying-glass"
              variant="ghost"
              size="sm"
              @click="showSearch = !showSearch"
              title="Search messages"
            />
            <UButton
              icon="i-heroicons-users"
              variant="ghost"
              size="sm"
              @click="toggleActiveUsers"
              title="Active users"
            />
            <UButton
              icon="i-heroicons-cog-6-tooth"
              variant="ghost"
              size="sm"
            />
          </div>
        </div>
        
        <!-- Search Bar -->
        <div v-if="showSearch" style="margin-top: 12px; padding: 12px; background: #f9fafb; border-radius: 8px;">
          <div style="display: flex; gap: 8px; align-items: center;">
            <input
              v-model="searchQuery"
              placeholder="Search messages..."
              style="flex: 1; padding: 8px 12px; border: 1px solid #d1d5db; border-radius: 6px; font-size: 14px;"
              @keyup.enter="performSearch"
            />
            <button
              @click="performSearch"
              style="padding: 8px 12px; background: #3b82f6; color: white; border: none; border-radius: 6px; font-size: 14px; cursor: pointer;"
            >
              Search
            </button>
            <button
              @click="clearSearch"
              style="padding: 8px 12px; background: #6b7280; color: white; border: none; border-radius: 6px; font-size: 14px; cursor: pointer;"
            >
              Clear
            </button>
          </div>
          
          <!-- Search Filters -->
          <div v-if="showSearch" style="margin-top: 8px; display: flex; gap: 8px; flex-wrap: wrap;">
            <input
              v-model="searchFilters.startDate"
              type="date"
              placeholder="Start date"
              style="padding: 6px 8px; border: 1px solid #d1d5db; border-radius: 4px; font-size: 12px;"
            />
            <input
              v-model="searchFilters.endDate"
              type="date"
              placeholder="End date"
              style="padding: 6px 8px; border: 1px solid #d1d5db; border-radius: 4px; font-size: 12px;"
            />
            <select
              v-model="searchFilters.userId"
              style="padding: 6px 8px; border: 1px solid #d1d5db; border-radius: 4px; font-size: 12px;"
            >
              <option value="">All users</option>
              <option v-for="participant in selectedRoom.participants" :key="participant.user_id" :value="participant.user_id">
                {{ participant.user?.username || 'Unknown User' }}
              </option>
            </select>
          </div>
        </div>
        
        <!-- Active Users Panel -->
        <div v-if="showActiveUsers" style="position: absolute; top: 100%; right: 0; width: 280px; background: white; border: 1px solid #e5e7eb; border-radius: 8px; box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1); z-index: 50; margin-top: 4px;">
          <div style="padding: 16px; border-bottom: 1px solid #e5e7eb;">
            <div style="display: flex; justify-content: space-between; align-items: center;">
              <h4 style="font-size: 14px; font-weight: 600; color: #111827;">Active Users</h4>
              <button
                @click="showActiveUsers = false"
                style="background: none; border: none; color: #6b7280; cursor: pointer; font-size: 16px;"
              >
                ×
              </button>
            </div>
          </div>
          
          <div style="max-height: 300px; overflow-y: auto;">
            <div v-if="activeUsersLoading" style="padding: 16px; text-align: center; color: #6b7280;">
              <div style="display: inline-block; width: 16px; height: 16px; border: 2px solid #e5e7eb; border-top: 2px solid #3b82f6; border-radius: 50%; animation: spin 1s linear infinite;"></div>
              <span style="margin-left: 8px; font-size: 14px;">Loading users...</span>
            </div>
            
            <div v-else-if="activeUsers.length === 0" style="padding: 16px; text-align: center; color: #6b7280; font-size: 14px;">
              No active users
            </div>
            
            <div v-else style="padding: 8px;">
              <div
                v-for="user in activeUsers"
                :key="user.id"
                style="display: flex; align-items: center; gap: 12px; padding: 8px; border-radius: 6px; transition: background-color 0.2s;"
                :style="{ backgroundColor: user.id === currentUser?.id ? '#f0f9ff' : 'transparent' }"
              >
                <div style="width: 32px; height: 32px; background: linear-gradient(135deg, #3b82f6, #8b5cf6); border-radius: 50%; display: flex; align-items: center; justify-content: center; color: white; font-size: 12px; font-weight: 600;">
                  {{ user.username?.charAt(0).toUpperCase() || '?' }}
                </div>
                <div style="flex: 1; min-width: 0;">
                  <p style="font-size: 14px; font-weight: 500; color: #111827; margin: 0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
                    {{ user.username }}
                    <span v-if="user.id === currentUser?.id" style="color: #3b82f6; font-size: 12px; margin-left: 4px;">(You)</span>
                  </p>
                  <p style="font-size: 12px; color: #6b7280; margin: 0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
                    {{ user.email }}
                  </p>
                </div>
                <div style="width: 8px; height: 8px; background: #10b981; border-radius: 50%; flex-shrink: 0;"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Welcome Screen -->
      <div v-else style="flex: 1; display: flex; align-items: center; justify-content: center; padding: 32px;">
        <div style="text-align: center; max-width: 448px;">
          <div style="width: 96px; height: 96px; background: linear-gradient(135deg, #dbeafe, #e9d5ff); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 24px;">
            <span style="font-size: 36px;">💬</span>
          </div>
          <h3 style="font-size: 24px; font-weight: bold; color: #111827; margin-bottom: 12px;">
            Welcome to KT Chat
          </h3>
          <p style="color: #ffffff; margin-bottom: 24px; font-size: 18px;">
            Select a chat room to start messaging with your team
          </p>
        </div>
      </div>

      <!-- Messages Area -->
      <div v-if="selectedRoom" style="flex: 1; display: flex; flex-direction: column; min-height: 0;">
        <!-- Messages List -->
        <div ref="messagesContainer" style="flex: 1; overflow-y: auto; padding: 16px; min-height: 0; scroll-behavior: smooth;">
          <div v-if="messagesLoading" class="flex justify-center">
            <UIcon name="i-heroicons-arrow-path" class="w-6 h-6 animate-spin text-gray-400" />
          </div>
          
          <div v-if="searchLoading" class="flex justify-center">
            <div style="display: flex; align-items: center; gap: 8px; padding: 16px; color: #6b7280;">
              <div style="width: 16px; height: 16px; border: 2px solid #e5e7eb; border-top: 2px solid #3b82f6; border-radius: 50%; animation: spin 1s linear infinite;"></div>
              <span style="font-size: 14px;">Searching messages...</span>
            </div>
          </div>
          
          <div v-else-if="messages.length === 0" class="text-center text-gray-500 dark:text-gray-400">
            <p>No messages yet. Start the conversation!</p>
          </div>

          <div v-else style="display: flex; flex-direction: column; gap: 16px;">
            <!-- Search Results -->
            <div v-if="searchResults.length > 0" style="margin-bottom: 16px; padding: 12px; background: #f0f9ff; border: 1px solid #0ea5e9; border-radius: 8px;">
              <h4 style="font-size: 14px; font-weight: 600; color: #0369a1; margin-bottom: 8px;">
                Search Results ({{ searchResults.length }} messages)
              </h4>
              <div style="display: flex; flex-direction: column; gap: 12px;">
                <div
                  v-for="message in searchResults"
                  :key="message.id || message.timestamp"
                  style="display: flex; gap: 12px; padding: 8px; background: white; border-radius: 6px; border: 1px solid #e5e7eb;"
                  :style="{ flexDirection: (message.userId === currentUser?.id || message.user_id === currentUser?.id) ? 'row-reverse' : 'row' }"
                >
                  <div style="width: 24px; height: 24px; background: linear-gradient(135deg, #3b82f6, #8b5cf6); border-radius: 50%; display: flex; align-items: center; justify-content: center; color: white; font-size: 10px;">
                    {{ getMessageUserInitial(message) }}
                  </div>
                  <div
                    style="max-width: 280px; display: flex; flex-direction: column;"
                    :style="{ alignItems: (message.userId === currentUser?.id || message.user_id === currentUser?.id) ? 'flex-end' : 'flex-start' }"
                  >
                    <p 
                      style="font-size: 10px; margin-bottom: 2px; font-weight: 600; text-transform: capitalize; padding: 1px 4px; border-radius: 3px; display: inline-block;"
                      :style="{ 
                        color: (message.userId === currentUser?.id || message.user_id === currentUser?.id) ? '#3b82f6' : '#374151',
                        background: (message.userId === currentUser?.id || message.user_id === currentUser?.id) ? 'rgba(59, 130, 246, 0.1)' : 'rgba(0, 0, 0, 0.05)'
                      }"
                    >
                      {{ (message.userId === currentUser?.id || message.user_id === currentUser?.id) ? 'You' : getMessageUsername(message) }}
                    </p>
                    <div
                      style="border-radius: 6px; padding: 6px 10px; font-size: 13px;"
                      :style="{ 
                        backgroundColor: (message.userId === currentUser?.id || message.user_id === currentUser?.id) ? '#3b82f6' : '#f3f4f6',
                        color: (message.userId === currentUser?.id || message.user_id === currentUser?.id) ? 'white' : '#111827'
                      }"
                    >
                      <p style="margin: 0;" v-html="highlightSearchTerm(message.content, searchQuery)"></p>
                    </div>
                    <p style="font-size: 10px; color: #6b7280; margin-top: 2px;">
                      {{ formatTime(message.createdAt || message.timestamp || message.created_at) }}
                    </p>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- Regular Messages -->
            <div
              v-for="message in allMessages"
              :key="message.id || message.timestamp"
              style="display: flex; gap: 12px;"
              :style="{ flexDirection: (message.userId === currentUser?.id || message.user_id === currentUser?.id) ? 'row-reverse' : 'row' }"
            >
              <div style="width: 32px; height: 32px; background: linear-gradient(135deg, #3b82f6, #8b5cf6); border-radius: 50%; display: flex; align-items: center; justify-content: center; color: white; font-size: 12px;">
                {{ getMessageUserInitial(message) }}
              </div>
              <div
                style="max-width: 350px; display: flex; flex-direction: column;"
                :style="{ alignItems: (message.userId === currentUser?.id || message.user_id === currentUser?.id) ? 'flex-end' : 'flex-start' }"
              >
                <!-- Username display -->
                <p 
                  style="font-size: 11px; margin-bottom: 4px; font-weight: 600; text-transform: capitalize; padding: 2px 6px; border-radius: 4px; display: inline-block;"
                  :style="{ 
                    color: (message.userId === currentUser?.id || message.user_id === currentUser?.id) ? '#3b82f6' : '#374151',
                    background: (message.userId === currentUser?.id || message.user_id === currentUser?.id) ? 'rgba(59, 130, 246, 0.1)' : 'rgba(0, 0, 0, 0.05)'
                  }"
                >
                  {{ (message.userId === currentUser?.id || message.user_id === currentUser?.id) ? 'You' : getMessageUsername(message) }}
                </p>
                <div
                  style="border-radius: 8px; padding: 8px 12px;"
                  :style="{ 
                    backgroundColor: (message.userId === currentUser?.id || message.user_id === currentUser?.id) ? '#3b82f6' : '#f3f4f6',
                    color: (message.userId === currentUser?.id || message.user_id === currentUser?.id) ? 'white' : '#111827'
                  }"
                >
                  <p style="font-size: 14px; margin: 0;">{{ message.content }}</p>
                </div>
                <p style="font-size: 12px; color: #6b7280; margin-top: 4px;">
                  {{ formatTime(message.createdAt || message.timestamp || message.created_at) }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Message Input -->
        <div style="border-top: 1px solid #e5e7eb; padding: 16px;">
          <!-- Connection Status -->
          <div v-if="!isConnected" style="margin-bottom: 8px; padding: 8px; background: #fef3c7; border: 1px solid #f59e0b; border-radius: 4px; font-size: 12px; color: #92400e; display: flex; justify-content: space-between; align-items: center;">
            <span>🔄 Connecting to chat...</span>
            <button 
              @click="reconnectWebSocket" 
              style="background: #f59e0b; color: white; border: none; border-radius: 4px; padding: 4px 8px; font-size: 11px; cursor: pointer;"
            >
              Reconnect
            </button>
          </div>
          
          <form @submit.prevent="sendMessage" style="display: flex; gap: 8px;">
            <input
              v-model="newMessage"
              placeholder="Type your message..."
              style="flex: 1; padding: 8px 12px; border: 1px solid #d1d5db; border-radius: 8px; background: white; color: #111827; font-size: 14px;"
              :disabled="!isConnected"
            />
            <button
              type="submit"
              :disabled="!newMessage.trim() || !isConnected"
              style="padding: 8px 12px; background: #3b82f6; color: white; border: none; border-radius: 8px; font-size: 14px; cursor: pointer; transition: background-color 0.2s;"
            >
              📤
            </button>
          </form>
        </div>
      </div>
    </div>


  </div>
</template>

<style scoped>
@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>

<script setup>
defineOptions({
  name: 'ChatPage'
})

// Apply authentication middleware
definePageMeta({
  middleware: 'auth'
})

// Composables
const config = useRuntimeConfig()
const { isAuthenticated, user: currentUser, getAuthHeaders, isInitializing, logout } = useAuth()
const { 
  isConnected, 
  messages: wsMessages, 
  connect: connectWebSocket, 
  disconnect: disconnectWebSocket, 
  clearMessages: clearWebSocketMessages,
  reconnect: reconnectWebSocket
} = useWebSocket()

// Reactive data
const rooms = ref([])
const selectedRoom = ref(null)
const messages = ref([])
const newMessage = ref('')
const loading = ref(false)
const messagesLoading = ref(false)
const messagesContainer = ref(null)

// Search state
const showSearch = ref(false)
const searchQuery = ref('')
const searchResults = ref([])
const searchLoading = ref(false)
const searchFilters = ref({
  startDate: '',
  endDate: '',
  userId: ''
})

// Active users state
const showActiveUsers = ref(false)
const activeUsers = ref([])
const activeUsersLoading = ref(false)

// Load rooms on mount
onMounted(() => {
  console.log('Current user data:', currentUser.value)
  loadRooms()
})

// Cleanup on unmount
onUnmounted(async () => {
  // Leave the current room if any
  if (selectedRoom.value?.id) {
    try {
      await $fetch(`/api/v1/chat/rooms/${selectedRoom.value.id}/leave`, {
        baseURL: config.public.apiBaseUrl,
        method: 'POST',
        headers: getAuthHeaders()
      })
    } catch (error) {
      console.error('Failed to leave room:', error)
    }
  }
  
  disconnectWebSocket()
})

// Load chat rooms
const loadRooms = async () => {
  if (!isAuthenticated.value) return
  
  loading.value = true
  try {
    const response = await $fetch('/api/v1/chat/rooms', {
      baseURL: config.public.apiBaseUrl,
      headers: getAuthHeaders()
    })
    rooms.value = response.rooms || []
  } catch (error) {
    console.error('Failed to load rooms:', error)
  } finally {
    loading.value = false
  }
}

// Select a room
const selectRoom = async (room) => {
  // Leave previous room if exists
  if (selectedRoom.value?.id && selectedRoom.value.id !== room.id) {
    await leaveRoom(selectedRoom.value.id)
  }
  
  // Disconnect from previous room
  disconnectWebSocket()
  clearWebSocketMessages()
  
  selectedRoom.value = room
  
  // Clear search results when switching rooms
  clearSearch()
  
  // Automatically join the room
  await joinRoom(room.id)
  
  // Load room details with updated participant count
  await loadRoomDetails(room.id)
  
  await loadMessages(room.id)
  
  // Connect to WebSocket for real-time messages
  connectWebSocket(room.id)
}

// Load messages for a room
const loadMessages = async (roomId) => {
  if (!roomId) return
  
  messagesLoading.value = true
  try {
    const response = await $fetch(`/api/v1/chat/rooms/${roomId}/messages`, {
      baseURL: config.public.apiBaseUrl,
      headers: getAuthHeaders()
    })
    console.log('API messages loaded:', response.messages)
    
    // Debug: Log first message to see user data
    if (response.messages && response.messages.length > 0) {
      console.log('First API message:', response.messages[0])
      console.log('First message user:', response.messages[0].user)
    }
    
    messages.value = response.messages || []
    
    // Auto-scroll to bottom after loading messages
    scrollToBottom()
  } catch (error) {
    console.error('Failed to load messages:', error)
  } finally {
    messagesLoading.value = false
  }
}

// Leave a room
const leaveRoom = async (roomId) => {
  if (!roomId) return
  
  try {
    await $fetch(`/api/v1/chat/rooms/${roomId}/leave`, {
      baseURL: config.public.apiBaseUrl,
      method: 'POST',
      headers: getAuthHeaders()
    })
    console.log('Left room:', roomId)
  } catch (error) {
    console.error('Failed to leave room:', error)
  }
}

// Join a room
const joinRoom = async (roomId) => {
  if (!roomId) return
  
  try {
    await $fetch(`/api/v1/chat/rooms/${roomId}/join`, {
      baseURL: config.public.apiBaseUrl,
      method: 'POST',
      headers: getAuthHeaders()
    })
    console.log('Joined room:', roomId)
  } catch (error) {
    console.error('Failed to join room:', error)
  }
}

// Load room details with participant count
const loadRoomDetails = async (roomId) => {
  if (!roomId) return
  
  try {
    const response = await $fetch(`/api/v1/chat/rooms/${roomId}`, {
      baseURL: config.public.apiBaseUrl,
      headers: getAuthHeaders()
    })
    
    // Update the selected room with fresh data
    if (response.room) {
      selectedRoom.value = response.room
      console.log('Room details loaded:', response.room)
      console.log('Participant count:', response.room.participants?.length || 0)
      
      // Update active users if the panel is open
      if (showActiveUsers.value) {
        await loadActiveUsers()
      }
    }
  } catch (error) {
    console.error('Failed to load room details:', error)
  }
}

// Perform search
const performSearch = async () => {
  if (!searchQuery.value.trim() || !selectedRoom.value) return
  
  searchLoading.value = true
  try {
    const response = await $fetch(`/api/v1/chat/rooms/${selectedRoom.value.id}/search`, {
      baseURL: config.public.apiBaseUrl,
      method: 'POST',
      headers: getAuthHeaders(),
      body: {
        query: searchQuery.value.trim(),
        limit: 50,
        offset: 0,
        start_date: searchFilters.value.startDate,
        end_date: searchFilters.value.endDate,
        user_id: searchFilters.value.userId
      }
    })
    
    searchResults.value = response.messages || []
    console.log('Search results:', searchResults.value)
  } catch (error) {
    console.error('Failed to search messages:', error)
    searchResults.value = []
  } finally {
    searchLoading.value = false
  }
}

// Highlight search terms in text
const highlightSearchTerm = (text, searchTerm) => {
  if (!searchTerm || !text) return text
  
  const regex = new RegExp(`(${searchTerm})`, 'gi')
  return text.replace(regex, '<mark style="background-color: #fef08a; padding: 1px 2px; border-radius: 2px;">$1</mark>')
}

// Load active users for current room
const loadActiveUsers = async () => {
  if (!selectedRoom.value?.id) return
  
  activeUsersLoading.value = true
  try {
    const response = await $fetch(`/api/v1/chat/rooms/${selectedRoom.value.id}`, {
      baseURL: config.public.apiBaseUrl,
      headers: getAuthHeaders()
    })
    
    if (response.room?.participants) {
      activeUsers.value = response.room.participants
        .filter(participant => participant.left_at === null)
        .map(participant => participant.user)
        .filter(user => user !== null)
    }
  } catch (error) {
    console.error('Failed to load active users:', error)
    activeUsers.value = []
  } finally {
    activeUsersLoading.value = false
  }
}

// Toggle active users panel
const toggleActiveUsers = async () => {
  showActiveUsers.value = !showActiveUsers.value
  if (showActiveUsers.value) {
    await loadActiveUsers()
  }
}

// Clear search
const clearSearch = () => {
  searchQuery.value = ''
  searchResults.value = []
  searchFilters.value = {
    startDate: '',
    endDate: '',
    userId: ''
  }
  showSearch.value = false
}

// Send a message
const sendMessage = async () => {
  if (!newMessage.value.trim() || !selectedRoom.value) return
  
  const message = newMessage.value.trim()
  newMessage.value = ''
  
  // Save to database via API (which also broadcasts via WebSocket)
  try {
    await $fetch(`/api/v1/chat/rooms/${selectedRoom.value.id}/messages`, {
      baseURL: config.public.apiBaseUrl,
      method: 'POST',
      headers: getAuthHeaders(),
      body: {
        content: message
      }
    })
    
    // Auto-scroll to bottom after sending message
    scrollToBottom()
  } catch (error) {
    console.error('Failed to save message to database:', error)
  }
}

// Computed property to combine API messages with WebSocket messages
const allMessages = computed(() => {
  const apiMessages = messages.value || []
  const wsMessagesList = wsMessages.value || []
  
  // Combine and sort by timestamp
  const combined = [...apiMessages, ...wsMessagesList]
  return combined.sort((a, b) => {
    const timeA = a.timestamp || a.created_at || a.createdAt
    const timeB = b.timestamp || b.created_at || b.createdAt
    
    const dateA = new Date(timeA)
    const dateB = new Date(timeB)
    
    // Handle invalid dates by putting them at the end
    if (isNaN(dateA.getTime()) && isNaN(dateB.getTime())) return 0
    if (isNaN(dateA.getTime())) return 1
    if (isNaN(dateB.getTime())) return -1
    
    return dateA.getTime() - dateB.getTime()
  })
})

// Auto-scroll to bottom of messages
const scrollToBottom = () => {
  if (messagesContainer.value) {
    nextTick(() => {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    })
  }
}

// Watch for WebSocket messages to update participant count and auto-scroll
watch(wsMessages, (newMessages) => {
  if (newMessages && newMessages.length > 0) {
    const latestMessage = newMessages[newMessages.length - 1]
    
    // Debug: Log message data to see username
    console.log('Latest WebSocket message:', latestMessage)
    console.log('Message username:', latestMessage.username)
    console.log('Message user:', latestMessage.user)
    
    // Handle participant count updates
    if (latestMessage.participant_count !== undefined && selectedRoom.value?.id === latestMessage.room_id) {
      // Update the participant count
      selectedRoom.value.participants = Array(latestMessage.participant_count).fill({})
      
      // Refresh active users if panel is open
      if (showActiveUsers.value) {
        loadActiveUsers()
      }
    }
    
    // Handle dedicated participant count messages
    if (latestMessage.type === 'participant_count' && selectedRoom.value?.id === latestMessage.room_id) {
      selectedRoom.value.participants = Array(latestMessage.participant_count).fill({})
      
      // Refresh active users if panel is open
      if (showActiveUsers.value) {
        loadActiveUsers()
      }
    }
    
    // Auto-scroll to bottom when new messages arrive
    scrollToBottom()
  }
}, { deep: true })

// Watch for changes in allMessages to auto-scroll
watch(allMessages, () => {
  scrollToBottom()
}, { deep: true })



// Get full username for message display
const getMessageUsername = (message) => {
  // Try different possible username sources
  const username = message.user?.username || message.username || currentUser.value?.username || 'Unknown User'
  
  // Debug: Log the username resolution
  console.log('Message username resolution:', {
    messageUser: message.user?.username,
    messageUsername: message.username,
    currentUser: currentUser.value?.username,
    final: username,
    messageId: message.id,
    messageType: message.type
  })
  
  // Handle empty or invalid username
  if (!username || username === '?' || username === '') {
    return 'Unknown User'
  }
  
  // Truncate very long usernames
  if (username.length > 20) {
    return username.substring(0, 17) + '...'
  }
  
  return username
}

// Get user initial for message display
const getMessageUserInitial = (message) => {
  // Try different possible username sources
  const username = message.user?.username || message.username || currentUser.value?.username || '?'
  
  // Debug: Log the username resolution
  console.log('Message username resolution:', {
    messageUser: message.user?.username,
    messageUsername: message.username,
    currentUser: currentUser.value?.username,
    final: username,
    messageId: message.id,
    messageType: message.type
  })
  
  // Handle empty or invalid username
  if (!username || username === '?' || username === '') {
    return '?'
  }
  
  return username.charAt(0).toUpperCase()
}

// Format time
const formatTime = (timestamp) => {
  if (!timestamp) return 'Invalid Date'
  
  try {
    // Handle both string timestamps and Date objects
    const date = timestamp instanceof Date ? timestamp : new Date(timestamp)
    
    // Check if the date is valid
    if (isNaN(date.getTime())) {
      console.warn('Invalid timestamp received:', timestamp, typeof timestamp)
      return 'Invalid Date'
    }
    
    return date.toLocaleTimeString([], { 
      hour: '2-digit', 
      minute: '2-digit' 
    })
  } catch (error) {
    console.error('Error formatting timestamp:', error, timestamp)
    return 'Invalid Date'
  }
}

// Handle logout
const handleLogout = async () => {
  try {
    // Disconnect WebSocket if connected
    if (isConnected.value) {
      disconnectWebSocket()
    }
    
    // Clear any pending messages
    clearWebSocketMessages()
    
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
</script> 