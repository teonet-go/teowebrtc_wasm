module github.com/kirill-scherba/teowebrtc/teowebrtc_server

replace github.com/kirill-scherba/teowebrtc/teowebrtc_signal_client => ../teowebrtc_signal_client

go 1.16

require (
	github.com/kirill-scherba/teowebrtc/teowebrtc_signal_client v0.0.0-00010101000000-000000000000
	github.com/pion/webrtc/v3 v3.0.30
)
