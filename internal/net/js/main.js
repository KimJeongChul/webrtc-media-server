let server_addr = 'ws://localhost:8080/ws';
var mediaserver = new WebSocket(server_addr)
var media_room_id = null;
var user_id = null;
var handle_id = null;

mediaserver.onopen = () => {
    console.log('Media server ws open');
}
mediaserver.onclose = () => {
    console.log('Media server ws closed');
};
mediaserver.onmessage = (msg) => {
    const parsedMsg = JSON.parse(msg.data);
    console.log("RECV msg : ", parsedMsg);
    switch(parsedMsg.method) {
        case "resCreateRoom": 
        {
            let inputMediaRoomID = document.getElementById('media_room_id');
            inputMediaRoomID.value = parsedMsg.roomID;
            media_room_id = parsedMsg.roomID;
        }
            break;
        case "resCandidate":
        {
            setCandidate(parsedMsg);
        }
            break;
        case "resAnswerSDP":
        {
            setSDP(parsedMsg);
        }
            break;
    }
}

//createHandle Create user and handle id
function createHandle() {
    user_id = uuidv4(); // Group
    handle_id = makeid(8);// UserHandle
    let inputUserID = document.getElementById('user_id');
    inputUserID.value = user_id;
    let inputHandleID = document.getElementById('handle_id');
    inputHandleID.value = handle_id;
}

//makeCall Make a Call
function makeCall() {
    if(!checkMediaRoomID()){
        alert("Please check media room id");
        return
    }

    createHandle();

    addRTCSession();
}

//releaseRoom Request create room
function releaseRoom() {
    const msgReleaseRoom = {
        method: 'releaseRoom'
    };
    sendMsg(msgReleaseRoom);
}

//createRoom Request create room
function createRoom() {
    const msgCreateRoom = {
        method: 'createRoom'
    };
    sendMsg(msgCreateRoom);
}

//checkMediaRoomID 
function checkMediaRoomID() {
    let inputMediaRoomID = document.getElementById('media_room_id');
    console.log(inputMediaRoomID.value);
    if(inputMediaRoomID.value === "") {
        return false;
    }
    return true;
}

//sendMsg Send message websocket
function sendMsg(msg) {
    console.log("SEND msg : ", JSON.stringify(msg));
    mediaserver.send(JSON.stringify(msg));
}