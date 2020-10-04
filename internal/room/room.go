package room

import (
	"sync"

	"github.com/KimJeongChul/webrtc-media-server/internal/types"
)

type Room struct {
	roomID string
	group  *sync.Map
}

func NewRoom(roomID string) types.Room {
	r := &Room{
		roomID: roomID,
		group:  new(sync.Map),
	}

	return r
}

func (r *Room) Register(userID string, handleID string) {
	r.group.Store(userID, new(sync.Map))
}
