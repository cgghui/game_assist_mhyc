package mhyc

import (
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
	go func() {
		lg := Receive.CreateChannel(&S2CMultiBossGetDamageLog{})
		ls := Receive.CreateChannel(&S2CMultiBossLeaveScene{})
		bi := Receive.CreateChannel(&S2CMultiBossInfo{})
		for {
			select {
			case <-lg.Wait():
				lg.Call.Message(nil)
			case <-ls.Wait():
				ls.Call.Message(nil)
			case <-ls.Wait():
				bi.Call.Message(nil)
			}
		}
	}()
	id := 8
	t := time.NewTimer(ms100)
	f := func() time.Duration {
		if RoleInfo.Get("MultiBoss_Times").Int64() == 0 {
			if RoleInfo.Get("MultiBoss_Add_Times").Int64() == 10 {
				return TomorrowDuration(9 * time.Hour)
			}
			mn := time.Unix(RoleInfo.Get("MultiBoss_NextTime").Int64(), 0).Local().Add(time.Minute)
			cur := time.Now()
			if cur.Before(mn) {
				return mn.Add(time.Minute).Sub(cur)
			}
		}
		info := &S2CMultiBossInfo{}
		Receive.Action(CLI.MultiBossInfo)
		if err := Receive.Wait(info, s3); err != nil {
			return ms500
		}
		for _, item := range info.Items {
			if item.Id == int32(id) {
				rt := time.Unix(item.ReliveTimestamp, 0).Local().Add(time.Minute)
				cur := time.Now()
				if cur.Before(rt) {
					return rt.Add(time.Minute).Sub(cur)
				}
			}
		}
		go func() {
			_ = CLI.MultiBossJoinScene(&C2SMultiBossJoinScene{Id: int32(id)})
		}()
		go func() {
			_ = Receive.Wait(&S2CMultiBossJoinScene{}, s30)
		}()
		enter := &S2CMonsterEnterMap{}
		_ = Receive.Wait(enter, s30)
		tc := time.NewTimer(0)
		for range tc.C {
			ret := &S2CStartFight{}
			go func() {
				_ = CLI.StartFight(&C2SStartFight{Id: enter.Id, Type: int64(id)})
			}()
			_ = Receive.Wait(ret, s3)
			if ret.Tag == 4022 {
				tc.Stop()
				break
			}
			tc.Reset(ms500)
		}
		return ms100
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

// MultiBossPlayerInBoss Boss - 本服BOSS - 多人
func (c *Connect) MultiBossPlayerInBoss() error {
	body, err := proto.Marshal(&C2SMultiBossPlayerInBoss{})
	if err != nil {
		return err
	}
	return c.send(1125, body)
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

// MultiBossInfo Boss - 本服BOSS - 多人 BOSS信息
func (c *Connect) MultiBossInfo() error {
	body, err := proto.Marshal(&C2SMultiBossInfo{})
	if err != nil {
		return err
	}
	return c.send(1130, body)
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

////////////////////////////////////////////////////////////

func (x *S2CMultiBossGetDamageLog) ID() uint16 {
	return 1122
}

// Message S2CMultiBossGetDamageLog 1122
func (x *S2CMultiBossGetDamageLog) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][MultiBossGetDamageLog] boss_id=%v my_damage=%v boss_state=%v item=%v", x.BossId, x.MyDamage, x.BossState, x.Items)
}

////////////////////////////////////////////////////////////

func (x *S2CMultiBossInfo) ID() uint16 {
	return 1111
}

// Message S2CMultiBossInfo 1111
func (x *S2CMultiBossInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][MultiBossInfo] items=%v", x.Items)
}

////////////////////////////////////////////////////////////

func (x *S2CMultiBossPlayerInBoss) ID() uint16 {
	return 1126
}

// Message S2CMultiBossPlayerInBoss 1126
func (x *S2CMultiBossPlayerInBoss) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][MultiBossPlayerInBoss] boss_id=%v damage=%v damage_order=%v", x.BossId, x.Damage, x.DamageOrder)
}

////////////////////////////////////////////////////////////

func (x *S2CMultiBossLeaveScene) ID() uint16 {
	return 1128
}

// Message S2CMultiBossPlayerInBoss 1126
func (x *S2CMultiBossLeaveScene) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][MultiBossLeaveScene] tag=%v", x.Tag)
}
