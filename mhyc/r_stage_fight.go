package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func StageFight(ctx context.Context) {
	tm := time.NewTimer(ms500)
	defer tm.Stop()
	ff := func() {
		Fight.Lock()
		defer Fight.Unlock()
		// 闯关
		for range tm.C {
			Receive.Action(CLI.StageFight)
			f := &S2CStageFight{}
			if _ = Receive.Wait(f, s3); f.Tag == 31 || f.Tag == 9012 || (f.Tag == 0 && f.Win == 0) {
				break
			}
			if f.Tag == 17003 {
				tm.Reset(time.Second)
				continue
			}
			if f.Tag == 0 {
				_ = Receive.Wait(&S2CBattlefieldReport{}, s3)
			}
			tm.Reset(ms100)
		}
		// 幸运转盘
		tm.Reset(ms100)
		for range tm.C {
			Receive.Action(CLI.GetStageDraw)
			_ = Receive.Wait(&S2CStageDraw{}, s3)
			if RoleInfo.Get("StageDrawTimes").Int64() <= 0 {
				break
			}
			tm.Reset(ms100)
		}
		tm.Reset(RandMillisecond(8, 30))
	}
	tc := time.NewTimer(ms100)
	defer tc.Stop()
	for {
		select {
		case <-tc.C:
			ff()
			tc.Reset(RandMillisecond(1800, 3600)) // 30 ~ 60 分钟
		case <-ctx.Done():
			return
		}
	}
}

// StageFight 闯关 - 开始
func (c *Connect) StageFight() error {
	body, err := proto.Marshal(DefineStageFight)
	if err != nil {
		return err
	}
	log.Println("[C][StageFight]")
	return c.send(103, body)
}

// GetStageDraw 闯关 幸运转盘
func (c *Connect) GetStageDraw() error {
	body, err := proto.Marshal(DefineStageDraw)
	if err != nil {
		return err
	}
	log.Println("[C][GetStageDraw]")
	return c.send(118, body)
}

////////////////////////////////////////////////////////////

func (x *S2CStageFight) ID() uint16 {
	return 104
}

func (x *S2CStageFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][StagePrize] tag=%v win=%v", x.Tag, x.Win)
}

////////////////////////////////////////////////////////////

func (x *S2CStageDraw) ID() uint16 {
	return 119
}

func (x *S2CStageDraw) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][CStageDraw] tag=%v id=%v", x.Tag, x.Id)
}
