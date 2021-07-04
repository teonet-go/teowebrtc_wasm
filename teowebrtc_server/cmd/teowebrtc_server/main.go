// Copyright 2021 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Webrts server sample application
package main

import (
	"flag"
	"log"

	"github.com/kirill-scherba/teowebrtc/teowebrtc_server"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	err := teowebrtc_server.Connect(*addr, "server-1", func(peer string) {
		log.Println("Connected to", peer)
	})
	if err != nil {
		log.Fatalln("connect error:", err)
	}

	select {}
}
