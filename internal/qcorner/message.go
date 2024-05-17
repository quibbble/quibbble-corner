package qcorner

const (
	ConnectionType = "connection"
	ChatType       = "chat"
)

type Message struct {
	Type    string      `json:"type"`
	Details interface{} `json:"details"`
}

type Player struct {
	Name string `json:"name"`
}

type ChatMessage struct {
	Name      string `json:"name"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
