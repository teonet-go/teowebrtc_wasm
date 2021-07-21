module github.com/kirill-scherba/teowebrtc/teowebrtc_server

replace github.com/kirill-scherba/teowebrtc/teowebrtc_client => ../teowebrtc_client

go 1.16

require (
	github.com/kirill-scherba/teowebrtc/teowebrtc_client v0.0.0-20210721095130-1e54d8193589
	github.com/kirill-scherba/teowebrtc/teowebrtc_signal_client v0.0.0-20210721095130-1e54d8193589
	github.com/pion/webrtc/v3 v3.0.31
)
