package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocket
type WebSocket struct {
	conn *websocket.Conn
	mu   sync.Mutex
	// 웹소켓에 전달되는 메시지 채널
	msg chan RequestMessage
}

// NewWebSocket
func NewWebSocket(conn *websocket.Conn, msg chan RequestMessage) *WebSocket {
	ws := &WebSocket{
		conn: conn,
		msg:  msg,
	}
	return ws
}

// Send
func (w *WebSocket) Send(v interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.WriteJSON(v)
}
