package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"sync"
	"time"
)

func init() {
	PCK[501] = &S2CUserBag{}
	PCK[520] = &S2CBagChange{}
}

var UserBag = &userBagBox{
	store: &sync.Map{},
}

var UserBagFirstWait = make(chan struct{})

var bagThread = make(chan *S2CUserBag)
var bagAction = make(chan struct{})

var waitUserBagItems = make([]*waitBagItemValue, 0)

// Bag 背包
func (c *Connect) Bag() {
	go func() {
		bagAction <- struct{}{}
		t := time.NewTimer(time.Minute)
		for range t.C {
			bagAction <- struct{}{}
			t.Reset(RandMillisecond(60, 120))
		}
	}()
	for {
		select {
		case <-bagAction:
			_ = c.userBag()
		case ret := <-bagThread:
			for _, item := range ret.Bag.Items {
				UserBag.Set(item.IId, item)
				for i := range waitUserBagItems {
					w := waitUserBagItems[i]
					if w.id == 0 && w.run == false {
						continue
					}
					if w.id == item.IId {
						w.C <- struct{}{}
						_ = w.Close()
					}
				}
			}
			// 是否有等待中的item
			for _, w := range waitUserBagItems {
				if w.id != 0 || w.run {
					go func() {
						bagAction <- struct{}{}
					}()
					break
				}
			}
			//
			UserBagFirstWait <- struct{}{}
		}
	}
}

// userBag 背包
func (c *Connect) userBag() error {
	body, err := proto.Marshal(&C2SUserBag{})
	if err != nil {
		return err
	}
	return c.send(500, body)
}

func (x *S2CUserBag) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Println("[S][CUserBag] success")
	bagThread <- x
	return
}

func (x *S2CBagChange) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][BagChange] %v", x)
	return
}

type userBagBox struct {
	store *sync.Map
}

func (u *userBagBox) Set(id int32, data *ItemData) {
	u.store.Store(id, data)
}

func (u *userBagBox) Has(id int32) bool {
	_, ok := u.store.Load(id)
	return ok
}

func (u *userBagBox) Get(id int32, timeout time.Duration) (*ItemData, bool) {
	ret, ok := u.store.Load(id)
	if !ok {
		if timeout == 0 {
			return nil, false
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		w := WaitUserBagItemValue(id)
		select {
		case <-w.C:
			{
				ret, ok = u.store.Load(id)
				if ok {
					return ret.(*ItemData), ok
				}
				return nil, false
			}
		case <-ctx.Done():
			_ = w.Close()
			return nil, false
		}
	}
	return ret.(*ItemData), ok
}

type waitBagItemValue struct {
	id  int32
	C   chan struct{}
	run bool
}

func (w *waitBagItemValue) Close() error {
	w.id = 0
	w.run = false
	return nil
}

func (w *waitBagItemValue) Open(id int32) {
	w.id = id
	w.run = true
}

func WaitUserBagItemValue(id int32) *waitBagItemValue {
	for i, w := range waitUserBagItems {
		if w.id == 0 && w.run == false {
			w.Open(id)
			return waitUserBagItems[i]
		}
	}
	w := waitBagItemValue{
		id:  id,
		C:   make(chan struct{}),
		run: true,
	}
	waitUserBagItems = append(waitUserBagItems, &w)
	bagAction <- struct{}{}
	return &w
}
