package types

type RoomManager interface {
	Destroy()
	Register() string
	Unregister(roomID string)
}
