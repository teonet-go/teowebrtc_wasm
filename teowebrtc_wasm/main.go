// Example from: https://github.com/golang/go/wiki/WebAssembly#getting-started
//
/* Build:

GOOS=js GOARCH=wasm go build -o main.wasm
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .

# install goexec: go get -u github.com/shurcooL/goexec
goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'

*/
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kirill-scherba/teowebrtc/teowebrtc_client"
)

func main() {

	addr := "localhost:8080"
	name := "web-1"
	server := "server-1"

	err := teowebrtc_client.Connect(addr, name, server, func(peer string, d *teowebrtc_client.DataChannel) {
		log.Println("Connected to", peer)
		// Send messages to created data channel
		var id = 0
		d.OnOpen(func() {
			for {
				id++
				msg := fmt.Sprintf("Hello from %s with id %d!", name, id)
				d.Send([]byte(msg))
				log.Printf("Send: %s", msg)
				time.Sleep(5 * time.Second)
			}
		})
		d.OnMessage(func(data []byte) {
			log.Printf("Got: %s", data)
		})
	})
	if err != nil {
		log.Fatalln("connect error:", err)
	}

	select {}
}
