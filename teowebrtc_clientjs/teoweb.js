/**
 * Teonet web application client class
 * 
 * @param {Object} ws Websocket connection
 * @returns {Teocli.teowebAnonym$0}
 */
function Teoweb() {

console.log("Teoweb loaded");

var url = "ws://localhost:8080/signal";
var ws = new WebSocket(url);
var server = "server-1";

const configuration = { iceServers: [{ urls: "stun:stun.l.google.com:19302" }] };
const pc = new RTCPeerConnection(configuration);
const dc = pc.createDataChannel("teo");

// sendSignal send signal to signal server
sendSignal = function(signal) {
    ws.send(JSON.stringify(signal))
}

// Show signaling state
pc.onsignalingstatechange = function(ev) {
    console.log("signaling state:", pc.signalingState)
  };

// Show ice connection state
pc.oniceconnectionstatechange = function(ev) {
    console.log("ICE connection state:", pc.iceConnectionState)
};

// Send any ice candidates to the other peer.
pc.onicecandidate = function(ev) {
    if (ev.candidate) {
        console.log("send candidate", ev.candidate);
        sendSignal({ signal: "candidate", peer: server, data: ev.candidate });
    } else {
        /* there are no more candidates coming during this negotiation */
    }
};  

// Let the "negotiationneeded" event trigger offer generation.
pc.onnegotiationneeded = async () => {
  try {
    offer = await pc.createOffer();
    pc.setLocalDescription(offer);
    console.log("send offer");
    sendSignal({ signal: "offer", peer: server, data: offer })
  } catch (err) {
    console.error(err);
  }
};

ws.onopen = function(ev)  { 
    console.log("ws.onopen");
    login = "web-1"
    console.log("send login", login);
    sendSignal({ signal: "login", login: login });
}
ws.onerror = function(ev) { 
    //
}
ws.onclose = function(ev) { 
    //
}
ws.onmessage = function(ev) { 
    obj = JSON.parse(ev.data);
    if (obj['signal'] == undefined) {
        console.log("Wrong signal received")
        return
    }

    switch (obj['signal']) {
        
        case "login":
            console.log("got login answer");
            // console.log("send offer");
            // sendOffer(server);
            break;

        case "answer":
            console.log("got answer signal", obj.data);
            pc.setRemoteDescription(obj.data)
            break;

        case "candidate":
            console.log("got candidate signal", obj.data);
            pc.addIceCandidate(obj.data);
            // .then(
            //     function(){ console.log("ok"); },
            //     function(){ console.log("err"); }
            // );
            break;    

        // default:    
    }
}

return {

/**
  * Send login command signal server
  *  
  * @param {string} addr Name of this client
  * @returns {undefined}
  */  
// login: function (addr) {
//     ws.send('{ "signal": "login", "login": "' + addr + '" }');
// },

};
};