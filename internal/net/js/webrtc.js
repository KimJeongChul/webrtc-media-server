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
            console.log(event.candidate);

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
    }).catch((err) => {
        console.log(err.stack);
    });
}