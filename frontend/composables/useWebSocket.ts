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

  // Connect to WebSocket
  const connect = (roomId: string) => {
    if (!user.value) return

    const wsUrl = config.public.wsUrl.replace('http', 'ws')
    const url = `${wsUrl}/api/v1/ws?user_id=${user.value.id}&username=${user.value.username}&room_id=${roomId}`
    
    socket.value = new WebSocket(url)

    socket.value.onopen = () => {
      console.log('WebSocket connected')
      isConnected.value = true
    }

    socket.value.onmessage = (event) => {
      try {
        const message: ChatMessage = JSON.parse(event.data)
        messages.value.push(message)
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }

    socket.value.onclose = () => {
      console.log('WebSocket disconnected')
      isConnected.value = false
    }

    socket.value.onerror = (error) => {
      console.error('WebSocket error:', error)
      isConnected.value = false
    }
  }

  // Disconnect from WebSocket
  const disconnect = () => {
    if (socket.value) {
      socket.value.close()
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
    clearMessages
  }
} 