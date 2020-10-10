package webrtc

import "github.com/pion/webrtc/v2"

// ResponseMessage ...
type ResponseMessage struct {
	Method    string                    `json:"method"`
	RoomID    string                    `json:"roomID"`
	UserID    string                    `json:"userID"`
	HandleID  string                    `json:"handleID"`
	Sdp       webrtc.SessionDescription `json:"sdp"`
	Candidate webrtc.ICECandidateInit   `json"candidate"`
	MediaDir  string                    `json:"mediaDir"`
	Status    int                       `json:"status"`
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

// createResMsgAnswerSDP ...
func createResMsgAnswerSDP(userID string, handleID string, sdp webrtc.SessionDescription) ResponseMessage {
	resMsg := ResponseMessage{
		Method:   "resAnswerSDP",
		UserID:   userID,
		HandleID: handleID,
		Sdp:      sdp,
		Status:   200,
	}
	return resMsg
}
