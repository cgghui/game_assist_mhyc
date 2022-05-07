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
		t := time.NewTimer(ms100)
		defer t.Stop()
		f := func() time.Duration {
			Fight.Lock()
			am := SetAction(ctx, "挂机-有尝奖励")
			defer func() {
				am.End()
				Fight.Unlock()
			}()
			return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
				Receive.Action(CLI.AFKGetBuyInfo)
				info := &S2CAFKGetBuyInfo{}
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, info, s3); err != nil {
					loop = 0
					next = RandMillisecond(6, 15)
					return
				}
				if info.Coin <= 0 {
					Receive.Action(CLI.AFKBuyTimes)
					buyTimes := &S2CAFKBuyTimes{}
					if err := Receive.WaitWithContextOrTimeout(am.Ctx, buyTimes, s3); err != nil {
						loop = 0
						next = RandMillisecond(6, 15)
						return
					}
					return ms100, 0
				}
				return 0, TomorrowDuration(RandMillisecond(1800, 3600))
			})
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
	defer t.Stop()
	f := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "挂机-定时领取挂机奖励")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			Receive.Action(CLI.GetAFKPrize)
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, info, s3)
			return 0, RandMillisecond(600, 1200)
		})
	}
	for {
		select {
		case <-t.C:
			t.Reset(f())
		case <-ctx.Done():
			return
		}
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
