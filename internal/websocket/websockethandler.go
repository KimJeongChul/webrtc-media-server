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
			ws.Send(resMsgCreateRoom)

		//Release room
		case msg.Method == "releaseRoom":
			roomID := msg.RoomID
			wsh.roomManager.Unregister(roomID)

			resMsgReleaseRoom := createResMsgReleaseRoom()
			ws.Send(resMsgReleaseRoom)

		//Add RTC session
		case msg.Method == "addRTCSession":
			roomID := msg.RoomID
			userID := msg.UserID
			handleID := msg.HandleID
			mediaDirection := msg.MediaDir
			sdp := msg.Sdp

			pc, err := wsh.webrtcManager.NewPeerConnection()
			if err != nil {
				continue
			}

			channel := wsh.webrtcManager.CreateChannel(roomID, userID, handleID, mediaDirection, ws, pc)
			log.Println(channel)

			switch {
			case mediaDirection == "recvonly":
				wsh.webrtcManager.AddPublisherRTCSession(channel, sdp.Sdp)
			case mediaDirection == "sendonly":
				// TODO create sub channel
			}

			room, err := wsh.roomManager.Load(roomID)
			if err != nil {
				log.Println("[ERROR] roomID:", roomID, err)
				continue
			}

			isExistUser := room.Exist(userID)
			if !isExistUser {
				room.Register(userID)
				room.RegisterHandle(userID, handleID, channel)
			} else {
				room.RegisterHandle(userID, handleID, channel)
			}

		//Set candidate
		case msg.Method == "candidate":
			// TODO Set Candidate

		//Request Subscribe
		case msg.Method == "requestSubscribe":
			// TODO Request Subscribe
		}
	}
}
