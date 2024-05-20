package qcorner

import (
	"encoding/json"
	"slices"
	"time"
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
	// remove messages older than storageTimeout
	idx := -1
	for i, msg := range qc.messages {
		t := time.Unix(msg.Timestamp, 0)
		if t.Add(storageTimeout).After(time.Now()) {
			idx = i
		}
	}
	if idx != -1 {
		qc.messages = qc.messages[idx:]
	}

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
