package teowebrtc_client

import (
	"errors"
	"sync"
)

// NewCmdType create new CmdType
func NewCmdType() *CmdType {
	return new(CmdType)
}

// CmdType command data structure
type CmdType struct {
	Cmd  byte
	Data []byte
}

func (c CmdType) MarshalBinary() (data []byte, err error) {
	data = make([]byte, len(c.Data)+1)
	data[0] = c.Cmd
	if len(c.Data) > 0 {
		copy(data[1:], c.Data)
	}
	return
}

func (c *CmdType) UnmarshalBinary(data []byte) (err error) {
	if len(data) < 1 {
		err = errors.New("to low packet size")
		return
	}
	c.Cmd = data[0]
	c.Data = nil
	if len(data) > 1 {
		c.Data = data[1:]
	}
	return
}

func NewSubscribe() (s *SubscrType) {
	s = new(SubscrType)
	s.m = make(map[uint]func(data []byte) bool)
	return
}

type SubscrType struct {
	id uint
	m  map[uint]func(data []byte) bool
	sync.RWMutex
}

func (s *SubscrType) Add(f func(data []byte) bool) (id uint) {
	s.Lock()
	defer s.Unlock()
	s.id++
	s.m[s.id] = f
	return s.id
}

func (s *SubscrType) Del(id uint) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, id)
}

func (s *SubscrType) Process(data []byte) bool {
	s.RLock()
	defer s.RUnlock()
	for _, f := range s.m {
		if f(data) {
			return true
		}
	}
	return false
}
