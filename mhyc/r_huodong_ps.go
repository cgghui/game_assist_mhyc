package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// HuoDongBusiness 活动<跑商>
func HuoDongBusiness(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	defer t1.Stop()
	f1 := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "活动-跑商")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		Receive.Action(CLI.GetBusinessPrize)
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CGetBusinessPrize{}, s3); err != nil {
			return RandMillisecond(0, 2)
		}
		//
		Receive.Action(CLI.StartBusiness)
		var r S2CStartBusiness
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &r, s3); err != nil {
			return RandMillisecond(0, 2)
		}
		if r.Tag == 50401 {
			Receive.Action(CLI.ContinueBusiness)
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CContinueBusiness{}, s3); err != nil {
				return RandMillisecond(0, 2)
			}
			return RandMillisecond(50, 60)
		}
		if r.Tag == 50402 {
			return TomorrowDuration(RandMillisecond(600, 1800))
		}
		if r.Tag == 0 {
			ListenMessageCallEx(&S2CBusinessData{}, func(data []byte) bool {
				b := &S2CBusinessData{}
				b.Message(data)
				if b.Data.State == 1 {
					// 领取奖励
					Receive.Action(CLI.GetBusinessPrize)
					_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CGetBusinessPrize{}, s3)
				}
				return b.Data.State != 1 // false cancel thread
			})
		}
		return time.Hour
	}
	for {
		select {
		case <-t1.C:
			t1.Reset(f1())
		case <-ctx.Done():
			return
		}
	}
}

////////////////////////////////////////////////////////////

// ContinueBusiness 跑商重连
func (c *Connect) ContinueBusiness() error {
	body, err := proto.Marshal(&C2SContinueBusiness{})
	if err != nil {
		return err
	}
	log.Printf("[C][ContinueBusiness]")
	return c.send(23303, body)
}

func (x *S2CContinueBusiness) ID() uint16 {
	return 23304
}

// Message S2CContinueBusiness Code:23304
func (x *S2CContinueBusiness) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ContinueBusiness] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

// StartBusiness 跑商 动作
func (c *Connect) StartBusiness() error {
	body, err := proto.Marshal(&C2SStartBusiness{Id: 2})
	if err != nil {
		return err
	}
	log.Printf("[C][StartBusiness] id=2")
	return c.send(23305, body)
}

func (x *S2CStartBusiness) ID() uint16 {
	return 23306
}

// Message S2CStartBusiness Code:23306
func (x *S2CStartBusiness) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][StartBusiness] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CBusinessData) ID() uint16 {
	return 23301
}

// Message S2CBusinessData Code:23301
func (x *S2CBusinessData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][BusinessData] data=%v", x.Data)
}

////////////////////////////////////////////////////////////

func (c *Connect) GetBusinessPrize() error {
	body, err := proto.Marshal(&C2SGetBusinessPrize{})
	if err != nil {
		return err
	}
	log.Printf("[C][GetBusinessPrize]")
	return c.send(23311, body)
}

func (x *S2CGetBusinessPrize) ID() uint16 {
	return 23312
}

// Message S2CGetBusinessPrize Code:23312
func (x *S2CGetBusinessPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetBusinessPrize] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}
