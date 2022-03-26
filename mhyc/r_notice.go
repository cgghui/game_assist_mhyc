package mhyc

import (
	"context"
	"errors"
	"log"
	"time"
)

var Receive *receiveMessageBox

func init() {
	Receive = &receiveMessageBox{
		ls: make([]*receiveMessage, 0),
	}
}

type receiveMessage struct {
	wait    chan []byte
	id      uint16
	running bool
}

func (r *receiveMessage) Wait() <-chan []byte {
	return r.wait
}

func (r *receiveMessage) Close() {
	r.id = 0
	r.running = false
}

func (r *receiveMessage) Open(id uint16) *receiveMessage {
	r.id = id
	r.running = true
	return r
}

// IsOpen open true, close false
func (r *receiveMessage) IsOpen() bool {
	return r.id > 0 && r.running
}

////////////////////////////////////////////////////////////////////////

var ErrReceiveMessageTimeout = errors.New("receive message timeout")

type receiveMessageBox struct {
	ls []*receiveMessage
}

func (r *receiveMessageBox) Action(act func() error) {
	go func(act func() error) {
		if err := act(); err != nil {
			log.Printf("Action Err: %v", err)
		}
	}(act)
}

func (r *receiveMessageBox) CreateChannel(call HandleMessage) *receiveMessage {
	id := call.ID()
	var rm *receiveMessage
	for _, w := range r.ls {
		if !w.IsOpen() {
			rm = w.Open(id)
			break
		}
	}
	if rm == nil {
		rm = &receiveMessage{
			wait:    make(chan []byte),
			id:      id,
			running: true,
		}
		r.ls = append(r.ls, rm)
	}
	return rm
}

func (r *receiveMessageBox) Wait(call HandleMessage, timeout ...time.Duration) error {
	if len(timeout) == 0 {
		return r.WaitWithContext(nil, call)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout[0])
	defer cancel()
	return r.WaitWithContext(ctx, call)
}

func (r *receiveMessageBox) WaitWithContext(ctx context.Context, call HandleMessage) error {
	rm := r.CreateChannel(call)
	defer rm.Close()
	if ctx == nil {
		call.Message(<-rm.wait)
		return nil
	}
	select {
	case data := <-rm.wait:
		call.Message(data)
		return nil
	case <-ctx.Done():
		return ErrReceiveMessageTimeout
	}
}

func (r *receiveMessageBox) Notify(id uint16, data []byte) {
	success := false
	for _, w := range r.ls {
		if w.running && w.id == id {
			w.wait <- data
			success = true
		}
	}
	if !success {
		log.Printf("receive: id[%d] manage func non-existent", id)
	}
}
