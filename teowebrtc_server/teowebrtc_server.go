// Copyright 2021 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Webrts server package
package teowebrtc_server

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kirill-scherba/teowebrtc/teowebrtc_signal_client"
	"github.com/pion/webrtc/v3"
)

func Connect(signalServerAddr, login string, connected func(peer string)) (err error) {

	// Create signal server client
	signal := teowebrtc_signal_client.New()

	// Connect to signal server
	err = signal.Connect(signalServerAddr, login)
	if err != nil {
		log.Println("can't connect to signal server, error:", err)
		return
	}
	log.Println()

	var skipRead = false

	var sig teowebrtc_signal_client.Signal
	for {

		// Wait offer signal
		if !skipRead {
			sig, err = signal.WaitOffer()
			if err != nil {
				log.Println("can't wait offer, error:", err)
				break
			}
		}
		skipRead = false

		// Unmarshal offer
		peer := sig.Peer
		sigData, _ := json.Marshal(sig.Data)
		var offer webrtc.SessionDescription
		json.Unmarshal(sigData, &offer)
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
			if connectionState.String() == "connected" {
				connected(peer)
			}
		})

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

		pc.OnDataChannel(func(d *webrtc.DataChannel) {
			log.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

			d.OnOpen(func() {
				log.Printf("Data channel opened")
			})

			// Register text message handling
			d.OnMessage(func(msg webrtc.DataChannelMessage) {
				log.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
				log.Println(d, msg)
				// Send echo answer
				data := []byte("Answer to: ")
				data = append(data, msg.Data...)
				d.Send(data)
			})

		})

		// Set the remote SessionDescription
		pc.SetRemoteDescription(offer)

		// Initiates answer and set local SessionDescription
		answer, _ := pc.CreateAnswer(nil)
		err = pc.SetLocalDescription(answer)
		if err != nil {
			fmt.Print("SetLocalDescription: ")
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
			sig, err = signal.WaitCandidate()
			if err != nil {
				break
			}

			// Unmarshal ICECandidate signal
			var i webrtc.ICECandidate
			sigData, _ := json.Marshal(sig.Data)
			err = json.Unmarshal(sigData, &i)
			if err != nil {
				log.Println("can't unmarshal candidate, error:", err)
				skipRead = true
				break
			}
			log.Printf("Got ICECandidatecandidate from %s", sig.Peer)

			// Add servers ICECandidate
			err = pc.AddICECandidate(i.ToJSON())
			if err != nil {
				log.Println("can't add ICECandidate, error:", err)
			}
		}
	}

	select {}

	return
}
