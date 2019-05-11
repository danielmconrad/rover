package marv

import (
	"sync"

	"github.com/gorilla/websocket"
)

// ClientMap NEEDSCOMMENT
type ClientMap struct {
	clients               map[*websocket.Conn]map[string]interface{}
	mutex                 sync.RWMutex
	initialBinaryMessages [][]byte
}

// NewClientMap NEEDSCOMMENT
func NewClientMap(initialBinaryMessages [][]byte) *ClientMap {
	return &ClientMap{
		clients:               map[*websocket.Conn]map[string]interface{}{},
		initialBinaryMessages: initialBinaryMessages,
	}
}

// Clients NEEDSCOMMENT
func (c *ClientMap) Clients() map[*websocket.Conn]map[string]interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.clients
}

// Pause NEEDSCOMMENT
func (c *ClientMap) Pause(client *websocket.Conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.clients[client]["started"] = false
}

// Start NEEDSCOMMENT
func (c *ClientMap) Start(client *websocket.Conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, alreadyHasClient := c.clients[client]
	c.clients[client]["started"] = true

	if !alreadyHasClient {
		for _, frame := range c.initialBinaryMessages {
			client.WriteMessage(websocket.BinaryMessage, frame)
		}
	}
}
