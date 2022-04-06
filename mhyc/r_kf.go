package mhyc

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func illusionSweep() time.Duration {
	Fight.Lock()
	defer Fight.Unlock()
	for _, ir := range []int32{1, 2, 11} {
		go func(ir int32) {
			_ = CLI.IllusionSweep(&C2SIllusionSweep{IllusionType: ir})
		}(ir)
		r := &S2CIllusionSweep{}
		_ = Receive.Wait(r)
		fmt.Println(r)
	}
	return TomorrowDuration(RandMillisecond(30000, 30600))
}

func yiJi() time.Duration {
	Receive.Action(CLI.YiJiInfo)
	var info S2CYiJiInfo
	_ = Receive.Wait(&info, s3)
	if len(info.Items) == 0 {
		return time.Second
	}
	id := int32(0)
	for _, item := range info.Items {
		if item.BossState == 1 {
			id = item.Id
			break
		}
	}
	// 进入
	go func() {
		_ = CLI.YiJiJoinScene(id)
	}()
	_ = Receive.Wait(&S2CYiJiJoinScene{}, s3)

	return TomorrowDuration(RandMillisecond(30000, 30600))
}

// KuaFu 跨服
func KuaFu(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	t4 := time.NewTimer(ms10)
	defer t1.Stop()
	defer t4.Stop()
	for {
		select {
		case <-t1.C:
			t1.Reset(yiJi())
		case <-t4.C:
			t4.Reset(illusionSweep())
		case <-ctx.Done():
			return
		}
	}
}

////////////////////////////////////////////////////////////

// YiJiJoinScene 进入遗迹场景
func (c *Connect) YiJiJoinScene(id int32) error {
	body, err := proto.Marshal(&C2SYiJiJoinScene{Id: id})
	if err != nil {
		return err
	}
	log.Printf("[C][进入遗迹场景] id=%v", id)
	return c.send(25315, body)
}

func (x *S2CYiJiJoinScene) ID() uint16 {
	return 25316
}

// Message S2CYiJiJoinScene Code:25316
func (x *S2CYiJiJoinScene) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][进入遗迹场景] tag=%v id=%v", x.Tag, x.Id)
}

////////////////////////////////////////////////////////////

// YiJiInfo 遗迹信息
func (c *Connect) YiJiInfo() error {
	body, err := proto.Marshal(&C2SYiJiInfo{})
	if err != nil {
		return err
	}
	log.Printf("[C][遗迹信息]")
	return c.send(25301, body)
}

func (x *S2CYiJiInfo) ID() uint16 {
	return 25302
}

// Message S2CYiJiInfo Code:25724
func (x *S2CYiJiInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][遗迹信息] items=%v", x.Items)
}

////////////////////////////////////////////////////////////

// IllusionSweep 扫荡
func (c *Connect) IllusionSweep(t *C2SIllusionSweep) error {
	body, err := proto.Marshal(t)
	if err != nil {
		return err
	}
	log.Printf("[C][IllusionSweep] illusion_type=%v", t.IllusionType)
	return c.send(25723, body)
}

func (x *S2CIllusionSweep) ID() uint16 {
	return 25724
}

// Message S2CIllusionSweep Code:25724
func (x *S2CIllusionSweep) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][IllusionSweep] %v", x)
}
