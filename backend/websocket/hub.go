package websocket

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Client struct {
	UserID uint
	Conn   chan []byte
}

type Hub struct {
	Clients    map[uint]*Client
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

type Message struct {
	UserID  uint        `json:"user_id"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

var GlobalHub = &Hub{
	Clients:    make(map[uint]*Client),
	Broadcast:  make(chan Message, 256),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.UserID] = client
			h.mu.Unlock()
			fmt.Printf("[WebSocket Hub] Client registered: UserID=%d, Total clients=%d\n", client.UserID, len(h.Clients))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Conn)
				fmt.Printf("[WebSocket Hub] Client unregistered: UserID=%d, Total clients=%d\n", client.UserID, len(h.Clients))
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.RLock()
			if client, ok := h.Clients[message.UserID]; ok {
				fmt.Printf("[WebSocket Hub] Sending message to UserID=%d, Type=%s, Payload=%+v\n", message.UserID, message.Type, message.Payload)
				select {
				case client.Conn <- encodeMessage(message):
					fmt.Printf("[WebSocket Hub] Message sent successfully to UserID=%d\n", message.UserID)
				default:
					fmt.Printf("[WebSocket Hub] Failed to send message to UserID=%d, channel full or closed\n", message.UserID)
					close(client.Conn)
					h.mu.RUnlock()
					h.mu.Lock()
					delete(h.Clients, message.UserID)
					h.mu.Unlock()
					h.mu.RLock()
				}
			} else {
				fmt.Printf("[WebSocket Hub] User not connected: UserID=%d\n", message.UserID)
			}
			h.mu.RUnlock()
		}
	}
}

func encodeMessage(msg Message) []byte {
	data := map[string]interface{}{
		"type":    msg.Type,
		"payload": msg.Payload,
	}
	bytes, _ := json.Marshal(data)
	return bytes
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID uint, msgType string, payload interface{}) {
	h.Broadcast <- Message{
		UserID:  userID,
		Type:    msgType,
		Payload: payload,
	}
}
