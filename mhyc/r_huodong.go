package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// HuoDong 活动
func HuoDong(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	defer t1.Stop()
	f1 := func(IllusionType int32) time.Duration {

		return RandMillisecond(300, 600)
	}

	for {
		select {
		case <-t1.C:
			t1.Reset(f1(2))
		case <-ctx.Done():
			return
		}
	}
}

////////////////////////////////////////////////////////////

// GetWestExp 护送 免费
func (c *Connect) GetWestExp() error {
	body, err := proto.Marshal(&C2SGetWestExp{GetType: 0})
	if err != nil {
		return err
	}
	log.Printf("[C][C2SGetWestExp] get_type=%v", 0)
	return c.send(463, body)
}

// GetWestExpRef 护送 用50元宝刷新
func (c *Connect) GetWestExpRef() error {
	body, err := proto.Marshal(&C2SGetWestExp{GetType: 1})
	if err != nil {
		return err
	}
	log.Printf("[C][C2SGetWestExp] get_type=%v", 1)
	return c.send(463, body)
}

func (x *S2CGetWestExp) ID() uint16 {
	return 464
}

// Message S2CGetWestExp Code:464
func (x *S2CGetWestExp) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetWestExp] %v", x)
}

////////////////////////////////////////////////////////////

// StartWestExp 护送动作
func (c *Connect) StartWestExp() error {
	body, err := proto.Marshal(&C2SStartWestExp{})
	if err != nil {
		return err
	}
	log.Printf("[C][StartWestExp]")
	return c.send(474, body)
}

func (x *S2CStartWestExp) ID() uint16 {
	return 475
}

// Message S2CStartWestExp Code:475
func (x *S2CStartWestExp) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][StartWestExp] tag=%v", x.Tag)
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
