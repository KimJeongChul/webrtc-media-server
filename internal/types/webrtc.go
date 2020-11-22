package types

import (
	"sync"
	"time"

	"github.com/pion/webrtc/v2"
)

// Channel ...
type Channel interface {
	GetRoomID() string
	GetUserID() string
	GetHandleID() string
	GetPeerConnection() *webrtc.PeerConnection
	GetVideoTrackLock() *sync.RWMutex
	GetAudioTrackLock() *sync.RWMutex
	GetVideoTrack() *webrtc.Track
	GetAudioTrack() *webrtc.Track
	SetVideoTrack(track *webrtc.Track)
	SetAudioTrack(track *webrtc.Track)
	GetRTCPPLIInterval() time.Duration
	GetVideoRTCPQuit() chan bool
	GetWebSocket() WebSocket
	GetIceQueue() []webrtc.ICECandidateInit
	SetIceQueue(iceQueue []webrtc.ICECandidateInit)
	GetIceLock() *sync.RWMutex
	SetIsSetRemoteSDP(isSetRemoteSDP bool)
	GetIsSetRemoteSDP() bool
}

// WebRTCManager ...
type WebRTCManager interface {
	NewPeerConnection() (*webrtc.PeerConnection, error)
	CreateChannel(roomID string, userID string, handleID string, mediaDirection string, ws WebSocket, pc *webrtc.PeerConnection) Channel
	AddPublisherRTCSession(pubChannel Channel, sdp string)
}
