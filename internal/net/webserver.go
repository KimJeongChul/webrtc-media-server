package net

import (
	"net/http"
	"time"

	"github.com/KimJeongChul/webrtc-media-server/internal/types"
	"github.com/go-chi/chi"
)

// WebServer
type WebServer struct {
	port string
	wsh  types.WebSocketHandler
}

// NewWebServer
func NewWebServer(port string, wsh types.WebSocketHandler) *WebServer {
	ws := &WebServer{
		port: port,
		wsh:  wsh,
	}
	return ws
}

// Start
func (ws *WebServer) Start() {

	//configureRLimit()

	// Router
	router := chi.NewRouter()

	router.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.wsh.Upgrade(w, r)
	})

	// WebPage Route
	router.Get("/", web)
	router.Handle("/js/*", http.StripPrefix("/js/", http.FileServer(http.Dir("internal/net/js"))))
	router.Handle("/css/*", http.StripPrefix("/css/", http.FileServer(http.Dir("internal/net/css"))))

	// WebServer
	srv := &http.Server{
		Addr:         ":" + ws.port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	// WebServer Listen and Serve
	panic(srv.ListenAndServe())
}
