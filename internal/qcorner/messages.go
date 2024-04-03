package qcorner

import "encoding/json"

func (qc *QCorner) broadcastConnectionMessage() {
	for p := range qc.connected {
		payload, _ := json.Marshal(Message{})
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
