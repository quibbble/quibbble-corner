package qcorner

const (
	ConnectionMessage = "connection"
	ChatMessage       = "chat"
	PongMessage       = "pong"
)

type Message struct {
	Type    string      `json:"type"`
	Details interface{} `json:"details"`
}

type Player struct {
	Name string `json:"name"`
}

type ChatDetails struct {
	Name      string `json:"name"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
