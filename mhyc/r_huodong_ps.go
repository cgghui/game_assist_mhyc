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
		Receive.Action(CLI.StartBusiness)
		var r S2CStartBusiness
		_ = Receive.Wait(&r, s3)
		if r.Tag == 50401 {
			Receive.Action(CLI.ContinueBusiness)
			_ = Receive.Wait(&S2CContinueBusiness{}, s3)
			return time.Minute
		}
		if r.Tag == 50402 {
			TomorrowDuration(RandMillisecond(600, 1800))
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
	log.Printf("[S][ContinueBusiness] tag=%v", x.Tag)
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
	log.Printf("[S][StartBusiness] tag=%v", x.Tag)
}
