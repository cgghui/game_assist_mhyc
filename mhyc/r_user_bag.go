package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"sync"
	"time"
)

// UserBag 背包
func (c *Connect) UserBag() error {
	body, err := proto.Marshal(&C2SUserBag{})
	if err != nil {
		return err
	}
	log.Println("[C][UserBag]")
	return c.send(500, body)
}

////////////////////////////////////////////////////////////

func (x *S2CUserBag) ID() uint16 {
	return 501
}

func (x *S2CUserBag) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("[S][UserBag] err=%v", err)
		return
	}
	if x.Bag.Items != nil {
		for _, item := range x.Bag.Items {
			UserBag.Set(item.IId, item)
		}
	}
}

////////////////////////////////////////////////////////////

func (x *S2CBagChange) ID() uint16 {
	return 520
}

func (x *S2CBagChange) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		log.Printf("[S][BagChange] err=%v", err)
		return
	}
	for _, item := range x.Change {
		UserBag.Set(item.Item.IId, item.Item)
	}
}

////////////////////////////////////////////////////////////

func (x *ItemFly) ID() uint16 {
	return 524
}

// Message ItemFly 524
func (x *ItemFly) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	for _, item := range x.Item {
		var dat *ItemData
		if val := UserBag.Get(item.IId); val != nil {
			dat = val
		} else {
			dat = &ItemData{}
		}
		dat.N = dat.N + item.N
		UserBag.Set(item.IId, dat)
	}
	log.Printf("[S][ItemFly] %v", x)
	return
}

////////////////////////////////////////////////////////////

var UserBag = &userBag{
	s: &sync.Map{},
}

type userBag struct{ s *sync.Map }

func (u *userBag) Set(id int32, data *ItemData) {
	u.s.Store(id, data)
}

func (u *userBag) Has(id int32) bool {
	_, ok := u.s.Load(id)
	return ok
}

func (u *userBag) GetAll() (ret map[string]interface{}) {
	ret = make(map[string]interface{})
	u.s.Range(func(key, value interface{}) bool {
		ret[key.(string)] = value
		return true
	})
	return
}

func (u *userBag) Get(id int32) *ItemData {
	ret, ok := u.s.Load(id)
	if !ok {
		return nil
	}
	return ret.(*ItemData)
}

const ts0 = 0
const ms10 = 10 * time.Millisecond
const ms100 = 100 * time.Millisecond
const ms500 = 500 * time.Millisecond
const s3 = 3 * time.Second
const s6 = 6 * time.Second
const s10 = 10 * time.Second
const s20 = 20 * time.Second
const s30 = 30 * time.Second
const s60 = 60 * time.Second
const s90 = 90 * time.Second

func (u *userBag) Wait(id int32, timeout time.Duration, c context.Context) *ItemData {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	tm := time.NewTimer(ms100)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-c.Done():
			return nil
		case <-tm.C:
			_ = CLI.UserBag()
			if err = Receive.WaitWithContext(ctx, &S2CUserBag{}); err != nil {
				return nil
			}
			ret, ok := u.s.Load(id)
			if !ok {
				tm.Reset(ms100)
				continue
			}
			return ret.(*ItemData)
		}
	}
}
