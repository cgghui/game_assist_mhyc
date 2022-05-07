package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func StageFight(ctx context.Context) {
	ff := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "闯关")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			sfChan := make(chan *S2CStageFight)
			go func() {
				defer close(sfChan)
				ret := &S2CStageFight{}
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, ret, s3); err != nil {
					sfChan <- nil
				} else {
					sfChan <- ret
				}
			}()
			Receive.Action(CLI.StageFight)
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CBattlefieldReport{}, s3)
			f := <-sfChan
			if f == nil {
				return 0, RandMillisecond(3, 8)
			}
			if f.Tag == 17003 {
				return 0, RandMillisecond(10, 20)
			}
			if f.Tag == 31 || f.Tag == 9012 || (f.Tag == 0 && f.Win == 0) {
				return 0, am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
					if RoleInfo.Get("StageDrawTimes").Int64() <= 0 {
						return 0, RandMillisecond(1800, 3600)
					}
					Receive.Action(CLI.GetStageDraw)
					if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CStageDraw{}, s3); err != nil {
						return 0, RandMillisecond(3, 8)
					}
					return ms100, 0
				})
			}
			return ms100, 0
		})
	}
	tc := time.NewTimer(ms100)
	defer tc.Stop()
	for {
		select {
		case <-tc.C:
			tc.Reset(ff())
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
