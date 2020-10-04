package main

import (
	"flag"
	"log"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"github.com/KimJeongChul/webrtc-media-server/internal/net"
	"github.com/KimJeongChul/webrtc-media-server/internal/room"
	"github.com/KimJeongChul/webrtc-media-server/internal/webrtc"
	"github.com/KimJeongChul/webrtc-media-server/internal/websocket"
)

func main() {
	// Write rotates log files from within the application.
	rl, _ := rotatelogs.New("./webrtc-media-server.%Y%m%d")

	log.Println("WebRTC Media Server start ...")

	mediaServerPort := flag.String("p", "8080", "media server port")
	rotateLog := flag.Bool("l", false, "rotate logs")
	flag.Parse()

	if *rotateLog {
		log.SetOutput(rl)
	}

	// WebRTCManager
	webrtcManager := webrtc.NewWebRTCManager()

	// RoomManager
	roomManager := room.NewRoomManager()

	// WebSocket Handler
	webSocketHandler := websocket.New(roomManager, webrtcManager)

	// WebServer
	ws := net.NewWebServer(*mediaServerPort, webSocketHandler)
	ws.Start()
}
