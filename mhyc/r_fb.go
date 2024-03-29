package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"strconv"
	"time"
)

// FuBen 副本
func FuBen(ctx context.Context) {
	// 材料扫荡
	t1 := time.NewTimer(ms100)
	defer t1.Stop()
	f1 := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "副本-材料扫荡")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		//
		isEnd := 0
		//
		ims := &S2CInstanceMaterialSweep{}
		Receive.Action(CLI.InstanceMaterialSweep)
		if _ = Receive.WaitWithContextOrTimeout(am.Ctx, ims, s3); ims.Tag != 0 {
			isEnd += 1
		}
		//
		isl1 := &S2CInstanceSLSweep{}
		Receive.Action(CLI.InstanceSLSweep1)
		if _ = Receive.WaitWithContextOrTimeout(am.Ctx, isl1, s3); isl1.Tag != 0 {
			isEnd += 1
		}
		//
		isl2 := &S2CInstanceSLSweep{}
		Receive.Action(CLI.InstanceSLSweep2)
		if _ = Receive.WaitWithContextOrTimeout(am.Ctx, isl2, s3); isl2.Tag != 0 {
			isEnd += 1
		}
		//
		if isEnd >= 3 {
			return TomorrowDuration(RandMillisecond(1800, 3600))
		}
		return time.Second
	}
	// 爬塔 宝石
	t2 := time.NewTimer(ms100)
	defer t2.Stop()
	f2 := func(ctx context.Context) time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "副本-爬塔-宝石")
		defer func() {
			go func() {
				_ = CLI.ClimbingTowerLeave(&C2SClimbingTowerLeave{TowerType: 1})
			}()
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CClimbingTowerLeave{}, s3)
			am.End()
			Fight.Unlock()
		}()
		enter := &S2CClimbingTowerEnter{}
		go func() {
			_ = CLI.ClimbingTowerEnter(&C2SClimbingTowerEnter{TowerType: 1})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, enter, s3); err != nil {
			return ms500
		}
		if enter.Tag != 0 {
			return TomorrowDuration(RandMillisecond(1800, 3600))
		}
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			go func() {
				_ = CLI.ClimbingTowerFight(&C2SClimbingTowerFight{TowerType: 1, Id: 0})
			}()
			r := &S2CClimbingTowerFight{}
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, r, s3); err != nil {
				return 0, RandMillisecond(0, 3)
			}
			if r.Tag != 0 {
				return 0, RandMillisecond(0, 3)
			}
			return ms10, 0
		})
	}
	// 爬塔 天仙
	t3 := time.NewTimer(ms100)
	defer t3.Stop()
	f3 := func(ctx context.Context) time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "副本-爬塔-天仙")
		defer func() {
			go func() {
				_ = CLI.ClimbingTowerLeave(&C2SClimbingTowerLeave{TowerType: 2})
			}()
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CClimbingTowerLeave{}, s3)
			am.End()
			Fight.Unlock()
		}()
		enter := &S2CClimbingTowerEnter{}
		go func() {
			_ = CLI.ClimbingTowerEnter(&C2SClimbingTowerEnter{TowerType: 2})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, enter, s3); err != nil {
			return ms500
		}
		if enter.Tag != 0 {
			return TomorrowDuration(RandMillisecond(1800, 3600))
		}
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			go func() {
				_ = CLI.ClimbingTowerFight(&C2SClimbingTowerFight{TowerType: 2, Id: 0})
			}()
			r := &S2CClimbingTowerFight{}
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, r, s3); err != nil {
				return 0, RandMillisecond(0, 3)
			}
			if r.Tag != 0 {
				return 0, RandMillisecond(0, 3)
			}
			return ms10, 0
		})
	}
	// 爬塔 战神
	t4 := time.NewTimer(ms100)
	defer t4.Stop()
	f4 := func(ctx context.Context) time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "副本-爬塔-战神")
		defer func() {
			go func() {
				_ = CLI.ClimbingTowerLeave(&C2SClimbingTowerLeave{TowerType: 3})
			}()
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CClimbingTowerLeave{}, s3)
			am.End()
			Fight.Unlock()
		}()
		enter := &S2CClimbingTowerEnter{}
		go func() {
			_ = CLI.ClimbingTowerEnter(&C2SClimbingTowerEnter{TowerType: 3})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, enter, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		if enter.Tag != 0 {
			return TomorrowDuration(RandMillisecond(1800, 3600))
		}
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			go func() {
				_ = CLI.ClimbingTowerFight(&C2SClimbingTowerFight{TowerType: 3, Id: 0})
			}()
			r := &S2CClimbingTowerFight{}
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, r, s3); err != nil {
				return 0, RandMillisecond(0, 3)
			}
			if r.Tag != 0 {
				return 0, RandMillisecond(0, 3)
			}
			return ms10, 0
		})
	}
	// 爬塔 仙童
	t5 := time.NewTimer(ms100)
	defer t5.Stop()
	f5 := func(ctx context.Context) time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "副本-爬塔-仙童")
		defer func() {
			go func() {
				_ = CLI.ClimbingTowerLeave(&C2SClimbingTowerLeave{TowerType: 4})
			}()
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CClimbingTowerLeave{}, s3)
			am.End()
			Fight.Unlock()
		}()
		enter := &S2CClimbingTowerEnter{}
		go func() {
			_ = CLI.ClimbingTowerEnter(&C2SClimbingTowerEnter{TowerType: 4})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, enter, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		if enter.Tag != 0 {
			return TomorrowDuration(RandMillisecond(1800, 3600))
		}
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			go func() {
				_ = CLI.ClimbingTowerFight(&C2SClimbingTowerFight{TowerType: 4, Id: 0})
			}()
			r := &S2CClimbingTowerFight{}
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, r, s3); err != nil {
				return 0, RandMillisecond(0, 3)
			}
			if r.Tag != 0 {
				return 0, RandMillisecond(0, 3)
			}
			return ms10, 0
		})
	}
	// 爬塔 剑魂
	t6 := time.NewTimer(ms100)
	defer t6.Stop()
	f6 := func(ctx context.Context) time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "副本-爬塔-剑魂")
		defer func() {
			go func() {
				_ = CLI.ClimbingTowerLeave(&C2SClimbingTowerLeave{TowerType: 5})
			}()
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CClimbingTowerLeave{}, s3)
			am.End()
			Fight.Unlock()
		}()
		enter := &S2CClimbingTowerEnter{}
		go func() {
			_ = CLI.ClimbingTowerEnter(&C2SClimbingTowerEnter{TowerType: 5})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, enter, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		if enter.Tag != 0 {
			// 分解
			Receive.Action(CLI.SwordSoulResolveJH)
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CSwordSoulResolve{}, s3); err != nil {
				return RandMillisecond(0, 3)
			}
			// 每日奖励
			Receive.Action(CLI.ClimbingTowerGetSwordSoulDayPrize)
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CClimbingTowerGetSwordSoulDayPrize{}, s3); err != nil {
				return RandMillisecond(0, 3)
			}
			//
			return TomorrowDuration(RandMillisecond(1800, 3600))
		}
		reTime := am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			go func() {
				_ = CLI.ClimbingTowerFight(&C2SClimbingTowerFight{TowerType: 5, Id: 0})
			}()
			r := &S2CClimbingTowerFight{}
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, r, s3); err != nil {
				return 0, RandMillisecond(0, 3)
			}
			if r.Tag != 0 {
				return 0, RandMillisecond(0, 3)
			}
			return ms10, 0
		})
		// 分解
		Receive.Action(CLI.SwordSoulResolveJH)
		_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CSwordSoulResolve{}, s3)
		// 每日奖励
		Receive.Action(CLI.ClimbingTowerGetSwordSoulDayPrize)
		_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CClimbingTowerGetSwordSoulDayPrize{}, s3)
		return reTime
	}
	// 组队 灵气
	t7 := time.NewTimer(ms100)
	defer t7.Stop()
	f7 := func() time.Duration {
		id := int32(241)
		Fight.Lock()
		am := SetAction(ctx, "副本-组队-灵气")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		matching := &S2CTeamInstanceMatching{}
		go func() {
			_ = CLI.TeamInstanceMatching(&C2STeamInstanceMatching{InstanceType: id})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, matching, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		if matching.Tag != 0 {
			return TomorrowDuration(RandMillisecond(1800, 3600))
		}
		go func() {
			user := make([]int64, 0)
			for _, player := range matching.Players {
				user = append(user, player.UserId)
			}
			_ = CLI.TeamInstanceStartFight(&C2STeamInstanceStartFight{InstanceType: id, UserIds: user})
		}()
		report := &S2CTeamInstanceGetReport{}
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, report, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		go func() {
			_ = CLI.TeamInstanceGetReport(&C2STeamInstanceGetReport{InstanceType: id})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, report, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		return ms100
	}
	// 组队 进阶
	t8 := time.NewTimer(ms100)
	defer t8.Stop()
	f8 := func() time.Duration {
		id := int32(242)
		Fight.Lock()
		am := SetAction(ctx, "副本-组队-进阶")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		matching := &S2CTeamInstanceMatching{}
		go func() {
			_ = CLI.TeamInstanceMatching(&C2STeamInstanceMatching{InstanceType: id})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, matching, s3); err != nil {
			return RandMillisecond(0, 3)
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
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, report, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		go func() {
			_ = CLI.TeamInstanceGetReport(&C2STeamInstanceGetReport{InstanceType: id})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, report, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		return ms100
	}
	// 组队 宠物装备
	t9 := time.NewTimer(ms100)
	defer t9.Stop()
	f9 := func() time.Duration {
		id := int32(340)
		Fight.Lock()
		am := SetAction(ctx, "副本-组队-宠物装备")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		matching := &S2CTeamInstanceMatching{}
		go func() {
			_ = CLI.TeamInstanceMatching(&C2STeamInstanceMatching{InstanceType: id})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, matching, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		if matching.Tag != 0 {
			return TomorrowDuration(RandMillisecond(1800, 3600))
		}
		go func() {
			user := make([]int64, 0)
			for _, player := range matching.Players {
				user = append(user, player.UserId)
			}
			_ = CLI.TeamInstanceStartFight(&C2STeamInstanceStartFight{InstanceType: id, UserIds: user})
		}()
		report := &S2CTeamInstanceGetReport{}
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, report, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		go func() {
			_ = CLI.TeamInstanceGetReport(&C2STeamInstanceGetReport{InstanceType: id})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, report, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		return ms100
	}
	// 组队 星图
	t10 := time.NewTimer(ms100)
	defer t10.Stop()
	f10 := func() time.Duration {
		id := int32(341)
		Fight.Lock()
		am := SetAction(ctx, "副本-组队-星图")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		matching := &S2CTeamInstanceMatching{}
		go func() {
			_ = CLI.TeamInstanceMatching(&C2STeamInstanceMatching{InstanceType: id})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, matching, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		if matching.Tag != 0 {
			return TomorrowDuration(RandMillisecond(1800, 3600))
		}
		go func() {
			user := make([]int64, 0)
			for _, player := range matching.Players {
				user = append(user, player.UserId)
			}
			_ = CLI.TeamInstanceStartFight(&C2STeamInstanceStartFight{InstanceType: id, UserIds: user})
		}()
		report := &S2CTeamInstanceGetReport{}
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, report, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		go func() {
			_ = CLI.TeamInstanceGetReport(&C2STeamInstanceGetReport{InstanceType: id})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, report, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		return ms100
	}
	// 仙林狩猎
	t11 := time.NewTimer(ms100)
	defer t11.Stop()
	f11 := func(ctx context.Context) time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "副本-仙林狩猎")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		nextFunc := func() time.Duration {
			Receive.Action(CLI.JungleHuntOpenBox)
			box := &S2CJungleHuntOpenBox{}
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CJungleHuntOpenBox{}, s3); err != nil {
				return RandMillisecond(0, 3)
			}
			if box.Tag != 0 {
				return RandMillisecond(0, 3)
			}
			if RoleInfo.Get("JungleHunt_LeftResetTimes").Int64() <= 0 {
				return TomorrowDuration(RandMillisecond(1800, 3600))
			}
			Receive.Action(CLI.JungleHuntReset)
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CJungleHuntReset{}, s3)
			return RandMillisecond(0, 3)
		}
		i := 1
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			if i >= 11 {
				return 0, nextFunc()
			}
			Receive.Action(CLI.JungleHuntData)
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CJungleHuntData{}, s3); err != nil {
				return 0, RandMillisecond(0, 3)
			}
			//
			go func(i int) {
				_ = CLI.JungleHuntFight(&C2SJungleHuntFight{CpId: int32(i)})
			}(i)
			r := &S2CJungleHuntFight{}
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, r, s3); err != nil {
				return 0, RandMillisecond(0, 3)
			}
			if r.Tag == 58871 || r.Tag == 58851 { // 全体阵亡 已通关
				return 0, nextFunc()
			}
			// 尝试阵亡复活
			// 第8层以下进行复活
			if r.CpId <= 8 && r.Tag == 0 && r.Win == 0 && RoleInfo.Get("Coin4").Int64() > 1000 {
				Receive.Action(CLI.JungleHuntTreat)
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CJungleHuntTreat{}, s3); err != nil {
					return 0, RandMillisecond(0, 3)
				}
				return time.Second, 0
			}
			if (r.Tag == 0 && r.Win == 1) || r.Tag == 58851 {
				i++
			}
			return time.Second, 0
		})
	}
	// 快捷挖宝
	t12 := time.NewTimer(ms100)
	defer t12.Stop()
	f12 := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "副本-快捷挖宝")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		//
		go func() {
			_ = CLI.DigTreasure10Times(1)
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CDigTreasure10Times{}, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		//
		go func() {
			_ = CLI.DigTreasure10Times(2)
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CDigTreasure10Times{}, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		//
		go func() {
			_ = CLI.DigTreasure10Times(3)
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CDigTreasure10Times{}, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		//
		return TomorrowDuration(RandMillisecond(1800, 3600))
	}
	// 秘境探险
	t13 := time.NewTimer(ms10)
	defer t13.Stop()
	f13 := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "副本-秘境探险")
		defer func() {
			Receive.Action(CLI.GetYJFBData)
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CGetYJFBData{}, s3)
			am.End()
			Fight.Unlock()
		}()
		// 扫荡
		Receive.Action(CLI.YJFBSweep)
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CYJFBSweep{}, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		// 进入探险
		go func() {
			_ = CLI.GetYJFBGuanQiaData(FuBenId, GuanQiaId)
		}()
		data := &S2CGetYJFBGuanQiaData{}
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, data, s3); err != nil {
			return RandMillisecond(0, 3)
		}
		//
		x := int32(1)
		y := int32(1)
		if data.PlayerGrid != nil {
			x = data.PlayerGrid.X
			y = data.PlayerGrid.Y
		}
		if x == 21 && y == 21 {
			x = 1
			y = 1
		}
		for ; x <= 21; x++ {
			if y >= 22 {
				y = 1
			}
			for ; y <= 21; y++ {
				go func() {
					_ = CLI.YJFBGuanQiaMove(&C2SYJFBGuanQiaMove{
						FuBenId:    FuBenId,
						GuanQiaId:  GuanQiaId,
						TargetGrid: &YJFBGrid{Y: y, X: x},
					})
				}()
				move := &S2CYJFBGuanQiaMove{}
				_ = Receive.WaitWithContextOrTimeout(am.Ctx, move, s3)
				if move.Tag != 0 || len(move.MovePath) == 0 {
					continue
				}
				for _, mp := range move.MovePath {
					if mp.MoveGrid.State == 1 && mp.MoveGrid.EventId != 0 {
						go func() {
							_ = CLI.YJFBGuanQiaMove(&C2SYJFBGuanQiaMove{
								FuBenId:    FuBenId,
								GuanQiaId:  GuanQiaId,
								TargetGrid: &YJFBGrid{Y: mp.MoveGrid.Y, X: mp.MoveGrid.X},
							})
						}()
						_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CYJFBGuanQiaMove{}, s3)
						go func() {
							_ = CLI.YJFBGuanQiaTriggerEvent(&C2SYJFBGuanQiaTriggerEvent{
								FuBenId:     FuBenId,
								GuanQiaId:   GuanQiaId,
								TriggerGrid: &YJFBGrid{Y: mp.MoveGrid.Y, X: mp.MoveGrid.X},
							})
						}()
						ret := &S2CYJFBGuanQiaTriggerEvent{}
						_ = Receive.WaitWithContextOrTimeout(am.Ctx, ret, s3)
						if ret.Tag == 57212 {
							go func() {
								_ = CLI.WareHouseReceiveItem(1)
							}()
							_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CWareHouseReceiveItem{}, s3)
							return TomorrowDuration(RandMillisecond(1800, 3600))
						}
					}
				}
			}
		}
		//
		return RandMillisecond(120, 300)
	}
	for {
		select {
		case <-t1.C:
			t1.Reset(f1())
		case <-t2.C:
			t2.Reset(f2(ctx))
		case <-t3.C:
			t3.Reset(f3(ctx))
		case <-t4.C:
			t4.Reset(f4(ctx))
		case <-t5.C:
			t5.Reset(f5(ctx))
		case <-t6.C:
			t6.Reset(f6(ctx))
		case <-t7.C:
			t7.Reset(f7())
		case <-t8.C:
			t8.Reset(f8())
		case <-t9.C:
			t9.Reset(f9())
		case <-t10.C:
			t10.Reset(f10())
		case <-t11.C:
			t11.Reset(f11(ctx))
		case <-t12.C:
			t12.Reset(f12())
		case <-t13.C:
			t13.Reset(f13())
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
	log.Printf("[S][InstanceMaterialSweep] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CInstanceSLSweep) ID() uint16 {
	return 23016
}

// Message S2CInstanceSLSweep 23016
func (x *S2CInstanceSLSweep) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][InstanceSLSweep] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CClimbingTowerEnter) ID() uint16 {
	return 22572
}

// Message S2CClimbingTowerEnter 22572
func (x *S2CClimbingTowerEnter) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ClimbingTowerEnter] tag=%v tag_msg=%s tower_type=%v", x.Tag, GetTagMsg(x.Tag), x.TowerType)
}

////////////////////////////////////////////////////////////

func (x *S2CClimbingTowerFight) ID() uint16 {
	return 22576
}

// Message S2CClimbingTowerFight 22576
func (x *S2CClimbingTowerFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ClimbingTowerFight] tag=%v tag_msg=%s id=%v tower_type=%v", x.Tag, GetTagMsg(x.Tag), x.Id, x.TowerType)
}

////////////////////////////////////////////////////////////

func (x *S2CClimbingTowerLeave) ID() uint16 {
	return 22574
}

// Message S2CClimbingTowerLeave 22574
func (x *S2CClimbingTowerLeave) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ClimbingTowerLeave] tag=%v tag_msg=%s tower_type=%v", x.Tag, GetTagMsg(x.Tag), x.TowerType)
}

////////////////////////////////////////////////////////////

func (x *S2CTeamInstanceMatching) ID() uint16 {
	return 24403
}

// Message S2CTeamInstanceMatching 24403
func (x *S2CTeamInstanceMatching) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][TeamInstanceMatching] tag=%v tag_msg=%s instance_type=%v players=%v", x.Tag, GetTagMsg(x.Tag), x.InstanceType, x.Players)
}

////////////////////////////////////////////////////////////

func (x *S2CTeamInstanceGetReport) ID() uint16 {
	return 24406
}

// Message S2CTeamInstanceGetReport 24406
func (x *S2CTeamInstanceGetReport) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	_ = CLI.EndFight(x.Report)
	log.Printf("[S][TeamInstanceGetReport] tag=%v tag_msg=%s instance_type=%v players=%v", x.Tag, GetTagMsg(x.Tag), x.InstanceType, x.Index)
}

////////////////////////////////////////////////////////////

func (x *S2CJungleHuntFight) ID() uint16 {
	return 28704
}

// Message S2CJungleHuntFight 28704
func (x *S2CJungleHuntFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][TeamInstanceGetReport] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CJungleHuntTreat) ID() uint16 {
	return 28706
}

// Message S2CJungleHuntTreat 28706
func (x *S2CJungleHuntTreat) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JungleHuntTreat] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CJungleHuntBattleArr) ID() uint16 {
	return 28724
}

// Message S2CJungleHuntBattleArr 28724
func (x *S2CJungleHuntBattleArr) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JungleHuntBattleArr] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CJungleHuntOpenBox) ID() uint16 {
	return 28710
}

// Message S2CJungleHuntOpenBox 28710
func (x *S2CJungleHuntOpenBox) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JungleHuntOpenBox] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
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

////////////////////////////////////////////////////////////

func (c *Connect) SwordSoulResolveJH() error {
	uid := strconv.FormatInt(RoleInfo.Get("UserId").Int64(), 10)
	items := []string{
		uid + "_229505",
		uid + "_229506",
	}
	body, err := proto.Marshal(&C2SSwordSoulResolve{ItemIds: items})
	if err != nil {
		return err
	}
	log.Printf("[C][SwordSoulResolve] item_ids=%v", items)
	return c.send(3205, body)
}

func (x *S2CSwordSoulResolve) ID() uint16 {
	return 3206
}

// Message S2CSwordSoulResolve 3206
func (x *S2CSwordSoulResolve) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SwordSoulResolve] tag=%v item_ids=%v", x.Tag, x.ItemIds)
}

////////////////////////////////////////////////////////////

func (c *Connect) ClimbingTowerGetSwordSoulDayPrize() error {
	body, err := proto.Marshal(&C2SClimbingTowerGetSwordSoulDayPrize{})
	if err != nil {
		return err
	}
	log.Println("[C][ClimbingTowerGetSwordSoulDayPrize]")
	return c.send(22581, body)
}

func (x *S2CClimbingTowerGetSwordSoulDayPrize) ID() uint16 {
	return 22582
}

// Message S2CClimbingTowerGetSwordSoulDayPrize 22582
func (x *S2CClimbingTowerGetSwordSoulDayPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SwordSoulResolve] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

// YJFBSweep 秘境探险
func (c *Connect) YJFBSweep() error {
	body, err := proto.Marshal(&C2SYJFBSweep{})
	if err != nil {
		return err
	}
	log.Println("[C][YJFBSweep]")
	return c.send(27207, body)
}

func (x *S2CYJFBSweep) ID() uint16 {
	return 27208
}

// Message S2CYJFBSweep 27208
func (x *S2CYJFBSweep) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][YJFBSweep] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

func (c *Connect) DigTreasure10Times(id int32) error {
	body, err := proto.Marshal(&C2SDigTreasure10Times{Id: id})
	if err != nil {
		return err
	}
	log.Printf("[C][DigTreasure10Times] id=%v", id)
	return c.send(22513, body)
}

func (x *S2CDigTreasure10Times) ID() uint16 {
	return 22514
}

// Message S2CDigTreasure10Times 22514
func (x *S2CDigTreasure10Times) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][DigTreasure10Times] tag=%v id=%v list=%v", x.Tag, x.Id, x.List)
}

////////////////////////////////////////////////////////////

func (c *Connect) GetYJFBGuanQiaData(fuBenId, guangQiaId int32) error {
	body, err := proto.Marshal(&C2SGetYJFBGuanQiaData{FuBenId: fuBenId, GuanQiaId: guangQiaId})
	if err != nil {
		return err
	}
	log.Printf("[C][GetYJFBGuanQiaData] fu_ben_id=%v guang_qia_id=%v", fuBenId, guangQiaId)
	return c.send(27201, body)
}

func (x *S2CGetYJFBGuanQiaData) ID() uint16 {
	return 27202
}

// Message S2CGetYJFBGuanQiaData 27202
func (x *S2CGetYJFBGuanQiaData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetYJFBGuanQiaData] tag=%v %v", x.Tag, x)
}

////////////////////////////////////////////////////////////

func (c *Connect) YJFBGuanQiaMove(m *C2SYJFBGuanQiaMove) error {
	body, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	log.Printf("[C][YJFBGuanQiaMove] fu_ben_id=%v guang_qia_id=%v x=%v y=%v", m.FuBenId, m.GuanQiaId, m.TargetGrid.X, m.TargetGrid.Y)
	return c.send(27203, body)
}

func (x *S2CYJFBGuanQiaMove) ID() uint16 {
	return 27204
}

// Message S2CYJFBGuanQiaMove 27204
func (x *S2CYJFBGuanQiaMove) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][YJFBGuanQiaMove] tag=%v %v", x.Tag, x)
}

////////////////////////////////////////////////////////////

func (c *Connect) YJFBGuanQiaTriggerEvent(e *C2SYJFBGuanQiaTriggerEvent) error {
	body, err := proto.Marshal(e)
	if err != nil {
		return err
	}
	log.Printf("[C][YJFBGuanQiaTriggerEvent] fu_ben_id=%v guang_qia_id=%v x=%v y=%v", e.FuBenId, e.GuanQiaId, e.TriggerGrid.X, e.TriggerGrid.Y)
	return c.send(27205, body)
}

func (x *S2CYJFBGuanQiaTriggerEvent) ID() uint16 {
	return 27206
}

// Message S2CYJFBGuanQiaTriggerEvent 27206
func (x *S2CYJFBGuanQiaTriggerEvent) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][YJFBGuanQiaTriggerEvent] tag=%v %v", x.Tag, x)
}

////////////////////////////////////////////////////////////

func (c *Connect) GetYJFBData() error {
	body, err := proto.Marshal(&C2SGetYJFBData{})
	if err != nil {
		return err
	}
	log.Println("[C][GetYJFBData]")
	return c.send(27209, body)
}

func (x *S2CGetYJFBData) ID() uint16 {
	return 27210
}

// Message S2CGetYJFBData 27210
func (x *S2CGetYJFBData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetYJFBData] tag=%v fb=%v", x.Tag, x.Fb)
}
