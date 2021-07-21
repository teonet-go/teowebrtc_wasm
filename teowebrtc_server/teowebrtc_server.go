// Copyright 2021 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Webrts server package
package teowebrtc_server

import (
	"encoding/json"
	"log"
	"time"

	"github.com/kirill-scherba/teowebrtc/teowebrtc_client"
	"github.com/kirill-scherba/teowebrtc/teowebrtc_signal_client"
	"github.com/pion/webrtc/v3"
)

func Connect(signalServerAddr, login string, connected func(peer string, dc *teowebrtc_client.DataChannel)) (err error) {

	// Create signal server client
	signal := teowebrtc_signal_client.New()

connect:
	// Connect to signal server
	err = signal.Connect(signalServerAddr, login)
	if err != nil {
		log.Println("Can't connect to signal server, error:", err)
		time.Sleep(5 * time.Second)
		goto connect
	}
	log.Println("Connected")

	var skipRead = false

	var sig teowebrtc_signal_client.Signal
	for {

		// Wait offer signal
		if !skipRead {
			sig, err = signal.WaitSignal()
			if err != nil {
				log.Println("can't wait offer, error:", err)
				goto connect
				// break
			}
		}
		skipRead = false

		// Unmarshal offer
		peer := sig.Peer
		var offer webrtc.SessionDescription
		json.Unmarshal(sig.Data, &offer)
		if err != nil {
			log.Println("can't unmarshal offer, error:", err)
			continue
		}
		log.Printf("got offer from %s", sig.Peer)

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
			continue
		}

		// Add handlers for setting up the connection.
		pc.OnSignalingStateChange(func(state webrtc.SignalingState) {
			log.Println("Signal changed:", state)
		})

		// Add handlers for setting up the connection.
		pc.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
			log.Printf("ICE Connection State has changed: %s\n", connectionState.String())
		})

		// Send AddICECandidate to remote peer
		pc.OnICECandidate(func(i *webrtc.ICECandidate) {
			if i != nil {
				candidateData, err := json.Marshal(i)
				if err != nil {
					log.Println("can't marshal ICECandidate, error:", err)
					return
				}
				signal.WriteCandidate(peer, candidateData)
			}
		})

		// Check ICEGathering state
		pc.OnICEGatheringStateChange(func(state webrtc.ICEGathererState) {
			switch state {
			case webrtc.ICEGathererStateGathering:
				log.Println("Collection of candidates has begin")

			case webrtc.ICEGathererStateComplete:
				log.Println("Collection of candidates is finished ")
				signal.WriteCandidate(peer, nil)
			}
		})

		pc.OnDataChannel(func(d *webrtc.DataChannel) {
			log.Printf("New DataChannel %s %d\n", d.Label(), d.ID())
			connected(peer, teowebrtc_client.NewDataChannel(d))
		})

		// Set the remote SessionDescription
		err = pc.SetRemoteDescription(offer)
		if err != nil {
			log.Print("SetRemoteDescription error: ", err)
			continue
		}

		// Initiates answer and set local SessionDescription
		answer, _ := pc.CreateAnswer(nil)
		err = pc.SetLocalDescription(answer)
		if err != nil {
			log.Print("SetLocalDescription error: ", err)
			continue
		}

		// Send answer to signal server
		answerData, err := json.Marshal(answer)
		if err != nil {
			log.Println("can't marshal answer, error:", err)
			continue
		}
		signal.WriteAnswer(peer, answerData)

		// Get client ICECandidate
		for {
			sig, err = signal.WaitSignal()
			if err != nil {
				break
			}

			// Unmarshal ICECandidate signal
			var i webrtc.ICECandidate
			if len(sig.Data) == 0 {
				log.Println("All ICECandidate processed")
				break
			}
			err = json.Unmarshal(sig.Data, &i)
			if err != nil {
				log.Println("can't unmarshal candidate, error:", err)
				skipRead = true
				break
			}
			log.Printf("Got ICECandidate from %s\n", sig.Peer)

			// Add servers ICECandidate
			err = pc.AddICECandidate(i.ToJSON())
			if err != nil {
				log.Println("can't add ICECandidate, error:", err)
			}
		}
	}

	// select {}
	// return
}
