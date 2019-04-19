package marv

// Message NEEDSCOMMENT
type Message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type messageHandler func(*Message) *Message
