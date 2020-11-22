package types

import (
	"sync"
)

type Handle interface {
	GetChannel() Channel
}

// Room ...
type Room interface {
	Register(userID string)
	Unregister(userID string)
	Exist(userID string) bool
	Load(userID string) (*sync.Map, error)
	RegisterHandle(userID string, handleID string, channel Channel) error
	LoadHandle(userID string, handleID string) (Handle, error)
}

// RoomManager ...
type RoomManager interface {
	Destroy()
	Register() string
	Unregister(roomID string)
	Load(roomID string) (Room, error)
	GetRooms() *sync.Map
}
