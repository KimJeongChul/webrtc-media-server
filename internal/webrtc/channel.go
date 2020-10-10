package webrtc

import (
	"sync"
	"time"

	"github.com/KimJeongChul/webrtc-media-server/internal/types"
	"github.com/pion/webrtc/v2"
)

// Channel ...
type Channel struct {
	// Media Room Information
	roomID         string
	userID         string
	handleID       string
	mediaDirection string

	// WebSocket
	ws types.WebSocket

	// PeerConnection
	pc *webrtc.PeerConnection

	// ICE connection state
	iceConnState string

	// RTCP
	rtcpPLIInterval time.Duration

	// Track
	videoTrack *webrtc.Track
	audioTrack *webrtc.Track

	// Mutex
	videoTrackLock *sync.RWMutex
	audioTrackLock *sync.RWMutex

	// Channel
	videoRTCPQuit chan bool
}

// GetRoomID ...
func (ch *Channel) GetRoomID() string {
	return ch.roomID
}

// GetUserID ...
func (ch *Channel) GetUserID() string {
	return ch.userID
}

// GetHandleID ...
func (ch *Channel) GetHandleID() string {
	return ch.handleID
}

// GetPeerConnection ...
func (ch *Channel) GetPeerConnection() *webrtc.PeerConnection {
	return ch.pc
}

// GetVideoTrack ...
func (ch *Channel) GetVideoTrack() *webrtc.Track {
	return ch.videoTrack
}

// GetAudioTrack ...
func (ch *Channel) GetAudioTrack() *webrtc.Track {
	return ch.audioTrack
}

// GetVideoTrackLock ...
func (ch *Channel) GetVideoTrackLock() *sync.RWMutex {
	return ch.videoTrackLock
}

// GetAudioTrackLock ...
func (ch *Channel) GetAudioTrackLock() *sync.RWMutex {
	return ch.audioTrackLock
}

// SetVideoTrack ...
func (ch *Channel) SetVideoTrack(track *webrtc.Track) {
	ch.videoTrack = track
}

// SetAudioTrack ...
func (ch *Channel) SetAudioTrack(track *webrtc.Track) {
	ch.audioTrack = track
}

// GetRTCPPLIInterval ...
func (ch *Channel) GetRTCPPLIInterval() time.Duration {
	return ch.rtcpPLIInterval
}

// GetVideoRTCPQuit ...
func (ch *Channel) GetVideoRTCPQuit() chan bool {
	return ch.videoRTCPQuit
}

// GetWebSocket ...
func (ch *Channel) GetWebSocket() types.WebSocket {
	return ch.ws
}
