package teowebrtc_client

import (
	"errors"
	"log"
	"syscall/js"
	"time"
)

var setDataCallback js.Value
var dc *DataChannel

func SetDataChannel(d *DataChannel) {
	dc = d
}

// SetData call vuejs setDatacallback to show data in page
func SetData(data []byte) {
	if setDataCallback.Type().String() == "undefined" {
		return
	}
	setDataCallback.Invoke(string(data))
}

// SetFuncs create webasm js functions
func SetFuncs(subscr *SubscrType) {

	// SetCallback set vuejs callback function
	js.Global().Set("SetCallback", js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		setDataCallback = inputs[len(inputs)-1:][0]
		setDataCallback.Invoke("Hello! Wait data received...")
		return js.Undefined()
	}))

	// Send to server
	js.Global().Set("Send", js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		if dc == nil {
			return errors.New("not initialised")
		}
		data := []byte(inputs[0].String())
		err := dc.Send(data)
		return err
	}))

	// SendCmd send command to server {cmd: byte, data: []byte, callback func(err, data)}
	js.Global().Set("SendCmd", js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		log.Println("SendCmd, inputs", inputs)
		var err error
		if dc == nil {
			err = errors.New("not initialised")
			return err
		}

		// Parse parameters
		c := NewCmdType()
		c.Cmd = byte(inputs[0].Int())
		c.Data = []byte(inputs[1].String())
		callback := inputs[len(inputs)-1:][0]
		data, err := c.MarshalBinary()
		if err != nil {
			return err
		}

		// Send data
		err = dc.Send(data)
		if err != nil {
			return err
		}

		// Process answer
		var id uint
		type waitType struct {
			err  error
			cmd  byte
			data []byte
		}
		var wait = make(chan waitType)
		id = subscr.Add(func(data []byte) (processed bool) {
			log.Println("Check data:", data)
			// "Got: Answer to: " // 11
			log.Println("process data + 11", data[11:])
			a := NewCmdType()
			err = a.UnmarshalBinary(data[11:])
			if err != nil {
				log.Println("unmarshal binary error:", err)
				return
			}
			if a.Cmd != c.Cmd {
				log.Println("skip cmd:", a.Cmd)
				return
			}
			log.Println("Got answer to cmd", a.Cmd, a.Data)
			wait <- waitType{nil, a.Cmd, a.Data}
			processed = true
			return
		})

		// Wait answer
		go func() {
			var res waitType
			select {
			case res = <-wait:
			case <-time.After(5 * time.Second):
				res.err = errors.New("timeout")
			}
			close(wait)
			subscr.Del(id)
			if res.err != nil {
				callback.Invoke(res.err, js.Null())
			} else {
				callback.Invoke(js.Null(), string(res.data))
			}
		}()
		return err
	}))
}
