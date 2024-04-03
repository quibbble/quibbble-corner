package qcorner

import (
	"context"
	"log"
	"net/http"

	"github.com/quibbble/quibbble-controller/pkg/auth"
	"nhooyr.io/websocket"
)

const ReadOnly = "readOnly"

func (qc *QCorner) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	qc.mux.ServeHTTP(w, r)
}

func (qc *QCorner) connectHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // allow origin checks are handled at the ingress
	})
	if err != nil {
		log.Println(err.Error())
		return
	}

	userId, ok := r.Context().Value(auth.UserID).(string)
	if !ok {
		userId = ReadOnly
	}
	username, ok := r.Context().Value(auth.Username).(string)
	if !ok {
		username = ReadOnly
	}

	p := NewConnection(&Player{userId, username}, conn, qc.inputCh)
	qc.joinCh <- p

	ctx := context.Background()
	go func() {
		if err := p.ReadPump(ctx); err != nil {
			log.Println(err.Error())
		}
		qc.leaveCh <- p
		p.conn.CloseNow()
	}()
	go func() {
		if err := p.WritePump(ctx); err != nil {
			log.Println(err.Error())
		}
	}()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}
