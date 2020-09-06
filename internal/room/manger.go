package room

import (
	"log"
	"sync"

	"github.com/KimJeongChul/webrtc-media-server/internal/utils"
	"github.com/google/uuid"
)

// RoomManager Manage to room session
type RoomManager struct {
	rooms *sync.Map
}

// NewRoomManager Create new RoomManager
func NewRoomManager() *RoomManager {
	rm := &RoomManager{
		rooms: new(sync.Map),
	}
	return rm
}

// Destroy RoomManager
func (rm *RoomManager) Destroy() {
	log.Println("[PION MEDIA SERVER] room manager destroy")
	utils.EraseSyncMap(rm.rooms)
	rm = nil
}

// Register Register new room
func (rm *RoomManager) Register() string {
	roomID := uuid.New().String()
	if _, ok := rm.rooms.Load(roomID); !ok {
		room := NewRoom(roomID)
		rm.rooms.Store(roomID, room)
	}
	return roomID
}

// Unregister Unregister new room
func (rm *RoomManager) Unregister(roomID string) {
	utils.EraseKeyInSyncMap(roomID, rm.rooms)
}
