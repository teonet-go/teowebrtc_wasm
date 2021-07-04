// Copyright 2021 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Webretc signal server test client application (for teonet network)
package main

import (
	"flag"
	"log"

	"github.com/kirill-scherba/teowebrtc/teowebrtc_signal_client"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)
	signal := teowebrtc_signal_client.New()

	// Connect to signal server
	err := signal.Connect(*addr)
	if err != nil {
		log.Println("can't connect to signal server")
		return
	}
	log.Println()

	// Send offer signal to webrtc server
	var offer = []byte("offer data ...")
	answer, err := signal.WriteOffer("server-1", offer)
	if err != nil {
		log.Println("can't sent offer to signal server")
		return
	}
	log.Println("got answer:", string(answer))
	log.Println()

	// Send answer signal to webrtc client
	answer = []byte("answer data ...")
	res, err := signal.WriteAnswer("client-1", answer)
	if err != nil {
		log.Println("can't sent answer to signal server")
		return
	}
	log.Println("got responce:", string(res))
	log.Println()

	// Send candidate signal to webrtc peer
	candidate := []byte("candidate data ...")
	res, err = signal.WriteCandidate("server-1", candidate)
	if err != nil {
		log.Println("can't sent candidate to signal server")
		return
	}
	log.Println("got responce:", string(res))
	log.Println()

	select {}
}
