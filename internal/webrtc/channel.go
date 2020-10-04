package webrtc

type Channel struct {
	// Media Room Information
	roomID   string
	userID   string
	handleID string
}

type ChannelOption func(*Channel)

func NewChannel(options ...func(*Channel)) *Channel {
	channel := &Channel{}

	for _, option := range options {
		option(channel)
	}

	//api.NewPeerConnection(peerConnectionConfig)
	return channel
}
