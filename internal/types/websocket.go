package types

import "net/http"

// WebSocketHandler ...
type WebSocketHandler interface {
	Upgrade(w http.ResponseWriter, r *http.Request) error
}
