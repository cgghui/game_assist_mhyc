package mhyc

import (
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func init() {
	PCK[22152] = &S2CAFKGetBuyInfo{}
	PCK[22154] = &S2CAFKBuyTimes{}
	PCK[22156] = &S2CGetAFKPrize{}
}

var afkThread = make(chan interface{})
var afkAction = make(chan struct{})

// AFK 挂机
func (c *Connect) AFK() {
	go func() {
		afkAction <- struct{}{}
		t := time.NewTimer(time.Minute)
		for range t.C {
			afkAction <- struct{}{}
			t.Reset(10 * time.Minute)
		}
	}()
	for range afkAction {
		for {
			_ = c.afkGetBuyInfo()
			if r := (<-afkThread).(*S2CAFKGetBuyInfo); r.Coin > 0 {
				break
			}
			_ = c.afkBuyTimes()
			<-afkThread
		}
		_ = c.getAFKPrize()
		<-afkThread
	}
}

func (c *Connect) afkGetBuyInfo() error {
	body, err := proto.Marshal(&C2SAFKGetBuyInfo{})
	if err != nil {
		return err
	}
	return c.send(22151, body)
}

// getAFKPrize 挂机收益
func (c *Connect) getAFKPrize() error {
	body, err := proto.Marshal(&C2SGetAFKPrize{})
	if err != nil {
		return err
	}
	return c.send(22155, body)
}

// afkBuyTimes 通过购买获取挂机奖励
func (c *Connect) afkBuyTimes() error {
	body, err := proto.Marshal(&C2SAFKBuyTimes{})
	if err != nil {
		return err
	}
	return c.send(22153, body)
}

func (x *S2CAFKGetBuyInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][AFKGetBuyInfo] %v", x)
	afkThread <- x
}

func (x *S2CGetAFKPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetAFKPrize] tag=%v", x.Tag)
	afkThread <- x
}

func (x *S2CAFKBuyTimes) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][AFKBuyTimes] tag=%v", x.Tag)
	afkThread <- x
}
