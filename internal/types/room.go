package types

// Room ...
type Room interface {
	Register(userID string, handleID string)
}

// RoomManager ...
type RoomManager interface {
	Destroy()
	Register() string
	Unregister(roomID string)
	Load(roomID string) (*Room, error)
}
