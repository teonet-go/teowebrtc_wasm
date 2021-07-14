# Teonet webrtc

This project contain webrtc signal server, webrtc server and client packages
and applictions. It organize transfer data packages between webrtc server and
clients.

This code is currently under development and we do not recommend using it in
production environment.

## Packages description

`teowebrtc_signal` - signal server packet and application. It start websocket
server, wait client connection and resend webrtc signals between webrtc clients
and server

`teowebrtc_signal_client` - signal server client, used to establish webrtc
connection

`teowebrtc_server` - webrtc server package and sample application to connect
clients

`teowebrtc_client` - webrtc client package and sample application which
connected to webrtc server

`teowebrtc_wasm` - webrtc web (wasm) sample application which connected to
webrtc server from browser.
