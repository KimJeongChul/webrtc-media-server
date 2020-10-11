let ice = {
    iceServers: [{
        url: "stun:stun.l.google.com:19302"
    }]
};

let media_constraints = {
    "audio": true,
    "video": {
        mandatory: {
            minWidth: 1280, 
            minHeight: 720,
        }
    }
};

let pc_constraints = {
    offerToReceiveAudio: 0, 
    offerToReceiveVideo: 0
};

var pc = null;

// Candidate
let queueCandidate = {};

// SDP
let recvSDP = {};

// Subscriber
let videoSub = {};
let pcSub = {};

//addRTCSession
function addRTCSession() {
    let button_call = document.getElementById('button_call');
    button_call.disabled = true;
    button_call.style.backgroundColor = '#212529';

    navigator.mediaDevices.getUserMedia(media_constraints).then(gotStream).catch((err) => {
        console.log('getUserMedia Error:' + err.stack);
    }).catch(function(err) {
        console.log(err);
    });
}

//gotStream 
function gotStream(stream) {
    //Create PeerConnection
    pc = new webkitRTCPeerConnection(ice); 

    //Event handler which specifies a function to be called when the icecandidate event occurs on an RTCPeerConnection instance.
    pc.onicecandidate = function (event) {
        if (event.candidate) {
            //Send Ice
            const msgCandidate = {
                method: 'candidate',
                candidate: event.candidate,
                roomID: media_room_id,
                userID: user_id,
                handleID: handle_id
            };
            sendMsg(msgCandidate);
        }
    }

    //Event handler which specifies a function to be called when the iceconnectionstatechange event is fired on an RTCPeerConnection instance.
    pc.oniceconnectionstatechange = (event) => {
        console.log("[Publisher] Connection State : " + pc.iceConnectionState)
        if (pc.iceConnectionState === "connected") {
            // Request subscribe
            const msgRequestSubscribe = {
                method: 'requestSubscribe',
                roomID: media_room_id,
                userID: user_id,
                handleID: handle_id,
            }
            sendMsg(msgRequestSubscribe);
        }
    }
    
    videoLocal = document.createElement('video');
    videoLocal.id = 'video_local';
    videoLocal.srcObject = stream;
    videoLocal.autoplay = true;
    
    videoLocal.width = 640;
    videoLocal.height = 480;

    divLocal = document.getElementById('local-div');
    divLocal.appendChild(videoLocal);

    //Adds a MediaStream as a local source of audio or video.
    pc.addStream(stream);

    //Initiates the creation of an SDP offer for the purpose of starting a new WebRTC connection to a remote peer. 
    pc.createOffer(pc_constraints).then((offer) => {
        //Changes the local description associated with the connection.
        return pc.setLocalDescription(offer);
    }).then(() => {
        const msgAddRTCSession = {
            method: 'addRTCSession',
            mediaDir: 'recvonly',
            roomID: media_room_id,
            userID: user_id,
            handleID: handle_id,
            sdp: pc.localDescription
        }
        sendMsg(msgAddRTCSession);

        queueCandidate[user_id] = [];

    }).catch((err) => {
        console.log(err.stack);
    });
}

// Set SDP
function setSDP(msg) {
    const req_user_id = msg.userID;
    const req_handle_id = msg.handleID;

    if(user_id === req_user_id && handle_id === req_handle_id) {
        // Publisher channel
        // Sets the specified session description as the remote peer's current offer or answer.
        pc.setRemoteDescription(msg.sdp, () => {
            console.log('Publisher Set Remote Description completed');
            recvSDP[req_user_id] = true;

            // Set ICE candidate
            queueCandidate[req_user_id].forEach((candidate) => {
                if (candidate != undefined) {
                    console.log('add queued ice candidate:' + JSON.stringify(candidate));
                    pc.addIceCandidate(new RTCIceCandidate(candidate)).catch((e) => {
                        console.log('add Queued ICE Exception:' + e);
                    });
                }
            });
        });
    } else {
        // Subscriber channel
        if (pcSub[req_user_id] !== undefined) {
            pcSub[req_user_id].setRemoteDescription(msg.sdp, () => {
                recvSDP[req_user_id] = true;
                queueCandidate[req_user_id].forEach((candidate) => {
                    pcSub[req_user_id].addIceCandidate(new RTCIceCandidate(candidate)).catch((e) => {
                        console.log('add Queued ICE Exception:' + e);
                    });
                });
            });
        }
    }
}

// Set ICE candidate
function setCandidate(msg) {
    const req_user_id = msg.userID;
    const req_handle_id = msg.handleID;
    if (user_id === req_user_id && handle_id === req_handle_id) {
        // Publisher channel
        if (recvSDP[req_user_id]) {
            // Adds this new remote candidate to the RTCPeerConnection's remote description, which describes the state of the remote end of the connection.
            pc.addIceCandidate(new RTCIceCandidate(msg.candidate)).catch((e) => {
                console.log('addIceCandidate Exception:' + e);
            });
        } else {
            queueCandidate[req_user_id].push(msg.candidate);
        }
    } else {
        // Subscriber channel
        if (pcSub[req_user_id] !== undefined) {
            if (recvSDP[req_user_id]) {
                pcSub[req_user_id].addIceCandidate(new RTCIceCandidate(msg.candidate)).catch((e) => {
                    console.log('addIceCandidate Exception:' + e);
                });
            } else {
                queueCandidate[req_user_id].push(msg.candidate);
            }
        }
    }
} 