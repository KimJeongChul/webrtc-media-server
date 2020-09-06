package net

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// handleMessage
func handleMessages(ws *WebSocket) {
	for {
		msg := <-ws.msg
		log.Println("[WebSocket] ws request message : ", msg)

		switch {
		case msg.Method == "createRoom":
		}
	}
}

// handler
func handler(w http.ResponseWriter, r *http.Request) {

	// Websocket upgrader
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR] Websocket upgrader", err)
	}

	msg := make(chan RequestMessage, 1024)
	ws := NewWebSocket(conn, msg)

	go handleMessages(ws)

	for {
		var req RequestMessage
		err := ws.conn.ReadJSON(&req)
		if err != nil {
			log.Println("[ERROR] Websocket Read Message JSON : ", err)
			break
		}

		ws.msg <- req
	}
}
