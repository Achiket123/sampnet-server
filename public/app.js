const signalingServer = new WebSocket("ws://localhost:8000/api/v1/call/2"); // Room ID: 1

let localStream;
let remoteStream = new MediaStream();
let peerConnection;
const configuration = {
  iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
};
  
// Initialize WebRTC
async function startWebRTC() {
  localStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
  localVideo.srcObject = localStream;

  peerConnection = new RTCPeerConnection(configuration);

  // Add local tracks to the connection
  localStream.getTracks().forEach((track) => peerConnection.addTrack(track, localStream));

  // Handle remote tracks
  peerConnection.ontrack = (event) => {
    remoteStream.addTrack(event.track);
    remoteVideo.srcObject = remoteStream;
    
  };

  // Handle ICE candidates
  peerConnection.onicecandidate = (event) => {
    if (event.candidate) {
      signalingServer.send(JSON.stringify({ type: "candidate", candidate: event.candidate }));
    }
  };

  signalingServer.onmessage = async (event) => {
    const message = JSON.parse(event.data);

    switch (message.type) {
      case "offer":
        await peerConnection.setRemoteDescription(new RTCSessionDescription(message.offer));
        const answer = await peerConnection.createAnswer();
        await peerConnection.setLocalDescription(answer);
        signalingServer.send(JSON.stringify({ type: "answer", answer }));
        break;

      case "answer":
        await peerConnection.setRemoteDescription(new RTCSessionDescription(message.answer));
        break;

      case "candidate":
        if (message.candidate) {
          await peerConnection.addIceCandidate(new RTCIceCandidate(message.candidate));
        }
        break;

      default:
        console.error("Unknown message type:", message.type);
    }
  };

  // Create an offer if you are the first participant
  const offer = await peerConnection.createOffer();
  await peerConnection.setLocalDescription(offer);
  signalingServer.send(JSON.stringify({ type: "offer", offer }));
}

// Start the WebRTC process
startCall.addEventListener("click", startWebRTC);
