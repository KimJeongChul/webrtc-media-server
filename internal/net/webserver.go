package net

import (
	"net/http"
	"time"

	"github.com/KimJeongChul/webrtc-media-server/internal/types"
)

// WebServer
type WebServer struct {
	port string
	rm   types.RoomManager
}

// NewWebServer
func NewWebServer(port string, rm types.RoomManager) *WebServer {
	ws := &WebServer{
		port: port,
		rm:   rm,
	}
	return ws
}

// Start
func (ws *WebServer) Start() {

	//configureRLimit()

	// Websocket Handler
	http.HandleFunc("/ws", handler)

	// WebPage Route
	http.HandleFunc("/", web)
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	// WebServer
	srv := &http.Server{
		Addr:         ":" + ws.port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// WebServer Listen and Serve
	panic(srv.ListenAndServe())
}
