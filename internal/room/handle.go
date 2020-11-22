package room

import (
	"github.com/KimJeongChul/webrtc-media-server/internal/types"
)

type Handle struct {
	Channel types.Channel
}

// NewHandle Create new handle
func NewHandle(channel types.Channel) *Handle {
	r := &Handle{
		Channel: channel,
	}
	return r
}

// GetChannel Get channel
func (h *Handle) GetChannel() types.Channel {
	return h.Channel
}
