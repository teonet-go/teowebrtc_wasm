// Copyright 2021 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Webrts client package
package teowebrtc_client

import (
	"encoding/json"
	"log"

	"github.com/kirill-scherba/teowebrtc/teowebrtc_signal_client"
	"github.com/pion/webrtc/v3"
)

func Connect(signalServerAddr, login, server string, connected func(peer string, dc *DataChannel)) (err error) {

	// Create signal server client
	signal := teowebrtc_signal_client.New()

	// Connect to signal server
	err = signal.Connect(signalServerAddr, login)
	if err != nil {
		log.Println("can't connect to signal server")
		return
	}
	log.Println()

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return
	}

	// Create DataChannel
	dc, err := pc.CreateDataChannel("teo", nil)
	if err != nil {
		return
	}

	pc.OnSignalingStateChange(func(state webrtc.SignalingState) {
		log.Println("Signal changed:", state)
	})

	// Add handlers for setting up the connection.
	pc.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("ICE Connection State has changed: %s\n", connectionState.String())
		if connectionState.String() == "connected" {
			connected(server, &DataChannel{dc})
		}
	})

	// Initiates the offer.
	offer, _ := pc.CreateOffer(nil)

	// Send offer and get answer
	offerData, err := json.Marshal(offer)
	if err != nil {
		return
	}
	message, err := signal.WriteOffer(server, offerData)
	if err != nil {
		return
	}

	// Unmarshal answer
	var errMsg = "can't unmarshal answer, error:"
	var sig teowebrtc_signal_client.Signal
	err = json.Unmarshal(message, &sig)
	if err != nil {
		log.Println(errMsg, err)
		return
	}
	peer := sig.Peer
	sigData, _ := json.Marshal(sig.Data)
	var answer webrtc.SessionDescription
	err = json.Unmarshal(sigData, &answer)
	if err != nil {
		log.Println(errMsg, err)
		return
	}
	log.Printf("Got answer from %s", sig.Peer)

	// Send AddICECandidate to remote peer
	pc.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i != nil {
			// check(offerPC.AddICECandidate(i.ToJSON()))
			log.Println("ICECandidate:", i)
			candidateData, err := json.Marshal(i)
			if err != nil {
				log.Panicln("can't marshal ICECandidate, error:", err)
				return
			}
			signal.WriteCandidate(peer, candidateData)
		}
	})

	// Set local SessionDescription
	err = pc.SetLocalDescription(offer)
	if err != nil {
		log.Println("SetLocalDescription error, err:", err)
		return
	}

	// Set remote SessionDescription
	err = pc.SetRemoteDescription(answer)
	if err != nil {
		log.Println("SetRemoteDescription error, err:", err)
		return
	}

	// Get servers ICECandidate
	for {
		sig, err := signal.WaitCandidate()
		if err != nil {
			break
		}

		// Unmarshal ICECandidate signal
		var i webrtc.ICECandidate
		sigData, _ := json.Marshal(sig.Data)
		err = json.Unmarshal(sigData, &i)
		if err != nil {
			log.Println("can't unmarshal candidate, error:", err)
			continue
		}
		log.Printf("Got ICECandidatecandidate from %s", sig.Peer)

		// Add servers ICECandidate
		err = pc.AddICECandidate(i.ToJSON())
		if err != nil {
			log.Println("can't add ICECandidate, error:", err)
		}
	}

	select {}

	return
}

type DataChannel struct {
	dc *webrtc.DataChannel
}

func (d *DataChannel) OnOpen(f func()) {
	d.dc.OnOpen(f)
}

func (d *DataChannel) Send(data []byte) error {
	return d.dc.Send(data)
}
