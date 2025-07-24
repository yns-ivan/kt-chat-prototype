<template>
  <div style="height: 100vh; display: flex;">
    <!-- Sidebar - Chat Rooms -->
    <div style="width: 320px; background: rgba(255, 255, 255, 0.9); backdrop-filter: blur(10px); border-right: 1px solid rgba(0, 0, 0, 0.1); display: flex; flex-direction: column; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);">
      <!-- Header -->
      <div style="padding: 24px; border-bottom: 1px solid rgba(0, 0, 0, 0.1);">
        <div style="display: flex; align-items: center; justify-content: space-between;">
          <div style="display: flex; align-items: center; gap: 12px;">
            <div style="width: 40px; height: 40px; background: linear-gradient(135deg, #3b82f6, #8b5cf6); border-radius: 12px; display: flex; align-items: center; justify-content: center;">
              <span style="color: white; font-size: 18px;">📋</span>
            </div>
            <h2 style="font-size: 18px; font-weight: bold; color: #111827;">Chat Rooms</h2>
          </div>
          <button
            @click="showCreateRoom = true"
            style="width: 32px; height: 32px; background: transparent; border: none; border-radius: 6px; display: flex; align-items: center; justify-content: center; cursor: pointer; font-size: 16px;"
          >
            ➕
          </button>
        </div>
      </div>

      <!-- Room List -->
      <div style="flex: 1; overflow-y: auto;">
        <div v-if="loading" style="padding: 16px;">
          <div style="height: 48px; width: 100%; background: #f3f4f6; border-radius: 8px; margin-bottom: 8px;"></div>
          <div style="height: 48px; width: 100%; background: #f3f4f6; border-radius: 8px; margin-bottom: 8px;"></div>
          <div style="height: 48px; width: 100%; background: #f3f4f6; border-radius: 8px;"></div>
        </div>
        
        <div v-else-if="rooms.length === 0" style="padding: 32px; text-align: center; color: #6b7280;">
          <div style="width: 64px; height: 64px; background: linear-gradient(135deg, #dbeafe, #e9d5ff); border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 16px;">
            <span style="font-size: 24px;">💬</span>
          </div>
          <h3 style="font-size: 18px; font-weight: 600; color: #111827; margin-bottom: 8px;">No chat rooms yet</h3>
          <p style="margin-bottom: 16px;">Create your first chat room to start messaging with others</p>
          <button
            @click="showCreateRoom = true"
            style="background: linear-gradient(135deg, #3b82f6, #8b5cf6); color: white; padding: 8px 16px; border-radius: 8px; border: none; font-size: 14px; cursor: pointer; box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1); transition: all 0.2s;"
          >
            ➕ Create your first room
          </button>
        </div>

        <div v-else style="padding: 16px; display: flex; flex-direction: column; gap: 8px;">
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
    <div style="flex: 1; display: flex; flex-direction: column;">
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
              icon="i-heroicons-users"
              variant="ghost"
              size="sm"
              @click="showParticipants = true"
            />
            <UButton
              icon="i-heroicons-cog-6-tooth"
              variant="ghost"
              size="sm"
            />
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
          <p style="color: #6b7280; margin-bottom: 24px; font-size: 18px;">
            Select a chat room to start messaging with your team
          </p>
          <button
            @click="showCreateRoom = true"
            style="background: linear-gradient(135deg, #3b82f6, #8b5cf6); color: white; padding: 12px 24px; border-radius: 8px; border: none; font-size: 16px; font-weight: 600; cursor: pointer; box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1); transition: all 0.2s;"
          >
            ➕ Create New Room
          </button>
        </div>
      </div>

      <!-- Messages Area -->
      <div v-if="selectedRoom" class="flex-1 flex flex-col">
        <!-- Messages List -->
        <div class="flex-1 overflow-y-auto p-4 space-y-4">
          <div v-if="messagesLoading" class="flex justify-center">
            <UIcon name="i-heroicons-arrow-path" class="w-6 h-6 animate-spin text-gray-400" />
          </div>
          
          <div v-else-if="messages.length === 0" class="text-center text-gray-500 dark:text-gray-400">
            <p>No messages yet. Start the conversation!</p>
          </div>

          <div v-else class="space-y-4">
            <div
              v-for="message in allMessages"
              :key="message.id || message.timestamp"
              class="flex space-x-3"
              :class="{ 'flex-row-reverse space-x-reverse': message.userId === currentUser?.id || message.user_id === currentUser?.id }"
            >
              <div style="width: 32px; height: 32px; background: linear-gradient(135deg, #3b82f6, #8b5cf6); border-radius: 50%; display: flex; align-items: center; justify-content: center; color: white; font-size: 12px;">
                {{ (message.user?.username || message.username || '?').charAt(0).toUpperCase() }}
              </div>
              <div
                style="max-width: 320px;"
                :style="{ display: 'flex', flexDirection: 'column', alignItems: (message.userId === currentUser?.id || message.user_id === currentUser?.id) ? 'flex-end' : 'flex-start' }"
              >
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
                  {{ formatTime(message.createdAt || message.timestamp) }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Message Input -->
        <div style="border-top: 1px solid #e5e7eb; padding: 16px;">
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

    <!-- Create Room Modal -->
    <div v-if="showCreateRoom" style="position: fixed; inset: 0; background: rgba(0, 0, 0, 0.5); display: flex; align-items: center; justify-content: center; z-index: 50;">
      <div style="background: white; border-radius: 16px; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04); max-width: 448px; width: 100%; margin: 0 16px;">
        <div style="padding: 24px; border-bottom: 1px solid #e5e7eb;">
          <h3 style="font-size: 18px; font-weight: 600; color: #111827;">Create New Chat Room</h3>
        </div>

        <form @submit.prevent="createRoom" style="padding: 24px; display: flex; flex-direction: column; gap: 16px;">
          <div>
            <label style="display: block; font-size: 14px; font-weight: 500; color: #374151; margin-bottom: 8px;">
              Room Name *
            </label>
            <input
              v-model="newRoom.name"
              type="text"
              placeholder="Enter room name"
              required
              style="width: 100%; padding: 8px 12px; border: 1px solid #d1d5db; border-radius: 8px; background: white; color: #111827; font-size: 14px;"
            />
          </div>
          
          <div>
            <label style="display: block; font-size: 14px; font-weight: 500; color: #374151; margin-bottom: 8px;">
              Description
            </label>
            <textarea
              v-model="newRoom.description"
              placeholder="Enter room description"
              rows="3"
              style="width: 100%; padding: 8px 12px; border: 1px solid #d1d5db; border-radius: 8px; background: white; color: #111827; font-size: 14px; resize: none;"
            ></textarea>
          </div>

          <div style="display: flex; align-items: center;">
            <input
              v-model="newRoom.isPrivate"
              type="checkbox"
              id="private-room"
              style="width: 16px; height: 16px; margin-right: 8px;"
            />
            <label for="private-room" style="font-size: 14px; font-weight: 500; color: #374151;">
              Private Room
            </label>
          </div>
        </form>

        <div style="padding: 24px; border-top: 1px solid #e5e7eb; display: flex; justify-content: flex-end; gap: 8px;">
          <button
            @click="showCreateRoom = false"
            style="padding: 8px 16px; color: #374151; background: transparent; border: none; border-radius: 8px; font-size: 14px; cursor: pointer; transition: background-color 0.2s;"
          >
            Cancel
          </button>
          <button
            @click="createRoom"
            :disabled="creatingRoom"
            style="padding: 8px 16px; background: #3b82f6; color: white; border: none; border-radius: 8px; font-size: 14px; cursor: pointer; transition: background-color 0.2s;"
          >
            {{ creatingRoom ? 'Creating...' : 'Create Room' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineOptions({
  name: 'ChatPage'
})

// Composables
const config = useRuntimeConfig()
const { isAuthenticated, user: currentUser, getAuthHeaders } = useAuth()
const { 
  isConnected, 
  messages: wsMessages, 
  connect: connectWebSocket, 
  disconnect: disconnectWebSocket, 
  sendMessage: sendWebSocketMessage,
  clearMessages: clearWebSocketMessages 
} = useWebSocket()

// Reactive data
const rooms = ref([])
const selectedRoom = ref(null)
const messages = ref([])
const newMessage = ref('')
const loading = ref(false)
const messagesLoading = ref(false)
const creatingRoom = ref(false)
const showCreateRoom = ref(false)
const showParticipants = ref(false)

// New room form
const newRoom = ref({
  name: '',
  description: '',
  isPrivate: false
})

// Load rooms on mount
onMounted(() => {
  loadRooms()
})

// Cleanup on unmount
onUnmounted(() => {
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
  // Disconnect from previous room
  disconnectWebSocket()
  clearWebSocketMessages()
  
  selectedRoom.value = room
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
    messages.value = response.messages || []
  } catch (error) {
    console.error('Failed to load messages:', error)
  } finally {
    messagesLoading.value = false
  }
}

// Send a message
const sendMessage = async () => {
  if (!newMessage.value.trim() || !selectedRoom.value) return
  
  const message = newMessage.value.trim()
  newMessage.value = ''
  
  // Send via WebSocket for real-time
  sendWebSocketMessage(message, selectedRoom.value.id)
  
  // Also save to database via API
  try {
    await $fetch(`/api/v1/chat/rooms/${selectedRoom.value.id}/messages`, {
      baseURL: config.public.apiBaseUrl,
      method: 'POST',
      headers: getAuthHeaders(),
      body: {
        content: message
      }
    })
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
  return combined.sort((a, b) => new Date(a.timestamp || a.created_at).getTime() - new Date(b.timestamp || b.created_at).getTime())
})

// Create a new room
const createRoom = async () => {
  if (!newRoom.value.name.trim()) return
  
  creatingRoom.value = true
  try {
    const response = await $fetch('/api/v1/chat/rooms', {
      baseURL: config.public.apiBaseUrl,
      method: 'POST',
      headers: getAuthHeaders(),
      body: newRoom.value
    })
    
    rooms.value.push(response.room)
    showCreateRoom.value = false
    newRoom.value = { name: '', description: '', isPrivate: false }
  } catch (error) {
    console.error('Failed to create room:', error)
  } finally {
    creatingRoom.value = false
  }
}

// Format time
const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleTimeString([], { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}
</script> 