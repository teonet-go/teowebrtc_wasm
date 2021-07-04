// Copyright 2021 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Webrtc signal server
package main

import (
	"flag"
	"log"

	"github.com/kirill-scherba/teowebrtc/teowebrtc_signal"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)
	teowebrtc_signal.New(*addr)
}
