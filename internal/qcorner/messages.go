package qcorner

import (
	"encoding/json"
	"slices"
)

func (qc *QCorner) broadcastConnectionMessage() {
	names := make([]string, 0)
	for p := range qc.connected {
		if !slices.Contains(names, p.player.Name) {
			names = append(names, p.player.Name)
		}
	}
	payload, _ := json.Marshal(Message{
		Type: ConnectionType,
		Details: struct {
			Names []string `json:"names"`
		}{
			Names: names,
		},
	})
	for p := range qc.connected {
		qc.sendMessage(p, payload)
	}
}

func (qc *QCorner) broadcastChatMessage(msg *ChatMessage) {
	payload, _ := json.Marshal(Message{
		Type:    ChatType,
		Details: msg,
	})
	for p := range qc.connected {
		qc.sendMessage(p, payload)
	}
}

func (qc *QCorner) sendChatMessages(player *Connection) {
	qc.mu.Lock()
	defer qc.mu.Unlock()
	for _, msg := range qc.messages {
		payload, _ := json.Marshal(Message{
			Type:    ChatType,
			Details: msg,
		})
		qc.sendMessage(player, payload)
	}
}

func (qc *QCorner) sendMessage(player *Connection, payload []byte) {
	select {
	case player.outputCh <- payload:
	default:
		delete(qc.connected, player)
		go player.Close()
	}
}
