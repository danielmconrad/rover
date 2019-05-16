package rover

import (
	"sync"

	"github.com/gorilla/websocket"
)

// ClientMap NEEDSCOMMENT
type ClientMap struct {
	clients               map[*websocket.Conn]bool
	mutex                 sync.RWMutex
	initialBinaryMessages []Frame
}

// NewClientMap NEEDSCOMMENT
func NewClientMap(initialBinaryMessages []Frame) *ClientMap {
	return &ClientMap{
		clients:               map[*websocket.Conn]bool{},
		initialBinaryMessages: initialBinaryMessages,
	}
}

// Clients NEEDSCOMMENT
func (c *ClientMap) Clients() map[*websocket.Conn]bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.clients
}

// Pause NEEDSCOMMENT
func (c *ClientMap) Pause(client *websocket.Conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.clients[client] = false
}

// Start NEEDSCOMMENT
func (c *ClientMap) Start(client *websocket.Conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, alreadyHasClient := c.clients[client]
	c.clients[client] = true

	if !alreadyHasClient {
		for _, frame := range c.initialBinaryMessages {
			client.WriteMessage(websocket.BinaryMessage, frame)
		}
	}
}
