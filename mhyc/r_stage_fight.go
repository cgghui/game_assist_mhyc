package mhyc

import (
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func init() {
	PCK[101] = &S2CBattlefieldReport{}
	PCK[104] = &S2CStageFight{}
	PCK[119] = &S2CStageDraw{}
	PCK[714] = &S2CGetHistoryTaskPrize{}
}

var stageFightThread = make(chan interface{})
var stageFightAction = make(chan struct{})

// StageFight 闯关
func (c *Connect) StageFight() {
	time.Sleep(time.Second * 3)
	go func() {
		stageFightAction <- struct{}{}
		t := time.NewTimer(time.Minute)
		for range t.C {
			stageFightAction <- struct{}{}
			t.Reset(RandMillisecond(60, 120))
		}
	}()
	run := func(val interface{}) {
		switch ret := val.(type) {
		case *S2CStageFight:
			{
				//sdt, ok := RoleInfo.Load("StageDrawTimes")
				//fmt.Println(sdt, ok)
				if ret.Tag == 31 {
					_ = c.getStageDraw()
				}
			}
		case *S2CBattlefieldReport:
			if ret.Win == 1 {
				_ = c.EndFight(ret)
				stageFightAction <- struct{}{}
			}
		case *S2CStageDraw:
			if ret.Tag == 0 {
				_ = c.getStageDraw()
				return
			}
			_ = c.getHistoryTaskPrize()
		case *S2CGetHistoryTaskPrize:
			if ret.Tag == 0 {
				_ = c.getHistoryTaskPrize()
			}
		}
	}
	for {
		select {
		case <-stageFightAction:
			_ = c.stageFight()
		case val := <-stageFightThread:
			go run(val)
		}
	}
}

// stageFight 闯关 - 开始
func (c *Connect) stageFight() error {
	body, err := proto.Marshal(DefineStageFight)
	if err != nil {
		return err
	}
	return c.send(103, body)
}

// getStageDraw 闯关 幸运转盘
func (c *Connect) getStageDraw() error {
	body, err := proto.Marshal(DefineStageDraw)
	if err != nil {
		return err
	}
	return c.send(118, body)
}

// getHistoryTaskPrize 主线任务奖励
func (c *Connect) getHistoryTaskPrize() error {
	body, err := proto.Marshal(DefineGetHistoryTaskPrize)
	if err != nil {
		return err
	}
	return c.send(713, body)
}

func (x *S2CStageFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][StagePrize] tag=%v win=%v", x.Tag, x.Win)
	stageFightThread <- x
}

func (x *S2CBattlefieldReport) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][BattlefieldReport] win=%v idx=%v", x.Win, x.Idx)
	stageFightThread <- x
}

func (x *S2CStageDraw) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][CStageDraw] tag=%v id=%v", x.Tag, x.Id)
	stageFightThread <- x
}

func (x *S2CGetHistoryTaskPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetHistoryTaskPrize] tag=%v raw=%v", x.Tag, x)
	stageFightThread <- x
}
