package mhyc

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// BossPersonal 个人BOSS
func BossPersonal() {
	t := time.NewTimer(ms100)
	f := func() time.Duration {
		ret := &S2CBossPersonalSweep{}
		Receive.Action(CLI.BossPersonalSweep)
		_ = Receive.Wait(ret, s3)
		if ret.Tag == 4055 { // end
			return TomorrowDuration(RandMillisecond(600, 1800))
		}
		return time.Minute
	}
	for range t.C {
		t.Reset(f())
	}
}

// BossVIP 至尊BOSS
func BossVIP() {
	t := time.NewTimer(ms100)
	f := func() time.Duration {
		ret := &S2CBossVipSweep{}
		Receive.Action(CLI.BossVipSweep)
		_ = Receive.Wait(ret, s3)
		if ret.Tag == 4055 { // end
			return TomorrowDuration(RandMillisecond(600, 1800))
		}
		return time.Minute
	}
	for range t.C {
		t.Reset(f())
	}
}

// BossMulti 多人BOSS
func BossMulti() {
	t := time.NewTimer(ms100)
	f := func() time.Duration {
		go func() {
			_ = CLI.MultiBossJoinScene(&C2SMultiBossJoinScene{Id: 8})
		}()
		enter := &S2CMonsterEnterMap{}
		_ = Receive.Wait(enter, s30)
		_ = Receive.Wait(&S2CMultiBossJoinScene{}, s30)
		fmt.Println()
		//ret := &S2CStartFight{}
		//go func() {
		//	_ = CLI.StartFight(&C2SStartFight{Id: 315, Type: 8})
		//}()
		//_ = Receive.Wait(ret, s3)
		return time.Hour
	}
	for range t.C {
		t.Reset(f())
	}
}

// BossPersonalSweep Boss - 本服BOSS - 个人BOSS 一键扫荡
func (c *Connect) BossPersonalSweep() error {
	body, err := proto.Marshal(&C2SBossPersonalSweep{})
	if err != nil {
		return err
	}
	return c.send(604, body)
}

// BossVipSweep Boss - 本服BOSS - 至尊BOSS 一键扫荡
func (c *Connect) BossVipSweep() error {
	body, err := proto.Marshal(&C2SBossVipSweep{})
	if err != nil {
		return err
	}
	return c.send(664, body)
}

// MultiBossJoinScene Boss - 本服BOSS - 多人
func (c *Connect) MultiBossJoinScene(i *C2SMultiBossJoinScene) error {
	body, err := proto.Marshal(i)
	if err != nil {
		return err
	}
	log.Printf("[C][MultiBossJoinScene] id=%v", i.Id)
	return c.send(1123, body)
}

////////////////////////////////////////////////////////////

func (x *S2CBossPersonalSweep) ID() uint16 {
	return 605
}

// Message S2CBossPersonalSweep 605
func (x *S2CBossPersonalSweep) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][BossPersonalSweep] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CBossVipSweep) ID() uint16 {
	return 665
}

// Message S2CBossVipSweep 665
func (x *S2CBossVipSweep) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][BossVipSweep] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CMultiBossJoinScene) ID() uint16 {
	return 1124
}

// Message S2CMultiBossJoinScene 1124
func (x *S2CMultiBossJoinScene) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][MultiBossJoinScene] tag=%v id=%v", x.Tag, x.Id)
}
