package qcorner

const (
	Ping = "ping"
	Chat = "chat"
)

type Action struct {
	Type    string      `json:"type"`
	Details interface{} `json:"details,omitempty"`
}

type ServerAction struct {
	*Connection
	*Action
}
