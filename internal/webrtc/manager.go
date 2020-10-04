package webrtc

import (
	"github.com/pion/webrtc/v2"
)

// WebRTCManager ...
type WebRTCManager struct {
	SettingEngine webrtc.SettingEngine
	MediaEngine   webrtc.MediaEngine
	Configuration webrtc.Configuration
	API           *webrtc.API
}

// NewWebRTCManager ...
func NewWebRTCManager() *WebRTCManager {
	wm := &WebRTCManager{}
	wm.SettingEngine = webrtc.SettingEngine{}
	wm.SettingEngine.SetEphemeralUDPPortRange(10000, 20000)

	wm.MediaEngine = webrtc.MediaEngine{}

	opusCodec := webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000)
	h264Codec := webrtc.NewRTPH264Codec(webrtc.DefaultPayloadTypeH264, 90000)

	wm.MediaEngine.RegisterCodec(opusCodec)
	wm.MediaEngine.RegisterCodec(h264Codec)

	wm.Configuration = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
	}

	wm.API = webrtc.NewAPI(webrtc.WithMediaEngine(wm.MediaEngine),
		webrtc.WithSettingEngine(wm.SettingEngine))

	return wm
}

// NewPeerConnection
func (wm *WebRTCManager) NewPeerConnection() (*webrtc.PeerConnection, error) {
	pc, err := wm.API.NewPeerConnection(wm.Configuration)
	if err != nil {
		return nil, err
	}
	return pc, nil
}
