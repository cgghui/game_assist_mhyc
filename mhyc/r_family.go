package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// FamilyJJC 竞技
// [RoleInfo] FamilyJJC_TimesLeft 剩余次数
// [RoleInfo] FamilyJJC_Times	  使用次数
// [RoleInfo] FamilyJJC_Score	  积分
func FamilyJJC(ctx context.Context) {
	t := time.NewTimer(ms10)
	defer t.Stop()
	f := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "家族竞技")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			if val := RoleInfo.Get("FamilyJJC_Times"); val != nil && val.Int64() >= 10 {
				i := 0
				return 0, am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
					go func(i int) {
						_ = CLI.FamilyJJCRecieveAward(int32(i))
					}(i)
					_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CFamilyJJCRecieveAward{}, s3)
					i++
					if i >= 4 {
						return 0, TomorrowDuration(RandMillisecond(1800, 3600))
					}
					return ms100, 0
				})
			}
			// 匹配
			Receive.Action(CLI.FamilyJJCJoin)
			ret := &S2CFamilyJJCJoin{}
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, ret, s3); err != nil {
				loop = 0
				next = RandMillisecond(6, 10)
				return
			}
			if ret.Tag == 17003 {
				loop = RandMillisecond(10, 20)
				next = 0
				return
			}
			if ret.Tag == 57606 { // 挑战次数不足
				i := 0
				return 0, am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
					go func(i int) {
						_ = CLI.FamilyJJCRecieveAward(int32(i))
					}(i)
					_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CFamilyJJCRecieveAward{}, s3)
					i++
					if i >= 4 {
						return 0, TomorrowDuration(RandMillisecond(1800, 3600))
					}
					return ms100, 0
				})
			}
			// 战斗
			go func() {
				_ = CLI.FamilyJJCFight(ret)
			}()
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CFamilyJJCFight{}, s3); err != nil {
				loop = 0
				next = RandMillisecond(6, 10)
				return
			}
			return time.Second, 0
		})
	}
	for {
		select {
		case <-t.C:
			t.Reset(f())
		case <-ctx.Done():
			return
		}
	}
}

func (c *Connect) FamilyJJCJoin() error {
	body, err := proto.Marshal(&C2SFamilyJJCJoin{})
	if err != nil {
		return err
	}
	log.Println("[C][FamilyJJCJoin]")
	return c.send(27357, body)
}

func (c *Connect) FamilyJJCRecieveAward(id int32) error {
	body, err := proto.Marshal(&C2SFamilyJJCRecieveAward{Id: id})
	if err != nil {
		return err
	}
	log.Printf("[C][FamilyJJCRecieveAward] id=%v", id)
	return c.send(27355, body)
}

func (c *Connect) FamilyJJCFight(act *S2CFamilyJJCJoin) error {
	dat := C2SFamilyJJCFight{
		UserId: make([]int64, 0, 0),
	}
	for _, self := range act.Self {
		dat.UserId = append(dat.UserId, self.UserId)
	}
	body, err := proto.Marshal(&dat)
	if err != nil {
		return err
	}
	log.Printf("[C][FamilyJJCRecieveAward] user=%v", dat.UserId)
	return c.send(27359, body)
}

////////////////////////////////////////////////////////////

func (x *S2CFamilyInfo) ID() uint16 {
	return 20002
}

func (x *S2CFamilyInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][FamilyInfo] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CFamilyJJCJoin) ID() uint16 {
	return 27358
}

func (x *S2CFamilyJJCJoin) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][FamilyJJCJoin] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CFamilyJJCFight) ID() uint16 {
	return 27363
}

func (x *S2CFamilyJJCFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][FamilyJJCFight] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CFamilyJJCRecieveAward) ID() uint16 {
	return 27356
}

func (x *S2CFamilyJJCRecieveAward) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][FamilyJJCRecieveAward] tag=%v", x.Tag)
}
