package qcorner

import (
	"log"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	goaway "github.com/TwiN/go-away"
)

const (
	storageLimit   = 50
	storageTimeout = time.Minute * 30
)

type QCorner struct {
	// admin user info
	adminUsername, adminPassword string

	// mux routes the various endpoints to the appropriate handler.
	mux http.ServeMux

	// connected represents all players currently connected to the server.
	connected map[*Connection]struct{}

	// joinCh and leaveCh adds/remove a player from the server.
	joinCh, leaveCh chan *Connection

	// inputCh sends actions to the server to be processed.
	inputCh chan *ServerAction

	// messages is the list of past messages
	messages []*ChatDetails

	// mu on messages
	mu sync.Mutex
}

func NewQCorner(adminUsername, adminPassword string) *QCorner {
	qc := &QCorner{
		adminUsername: adminUsername,
		adminPassword: adminPassword,
		connected:     make(map[*Connection]struct{}),
		joinCh:        make(chan *Connection),
		leaveCh:       make(chan *Connection),
		inputCh:       make(chan *ServerAction),
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
		case a := <-qc.inputCh:
			switch a.Type {
			case Ping:
				qc.sendPongMessage(a.Connection)
			case Chat:
				message, ok := a.Details.(string)
				if !ok {
					continue
				}
				message = goaway.Censor(message)
				msg := &ChatDetails{
					Name:      a.player.Name,
					Message:   message,
					Timestamp: time.Now().Unix(),
				}
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
}

func (qc *QCorner) clean() {
	// remove messages older than storageTimeout
	for range time.Tick(time.Minute) {
		qc.mu.Lock()
		idx := -1
		for i, msg := range qc.messages {
			t := time.Unix(msg.Timestamp, 0)
			if t.Add(storageTimeout).Before(time.Now()) {
				idx = i
			}
		}
		if idx != -1 {
			qc.messages = qc.messages[idx+1:]
		}
		qc.mu.Unlock()
	}
}
