package webrtc

import (
	"io"
	"log"
	"sync"
	"time"

	"github.com/KimJeongChul/webrtc-media-server/internal/types"
	"github.com/pion/rtcp"
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

// CreateChannel
func (wm *WebRTCManager) CreateChannel(roomID string, userID string, handleID string, mediaDirection string, ws types.WebSocket, pc *webrtc.PeerConnection) types.Channel {
	channel := &Channel{
		roomID:         roomID,
		userID:         userID,
		handleID:       handleID,
		mediaDirection: mediaDirection,
		pc:             pc,
		ws:             ws,
	}

	// Publisher Channel
	if channel.mediaDirection == "recvonly" {
		// Create a new RtpTransceiver and add it to the set of transceivers.
		_, err := channel.pc.AddTransceiver(webrtc.RTPCodecTypeAudio)
		if err != nil {
			log.Println("[ERROR] webrtc pc AddTransceiver RTPCodecTypeAudio error:", err)
		}

		_, err = channel.pc.AddTransceiver(webrtc.RTPCodecTypeVideo)
		if err != nil {
			log.Println("[ERROR] webrtc pc AddTransceiver RTPCodecTypeVideo error:", err)
		}
	}

	// Sets an event handler which is invoked when a new ICE candidate is found.
	channel.pc.OnICECandidate(func(ice *webrtc.ICECandidate) {
		candidate := ice.ToJSON()
		resMsg := createResMsgCandidate(channel.GetUserID(), channel.GetHandleID(), candidate)
		log.Println(resMsg)
		channel.ws.Send(resMsg)
	})

	// Sets an event handler which is called when an ICE connection state is changed.
	channel.pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		channel.iceConnState = state.String()
		log.Printf("User ID : %s Handle ID : %s ICE connection state : %s\n", channel.GetUserID(), channel.GetHandleID(), channel.iceConnState)
	})

	channel.rtcpPLIInterval = time.Second

	channel.videoTrack = &webrtc.Track{}
	channel.audioTrack = &webrtc.Track{}

	channel.videoTrackLock = &sync.RWMutex{}
	channel.audioTrackLock = &sync.RWMutex{}

	channel.videoRTCPQuit = make(chan bool)

	return channel
}

// AddPublisherRTCSession
func (wm *WebRTCManager) AddPublisherRTCSession(pubChannel types.Channel, sdp string) {
	//roomID := pubChannel.GetRoomID()
	userID := pubChannel.GetUserID()
	handleID := pubChannel.GetHandleID()
	pubPC := pubChannel.GetPeerConnection()
	ws := pubChannel.GetWebSocket()
	var err error

	// Sets an event handler which is called when remote track arrives from a remote peer.
	pubPC.OnTrack(func(remoteTrack *webrtc.Track, receiver *webrtc.RTPReceiver) {
		if remoteTrack.PayloadType() == webrtc.DefaultPayloadTypeVP8 || remoteTrack.PayloadType() == webrtc.DefaultPayloadTypeVP9 || remoteTrack.PayloadType() == webrtc.DefaultPayloadTypeH264 {
			codec := remoteTrack.Codec()
			log.Println("[TRACK] connected remote track:", remoteTrack.Kind(), codec)

			videoTrackLock := pubChannel.GetVideoTrackLock()

			videoTrackLock.Lock()
			videoTrack, err := pubPC.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "video", "pion")
			if err != nil {
				log.Println("[ERROR] WebRTC Peer Connection NewTrack(video) error:", err)
			}
			pubChannel.SetVideoTrack(videoTrack)
			videoTrackLock.Unlock()

			go func() {
				rtcpPLIInterval := pubChannel.GetRTCPPLIInterval()
				ticker := time.NewTicker(rtcpPLIInterval)
				pubVideoTrack := pubChannel.GetVideoTrack()
				videoRTCPQuit := pubChannel.GetVideoRTCPQuit()
				for range ticker.C {
					select {
					case <-videoRTCPQuit:
						return
					default:
						if pubPC == nil {
							return
						}
						if videoTrack == nil {
							return
						}
						rtcpSendErr := pubPC.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{
							MediaSSRC: pubVideoTrack.SSRC(),
						}})

						if rtcpSendErr != nil {
							if rtcpSendErr == io.ErrClosedPipe {
								return
							}
							log.Println(rtcpSendErr)
						}
					}
				}
			}()

			pubVideoTrack := pubChannel.GetVideoTrack()

			for {
				if videoTrack == nil {
					return
				}

				videoRtp, _ := remoteTrack.ReadRTP()
				if videoRtp == nil {
					continue
				}

				videoTrackLock.RLock()
				pubVideoTrack.WriteRTP(videoRtp)
				videoTrackLock.RUnlock()
			}

		} else if remoteTrack.PayloadType() == webrtc.DefaultPayloadTypeOpus {
			codec := remoteTrack.Codec()
			log.Println("[TRACK] connected remote track:", remoteTrack.Kind(), codec)

			audioTrackLock := pubChannel.GetAudioTrackLock()

			audioTrackLock.Lock()
			audioTrack, err := pubPC.NewTrack(remoteTrack.PayloadType(), remoteTrack.SSRC(), "audio", "pion")
			if err != nil {
				log.Println("[ERROR] WebRTC Peer Connection NewTrack(video) error:", err)
			}
			pubChannel.SetAudioTrack(audioTrack)
			audioTrackLock.Unlock()

			pubAudioTrack := pubChannel.GetAudioTrack()

			for {
				if pubAudioTrack == nil {
					return
				}

				audioRtp, _ := remoteTrack.ReadRTP()
				if audioRtp == nil {
					continue
				}
				audioTrackLock.RLock()
				pubAudioTrack.WriteRTP(audioRtp)
				audioTrackLock.RUnlock()
			}
		}
	})

	err = pubPC.SetRemoteDescription(
		webrtc.SessionDescription{
			SDP:  string(sdp),
			Type: webrtc.SDPTypeOffer})

	answer, err := pubPC.CreateAnswer(nil)
	log.Println("[ERROR] peerconnection CreateAnswer error:", err)

	pubPC.OnICECandidate(func(ice *webrtc.ICECandidate) {
		if ice != nil {
			candidate := ice.ToJSON()
			resMsg := createResMsgCandidate(userID, handleID, candidate)
			log.Println(resMsg)
			ws.Send(resMsg)
		}
	})

	pubPC.SetLocalDescription(answer)

	resMsgAnswerSDP := createResMsgAnswerSDP(userID, handleID, answer)
	log.Println(resMsgAnswerSDP)
	ws.Send(resMsgAnswerSDP)
}
