import { ref } from 'vue'

export interface WebSocketMessage {
  type: string
  payload: any
}

class WebSocketService {
  private ws: WebSocket | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 3000
  private listeners: Map<string, Array<(payload: any) => void>> = new Map()
  private pingInterval: number | null = null
  private token: string = ''

  public isConnected = ref(false)

  connect(token: string) {
    if (this.ws?.readyState === WebSocket.OPEN) {
      // console.log('WebSocket already connected')
      return
    }

    this.token = token

    // Determine WebSocket URL based on current location
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.hostname
    const port = window.location.hostname === 'localhost' ? ':8080' : ''
    const wsUrl = `${protocol}//${host}${port}/ws?token=${token}`

    try {
      this.ws = new WebSocket(wsUrl)

      this.ws.onopen = () => {
        // console.log('WebSocket connected')
        this.isConnected.value = true
        this.reconnectAttempts = 0
        this.startPingInterval()
      }

      this.ws.onmessage = (event) => {
        try {
          // console.log('[WebSocket Service] Raw message received:', event.data)
          const message: WebSocketMessage = JSON.parse(event.data)
          // console.log('[WebSocket Service] Parsed message:', message)
          this.handleMessage(message)
        } catch (error) {
          // console.error('[WebSocket Service] Failed to parse message:', error)
        }
      }

      this.ws.onerror = (error) => {
        // console.error('WebSocket error:', error)
      }

      this.ws.onclose = () => {
        // console.log('WebSocket disconnected')
        this.isConnected.value = false
        this.stopPingInterval()
        this.attemptReconnect(this.token)
      }
    } catch (error) {
      // console.error('Failed to connect WebSocket:', error)
      this.attemptReconnect(this.token)
    }
  }

  private startPingInterval() {
    this.stopPingInterval()
    // Send ping every 30 seconds to keep connection alive
    this.pingInterval = window.setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify({ type: 'ping' }))
      }
    }, 30000)
  }

  private stopPingInterval() {
    if (this.pingInterval !== null) {
      clearInterval(this.pingInterval)
      this.pingInterval = null
    }
  }

  private attemptReconnect(token: string) {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      // console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`)
      setTimeout(() => {
        this.connect(token)
      }, this.reconnectDelay)
    } else {
      // console.error('Max reconnect attempts reached')
    }
  }

  private handleMessage(message: WebSocketMessage) {
    // console.log('[WebSocket Service] Handling message type:', message.type)
    const listeners = this.listeners.get(message.type)
    if (listeners) {
      // console.log('[WebSocket Service] Found', listeners.length, 'listeners for type:', message.type)
      listeners.forEach(callback => callback(message.payload))
    } else {
      // console.log('[WebSocket Service] No listeners found for type:', message.type)
    }
  }

  on(type: string, callback: (payload: any) => void) {
    if (!this.listeners.has(type)) {
      this.listeners.set(type, [])
    }
    this.listeners.get(type)!.push(callback)
  }

  off(type: string, callback: (payload: any) => void) {
    const listeners = this.listeners.get(type)
    if (listeners) {
      const index = listeners.indexOf(callback)
      if (index > -1) {
        listeners.splice(index, 1)
      }
    }
  }

  disconnect() {
    this.stopPingInterval()
    if (this.ws) {
      this.ws.close()
      this.ws = null
      this.isConnected.value = false
    }
  }
}

export const wsService = new WebSocketService()
