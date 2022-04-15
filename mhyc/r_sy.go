package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// ShenYu 神域
func ShenYu(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	defer t1.Stop()
	f1 := func() time.Duration {
		time.Sleep(s10)
		Fight.Lock()
		defer func() {
			Receive.Action(CLI.LeaveShenYu)
			_ = Receive.Wait(&S2CLeaveShenYu{}, s3)
			Fight.Unlock()
		}()
		//
		Receive.Action(CLI.ShenYuData)
		_ = Receive.Wait(&S2CShenYuData{}, s3)
		//
		Receive.Action(CLI.EnterShenYu)
		var sy S2CEnterShenYu
		_ = Receive.Wait(&sy, s3)
		//

		//
		time.Sleep(s60)
		return TomorrowDuration(RandMillisecond(30000, 30600))
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

func (c *Connect) ShenYuData() error {
	body, err := proto.Marshal(&C2SShenYuData{})
	if err != nil {
		return err
	}
	log.Println("[C][ShenYuData]")
	return c.send(28270, body)
}

func (x *S2CShenYuData) ID() uint16 {
	return 28271
}

// Message S2CShenYuData 28271
func (x *S2CShenYuData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ShenYuData] tag=%v season=%v first_reward=%v", x.Tag, x.Season, x.FirstReward)
}

////////////////////////////////////////////////////////////

func (c *Connect) EnterShenYu() error {
	body, err := proto.Marshal(&C2SEnterShenYu{})
	if err != nil {
		return err
	}
	log.Println("[C][进入神域]")
	return c.send(51001, body)
}

func (x *S2CEnterShenYu) ID() uint16 {
	return 51002
}

// Message S2CEnterShenYu 51002
func (x *S2CEnterShenYu) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][进入神域] tag=%v %v", x.Tag, x)
}

////////////////////////////////////////////////////////////

func (c *Connect) LeaveShenYu() error {
	body, err := proto.Marshal(&C2SLeaveShenYu{})
	if err != nil {
		return err
	}
	log.Println("[C][离开神域]")
	return c.send(51003, body)
}

func (x *S2CLeaveShenYu) ID() uint16 {
	return 51004
}

// Message S2CLeaveShenYu 51004
func (x *S2CLeaveShenYu) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][离开神域] tag=%v", x.Tag)
}
