package websocket

// ResponseMessage ...
type ResponseMessage struct {
	Method string `json:"method"`
	RoomID string `json:"roomID"`
	Status int    `json:"status"`
}

// createResMsgCreateRoom ...
func createResMsgCreateRoom(roomID string) ResponseMessage {
	resMsg := ResponseMessage{
		Method: "resCreateRoom",
		RoomID: roomID,
		Status: 200,
	}
	return resMsg
}

// createResMsgCreateRoom ...
func createResMsgReleaseRoom() ResponseMessage {
	resMsg := ResponseMessage{
		Status: 200,
	}
	return resMsg
}
