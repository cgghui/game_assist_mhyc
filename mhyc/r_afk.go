package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// AFK 挂机
func AFK(ctx context.Context) {
	// 定时领取有尝奖励
	go func() {
		info := &S2CAFKGetBuyInfo{}
		buyTimes := &S2CAFKBuyTimes{}
		t := time.NewTimer(ms100)
		defer t.Stop()
		f := func() time.Duration {
			Fight.Lock()
			defer Fight.Unlock()
			for {
				Receive.Action(CLI.AFKGetBuyInfo)
				if err := Receive.Wait(info, s3); err != nil {
					continue
				}
				if CloseConn {
					return s3
				}
				if info.Coin <= 0 {
					Receive.Action(CLI.AFKBuyTimes)
					_ = Receive.Wait(buyTimes, s3)
					continue
				}
				return TomorrowDuration(RandMillisecond(30000, 30600))
			}
		}
		for {
			select {
			case <-t.C:
				t.Reset(f())
			case <-ctx.Done():
				return
			}
		}
	}()
	// 定时领取挂机奖励
	info := &S2CGetAFKPrize{}
	t := time.NewTimer(ms100)
	for range t.C {
		Receive.Action(CLI.GetAFKPrize)
		_ = Receive.Wait(info, s3)
		t.Reset(RandMillisecond(60, 180)) // 1 ~ 3 分钟
	}
}

func (c *Connect) AFKGetBuyInfo() error {
	body, err := proto.Marshal(&C2SAFKGetBuyInfo{})
	if err != nil {
		return err
	}
	log.Println("[C][AFKGetBuyInfo]")
	return c.send(22151, body)
}

// GetAFKPrize 挂机收益
func (c *Connect) GetAFKPrize() error {
	body, err := proto.Marshal(&C2SGetAFKPrize{})
	if err != nil {
		return err
	}
	log.Println("[C][GetAFKPrize]")
	return c.send(22155, body)
}

// AFKBuyTimes 通过购买获取挂机奖励
func (c *Connect) AFKBuyTimes() error {
	body, err := proto.Marshal(&C2SAFKBuyTimes{})
	if err != nil {
		return err
	}
	log.Println("[C][AFKBuyTimes]")
	return c.send(22153, body)
}

////////////////////////////////////////////////////////////

func (x *S2CAFKGetBuyInfo) ID() uint16 {
	return 22152
}

func (x *S2CAFKGetBuyInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][AFKGetBuyInfo] %v", x)
}

////////////////////////////////////////////////////////////

func (x *S2CGetAFKPrize) ID() uint16 {
	return 22156
}

func (x *S2CGetAFKPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetAFKPrize] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CAFKBuyTimes) ID() uint16 {
	return 22154
}

func (x *S2CAFKBuyTimes) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][AFKBuyTimes] tag=%v", x.Tag)
}
