package qcorner

const (
	ConnectionType = "connection"
	ChatType       = "chat"
)

type Message struct {
	Type    string      `json:"type"`
	Details interface{} `json:"details"`
}

type ConnectionMessage struct {
	Players []*Player `json:"players"`
}

type Player struct {
	UserID   string
	Username string
}

type ChatMessage struct {
	*Player
	Message string
}
