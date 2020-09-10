package websocket

type RequestMessage struct {
	Method string `json:"method"`
	RoomID string `json:"roomID"`
}
