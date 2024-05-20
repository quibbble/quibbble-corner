package qcorner

import (
	"log"
	"net/http"
	"runtime/debug"
	"sync"
	"time"
)

const (
	storageLimit   = 50
	storageTimeout = time.Minute * 30
)

type QCorner struct {
	// mux routes the various endpoints to the appropriate handler.
	mux http.ServeMux

	// connected represents all players currently connected to the server.
	connected map[*Connection]struct{}

	// joinCh and leaveCh adds/remove a player from the server.
	joinCh, leaveCh chan *Connection

	// inputCh sends actions to the server to be processed.
	inputCh chan *ChatMessage

	// messages is the list of past messages
	messages []*ChatMessage

	// mu on messages
	mu sync.Mutex
}

func NewQCorner() *QCorner {
	qc := &QCorner{
		connected: make(map[*Connection]struct{}),
		joinCh:    make(chan *Connection),
		leaveCh:   make(chan *Connection),
		inputCh:   make(chan *ChatMessage),
	}
	go qc.start()
	go qc.clean()
	qc.mux.Handle("GET /qcorner", http.HandlerFunc(qc.connectHandler))
	qc.mux.HandleFunc("GET /health", healthHandler)
	return qc
}

func (qc *QCorner) start() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(string(debug.Stack()))
		}
	}()

	for {
		select {
		case p := <-qc.joinCh:
			qc.connected[p] = struct{}{}
			qc.broadcastConnectionMessage()
			qc.sendChatMessages(p)
		case p := <-qc.leaveCh:
			delete(qc.connected, p)
			go p.Close()
			qc.broadcastConnectionMessage()
		case msg := <-qc.inputCh:
			qc.mu.Lock()
			qc.messages = append(qc.messages, msg)
			if len(qc.messages) > storageLimit {
				qc.messages = qc.messages[1:]
			}
			qc.mu.Unlock()
			qc.broadcastChatMessage(msg)
		}
	}
}

func (qc *QCorner) clean() {
	// remove messages older than storageTimeout
	for range time.Minute {
		qc.mu.Lock()
		idx := -1
		for i, msg := range qc.messages {
			t := time.Unix(msg.Timestamp, 0)
			if t.Add(storageTimeout).After(time.Now()) {
				idx = i
			}
		}
		if idx != -1 {
			qc.messages = qc.messages[idx+1:]
		}
		qc.mu.Unlock()
	}
}
