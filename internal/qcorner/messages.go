package qcorner

import (
	"encoding/json"
	"slices"
)

func (qc *QCorner) broadcastConnectionMessage() {
	players := make([]string, 0)
	usernames := make(map[string]string)
	for p := range qc.connected {
		if !slices.Contains(players, p.player.UserID) {
			players = append(players, p.player.UserID)
			usernames[p.player.UserID] = p.player.Username
		}
	}
	payload, _ := json.Marshal(Message{
		Type: ConnectionType,
		Details: struct {
			Players   []string          `json:"players"`
			Usernames map[string]string `json:"usernames"`
		}{
			Players:   players,
			Usernames: usernames,
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
