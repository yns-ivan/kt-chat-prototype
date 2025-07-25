interface ChatMessage {
  type: string
  room_id: string
  user_id: string
  username: string
  content: string
  timestamp: string
  files?: unknown[]
}

export const useWebSocket = () => {
  const socket = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const messages = ref<ChatMessage[]>([])
  const config = useRuntimeConfig()
  const { user } = useAuth()
  
  // Reconnection logic
  const reconnectAttempts = ref(0)
  const maxReconnectAttempts = 5
  const reconnectDelay = 1000 // 1 second
  const currentRoomId = ref<string | null>(null)
  const reconnectTimeout = ref<NodeJS.Timeout | null>(null)

  // Connect to WebSocket
  const connect = (roomId: string) => {
    if (!user.value) return

    // Clear any existing reconnection timeout
    if (reconnectTimeout.value) {
      clearTimeout(reconnectTimeout.value)
      reconnectTimeout.value = null
    }

    // Reset reconnection attempts for new room
    reconnectAttempts.value = 0
    currentRoomId.value = roomId

    const wsUrl = config.public.apiBaseUrl.replace('http', 'ws')
    const url = `${wsUrl}/api/v1/ws?user_id=${user.value.id}&username=${user.value.username}&room_id=${roomId}`
    
    socket.value = new WebSocket(url)

    socket.value.onopen = () => {
      console.log('WebSocket connected')
      isConnected.value = true
      reconnectAttempts.value = 0 // Reset attempts on successful connection
    }

    socket.value.onmessage = (event) => {
      try {
        const message: ChatMessage = JSON.parse(event.data)
        messages.value.push(message)
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }

    socket.value.onclose = (event) => {
      console.log('WebSocket disconnected', event.code, event.reason)
      isConnected.value = false
      
      // Attempt to reconnect if not a manual disconnect
      if (event.code !== 1000 && currentRoomId.value && reconnectAttempts.value < maxReconnectAttempts) {
        attemptReconnect()
      }
    }

    socket.value.onerror = (error) => {
      console.error('WebSocket error:', error)
      isConnected.value = false
    }
  }

  // Attempt to reconnect
  const attemptReconnect = () => {
    if (!currentRoomId.value || reconnectAttempts.value >= maxReconnectAttempts) return

    reconnectAttempts.value++
    console.log(`Attempting to reconnect (${reconnectAttempts.value}/${maxReconnectAttempts})...`)

    reconnectTimeout.value = setTimeout(() => {
      if (currentRoomId.value) {
        connect(currentRoomId.value)
      }
    }, reconnectDelay * reconnectAttempts.value) // Exponential backoff
  }

  // Disconnect from WebSocket
  const disconnect = () => {
    // Clear any pending reconnection attempts
    if (reconnectTimeout.value) {
      clearTimeout(reconnectTimeout.value)
      reconnectTimeout.value = null
    }
    
    // Reset reconnection state
    reconnectAttempts.value = 0
    currentRoomId.value = null
    
    if (socket.value) {
      socket.value.close(1000, 'Manual disconnect') // Use code 1000 for normal closure
      socket.value = null
      isConnected.value = false
    }
  }

  // Send a message
  const sendMessage = (content: string, roomId: string) => {
    if (!socket.value || !isConnected.value || !user.value) return

    const message: ChatMessage = {
      type: 'message',
      room_id: roomId,
      user_id: String(user.value.id),
      username: user.value.username,
      content,
      timestamp: new Date().toISOString()
    }

    socket.value.send(JSON.stringify(message))
  }

  // Clear messages
  const clearMessages = () => {
    messages.value = []
  }

  return {
    socket: readonly(socket),
    isConnected: readonly(isConnected),
    messages: readonly(messages),
    connect,
    disconnect,
    sendMessage,
    clearMessages,
    reconnect: () => {
      if (currentRoomId.value) {
        connect(currentRoomId.value)
      }
    }
  }
} 