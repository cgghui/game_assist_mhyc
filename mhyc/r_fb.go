package mhyc

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// FuBen 副本
func FuBen(ctx context.Context) {
	// 材料扫荡
	t1 := time.NewTimer(ms100)
	defer t1.Stop()
	f1 := func() time.Duration {
		Fight.Lock()
		defer Fight.Unlock()
		//
		isEnd := 0
		//
		ims := &S2CInstanceMaterialSweep{}
		Receive.Action(CLI.InstanceMaterialSweep)
		if _ = Receive.Wait(ims, s3); ims.Tag != 0 {
			isEnd += 1
		}
		//
		isl1 := &S2CInstanceSLSweep{}
		Receive.Action(CLI.InstanceSLSweep1)
		if _ = Receive.Wait(isl1, s3); isl1.Tag != 0 {
			isEnd += 1
		}
		//
		isl2 := &S2CInstanceSLSweep{}
		Receive.Action(CLI.InstanceSLSweep2)
		if _ = Receive.Wait(isl2, s3); isl2.Tag != 0 {
			isEnd += 1
		}
		//
		if isEnd >= 3 {
			return TomorrowDuration(RandMillisecond(30000, 30600))
		}
		return time.Second
	}
	// 爬塔 宝石
	t2 := time.NewTimer(ms100)
	defer t2.Stop()
	f2 := func() time.Duration {
		Fight.Lock()
		defer Fight.Unlock()
		enter := &S2CClimbingTowerEnter{}
		go func() {
			_ = CLI.ClimbingTowerEnter(&C2SClimbingTowerEnter{TowerType: 1})
		}()
		if err := Receive.Wait(enter, s3); err != nil {
			return ms500
		}
		if enter.Tag != 0 {
			return TomorrowDuration(RandMillisecond(30000, 30600))
		}
		for {
			go func() {
				_ = CLI.ClimbingTowerFight(&C2SClimbingTowerFight{TowerType: 1, Id: 0})
			}()
			r := &S2CClimbingTowerFight{}
			if err := Receive.Wait(r, s3); err != nil {
				return ms100
			}
			if r.Tag != 0 {
				break
			}
		}
		go func() {
			_ = CLI.ClimbingTowerLeave(&C2SClimbingTowerLeave{TowerType: 1})
		}()
		_ = Receive.Wait(&S2CClimbingTowerLeave{}, s3)
		return ms100
	}
	// 爬塔 天仙
	t3 := time.NewTimer(ms100)
	defer t3.Stop()
	f3 := func() time.Duration {
		Fight.Lock()
		defer Fight.Unlock()
		enter := &S2CClimbingTowerEnter{}
		go func() {
			_ = CLI.ClimbingTowerEnter(&C2SClimbingTowerEnter{TowerType: 2})
		}()
		if err := Receive.Wait(enter, s3); err != nil {
			return ms500
		}
		if enter.Tag != 0 {
			return TomorrowDuration(RandMillisecond(30000, 30600))
		}
		for {
			go func() {
				_ = CLI.ClimbingTowerFight(&C2SClimbingTowerFight{TowerType: 2, Id: 0})
			}()
			r := &S2CClimbingTowerFight{}
			if err := Receive.Wait(r, s3); err != nil {
				return ms100
			}
			if r.Tag != 0 {
				break
			}
		}
		go func() {
			_ = CLI.ClimbingTowerLeave(&C2SClimbingTowerLeave{TowerType: 2})
		}()
		_ = Receive.Wait(&S2CClimbingTowerLeave{}, s3)
		return ms100
	}
	// 爬塔 战神
	t4 := time.NewTimer(ms100)
	defer t4.Stop()
	f4 := func() time.Duration {
		Fight.Lock()
		defer Fight.Unlock()
		enter := &S2CClimbingTowerEnter{}
		go func() {
			_ = CLI.ClimbingTowerEnter(&C2SClimbingTowerEnter{TowerType: 3})
		}()
		if err := Receive.Wait(enter, s3); err != nil {
			return ms500
		}
		if enter.Tag != 0 {
			return TomorrowDuration(RandMillisecond(30000, 30600))
		}
		for {
			go func() {
				_ = CLI.ClimbingTowerFight(&C2SClimbingTowerFight{TowerType: 3, Id: 0})
			}()
			r := &S2CClimbingTowerFight{}
			if err := Receive.Wait(r, s3); err != nil {
				return ms100
			}
			if r.Tag != 0 {
				break
			}
		}
		go func() {
			_ = CLI.ClimbingTowerLeave(&C2SClimbingTowerLeave{TowerType: 3})
		}()
		_ = Receive.Wait(&S2CClimbingTowerLeave{}, s3)
		return ms100
	}
	// 爬塔 仙童
	t5 := time.NewTimer(ms100)
	defer t5.Stop()
	f5 := func() time.Duration {
		Fight.Lock()
		defer Fight.Unlock()
		enter := &S2CClimbingTowerEnter{}
		go func() {
			_ = CLI.ClimbingTowerEnter(&C2SClimbingTowerEnter{TowerType: 4})
		}()
		if err := Receive.Wait(enter, s3); err != nil {
			return ms500
		}
		if enter.Tag != 0 {
			return TomorrowDuration(RandMillisecond(30000, 30600))
		}
		for {
			go func() {
				_ = CLI.ClimbingTowerFight(&C2SClimbingTowerFight{TowerType: 4, Id: 0})
			}()
			r := &S2CClimbingTowerFight{}
			if err := Receive.Wait(r, s3); err != nil {
				return ms100
			}
			if r.Tag != 0 {
				break
			}
		}
		go func() {
			_ = CLI.ClimbingTowerLeave(&C2SClimbingTowerLeave{TowerType: 4})
		}()
		_ = Receive.Wait(&S2CClimbingTowerLeave{}, s3)
		return ms100
	}
	// 爬塔 剑魂
	t6 := time.NewTimer(ms100)
	defer t6.Stop()
	f6 := func() time.Duration {
		Fight.Lock()
		defer Fight.Unlock()
		enter := &S2CClimbingTowerEnter{}
		go func() {
			_ = CLI.ClimbingTowerEnter(&C2SClimbingTowerEnter{TowerType: 5})
		}()
		if err := Receive.Wait(enter, s3); err != nil {
			return ms500
		}
		if enter.Tag != 0 {
			return TomorrowDuration(RandMillisecond(30000, 30600))
		}
		for {
			go func() {
				_ = CLI.ClimbingTowerFight(&C2SClimbingTowerFight{TowerType: 5, Id: 0})
			}()
			r := &S2CClimbingTowerFight{}
			if err := Receive.Wait(r, s3); err != nil {
				return ms100
			}
			if r.Tag != 0 {
				break
			}
		}
		go func() {
			_ = CLI.ClimbingTowerLeave(&C2SClimbingTowerLeave{TowerType: 5})
		}()
		_ = Receive.Wait(&S2CClimbingTowerLeave{}, s3)
		return ms100
	}
	// 组队 灵气
	t7 := time.NewTimer(ms100)
	defer t7.Stop()
	f7 := func() time.Duration {
		id := int32(241)
		Fight.Lock()
		defer Fight.Unlock()
		matching := &S2CTeamInstanceMatching{}
		go func() {
			_ = CLI.TeamInstanceMatching(&C2STeamInstanceMatching{InstanceType: id})
		}()
		if err := Receive.Wait(matching, s3); err != nil {
			return ms500
		}
		if matching.Tag != 0 {
			return TomorrowDuration(RandMillisecond(30000, 30600))
		}
		go func() {
			user := make([]int64, 0)
			for _, player := range matching.Players {
				user = append(user, player.UserId)
			}
			_ = CLI.TeamInstanceStartFight(&C2STeamInstanceStartFight{InstanceType: id, UserIds: user})
		}()
		report := &S2CTeamInstanceGetReport{}
		_ = Receive.Wait(report, s3)
		go func() {
			_ = CLI.TeamInstanceGetReport(&C2STeamInstanceGetReport{InstanceType: id})
		}()
		_ = Receive.Wait(report, s3)
		fmt.Println(report)
		return ms100
	}
	// 组队 进阶
	t8 := time.NewTimer(ms100)
	defer t8.Stop()
	f8 := func() time.Duration {
		id := int32(242)
		Fight.Lock()
		defer Fight.Unlock()
		matching := &S2CTeamInstanceMatching{}
		go func() {
			_ = CLI.TeamInstanceMatching(&C2STeamInstanceMatching{InstanceType: id})
		}()
		if err := Receive.Wait(matching, s3); err != nil {
			return ms500
		}
		if matching.Tag != 0 {
			return TomorrowDuration(RandMillisecond(30000, 30600))
		}
		go func() {
			user := make([]int64, 0)
			for _, player := range matching.Players {
				user = append(user, player.UserId)
			}
			_ = CLI.TeamInstanceStartFight(&C2STeamInstanceStartFight{InstanceType: id, UserIds: user})
		}()
		report := &S2CTeamInstanceGetReport{}
		_ = Receive.Wait(report, s3)
		go func() {
			_ = CLI.TeamInstanceGetReport(&C2STeamInstanceGetReport{InstanceType: id})
		}()
		_ = Receive.Wait(report, s3)
		fmt.Println(report)
		return ms100
	}
	// 组队 宠物装备
	t9 := time.NewTimer(ms100)
	defer t9.Stop()
	f9 := func() time.Duration {
		id := int32(340)
		Fight.Lock()
		defer Fight.Unlock()
		matching := &S2CTeamInstanceMatching{}
		go func() {
			_ = CLI.TeamInstanceMatching(&C2STeamInstanceMatching{InstanceType: id})
		}()
		if err := Receive.Wait(matching, s3); err != nil {
			return ms500
		}
		if matching.Tag != 0 {
			return TomorrowDuration(RandMillisecond(30000, 30600))
		}
		go func() {
			user := make([]int64, 0)
			for _, player := range matching.Players {
				user = append(user, player.UserId)
			}
			_ = CLI.TeamInstanceStartFight(&C2STeamInstanceStartFight{InstanceType: id, UserIds: user})
		}()
		report := &S2CTeamInstanceGetReport{}
		_ = Receive.Wait(report, s3)
		go func() {
			_ = CLI.TeamInstanceGetReport(&C2STeamInstanceGetReport{InstanceType: id})
		}()
		_ = Receive.Wait(report, s3)
		return ms100
	}
	// 组队 星图
	t10 := time.NewTimer(ms100)
	defer t10.Stop()
	f10 := func() time.Duration {
		id := int32(341)
		Fight.Lock()
		defer Fight.Unlock()
		matching := &S2CTeamInstanceMatching{}
		go func() {
			_ = CLI.TeamInstanceMatching(&C2STeamInstanceMatching{InstanceType: id})
		}()
		if err := Receive.Wait(matching, s3); err != nil {
			return ms500
		}
		if matching.Tag != 0 {
			return TomorrowDuration(RandMillisecond(30000, 30600))
		}
		go func() {
			user := make([]int64, 0)
			for _, player := range matching.Players {
				user = append(user, player.UserId)
			}
			_ = CLI.TeamInstanceStartFight(&C2STeamInstanceStartFight{InstanceType: id, UserIds: user})
		}()
		report := &S2CTeamInstanceGetReport{}
		_ = Receive.Wait(report, s3)
		go func() {
			_ = CLI.TeamInstanceGetReport(&C2STeamInstanceGetReport{InstanceType: id})
		}()
		_ = Receive.Wait(report, s3)
		return ms100
	}
	//
	// cp 1 2
	t11 := time.NewTimer(ms100)
	defer t11.Stop()
	f11 := func() time.Duration {
		Fight.Lock()
		defer Fight.Unlock()
		i := 1
		ts := time.NewTimer(ms10)
		defer ts.Stop()
		for range ts.C {
			if i >= 11 {
				break
			}
			Receive.Action(CLI.JungleHuntData)
			_ = Receive.Wait(&S2CJungleHuntData{}, s3)
			//
			go func(i int) {
				_ = CLI.JungleHuntFight(&C2SJungleHuntFight{CpId: int32(i)})
			}(i)
			r := &S2CJungleHuntFight{}
			_ = Receive.Wait(r, s3)
			if r.Tag == 58871 { // 全体阵亡
				break
			}
			// 尝试阵亡复活
			// 第8层以下进行复活
			if r.CpId <= 8 && r.Tag == 0 && r.Win == 0 && RoleInfo.Get("Coin4").Int64() > 1000 {
				Receive.Action(CLI.JungleHuntTreat)
				_ = Receive.Wait(&S2CJungleHuntTreat{}, s3)
				ts.Reset(time.Second)
				continue
			}
			if (r.Tag == 0 && r.Win == 1) || r.Tag == 58851 {
				i++
			}
			ts.Reset(time.Second)
			continue
		}
		box := &S2CJungleHuntOpenBox{}
		Receive.Action(CLI.JungleHuntOpenBox)
		_ = Receive.Wait(&S2CJungleHuntOpenBox{}, s3)
		if box.Tag != 0 {
			return ms100
		}
		if RoleInfo.Get("JungleHunt_LeftResetTimes").Int64() <= 0 {
			return TomorrowDuration(18000)
		}
		Receive.Action(CLI.JungleHuntReset)
		_ = Receive.Wait(&S2CJungleHuntReset{}, s3)
		return ms100
	}
	//
	for {
		select {
		case <-t1.C:
			t1.Reset(f1())
		case <-t2.C:
			t2.Reset(f2())
		case <-t3.C:
			t3.Reset(f3())
		case <-t4.C:
			t4.Reset(f4())
		case <-t5.C:
			t5.Reset(f5())
		case <-t6.C:
			t6.Reset(f6())
		case <-t7.C:
			t7.Reset(f7())
		case <-t8.C:
			t8.Reset(f8())
		case <-t9.C:
			t9.Reset(f9())
		case <-t10.C:
			t10.Reset(f10())
		case <-t11.C:
			t11.Reset(f11())
		case <-ctx.Done():
			return
		}
	}
}

// InstanceMaterialSweep 副本 材料 一键扫荡
func (c *Connect) InstanceMaterialSweep() error {
	body, err := proto.Marshal(&C2SInstanceMaterialSweep{Id: 0})
	if err != nil {
		return err
	}
	log.Println("[C][InstanceMaterialSweep] ID: 0")
	return c.send(609, body)
}

// InstanceSLSweep1 副本 试炼 一键扫荡
func (c *Connect) InstanceSLSweep1() error {
	body, err := proto.Marshal(&C2SInstanceSLSweep{Id: 1})
	if err != nil {
		return err
	}
	log.Println("[C][InstanceSLSweep] ID: 1")
	return c.send(23015, body)
}

// InstanceSLSweep2 副本 试炼 一键扫荡
func (c *Connect) InstanceSLSweep2() error {
	body, err := proto.Marshal(&C2SInstanceSLSweep{Id: 2})
	if err != nil {
		return err
	}
	log.Println("[C][InstanceSLSweep] ID: 2")
	return c.send(23015, body)
}

// ClimbingTowerEnter 进入爬塔场景
func (c *Connect) ClimbingTowerEnter(enter *C2SClimbingTowerEnter) error {
	body, err := proto.Marshal(enter)
	if err != nil {
		return err
	}
	log.Println("[C][ClimbingTowerEnter] TowerType: 1")
	return c.send(22571, body)
}

func (c *Connect) ClimbingTowerLeave(enter *C2SClimbingTowerLeave) error {
	body, err := proto.Marshal(enter)
	if err != nil {
		return err
	}
	log.Printf("[C][ClimbingTowerLeave] TowerType: %v", enter.TowerType)
	return c.send(22573, body)
}

// ClimbingTowerFight 副本 - 爬塔 - 战斗
func (c *Connect) ClimbingTowerFight(act *C2SClimbingTowerFight) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(22575, body)
}

func (c *Connect) TeamInstanceMatching(act *C2STeamInstanceMatching) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(24402, body)
}

func (c *Connect) TeamInstanceStartFight(act *C2STeamInstanceStartFight) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(24404, body)
}

func (c *Connect) TeamInstanceGetReport(act *C2STeamInstanceGetReport) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(24405, body)
}

func (c *Connect) JungleHuntFight(act *C2SJungleHuntFight) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(28703, body)
}

func (c *Connect) JungleHuntTreat() error {
	body, err := proto.Marshal(&C2SJungleHuntTreat{})
	if err != nil {
		return err
	}
	return c.send(28705, body)
}

func (c *Connect) JungleHuntBattleArr(arr *C2SJungleHuntBattleArr) error {
	body, err := proto.Marshal(arr)
	if err != nil {
		return err
	}
	return c.send(28723, body)
}

func (c *Connect) JungleHuntOpenBox() error {
	body, err := proto.Marshal(&C2SJungleHuntBattleArr{CpId: 0})
	if err != nil {
		return err
	}
	return c.send(28709, body)
}

func (c *Connect) JungleHuntReset() error {
	body, err := proto.Marshal(&C2SJungleHuntReset{})
	if err != nil {
		return err
	}
	return c.send(28707, body)
}

func (c *Connect) JungleHuntData() error {
	body, err := proto.Marshal(&C2SJungleHuntData{})
	if err != nil {
		return err
	}
	return c.send(28701, body)
}

////////////////////////////////////////////////////////////

func (x *S2CInstanceMaterialSweep) ID() uint16 {
	return 610
}

// Message S2CInstanceMaterialSweep 610
func (x *S2CInstanceMaterialSweep) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][InstanceMaterialSweep] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CInstanceSLSweep) ID() uint16 {
	return 23016
}

// Message S2CInstanceSLSweep 23016
func (x *S2CInstanceSLSweep) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][InstanceSLSweep] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CClimbingTowerEnter) ID() uint16 {
	return 22572
}

// Message S2CClimbingTowerEnter 22572
func (x *S2CClimbingTowerEnter) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ClimbingTowerEnter] tag=%v tower_type=%v", x.Tag, x.TowerType)
}

////////////////////////////////////////////////////////////

func (x *S2CClimbingTowerFight) ID() uint16 {
	return 22576
}

// Message S2CClimbingTowerFight 22576
func (x *S2CClimbingTowerFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ClimbingTowerFight] tag=%v id=%v tower_type=%v", x.Tag, x.Id, x.TowerType)
}

////////////////////////////////////////////////////////////

func (x *S2CClimbingTowerLeave) ID() uint16 {
	return 22574
}

// Message S2CClimbingTowerLeave 22574
func (x *S2CClimbingTowerLeave) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ClimbingTowerLeave] tag=%v tower_type=%v", x.Tag, x.TowerType)
}

////////////////////////////////////////////////////////////

func (x *S2CTeamInstanceMatching) ID() uint16 {
	return 24403
}

// Message S2CTeamInstanceMatching 24403
func (x *S2CTeamInstanceMatching) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][TeamInstanceMatching] tag=%v instance_type=%v players=%v", x.Tag, x.InstanceType, x.Players)
}

////////////////////////////////////////////////////////////

func (x *S2CTeamInstanceGetReport) ID() uint16 {
	return 24406
}

// Message S2CTeamInstanceGetReport 24406
func (x *S2CTeamInstanceGetReport) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	_ = CLI.EndFight(x.Report)
	log.Printf("[S][TeamInstanceGetReport] tag=%v instance_type=%v players=%v", x.Tag, x.InstanceType, x.Index)
}

////////////////////////////////////////////////////////////

func (x *S2CJungleHuntFight) ID() uint16 {
	return 28704
}

// Message S2CJungleHuntFight 28704
func (x *S2CJungleHuntFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][TeamInstanceGetReport] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CJungleHuntTreat) ID() uint16 {
	return 28706
}

// Message S2CJungleHuntTreat 28706
func (x *S2CJungleHuntTreat) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JungleHuntTreat] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CJungleHuntBattleArr) ID() uint16 {
	return 28724
}

// Message S2CJungleHuntBattleArr 28724
func (x *S2CJungleHuntBattleArr) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JungleHuntBattleArr] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CJungleHuntOpenBox) ID() uint16 {
	return 28710
}

// Message S2CJungleHuntOpenBox 28710
func (x *S2CJungleHuntOpenBox) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JungleHuntOpenBox] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CJungleHuntReset) ID() uint16 {
	return 28708
}

// Message S2CJungleHuntReset 28708
func (x *S2CJungleHuntReset) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JungleHuntReset] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CJungleHuntData) ID() uint16 {
	return 28702
}

// Message S2CJungleHuntData 28702
func (x *S2CJungleHuntData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JungleHuntData] %v", x)
}
