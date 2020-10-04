package websocket

import "github.com/pion/webrtc/v2"

// SDP
type Sdp struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
}

// RequestMessage
type RequestMessage struct {
	Method    string                  `json:"method"`
	RoomID    string                  `json:"roomID"`
	UserID    string                  `json:"userID"`
	HandleID  string                  `json:"handleID"`
	Sdp       Sdp                     `json:"sdp"`
	Candidate webrtc.ICECandidateInit `json"candidate"`
	MediaDir  string                  `json:"mediaDir"`
}
