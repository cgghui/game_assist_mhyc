package mhyc

import (
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func StageFight() {
	go func() {
		t := time.NewTimer(ms100)
		for range t.C {
			h := &S2CGetHistoryTaskPrize{}
			Receive.Action(CLI.GetHistoryTaskPrize)
			if _ = Receive.Wait(714, h, s3); h.Tag == 0 {
				t.Reset(ms500)
				continue
			}
			t.Reset(RandMillisecond(1800, 3600)) // 30 ~ 60 分钟
		}
	}()
	tc := time.NewTimer(ms100)
	tm := time.NewTimer(ms500)
	for range tc.C {
		for range tm.C {
			Receive.Action(CLI.StageFight)
			f := &S2CStageFight{}
			if _ = Receive.Wait(104, f, s3); f.Tag == 31 {
				break
			}
			if f.Tag == 0 {
				_ = Receive.Wait(101, &S2CBattlefieldReport{}, s3)
			}
			tm.Reset(ms500)
		}
		tm.Reset(ms100)
		for range tm.C {
			r := &S2CStageDraw{}
			Receive.Action(CLI.GetStageDraw)
			if _ = Receive.Wait(119, r, s3); r.Tag == 0 {
				tm.Reset(ms100)
				continue
			}
			break
		}
		tc.Reset(RandMillisecond(1800, 3600)) // 30 ~ 60 分钟
	}
}

// StageFight 闯关 - 开始
func (c *Connect) StageFight() error {
	body, err := proto.Marshal(DefineStageFight)
	if err != nil {
		return err
	}
	return c.send(103, body)
}

// GetStageDraw 闯关 幸运转盘
func (c *Connect) GetStageDraw() error {
	body, err := proto.Marshal(DefineStageDraw)
	if err != nil {
		return err
	}
	return c.send(118, body)
}

// GetHistoryTaskPrize 主线任务奖励
func (c *Connect) GetHistoryTaskPrize() error {
	body, err := proto.Marshal(DefineGetHistoryTaskPrize)
	if err != nil {
		return err
	}
	return c.send(713, body)
}

func (x *S2CStageFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][StagePrize] tag=%v win=%v", x.Tag, x.Win)
}

func (x *S2CStageDraw) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][CStageDraw] tag=%v id=%v", x.Tag, x.Id)
}

func (x *S2CGetHistoryTaskPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetHistoryTaskPrize] tag=%v raw=%v", x.Tag, x)
}
