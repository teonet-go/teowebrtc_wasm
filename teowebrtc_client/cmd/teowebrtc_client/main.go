// Copyright 2021 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Webrts client sample application
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/kirill-scherba/teowebrtc/teowebrtc_client"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var name = flag.String("name", "client-1", "client name")
var server = flag.String("server", "server-1", "server name")

func main() {
	flag.Parse()
	log.SetFlags(0)

	var id = 0
connect:
	// Connect to teo webrtc application (server)
	err := teowebrtc_client.Connect(*addr, *name, *server, func(peer string, d *teowebrtc_client.DataChannel) {
		log.Println("Connected to", peer)
		var connected = true

		// On open Send messages to created data channel
		d.OnOpen(func() {
			for connected {
				id++
				msg := fmt.Sprintf("Hello from %s with id %d!", *name, id)
				err := d.Send([]byte(msg))
				if err != nil {
					log.Printf("Send error: %s\n", err)
					continue
				}
				log.Printf("Send: %s", msg)
				time.Sleep(5 * time.Second)
			}
		})

		d.OnClose(func() {
			log.Println("Connection closed")
			connected = false
		})

		d.OnMessage(func(data []byte) {
			log.Printf("Got: %s", data)
		})
	})
	if err != nil {
		log.Println("connect error:", err)
	}

	// Reconnect after five seconds
	time.Sleep(5 * time.Second)
	goto connect
}
