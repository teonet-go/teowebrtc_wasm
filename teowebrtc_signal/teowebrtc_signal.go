// Copyright 2021 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Webretc signal server (for teonet network)
package teowebrtc_signal

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// New create new teonet webrtc signal server
func New(addr string) (sig *Signal) {
	sig = new(Signal)
	sig.peers = make(map[string]*websocket.Conn)
	sig.serveWS(addr)
	return
}

type Signal struct {
	peers map[string]*websocket.Conn
	sync.RWMutex
}

type loginSignal struct {
	Signal string `json:"signal"`
	Login  string `json:"login"`
}

type signalSignal struct {
	Signal string      `json:"signal"`
	Peer   string      `json:"peer"`
	Data   interface{} `json:"data"`
}

func (sig *Signal) addPeer(address string, conn *websocket.Conn) {
	sig.Lock()
	defer sig.Unlock()
	sig.peers[address] = conn
}

func (sig *Signal) delPeer(address string) {
	sig.Lock()
	defer sig.Unlock()
	delete(sig.peers, address)
}

// getPeer get peer connection by address
func (sig *Signal) getPeer(address string) (conn *websocket.Conn, ok bool) {
	sig.RLock()
	defer sig.RUnlock()
	conn, ok = sig.peers[address]
	return
}

// address get peer address by connection
func (sig *Signal) address(conn *websocket.Conn) string {
	sig.RLock()
	defer sig.RUnlock()
	for addr, c := range sig.peers {
		if c == conn {
			return addr
		}
	}
	return ""
}

func (sig *Signal) serveWS(addr string) {
	http.HandleFunc("/signal", sig.signal)
	http.HandleFunc("/", sig.home)
	log.Fatal(http.ListenAndServe(addr, nil))
}

var upgrader = websocket.Upgrader{} // use default options

func (sig *Signal) signal(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade error:", err)
		return
	}
	defer c.Close()

	// Signals
	const (
		login     = "login"
		offer     = "offer"
		answer    = "answer"
		candidate = "candidate"
	)

	// This connection peer address
	var address string

	// writeErrMessage write json message with error
	writeErrMessage := func(messageType int, err error) {
		msg := []byte(fmt.Sprintf("{\"err\":\"%s\"}", err))
		log.Print("Error: ", err)
		c.WriteMessage(messageType, msg)
	}

	// Check connected
	connected := func(messageType int) bool {
		if address != "" {
			return true
		}
		err := errors.New("should login first")
		writeErrMessage(messageType, err)
		return false
	}

	// Befor close
	defer func() { log.Printf("Connection to %s closed", address) }()

	// Process messages
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Read message error:", err)
			return
		}

		// Parse json
		var j map[string]interface{}
		json.Unmarshal(message, &j)

		// Parse signal
		signal, ok := j["signal"]
		if !ok {
			err := errors.New("wrong command, signal must be present")
			writeErrMessage(mt, err)
			return
		}

		switch signal {

		case login:
			// Unmarshal signal
			var s loginSignal
			err = json.Unmarshal(message, &s)
			if err != nil {
				writeErrMessage(mt, err)
				return
			}

			// Save login
			address = s.Login

			log.Printf("Got %s from %s", signal, address)

			// Add peer to users map and remove it when connection closed
			sig.addPeer(s.Login, c)
			defer sig.delPeer(s.Login)

		case offer, answer, candidate:
			// Check login was done
			if !connected(mt) {
				return
			}

			log.Printf("Got %s from %s", signal, address)

			// Unmarshal signal
			var s signalSignal
			err = json.Unmarshal(message, &s)
			if err != nil {
				writeErrMessage(mt, err)
				return
			}

			// Resend signal to peer
			log.Printf("Resend %s to %s", s.Signal, s.Peer)
			conn, ok := sig.getPeer(s.Peer)
			if !ok {
				err := errors.New("peer does not connected")
				writeErrMessage(mt, err)
				continue
			}
			s.Peer = address
			message, err = json.Marshal(s)
			if err != nil {
				return
			}
			conn.WriteMessage(mt, message)
			continue
		}

		// Send answer
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write message error:", err)
			return
		}
	}
}

func (sig *Signal) home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/signal")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))
