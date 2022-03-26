package mhyc

import (
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// FamilyJJC 竞技
// [RoleInfo] FamilyJJC_TimesLeft 剩余次数
// [RoleInfo] FamilyJJC_Times	  使用次数
// [RoleInfo] FamilyJJC_Score	  积分
func FamilyJJC() {
	t := time.NewTimer(ms10)
	f := func() {
		Fight.Lock()
		defer Fight.Unlock()
		// 战斗
		for {
			if val := RoleInfo.Get("FamilyJJC_Times"); val != nil {
				if val.Int64() >= 10 {
					break
				}
			}
			ret := &S2CFamilyJJCJoin{}
			Receive.Action(CLI.FamilyJJCJoin)
			_ = Receive.Wait(ret, s3)
			if ret.Tag == 0 {
				go func() {
					_ = CLI.FamilyJJCFight(ret)
				}()
				_ = Receive.Wait(&S2CFamilyJJCFight{}, s3)
				time.Sleep(time.Second)
				continue
			}
			if ret.Tag == 17003 {
				time.Sleep(time.Second)
				continue
			}
			// end
			if ret.Tag == 57606 {
				break
			}
		}
		// 领取奖励
		for i := 0; i < 4; i++ {
			go func(i int) {
				_ = CLI.FamilyJJCRecieveAward(int32(i))
			}(i)
			_ = Receive.Wait(&S2CFamilyJJCRecieveAward{}, s3)
		}
	}
	for range t.C {
		f()
		t.Reset(RandMillisecond(1800, 3600)) // 30 ~ 60 分钟
	}
}

func (c *Connect) FamilyJJCJoin() error {
	body, err := proto.Marshal(&C2SFamilyJJCJoin{})
	if err != nil {
		return err
	}
	return c.send(27357, body)
}

func (c *Connect) FamilyJJCRecieveAward(id int32) error {
	body, err := proto.Marshal(&C2SFamilyJJCRecieveAward{Id: id})
	if err != nil {
		return err
	}
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
