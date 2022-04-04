package mhyc

import (
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func StageFight() {
	// 定时领取主线任务奖励
	go func() {
		t := time.NewTimer(ms100)
		for range t.C {
			h := &S2CGetHistoryTaskPrize{}
			Receive.Action(CLI.GetHistoryTaskPrize)
			if _ = Receive.Wait(h, s3); h.Tag == 0 {
				t.Reset(ms100)
				continue
			}
			t.Reset(RandMillisecond(1800, 3600)) // 30 ~ 60 分钟
		}
	}()
	tm := time.NewTimer(ms500)
	ff := func() {
		Fight.Lock()
		defer Fight.Unlock()
		// 闯关
		for range tm.C {
			Receive.Action(CLI.StageFight)
			f := &S2CStageFight{}
			if _ = Receive.Wait(f, s3); f.Tag == 31 || f.Tag == 9012 {
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
			r := &S2CStageDraw{}
			Receive.Action(CLI.GetStageDraw)
			if _ = Receive.Wait(r, s3); r.Tag == 0 {
				tm.Reset(ms100)
				continue
			}
			break
		}
		tm.Reset(RandMillisecond(8, 30))
	}
	tc := time.NewTimer(ms100)
	for range tc.C {
		ff()
		tc.Reset(RandMillisecond(1800, 3600)) // 30 ~ 60 分钟
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

// GetHistoryTaskPrize 主线任务奖励
func (c *Connect) GetHistoryTaskPrize() error {
	body, err := proto.Marshal(DefineGetHistoryTaskPrize)
	if err != nil {
		return err
	}
	log.Println("[C][GetHistoryTaskPrize]")
	return c.send(713, body)
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

////////////////////////////////////////////////////////////

func (x *S2CGetHistoryTaskPrize) ID() uint16 {
	return 714
}

func (x *S2CGetHistoryTaskPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetHistoryTaskPrize] tag=%v raw=%v", x.Tag, x)
}
