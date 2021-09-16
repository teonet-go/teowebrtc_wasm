module github.com/kirill-scherba/teowebrtc/teowebrtc_server

// replace github.com/kirill-scherba/teowebrtc/teowebrtc_client => /home/kirill/go/src/github.com/kirill-scherba/teowebrtc/teowebrtc_client
// replace github.com/kirill-scherba/teowebrtc/teowebrtc_signal_client => /home/kirill/go/src/github.com/kirill-scherba/teowebrtc/teowebrtc_signal_client

go 1.16

require (
	github.com/kirill-scherba/teowebrtc/teowebrtc_client v0.0.0-20210721104603-115f62a59ac6
	github.com/kirill-scherba/teowebrtc/teowebrtc_signal_client v0.0.0-20210721104603-115f62a59ac6
	github.com/pion/webrtc/v3 v3.0.31
)
