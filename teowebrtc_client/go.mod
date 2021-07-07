module github.com/kirill-scherba/teowebrtc/teowebrtc_client

replace github.com/kirill-scherba/teowebrtc/teowebrtc_signal_client => ../teowebrtc_signal_client

go 1.16

require (
	github.com/kirill-scherba/teowebrtc/teowebrtc_signal_client v0.0.0-00010101000000-000000000000
	github.com/pion/webrtc/v3 v3.0.30
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)
