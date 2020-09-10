package websocket

import (
	"log"
	"net/http"

	"github.com/KimJeongChul/webrtc-media-server/internal/types"
	"github.com/gorilla/websocket"
)

// WebSocketHandler ...
type WebSocketHandler struct {
	upgrader    websocket.Upgrader
	roomManager types.RoomManager
}

// New Create WebSocketHandler
func New(rm types.RoomManager) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		roomManager: rm,
	}
}

// Upgrade gorilla/websocket Upgrader
func (wsh *WebSocketHandler) Upgrade(w http.ResponseWriter, r *http.Request) error {
	conn, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[ERROR] Failed to upgrade websocket connection")
		return err
	}

	wsh.handle(conn)

	return nil
}

/// handle ...
func (wsh *WebSocketHandler) handle(conn *websocket.Conn) {
	msg := make(chan RequestMessage, 1024)
	ws := NewWebSocket(conn, msg)

	go wsh.handleMessage(ws)

	for {
		var reqMsg RequestMessage

		err := ws.conn.ReadJSON(&reqMsg)
		if err != nil {
			log.Printf("[WEBSOCKET] error: %v", err)
			break
		}
	}
}

// handleMessage ...
func (wsh *WebSocketHandler) handleMessage(ws *WebSocket) {
	for {
		msg := <-ws.msg
		log.Println("[WebSocket] ws request message : ", msg)

		switch {
		// Create room
		case msg.Method == "createRoom":
			roomID := wsh.roomManager.Register()

			resMsgCreateRoom := createResMsgCreateRoom(roomID)
			ws.send(resMsgCreateRoom)

		// Release room
		case msg.Method == "releaseRoom":
			roomID := msg.RoomID
			wsh.roomManager.Unregister(roomID)

			resMsgReleaseRoom := createResMsgReleaseRoom()
			ws.send(resMsgReleaseRoom)
		}
	}
}
