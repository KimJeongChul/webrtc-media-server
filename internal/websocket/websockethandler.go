package websocket

import (
	"log"
	"net/http"

	"github.com/KimJeongChul/webrtc-media-server/internal/types"
	"github.com/gorilla/websocket"
)

// WebSocketHandler ...
type WebSocketHandler struct {
	upgrader      websocket.Upgrader
	roomManager   types.RoomManager
	webrtcManager types.WebRTCManager
}

// New Create WebSocketHandler
func New(rm types.RoomManager, wm types.WebRTCManager) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		roomManager:   rm,
		webrtcManager: wm,
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
		ws.msg <- reqMsg
	}
}

// handleMessage ...
func (wsh *WebSocketHandler) handleMessage(ws *WebSocket) {
	for {
		msg := <-ws.msg
		log.Println("[WebSocket] ws request message : ", msg)

		switch {
		//Create room
		case msg.Method == "createRoom":
			roomID := wsh.roomManager.Register()

			resMsgCreateRoom := createResMsgCreateRoom(roomID)
			ws.send(resMsgCreateRoom)

		//Release room
		case msg.Method == "releaseRoom":
			roomID := msg.RoomID
			wsh.roomManager.Unregister(roomID)

			resMsgReleaseRoom := createResMsgReleaseRoom()
			ws.send(resMsgReleaseRoom)

		//Add RTC session
		case msg.Method == "addRTCSession":
			roomID := msg.RoomID
			userID := msg.UserID
			handleID := msg.HandleID

			pc, err := wsh.webrtcManager.NewPeerConnection()
			if err != nil {
				continue
			}

			// TODO Create Channel

			room, err := wsh.roomManager.Load(roomID)
			if err != nil {
				log.Println("[ERROR] roomID:", roomID, err)
				continue
			}

			if false {
				log.Println(userID, handleID, pc, room)
			}

		//Set candidate
		case msg.Method == "candidate":
			// TODO Set Candidate
		}
	}
}
