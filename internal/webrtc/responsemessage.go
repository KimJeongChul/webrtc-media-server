package webrtc

import "github.com/pion/webrtc/v2"

// SDP
type Sdp struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
}

// ResponseMessage ...
type ResponseMessage struct {
	Method    string                  `json:"method"`
	RoomID    string                  `json:"roomID"`
	UserID    string                  `json:"userID"`
	HandleID  string                  `json:"handleID"`
	Sdp       Sdp                     `json:"sdp"`
	Candidate webrtc.ICECandidateInit `json"candidate"`
	MediaDir  string                  `json:"mediaDir"`
	Status    int                     `json:"status"`
}

// createResMsgCandidate ...
func createResMsgCandidate(userID string, handleID string, candidate webrtc.ICECandidateInit) ResponseMessage {
	resMsg := ResponseMessage{
		Method:    "resCandidate",
		UserID:    userID,
		HandleID:  handleID,
		Candidate: candidate,
		Status:    200,
	}
	return resMsg
}
