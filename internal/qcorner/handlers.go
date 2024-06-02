package qcorner

import (
	"context"
	"log"
	"net/http"
	"strings"

	"nhooyr.io/websocket"
)

func (qc *QCorner) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	qc.mux.ServeHTTP(w, r)
}

func (qc *QCorner) connectHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// admin auth
	split := strings.Split(name, ":")
	if len(split) == 2 && split[0] == qc.adminUsername && split[1] == qc.adminPassword {
		name = split[0]
	} else if split[0] == qc.adminUsername {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // allow origin checks are handled at the ingress
	})
	if err != nil {
		log.Println(err.Error())
		return
	}

	p := NewConnection(&Player{name}, conn, qc.inputCh)
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
