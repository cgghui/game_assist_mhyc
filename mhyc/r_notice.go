package mhyc

import (
	"context"
	"errors"
	"log"
	"time"
)

var Receive *receiveMessageBox

func init() {
	Init()
}

func Init() {
	Receive = &receiveMessageBox{
		ls: make([]*receiveMessage, 0),
	}
}

type receiveMessage struct {
	wait    chan []byte
	id      uint16
	running bool
	Call    HandleMessage
}

func (r *receiveMessage) Wait() <-chan []byte {
	return r.wait
}

func (r *receiveMessage) Close() {
	r.id = 0
	r.running = false
	r.Call = nil
}

func (r *receiveMessage) Open(id uint16, call HandleMessage) *receiveMessage {
	r.id = id
	r.running = true
	r.Call = call
	return r
}

// IsOpen open true, close false
func (r *receiveMessage) IsOpen() bool {
	return r.id > 0 && r.running
}

////////////////////////////////////////////////////////////////////////

var ErrReceiveMessageTimeout = errors.New("receive message timeout")
var ErrReceiveClose = errors.New("receive channel close")

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
			rm = w.Open(id, call)
			break
		}
	}
	if rm == nil {
		rm = &receiveMessage{
			wait:    make(chan []byte),
			id:      id,
			running: true,
			Call:    call,
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

func (r *receiveMessageBox) WaitWithContextOrTimeout(c context.Context, call HandleMessage, timeout ...time.Duration) error {
	if len(timeout) == 0 {
		return r.WaitWithContext(nil, call)
	}
	ctx, cancel := context.WithTimeout(c, timeout[0])
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
	case data, ok := <-rm.wait:
		if ok == false {
			return ErrReceiveClose
		}
		call.Message(data)
		return nil
	case <-ctx.Done():
		return ErrReceiveMessageTimeout
	}
}

func (r *receiveMessageBox) Close() {
	for _, w := range r.ls {
		close(w.wait)
		w.Close()
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

func ListenMessage(ctx context.Context, hm HandleMessage) {
	channel := Receive.CreateChannel(hm)
	defer channel.Close()
	if ctx == nil {
		for data := range channel.Wait() {
			channel.Call.Message(data)
		}
	}
	for {
		select {
		case data := <-channel.Wait():
			channel.Call.Message(data)
		case <-ctx.Done():
			return
		}
	}
}

func ListenMessageCall(ctx context.Context, hm HandleMessage, call func(data []byte)) {
	channel := Receive.CreateChannel(hm)
	defer channel.Close()
	if ctx == nil {
		for data := range channel.Wait() {
			call(data)
		}
	}
	for {
		select {
		case data := <-channel.Wait():
			call(data)
		case <-ctx.Done():
			return
		}
	}
}

// ListenMessageCallEx call 返回false退出监听
func ListenMessageCallEx(hm HandleMessage, call func(data []byte) bool) {
	channel := Receive.CreateChannel(hm)
	defer channel.Close()
	for data := range channel.Wait() {
		if !call(data) { // call is false return
			return
		}
	}
}

func ListenMessageNotify(call HandleMessage, timeout ...time.Duration) <-chan struct{} {
	notify := make(chan struct{})
	go func() {
		_ = Receive.Wait(call, timeout...)
		notify <- struct{}{}
		close(notify)
	}()
	return notify
}

func FightAction(ctx context.Context, Id, Type int64) (*S2CStartFight, *S2CBattlefieldReport) {
	c := make(chan *S2CStartFight)
	defer close(c)
	go func() {
		sf := &S2CStartFight{}
		if err := Receive.WaitWithContextOrTimeout(ctx, sf, s3); err != nil {
			c <- nil
		} else {
			c <- sf
		}
	}()
	go func() {
		_ = CLI.StartFight(&C2SStartFight{Id: Id, Type: Type})
	}()
	r := &S2CBattlefieldReport{}
	if err := Receive.WaitWithContextOrTimeout(ctx, r, s3); err != nil {
		r = nil
	}
	return <-c, r
}
