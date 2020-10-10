package types

import "net/http"

// WebSocket
type WebSocket interface {
	Send(v interface{}) error
}

// WebSocketHandler ...
type WebSocketHandler interface {
	Upgrade(w http.ResponseWriter, r *http.Request) error
}
