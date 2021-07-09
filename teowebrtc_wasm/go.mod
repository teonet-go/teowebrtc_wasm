module github.com/kirill-scherba/teowebrtc/teowebrtc_wasm

replace github.com/kirill-scherba/teowebrtc/teowebrtc_client => ../teowebrtc_client

replace github.com/kirill-scherba/teowebrtc/teowebrtc_signal_client => ../teowebrtc_signal_client

go 1.16

require (
	github.com/kirill-scherba/teowebrtc/teowebrtc_client v0.0.0-00010101000000-000000000000
	github.com/kirill-scherba/teowebrtc/teowebrtc_signal_client v0.0.0-00010101000000-000000000000 // indirect
	github.com/shurcooL/go v0.0.0-20200502201357-93f07166e636 // indirect
	github.com/shurcooL/go-goon v0.0.0-20210110234559-7585751d9a17 // indirect
	github.com/shurcooL/goexec v0.0.0-20200425235707-36ff6d2d1adc // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/tools v0.1.4 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)
