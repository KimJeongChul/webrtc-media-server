package room

type Room struct {
	roomID string
}

func NewRoom(roomID string) *Room {
	r := &Room{
		roomID: roomID,
	}

	return r
}
