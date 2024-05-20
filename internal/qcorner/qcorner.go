package qcorner

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

const (
	storageLimit   = 10
	storageTimeout = time.Minute * 10
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
}

func NewQCorner() *QCorner {
	gs := &QCorner{
		connected: make(map[*Connection]struct{}),
		joinCh:    make(chan *Connection),
		leaveCh:   make(chan *Connection),
		inputCh:   make(chan *ChatMessage),
	}
	go gs.Start()
	gs.mux.Handle("GET /qcorner", http.HandlerFunc(gs.connectHandler))
	gs.mux.HandleFunc("GET /health", healthHandler)
	return gs
}

func (qc *QCorner) Start() {
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
			qc.messages = append(qc.messages, msg)
			if len(qc.messages) > storageLimit {
				qc.messages = qc.messages[1:]
			}
			qc.broadcastChatMessage(msg)
		}
	}
}
