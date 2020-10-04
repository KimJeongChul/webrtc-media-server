package types

import (
	"github.com/pion/webrtc/v2"
)

// Channel ...
type Channel interface {
}

// WebRTCManager ...
type WebRTCManager interface {
	NewPeerConnection() (*webrtc.PeerConnection, error)
}
