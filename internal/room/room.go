package room

import (
	"errors"
	"log"
	"sync"

	"github.com/KimJeongChul/webrtc-media-server/internal/types"
	"github.com/KimJeongChul/webrtc-media-server/internal/utils"
)

// Room ...
type Room struct {
	roomID    string
	userGroup *sync.Map
}

// NewRoom Create new room
func NewRoom(roomID string) *Room {
	r := &Room{
		roomID:    roomID,
		userGroup: new(sync.Map),
	}
	return r
}

// Register Register with userGroup
func (r *Room) Register(userID string) {
	r.userGroup.Store(userID, new(sync.Map))
}

// Unregsiter Unregister from userGroup
func (r *Room) Unregister(userID string) {
	utils.EraseKeyInSyncMap(userID, r.userGroup)
}

// Exist ...
func (r *Room) Exist(userID string) bool {
	if _, ok := r.userGroup.Load(userID); ok {
		return true
	} else {
		return false
	}
}

// Load Load userGroup
func (r *Room) Load(userID string) (*sync.Map, error) {
	if u, ok := r.userGroup.Load(userID); ok {
		user := u.(*sync.Map)
		return user, nil
	} else {
		return nil, errors.New("[ERROR] the userGroup doesn't exist")
	}
}

// RegisterHandle Register handle
func (r *Room) RegisterHandle(userID string, handleID string, channel types.Channel) error {
	user, err := r.Load(userID)
	if err != nil {
		log.Println("[ERROR] Room Load error:", err)
		return errors.New("[ERROR] register handle failed!")
	} else if user != nil {
		handle := NewHandle(channel)
		user.Store(handleID, handle)
		log.Println("유저 핸들 저장", handleID, handle)
	}
	return nil
}

// LoadHandle
func (r *Room) LoadHandle(userID string, handleID string) (types.Handle, error) {
	user, err := r.Load(userID)
	if err == nil {
		if h, ok := user.Load(handleID); ok {
			handle := h.(*Handle)
			return handle, nil
		} else {
			return nil, errors.New("[ERROR] the handle doesn't exist")
		}
	} else {
		log.Println("[ERROR] Room Load error:", err)
		return nil, errors.New("[ERROR] load handle failed!")
	}
}
