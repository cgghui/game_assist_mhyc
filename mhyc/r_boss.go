package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"sort"
	"time"
)

// BossPersonal 个人BOSS
func BossPersonal(ctx context.Context) {
	t := time.NewTimer(ms100)
	defer t.Stop()
	f := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "BOSS-个人BOSS")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		Receive.Action(CLI.BossPersonalSweep)
		ret := &S2CBossPersonalSweep{}
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, ret, s3); err != nil {
			return RandMillisecond(3, 6)
		}
		if ret.Tag == 4055 { // end
			return TomorrowDuration(RandMillisecond(600, 1800))
		}
		return RandMillisecond(600, 900)
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

// BossVIP 至尊BOSS
func BossVIP(ctx context.Context) {
	t := time.NewTimer(ms100)
	defer t.Stop()
	f := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "BOSS-BossVIP")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		Receive.Action(CLI.BossVipSweep)
		ret := &S2CBossVipSweep{}
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, ret, s3); err != nil {
			return RandMillisecond(3, 6)
		}
		if ret.Tag == 4055 { // end
			return TomorrowDuration(RandMillisecond(600, 1800))
		}
		return RandMillisecond(600, 900)
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

// BossXYCM 降妖除魔-奖励
func BossXYCM(ctx context.Context) {
	t := time.NewTimer(ms100)
	defer t.Stop()
	f := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "BOSS-降妖除魔-奖励")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		i := 0
		layer := []int32{10, 20, 30, 40, 50}
		count := len(layer)
		reTime := am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			go func() {
				_ = CLI.RecLimitFightSpeedReward(layer[i])
			}()
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CRecLimitFightSpeedReward{}, s3); err != nil {
				return 0, RandMillisecond(3, 6)
			}
			i++
			if i >= count {
				Receive.Action(CLI.RecLimitFightReward)
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CRecLimitFightReward{}, s3); err != nil {
					return 0, RandMillisecond(3, 6)
				}
				return 0, TomorrowDuration(RandMillisecond(1800, 3600))
			}
			return ms100, 0
		})
		return reTime
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

// BossXYCMGo 降妖除魔-战斗
func BossXYCMGo(ctx context.Context) {
	t := time.NewTimer(ms100)
	defer t.Stop()
	f := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "BOSS-降妖除魔-战斗")
		go func() {
			_ = CLI.JoinActive(&C2SJoinActive{AId: 101})
		}()
		join := &S2CJoinActive{}
		err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CJoinActive{}, s3)
		defer func() {
			go func() {
				_ = CLI.LeaveActive(&C2SLeaveActive{AId: 101})
			}()
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CLeaveActive{}, s3)
			am.End()
			Fight.Unlock()
		}()
		if err != nil || join.Tag != 0 {
			return RandMillisecond(3000, 3600)
		}
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			Receive.Action(CLI.ChallengeLimitFight)
			r := &S2CChallengeLimitFight{}
			if err = Receive.WaitWithContextOrTimeout(am.Ctx, r, s3); err != nil {
				return 0, RandMillisecond(600, 1800)
			}
			if r.Tag == 0 {
				return ms100, 0
			}
			return 0, RandMillisecond(3000, 3600)
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

// BossMulti 多人BOSS
func BossMulti(ctx context.Context) {
	t := time.NewTimer(ms100)
	defer t.Stop()
	f := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "BOSS-多人BOSS")
		defer func() {
			go func() {
				_ = CLI.MultiBossLeaveScene(&C2SMultiBossLeaveScene{Id: BossMultiID})
			}()
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CMultiBossLeaveScene{}, s3)
			am.End()
			Fight.Unlock()
		}()
		// 检测是否有挑战次数
		if RoleInfo.Get("MultiBoss_Times").Int64() == 0 {
			// 无
			if RoleInfo.Get("MultiBoss_Add_Times").Int64() == 10 {
				return TomorrowDuration(RandMillisecond(1800, 3600))
			}
			// 有
			mn := time.Unix(RoleInfo.Get("MultiBoss_NextTime").Int64(), 0).Local().Add(time.Minute)
			cur := time.Now()
			if cur.Before(mn) {
				return mn.Add(time.Second).Sub(cur)
			}
		}
		// 监听相关消息
		listens := []HandleMessage{&S2CMultiBossGetDamageLog{}, &S2CMultiBossLeaveScene{}, &S2CMultiBossInfo{}}
		for i := range listens {
			go ListenMessage(am.Ctx, listens[i])
		}
		// 获取BOSS信息
		info := &S2CMultiBossInfo{}
		Receive.Action(CLI.MultiBossInfo)
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, info, s3); err != nil || len(info.Items) == 0 {
			return RandMillisecond(3, 6)
		}
		for _, item := range info.Items {
			// 检测BOSS是否在冷却
			if item.Id == BossMultiID {
				rt := time.Unix(item.ReliveTimestamp, 0).Local().Add(time.Minute)
				cur := time.Now()
				if cur.Before(rt) {
					return rt.Add(time.Second).Sub(cur)
				}
			}
		}
		MonsterChan := make(chan *S2CMonsterEnterMap)
		go ListenMessageCall(am.Ctx, &S2CMonsterEnterMap{}, func(data []byte) {
			defer close(MonsterChan)
			var enter S2CMonsterEnterMap
			if err := proto.Unmarshal(data, &enter); err == nil {
				MonsterChan <- &enter
			}
		})
		go func() {
			_ = CLI.MultiBossJoinScene(&C2SMultiBossJoinScene{Id: BossMultiID})
		}()
		_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CMultiBossJoinScene{}, s30)
		for monster := range MonsterChan {
			return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
				ret, _ := FightAction(am.Ctx, monster.Id, 8)
				if ret == nil {
					return 0, time.Second
				}
				if ret.Tag == 4022 || ret.Tag == 17002 { // 逃跑
					return 0, time.Second
				}
				return ms500, 0
			})
		}
		return time.Second
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

// XuanShangBoss 悬赏BOSS
func XuanShangBoss(ctx context.Context) {
	t := time.NewTimer(ms100)
	defer t.Stop()
	f := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "BOSS-悬赏BOSS")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		info := &S2CXuanShangBossInfo{}
		Receive.Action(CLI.XuanShangBossInfo)
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, info, s3); err != nil {
			return ms100
		}
		// 没有挑战次数
		if info.LeftKillTimes <= 0 {
			Receive.Action(CLI.XuanShangBossScoreReward)
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CXuanShangBossScoreReward{}, s3); err != nil {
				return RandMillisecond(3, 6)
			}
			return TomorrowDuration(RandMillisecond(1800, 3600))
		}
		// BOSS级别小于3时，刷新BOSS级别
		if info.XuanShangID <= 3 {
			// 无刷新次数
			if info.LeftFreeRefreshTimes <= 0 {
				tt := time.Unix(info.NextFreeRefreshTimesTimeStamp, 0).Local().Add(s10)
				cur := time.Now()
				if cur.Before(tt) { // 等待刷新时间的到来
					return tt.Sub(cur)
				}
				return ms500
			}
			// 有刷新次数
			Receive.Action(CLI.XuanShangBossRefresh)
			refresh := &S2CXuanShangBossRefresh{}
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, refresh, s3); err != nil {
				return RandMillisecond(3, 6)
			}
			if refresh.XuanShangID <= 3 { // 刷新后乃然低于3，退出重来
				return ms100
			}
		}
		// 接受悬赏
		Receive.Action(CLI.XuanShangBossAccept)
		accept := &S2CXuanShangBossAccept{}
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, accept, s3); err != nil {
			return RandMillisecond(3, 6)
		}
		// 进入战场
		go func() {
			_ = CLI.XuanShangBossJoinScene(accept.BossID)
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CXuanShangBossJoinScene{}, s3); err != nil {
			return RandMillisecond(3, 6)
		}
		// 开始战斗
		FightAction(am.Ctx, accept.BossID, 8)
		return ms500
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

// worldBossActTime 世界BOSS活动时间
func worldBossActTime() time.Duration {
	cur := time.Now()
	y := cur.Year()
	m := cur.Month()
	d := cur.Day()
	actStartTime := []time.Time{
		time.Date(y, m, d, 10, 30, 0, 0, time.Local).Add(ms500),
		time.Date(y, m, d, 14, 30, 0, 0, time.Local).Add(ms500),
		time.Date(y, m, d, 16, 30, 0, 0, time.Local).Add(ms500),
		time.Date(y, m, d, 19, 30, 0, 0, time.Local).Add(ms500),
	}
	for _, ast := range actStartTime {
		if cur.Before(ast) {
			return ast.Sub(cur)
		}
		if cur.Before(ast.Add(time.Minute)) {
			return 0
		}
	}
	return TomorrowDuration(3 * time.Hour)
}

// WorldBoss 世界BOSS
func WorldBoss(ctx context.Context) {
	t := time.NewTimer(ms100)
	f := func() time.Duration {
		if td := worldBossActTime(); td != 0 {
			return td
		}
		Fight.Lock()
		am := SetAction(ctx, "BOSS-世界BOSS", 3*time.Minute)
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		// 进入前提前准备
		monster := make(chan *S2CMonsterEnterMap)
		go ListenMessageCall(am.Ctx, &S2CMonsterEnterMap{}, func(data []byte) {
			defer close(monster)
			var enter S2CMonsterEnterMap
			if err := proto.Unmarshal(data, &enter); err == nil {
				monster <- &enter
			} else {
				monster <- nil
			}
		})
		// 结束
		go ListenMessageCall(am.Ctx, &S2CWorldBossCloseScene{}, func(data []byte) {
			am.End()
		})
		// 等待 摇筛子
		go ListenMessageCall(am.Ctx, &S2CWorldBossBreakShieldInfo{}, func(data []byte) {
			r := &S2CWorldBossBreakShieldInfo{}
			r.Message(data)
			if r.MyState == 0 && r.MyPoints == 0 {
				go func() {
					go func() {
						_ = CLI.WorldBossStakePoints(1)
					}()
					_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CWorldBossBreakShieldInfo{}, s3)
				}()
			}
		})
		// 进入世界BOSS
		Receive.Action(CLI.BossGlobalJoinActive)
		join := &S2CJoinActive{}
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, join, s3); err != nil {
			return ms100
		}
		if join.Tag != 0 {
			return time.Second
		}
		// 离开
		defer func() {
			go func() {
				_ = CLI.LeaveActive(&C2SLeaveActive{AId: 2})
			}()
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CLeaveActive{}, s3)
		}()
		boss := <-monster
		if boss == nil {
			return ms300
		}
		am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			s, r := FightAction(am.Ctx, boss.Id, 8)
			if s == nil || r == nil {
				return ms100, 0
			}
			if r.Win == 1 {
				return 0, time.Second
			}
			return ms100, 0
		})
		i := int32(1)
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			go func() {
				_ = CLI.WorldBossReachGoalGetPrize(i)
			}()
			ret := &S2CWorldBossReachGoalGetPrize{}
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, ret, s3); err != nil {
				return 0, ms500
			}
			if ret.Tag != 0 {
				return 0, ms500
			}
			return ms100, 0
		})
	}
	//
	for {
		select {
		case <-t.C:
			t.Reset(f())
		case <-ctx.Done():
			return
		}
	}
}

func BossHome(ctx context.Context) {
	t := time.NewTimer(ms100)
	defer t.Stop()
	f := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "BOSS-BOSS之家")
		defer func() {
			Receive.Action(CLI.BossHomeLeaveScene)
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CBossHomeLeaveScene{}, s3)
			am.End()
			Fight.Unlock()
		}()
		// 地图怪
		monster := make([]*S2CMonsterEnterMap, 0)
		go monsterEnterMap(am.Ctx, &monster)
		// 怪信息
		bossInfoChan := make(chan *S2CHomeBossInfo)
		go func() {
			defer close(bossInfoChan)
			info := &S2CHomeBossInfo{}
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, info, s3); err != nil {
				bossInfoChan <- nil
			} else {
				bossInfoChan <- info
			}
		}()
		Receive.Action(CLI.BossHomeJoinScene)
		join := &S2CBossHomeJoinScene{}
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, join, s3); err != nil {
			return RandMillisecond(3, 6)
		}
		bossInfo := <-bossInfoChan // 等待BOSS信息返回
		// 挑战体力不足
		if join.Tag == 4049 {
			// 领取奖励，明天再战
			Receive.Action(CLI.HomeBossReceiveTempBag)
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CHomeBossReceiveTempBag{}, s3); err != nil {
				return RandMillisecond(3, 6)
			}
			return TomorrowDuration(RandMillisecond(1800, 3600))
		}
		<-time.After(2000 * time.Millisecond)
		// 地图内无怪时
		// 尝试等待所有怪冷却后再战
		if len(monster) == 0 {
			if bossInfo == nil {
				return RandMillisecond(3, 6)
			}
			var timeList = make([]int64, 0)
			for _, info := range bossInfo.Items {
				timeList = append(timeList, info.ReliveTimestamp)
			}
			if len(timeList) == 0 {
				return RandMillisecond(30, 60)
			}
			sort.Slice(timeList, func(i, j int) bool {
				return timeList[i] > timeList[j]
			})
			ttm := time.Unix(timeList[0], 0).Local()
			cur := time.Now()
			if cur.Before(ttm) {
				return ttm.Add(ms100).Sub(cur)
			}
			return RandMillisecond(20, 40)
		}
		// 打怪
		count := len(monster)
		i := 0
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			m := monster[i]
			_, r := FightAction(am.Ctx, m.Id, 8)
			if r == nil {
				return 0, ms100
			}
			if r.Win == 1 {
				i++
				if i >= count {
					return 0, RandMillisecond(3, 6)
				}
			}
			return ms100, 0
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

func BossXLD(ctx context.Context) {
	t := time.NewTimer(ms100)
	defer t.Stop()
	f := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "BOSS-BOSS凶灵岛")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		info := &S2CXLDBossInfo{}
		Receive.Action(CLI.XLDBossInfo)
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, info, s3); err != nil || len(info.Items) == 0 {
			return RandMillisecond(1, 3)
		}
		var timeList = make([]int64, 0)
		for _, item := range info.Items {
			timeList = append(timeList, item.NT)
		}
		sort.Slice(timeList, func(i, j int) bool {
			return timeList[i] > timeList[j]
		})
		ttm := time.Unix(timeList[0], 0).Local()
		cur := time.Now()
		if cur.Before(ttm) {
			return ttm.Add(ms100).Sub(cur)
		}
		bs := &S2CXLDBossSweep{}
		Receive.Action(CLI.XLDBossSweep)
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, bs, s3); err != nil {
			return RandMillisecond(1, 3)
		}
		if bs.Tag == 57015 {
			return TomorrowDuration(RandMillisecond(1800, 3600))
		}
		return RandMillisecond(1, 3)
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

// collectSC 采集仙草
func collectSC(ctx context.Context, field string, xsdID, bossID int32) time.Duration {
	Fight.Lock()
	name := ""
	if field == "XsdXsdDayCollectTimes" {
		name = "凶神岛"
	}
	if field == "XsdXmdDayCollectTimes" {
		name = "凶冥岛"
	}
	am := SetAction(ctx, "BOSS-采集仙草-"+name)
	defer func() {
		go func() {
			_ = CLI.XsdBossLeaveScene(&C2SXsdBossLeaveScene{XsdId: xsdID, BossId: bossID})
		}()
		_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CXsdBossLeaveScene{}, s3)
		am.End()
		Fight.Unlock()
	}()
	// 采集次数不足
	if RoleInfo.Get(field).Int64() >= 3 {
		return TomorrowDuration(RandMillisecond(1800, 3600))
	}
	// 怪信息
	bossInfoChan := make(chan *S2CXsdBossInfo)
	go func() {
		defer close(bossInfoChan)
		info := &S2CXsdBossInfo{}
		if err := Receive.Wait(info, s30); err != nil {
			bossInfoChan <- nil
		} else {
			bossInfoChan <- info
		}
	}()
	// 进入场景
	join := &S2CXsdBossJoinScene{}
	go func() {
		_ = CLI.XsdBossJoinScene(&C2SXsdBossJoinScene{XsdId: xsdID, BossId: bossID})
	}()
	if err := Receive.WaitWithContextOrTimeout(am.Ctx, join, s30); err != nil {
		return time.Second
	}
	//
	if field == "XsdXsdDayCollectTimes" {
		go func() {
			_ = CLI.DropItems(39051)
		}()
	}
	if field == "XsdXmdDayCollectTimes" {
		go func() {
			_ = CLI.DropItems(39053)
		}()
	}
	if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CGetDropItems{}, s3); err != nil {
		return time.Second
	}
	go func() {
		_ = CLI.StartMove(&C2SStartMove{P: []int32{13, 24}})
	}()
	if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CStartMove{}, s3); err != nil {
		return time.Second
	}
	//
	bossList := <-bossInfoChan // BOSS
	if bossList == nil {
		return ms100
	}
	for _, boss := range bossList.Items {
		if boss.BossId == bossID {
			if boss.State != 0 || boss.ReliveTimestamp == 0 {
				break
			}
			cur := time.Now()
			brt := time.Unix(boss.ReliveTimestamp, 0).Local()
			if cur.Before(brt) {
				return brt.Add(ms100).Sub(cur)
			}
			break
		}
	}
	// 采集仙草 视探
	go func() {
		_ = CLI.XsdCollect(&C2SXsdCollect{XsdId: xsdID, CollId: bossID, CollAct: 1})
	}()
	collect := &S2CXsdCollect{}
	if err := Receive.WaitWithContextOrTimeout(am.Ctx, collect, s3); err != nil {
		return time.Second
	}
	if collect.Tag == 0 && collect.CollState == 1 && RoleInfo.Get("UserId").Int64() != collect.CollUserId {
		return s30
	}
	// 采集仙草 采集
	go func() {
		_ = CLI.XsdCollect(&C2SXsdCollect{XsdId: xsdID, CollId: bossID, CollAct: 0})
	}()
	collect = &S2CXsdCollect{}
	if err := Receive.WaitWithContextOrTimeout(am.Ctx, collect, s90); err != nil {
		return time.Second
	}
	if field == "XsdXsdDayCollectTimes" {
		go func() {
			_ = CLI.RoutePath(&C2SRoutePath{MapId: 2555, FX: 13, FY: 24, TX: 38, TY: 73})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CRoutePath{}, s3); err != nil {
			return time.Second
		}
	}
	if field == "XsdXmdDayCollectTimes" {
		go func() {
			_ = CLI.RoutePath(&C2SRoutePath{MapId: 2566, FX: 13, FY: 24, TX: 14, TY: 70})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CRoutePath{}, s3); err != nil {
			return time.Second
		}
	}
	cur := time.Now()
	brt := time.Unix(collect.FinishTimestamp, 0).Local()
	if cur.Before(brt) {
		<-time.After(brt.Add(ms100).Sub(cur))
	}
	return ms100
}

// bossBattleScene BOSS战斗场景
func bossBattleScene(ctx context.Context, field string, xsdID, bossID int32) time.Duration {
	Fight.Lock()
	name := ""
	if field == "XsdXsdDayFightTimes" {
		name = "凶神岛之战"
	}
	if field == "XsdXmdDayFightTimes" {
		name = "凶冥岛之战"
	}
	am := SetAction(ctx, "BOSS-"+name)
	defer func() {
		go func() {
			_ = CLI.XsdBossLeaveScene(&C2SXsdBossLeaveScene{XsdId: xsdID, BossId: bossID})
		}()
		_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CXsdBossLeaveScene{}, s3)
		am.End()
		Fight.Unlock()
	}()
	if RoleInfo.Get(field).Int64() <= 0 {
		return TomorrowDuration(RandMillisecond(1800, 3600))
	}
	// 地图怪
	monster := make([]*S2CMonsterEnterMap, 0)
	go monsterEnterMap(am.Ctx, &monster)
	// 怪信息
	bossInfoChan := make(chan *S2CXsdBossInfo)
	go func() {
		defer close(bossInfoChan)
		info := &S2CXsdBossInfo{}
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, info, s3); err != nil {
			bossInfoChan <- nil
		} else {
			bossInfoChan <- info
		}
	}()
	join := &S2CXsdBossJoinScene{}
	go func() {
		_ = CLI.XsdBossJoinScene(&C2SXsdBossJoinScene{XsdId: xsdID, BossId: bossID})
	}()
	if err := Receive.WaitWithContextOrTimeout(am.Ctx, join, s3); err != nil {
		return RandMillisecond(1, 3)
	}
	<-bossInfoChan // 等待BOSS信息返回
	// 打怪
	<-time.After(time.Second)
	if len(monster) > 0 {
		// 按怪的血量排序，优先攻击血量多的怪（奖励多些）
		var HP = make([]int64, 0)
		for i := range monster {
			HP = append(HP, monster[i].Hp)
		}
		sort.Slice(HP, func(i, j int) bool {
			return HP[i] > HP[j]
		})
		count := len(HP)
		i := 0
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			if i >= count {
				return 0, time.Second
			}
			hp := HP[i]
			// 找到同等血量的怪
			idx := -1
			for j, m := range monster {
				if m.Hp == hp {
					idx = j
					break
				}
			}
			if idx == -1 {
				i++
				return ms100, 0
			}
			s, r := FightAction(am.Ctx, monster[idx].Id, 8)
			if s == nil || r == nil || s.Tag == 57006 || s.Tag == 57005 || s.Tag == 57016 { // 凶兽未解锁//
				i++
				return ms100, 0
			}
			if r.Win == 1 { // 斗报胜利
				i++
				return ms100, 0
			}
			return ms100, 0
		})
	}
	return s3
}

func BossXSD(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	t2 := time.NewTimer(ms100)
	defer t1.Stop()
	defer t2.Stop()
	for {
		select {
		case <-t1.C:
			t1.Reset(bossBattleScene(ctx, "XsdXsdDayFightTimes", 1, 1))
		case <-t2.C:
			t2.Reset(collectSC(ctx, "XsdXsdDayCollectTimes", 1, 7))
		case <-ctx.Done():
			return
		}
	}
}

func BossXMD(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	t2 := time.NewTimer(ms100)
	defer t1.Stop()
	defer t2.Stop()
	for {
		select {
		case <-t1.C:
			t1.Reset(bossBattleScene(ctx, "XsdXmdDayFightTimes", 2, 1))
		case <-t2.C:
			t1.Reset(collectSC(ctx, "XsdXmdDayCollectTimes", 2, 7))
		case <-ctx.Done():
			return
		}
	}
}

func BossHLTJ(ctx context.Context) {
	t1 := time.NewTimer(time.Second)
	f1 := func() time.Duration {
		if RoleInfo.Get("HLPower").Int64() <= 5 {
			return TomorrowDuration(RandMillisecond(1800, 3600))
		}
		Fight.Lock()
		am := SetAction(ctx, "BOSS-幻灵天界")
		defer func() {
			go func() {
				_ = CLI.LeaveHLFB(&C2SLeaveHLFB{InsId: HltjID})
			}()
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CLeaveHLFB{}, s3)
			am.End()
			Fight.Unlock()
		}()
		go func() {
			_ = CLI.C2SGetHLBossList(HltjID)
		}()
		bossList := &S2CGetHLBossList{}
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, bossList, s3); err != nil {
			if RoleInfo.Get("HLPower").Int64() <= 5 {
				return TomorrowDuration(RandMillisecond(1800, 3600))
			}
			return RandMillisecond(3, 6)
		}
		// 进入场景
		go func() {
			_ = CLI.EnterHLFB(&C2SEnterHLFB{InsId: HltjID, Type: 2})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CEnterHLFB{}, s3); err != nil {
			return RandMillisecond(1, 3)
		}
		// 组队
		go func() {
			_ = CLI.CreateTeam(&C2SCreateTeam{IsCross: 1, FuncId: 14105, Key1: 1, Key2: int64(HltjID), Key4: 0})
		}()
		var ct S2CCreateTeam
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &ct, s3); err != nil {
			return time.Second
		}
		if ct.Team == nil {
			return time.Second
		}
		defer func() {
			go func() {
				_ = CLI.LeaveTeam(ct.Team.TeamId)
			}()
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CLeaveTeam{}, s3)
		}()
		go func() {
			_ = CLI.Teams(&C2STeams{IsCross: 1, FuncId: 14105, Key1: 1, Key2: int64(HltjID), Key4: 0})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CTeams{}, s3); err != nil {
			return RandMillisecond(1, 3)
		}
		go func() {
			_ = CLI.InviteTeam(ct.Team.TeamId, 5)
		}()
		var teamInfo S2CTeamInfo // 等待成员加入
		ListenMessageCallEx(&S2CTeamInfo{}, func(data []byte) bool {
			teamInfo.Message(data)
			return len(teamInfo.Players) < HltjTeamRen
		})
		// 120402
		ReviveList := make([]int64, 0)
		for _, hl := range bossList.HLBossList {
			if hl.Revive == 0 {
				continue
			}
			ReviveList = append(ReviveList, hl.Revive)
		}
		if len(ReviveList) == 4 {
			sort.Slice(ReviveList, func(i, j int) bool {
				return ReviveList[i] > ReviveList[j]
			})
			ttm := time.Unix(ReviveList[0], 0).Local()
			cur := time.Now()
			if cur.Before(ttm) {
				return ttm.Add(time.Second).Sub(cur)
			}
			return s60
		}
		tc := time.NewTimer(ms500)
		defer tc.Stop()
		i := 0
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			if i >= len(bossList.HLBossList) {
				return 0, TomorrowDuration(RandMillisecond(1800, 3600))
			}
			boss := bossList.HLBossList[i]
			if boss.Revive != 0 {
				i++
				return ms500, 0
			}
			retChan := make(chan *S2CStartFightHLPVE)
			go func() {
				defer close(retChan)
				r := &S2CStartFightHLPVE{}
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, r, s3); err != nil {
					retChan <- nil
				} else {
					retChan <- r
				}
			}()
			go func() {
				_ = CLI.StartFightHLPVE(&C2SStartFightHLPVE{InsId: boss.InsId, BossId: int64(boss.Id)})
			}()
			r := &S2CBattlefieldReport{}
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, r, s3)
			p := <-retChan
			if p == nil {
				return 0, RandMillisecond(3, 6)
			}
			if p.Tag == 56713 { // 复活中
				ttm := time.Unix(RoleInfo.Get("ReviveTime").Int64(), 0).Local()
				cur := time.Now()
				if cur.Before(ttm) {
					return ttm.Add(time.Second).Sub(cur), 0
				}
				return 0, RandMillisecond(3, 6)
			}
			if p.Tag == 56714 {
				go func() {
					_ = CLI.WareHouseReceiveItem(2)
				}()
				_ = Receive.Wait(&S2CWareHouseReceiveItem{}, s3)
				return 0, TomorrowDuration(RandMillisecond(1800, 3600))
			}
			if r.Win == 1 && p.Tag == 0 {
				i++
			}
			return ms500, 0
		})
	}
	//
	for {
		select {
		case <-t1.C:
			t1.Reset(f1())
		case <-ctx.Done():
			return
		}
	}
}

func fightActionBDJJ(ctx context.Context, act func() error) (*S2CBangDanJJFight, *S2CBattlefieldReport) {
	c := make(chan *S2CBangDanJJFight)
	go func() {
		defer close(c)
		sf := &S2CBangDanJJFight{}
		if err := Receive.WaitWithContextOrTimeout(ctx, sf, s3); err != nil {
			c <- nil
		} else {
			c <- sf
		}
	}()
	Receive.Action(act)
	r := &S2CBattlefieldReport{}
	if err := Receive.WaitWithContextOrTimeout(ctx, r, s3); err != nil {
		r = nil
	}
	return <-c, r
}

func BossBDJJ(ctx context.Context) {
	t1 := time.NewTimer(time.Second)
	f1 := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "BOSS-榜单竞技")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		reTime := am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			f, _ := fightActionBDJJ(am.Ctx, CLI.C2SBangDanJJFight1)
			if f == nil {
				return 0, RandMillisecond(1, 3)
			}
			if f.Tag != 0 { // 60106 战斗次数不足
				return 0, 0
			}
			return ms500, 0
		})
		if reTime != 0 {
			return reTime
		}
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			f, _ := fightActionBDJJ(am.Ctx, CLI.C2SBangDanJJFight2)
			if f == nil {
				return 0, RandMillisecond(1, 3)
			}
			if f.Tag != 0 { // 60106 战斗次数不足
				return 0, TomorrowDuration(RandMillisecond(1800, 3600))
			}
			return ms500, 0
		})
	}
	//
	for {
		select {
		case <-t1.C:
			t1.Reset(f1())
		case <-ctx.Done():
			return
		}
	}
}

func monsterEnterMap(ctx context.Context, result *[]*S2CMonsterEnterMap) {
	ListenMessageCall(ctx, &S2CMonsterEnterMap{}, func(data []byte) {
		var enter S2CMonsterEnterMap
		if err := proto.Unmarshal(data, &enter); err == nil {
			*result = append(*result, &enter)
		}
	})
}

func (c *Connect) EnterHLFB(s *C2SEnterHLFB) error {
	body, err := proto.Marshal(s)
	if err != nil {
		return err
	}
	log.Printf("[C][EnterHLFB] ins_id=%v type=%v", s.InsId, s.Type)
	return c.send(27133, body)
}

func (c *Connect) LeaveHLFB(s *C2SLeaveHLFB) error {
	body, err := proto.Marshal(s)
	if err != nil {
		return err
	}
	log.Printf("[C][LeaveHLFB] ins_id=%v", s.InsId)
	return c.send(27135, body)
}

func (c *Connect) XsdBossLeaveScene(s *C2SXsdBossLeaveScene) error {
	body, err := proto.Marshal(s)
	if err != nil {
		return err
	}
	log.Printf("[C][XsdBossLeaveScene] xsd_id=%v boss_id=%v", s.XsdId, s.BossId)
	return c.send(15033, body)
}

func (c *Connect) BossHomeLeaveScene() error {
	body, err := proto.Marshal(&C2SBossHomeLeaveScene{HomeId: BossHomeID})
	if err != nil {
		return err
	}
	log.Println("[C][BossHomeLeaveScene]")
	return c.send(15033, body)
}

// BossPersonalSweep Boss - 本服BOSS - 个人BOSS 一键扫荡
func (c *Connect) BossPersonalSweep() error {
	body, err := proto.Marshal(&C2SBossPersonalSweep{})
	if err != nil {
		return err
	}
	log.Println("[C][BossPersonalSweep]")
	return c.send(604, body)
}

// BossVipSweep Boss - 本服BOSS - 至尊BOSS 一键扫荡
func (c *Connect) BossVipSweep() error {
	body, err := proto.Marshal(&C2SBossVipSweep{})
	if err != nil {
		return err
	}
	log.Println("[C][BossVipSweep]")
	return c.send(664, body)
}

// MultiBossPlayerInBoss Boss - 本服BOSS - 多人
func (c *Connect) MultiBossPlayerInBoss() error {
	body, err := proto.Marshal(&C2SMultiBossPlayerInBoss{})
	if err != nil {
		return err
	}
	log.Println("[C][MultiBossPlayerInBoss]")
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

func (c *Connect) MultiBossLeaveScene(i *C2SMultiBossLeaveScene) error {
	body, err := proto.Marshal(i)
	if err != nil {
		return err
	}
	log.Printf("[C][MultiBossLeaveScene] id=%v", i.Id)
	return c.send(1127, body)
}

// MultiBossInfo Boss - 本服BOSS - 多人 BOSS信息
func (c *Connect) MultiBossInfo() error {
	body, err := proto.Marshal(&C2SMultiBossInfo{})
	if err != nil {
		return err
	}
	log.Println("[C][MultiBossInfo]")
	return c.send(1130, body)
}

// XuanShangBossInfo Boss - 本服BOSS - 悬赏BOSS
func (c *Connect) XuanShangBossInfo() error {
	body, err := proto.Marshal(&C2SXuanShangBossInfo{})
	if err != nil {
		return err
	}
	log.Println("[C][XuanShangBossInfo]")
	return c.send(12451, body)
}

func (c *Connect) XuanShangBossRefresh() error {
	body, err := proto.Marshal(&C2SXuanShangBossRefresh{RefreshType: 0})
	if err != nil {
		return err
	}
	log.Println("[C][XuanShangBossRefresh]")
	return c.send(12455, body)
}

func (c *Connect) XuanShangBossAccept() error {
	body, err := proto.Marshal(&C2SXuanShangBossAccept{})
	if err != nil {
		return err
	}
	log.Println("[C][XuanShangBossAccept]")
	return c.send(12457, body)
}

func (c *Connect) XuanShangBossJoinScene(bossID int64) error {
	body, err := proto.Marshal(&C2SXuanShangBossJoinScene{BossID: bossID})
	if err != nil {
		return err
	}
	log.Printf("[C][XuanShangBossJoinScene] boss_id=%v", bossID)
	return c.send(12459, body)
}

func (c *Connect) XuanShangBossScoreReward() error {
	body, err := proto.Marshal(&C2SXuanShangBossScoreReward{})
	if err != nil {
		return err
	}
	log.Println("[C][XuanShangBossScoreReward]")
	return c.send(12466, body)
}

func (c *Connect) BossHomeJoinScene() error {
	body, err := proto.Marshal(&C2SBossHomeJoinScene{HomeId: BossHomeID})
	if err != nil {
		return err
	}
	log.Println("[C][BossHomeJoinScene]")
	return c.send(15031, body)
}

func (c *Connect) BossGlobalJoinActive() error {
	body, err := proto.Marshal(&C2SJoinActive{AId: 2})
	if err != nil {
		return err
	}
	log.Println("[C][GlobalJoinActive]")
	return c.send(1507, body)
}

func (c *Connect) HomeBossReceiveTempBag() error {
	body, err := proto.Marshal(&C2SHomeBossReceiveTempBag{})
	if err != nil {
		return err
	}
	log.Println("[C][HomeBossReceiveTempBag]")
	return c.send(15041, body)
}

func (c *Connect) XLDBossSweep() error {
	body, err := proto.Marshal(&C2SXLDBossSweep{Id: 1})
	if err != nil {
		return err
	}
	log.Println("[C][XLDBossSweep] XLD:1")
	return c.send(26205, body)
}

func (c *Connect) DropItems(id int32) error {
	body, err := proto.Marshal(&C2SGetDropItems{DropId: id})
	if err != nil {
		return err
	}
	log.Printf("[C][DropItemXLD] drop_id=%v", id)
	return c.send(26401, body)
}

func (c *Connect) XLDBossInfo() error {
	body, err := proto.Marshal(&C2SXLDBossInfo{})
	if err != nil {
		return err
	}
	log.Println("[C][XLDBossInfo]")
	return c.send(26201, body)
}

func (c *Connect) XsdBossJoinScene(join *C2SXsdBossJoinScene) error {
	body, err := proto.Marshal(join)
	if err != nil {
		return err
	}
	log.Println("[C][XsdBossJoinScene]")
	return c.send(26233, body)
}

func (c *Connect) XsdCollect(collect *C2SXsdCollect) error {
	body, err := proto.Marshal(collect)
	if err != nil {
		return err
	}
	log.Println("[C][XsdCollect]")
	return c.send(26245, body)
}

////////////////////////////////////////////////////////////

func (c *Connect) C2SBangDanJJFight1() error {
	body, err := proto.Marshal(&C2SBangDanJJFight{JJId: 1})
	if err != nil {
		return err
	}
	log.Println("[C][BangDanJJFight] jj_id=1")
	return c.send(29604, body)
}

func (c *Connect) C2SBangDanJJFight2() error {
	body, err := proto.Marshal(&C2SBangDanJJFight{JJId: 2})
	if err != nil {
		return err
	}
	log.Println("[C][BangDanJJFight] jj_id=2")
	return c.send(29604, body)
}

func (x *S2CBangDanJJFight) ID() uint16 {
	return 29609
}

// Message S2CBangDanJJFight 29609
func (x *S2CBangDanJJFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][BangDanJJFight] tag=%v tag_msg=%s %v", x.Tag, GetTagMsg(x.Tag), x)
}

////////////////////////////////////////////////////////////

func (x *S2CBossPersonalSweep) ID() uint16 {
	return 605
}

// Message S2CBossPersonalSweep 605
func (x *S2CBossPersonalSweep) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][BossPersonalSweep] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CBossVipSweep) ID() uint16 {
	return 665
}

// Message S2CBossVipSweep 665
func (x *S2CBossVipSweep) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][BossVipSweep] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CMultiBossJoinScene) ID() uint16 {
	return 1124
}

// Message S2CMultiBossJoinScene 1124
func (x *S2CMultiBossJoinScene) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][MultiBossJoinScene] tag=%v tag_msg=%s id=%v", x.Tag, GetTagMsg(x.Tag), x.Id)
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
	log.Printf("[S][MultiBossLeaveScene] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CXuanShangBossInfo) ID() uint16 {
	return 12452
}

// Message S2CXuanShangBossInfo 12452
// LeftKillTimes 剩余讨伐次数
// CurScore 当前积分
// LeftFreeRefreshTimes 剩余免费刷新
// NextFreeRefreshTimesTimeStamp 免费刷新品质（次数恢复）
func (x *S2CXuanShangBossInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XuanShangBossInfo] %v", x)
}

////////////////////////////////////////////////////////////

func (x *S2CXuanShangBossRefresh) ID() uint16 {
	return 12456
}

// Message S2CXuanShangBossRefresh 12456
func (x *S2CXuanShangBossRefresh) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XuanShangBossRefresh] tag=%v tag_msg=%s xuan_shang_id=%v", x.Tag, GetTagMsg(x.Tag), x.XuanShangID)
}

////////////////////////////////////////////////////////////

func (x *S2CXuanShangBossAccept) ID() uint16 {
	return 12458
}

// Message S2CXuanShangBossAccept 12458
func (x *S2CXuanShangBossAccept) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XuanShangBossRefresh] tag=%v tag_msg=%s boss_id=%v", x.Tag, GetTagMsg(x.Tag), x.BossID)
}

////////////////////////////////////////////////////////////

func (x *S2CXuanShangBossJoinScene) ID() uint16 {
	return 12460
}

// Message S2CXuanShangBossJoinScene 12460
func (x *S2CXuanShangBossJoinScene) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XuanShangBossJoinScene] tag=%v tag_msg=%s boss_id=%v", x.Tag, GetTagMsg(x.Tag), x.BossID)
}

////////////////////////////////////////////////////////////

func (x *S2CXuanShangBossScoreReward) ID() uint16 {
	return 12467
}

// Message S2CXuanShangBossScoreReward 12467
func (x *S2CXuanShangBossScoreReward) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XuanShangBossScoreReward] tag=%v tag_msg=%s score_reward=%v", x.Tag, GetTagMsg(x.Tag), x.ScoreRewardGet)
}

////////////////////////////////////////////////////////////

func (x *S2CBossHomeJoinScene) ID() uint16 {
	return 15032
}

// Message S2CBossHomeJoinScene 15032
func (x *S2CBossHomeJoinScene) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][BossHomeJoinScene] tag=%v tag_msg=%s home_id=%v", x.Tag, GetTagMsg(x.Tag), x.HomeId)
}

////////////////////////////////////////////////////////////

func (x *S2CHomeBossInfo) ID() uint16 {
	return 15030
}

// Message S2CHomeBossInfo 15030
func (x *S2CHomeBossInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][HomeBossInfo] %v", x)
}

////////////////////////////////////////////////////////////

func (x *S2CHomeBossReceiveTempBag) ID() uint16 {
	return 15042
}

// Message S2CHomeBossReceiveTempBag 15042
func (x *S2CHomeBossReceiveTempBag) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][HomeBossReceiveTempBag] %v", x)
}

////////////////////////////////////////////////////////////

func (x *S2CXLDBossSweep) ID() uint16 {
	return 26206
}

// Message S2CXLDBossSweep 26206
func (x *S2CXLDBossSweep) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XLDBossSweep] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CGetDropItems) ID() uint16 {
	return 26402
}

// Message S2CGetDropItems 26402
func (x *S2CGetDropItems) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XLDBossSweep] drop_id=%v items=%v", x.DropId, x.ItemData)
}

////////////////////////////////////////////////////////////

func (x *S2CXLDBossInfo) ID() uint16 {
	return 26202
}

// Message S2CXLDBossInfo 26202
func (x *S2CXLDBossInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XLDBossInfo] items=%v", x.Items)
}

////////////////////////////////////////////////////////////

func (x *S2CXsdBossJoinScene) ID() uint16 {
	return 26234
}

// Message S2CXsdBossJoinScene 26234
func (x *S2CXsdBossJoinScene) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XsdBossJoinScene] tag=%v tag_msg=%s xsd_id=%v boss_id=%v", x.Tag, GetTagMsg(x.Tag), x.XsdId, x.BossId)
}

////////////////////////////////////////////////////////////

func (x *S2CXsdBossInfo) ID() uint16 {
	return 26232
}

// Message S2CXsdBossInfo 26232
func (x *S2CXsdBossInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XsdBossInfo] items=%v", x.Items)
}

////////////////////////////////////////////////////////////

func (x *S2CXsdCollect) ID() uint16 {
	return 26246
}

// Message S2CXsdCollect 26246
func (x *S2CXsdCollect) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XsdCollect] tag=%v tag_msg=%s %v", x.Tag, GetTagMsg(x.Tag), x)
}

////////////////////////////////////////////////////////////

func (x *S2CBossHomeLeaveScene) ID() uint16 {
	return 15034
}

// Message S2CBossHomeLeaveScene 15034
func (x *S2CBossHomeLeaveScene) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][BossHomeLeaveScene] tag=%v tag_msg=%s %v", x.Tag, GetTagMsg(x.Tag), x)
}

////////////////////////////////////////////////////////////

func (x *S2CXsdBossLeaveScene) ID() uint16 {
	return 26236
}

// Message S2CXsdBossLeaveScene 15034
func (x *S2CXsdBossLeaveScene) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XsdBossLeaveScene] tag=%v tag_msg=%s %v", x.Tag, GetTagMsg(x.Tag), x)
}

////////////////////////////////////////////////////////////

func (x *S2CEnterHLFB) ID() uint16 {
	return 27134
}

// Message S2CEnterHLFB 27134
func (x *S2CEnterHLFB) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][EnterHLFB] tag=%v tag_msg=%s %v", x.Tag, GetTagMsg(x.Tag), x)
}

////////////////////////////////////////////////////////////

func (x *S2CLeaveHLFB) ID() uint16 {
	return 27136
}

// Message S2CLeaveHLFB 27136
func (x *S2CLeaveHLFB) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][LeaveHLFB] tag=%v tag_msg=%s %v", x.Tag, GetTagMsg(x.Tag), x)
}

////////////////////////////////////////////////////////////

func (c *Connect) StartFightHLPVE(s *C2SStartFightHLPVE) error {
	body, err := proto.Marshal(s)
	if err != nil {
		return err
	}
	log.Printf("[C][StartFightHLPVE] ins_id=%v boss_id=%v", s.InsId, s.BossId)
	return c.send(27137, body)
}

func (x *S2CStartFightHLPVE) ID() uint16 {
	return 27138
}

// Message S2CStartFightHLPVE 27138
func (x *S2CStartFightHLPVE) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][StartFightHLPVE] tag=%v tag_msg=%s tag_p=%v boss_id=%v", x.Tag, GetTagMsg(x.Tag), x.TagP, x.BossId)
}

////////////////////////////////////////////////////////////

func (c *Connect) C2SGetHLBossList(insID int32) error {
	body, err := proto.Marshal(&C2SGetHLBossList{InsId: insID})
	if err != nil {
		return err
	}
	log.Printf("[C][GetHLBossList] ins_id=%v", insID)
	return c.send(27131, body)
}

func (x *S2CGetHLBossList) ID() uint16 {
	return 27132
}

// Message S2CGetHLBossList 27132
func (x *S2CGetHLBossList) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetHLBossList] tag=%v tag_msg=%s hl_boss_list=%v", x.Tag, GetTagMsg(x.Tag), x.HLBossList)
}

////////////////////////////////////////////////////////////

func (c *Connect) RecLimitFightSpeedReward(layer int32) error {
	body, err := proto.Marshal(&C2SRecLimitFightSpeedReward{Layer: layer})
	if err != nil {
		return err
	}
	log.Printf("[C][RecLimitFightSpeedReward] layer=%v", layer)
	return c.send(24707, body)
}

func (x *S2CRecLimitFightSpeedReward) ID() uint16 {
	return 24708
}

// Message S2CRecLimitFightSpeedReward 24708
func (x *S2CRecLimitFightSpeedReward) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][RecLimitFightSpeedReward] tag=%v tag_msg=%s %v", x.Tag, GetTagMsg(x.Tag), x)
}

////////////////////////////////////////////////////////////

func (c *Connect) RecLimitFightReward() error {
	body, err := proto.Marshal(&C2SRecLimitFightReward{})
	if err != nil {
		return err
	}
	log.Println("[C][RecLimitFightReward]")
	return c.send(24705, body)
}

func (x *S2CRecLimitFightReward) ID() uint16 {
	return 24706
}

// Message S2CRecLimitFightReward 24706
func (x *S2CRecLimitFightReward) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][RecLimitFightReward] tag=%v tag_msg=%s %v", x.Tag, GetTagMsg(x.Tag), x)
}

////////////////////////////////////////////////////////////

func (c *Connect) RoutePath(p *C2SRoutePath) error {
	body, err := proto.Marshal(p)
	if err != nil {
		return err
	}
	log.Println("[C][RoutePath]")
	return c.send(154, body)
}

func (x *S2CRoutePath) ID() uint16 {
	return 155
}

// Message S2CRoutePath 155
func (x *S2CRoutePath) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][RoutePath] tag=%v tag_msg=%s map_id=%v points=%v", x.Tag, GetTagMsg(x.Tag), x.MapId, x.Points)
}

////////////////////////////////////////////////////////////

func (x *S2CWorldBossLevel) ID() uint16 {
	return 15017
}

// Message S2CWorldBossLevel 15017
func (x *S2CWorldBossLevel) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][WorldBossLevel] level=%v", x.Level)
}

////////////////////////////////////////////////////////////

func (x *S2CWorldBossEnd) ID() uint16 {
	return 15018
}

// Message S2CWorldBossEnd 15018
func (x *S2CWorldBossEnd) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][WorldBossEnd] tag=%v tag_msg=%s scene_close_time=%v", x.Tag, GetTagMsg(x.Tag), x.SceneCloseTime)
}

////////////////////////////////////////////////////////////

func (x *S2CWorldBossCloseScene) ID() uint16 {
	return 15019
}

// Message S2CWorldBossCloseScene 15019
func (x *S2CWorldBossCloseScene) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][WorldBossCloseScene] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CWorldBossBreakShieldInfo) ID() uint16 {
	return 15014
}

// Message S2CWorldBossBreakShieldInfo 15014
func (x *S2CWorldBossBreakShieldInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	t := time.Unix(x.OverTimestamp, 0).Local().Format("2006-01-02 15:04:05")
	log.Printf("[S][WorldBossBreakShieldInfo] over_time=%v", t)
}

////////////////////////////////////////////////////////////

func (c *Connect) WorldBossStakePoints(op int32) error {
	body, err := proto.Marshal(&C2SWorldBossStakePoints{Op: op})
	if err != nil {
		return err
	}
	log.Printf("[C][WorldBossStakePoints] op=%d", op)
	return c.send(15015, body)
}

func (x *S2CWorldBossStakePoints) ID() uint16 {
	return 15016
}

// Message S2CWorldBossStakePoints 15016
func (x *S2CWorldBossStakePoints) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][WorldBossStakePoints] tag=%v tag_msg=%s state=%v points=%v", x.Tag, GetTagMsg(x.Tag), x.State, x.Points)
}

////////////////////////////////////////////////////////////

func (c *Connect) WorldBossReachGoalGetPrize(id int32) error {
	body, err := proto.Marshal(&C2SWorldBossReachGoalGetPrize{Id: id})
	if err != nil {
		return err
	}
	log.Printf("[C][WorldBossReachGoalGetPrize] id=%d", id)
	return c.send(15012, body)
}

func (x *S2CWorldBossReachGoalGetPrize) ID() uint16 {
	return 15013
}

// Message S2CWorldBossReachGoalGetPrize 15013
func (x *S2CWorldBossReachGoalGetPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][WorldBossReachGoalGetPrize] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (c *Connect) ChallengeLimitFight() error {
	body, err := proto.Marshal(&C2SChallengeLimitFight{})
	if err != nil {
		return err
	}
	log.Println("[C][ChallengeLimitFight]")
	return c.send(24703, body)
}

func (x *S2CChallengeLimitFight) ID() uint16 {
	return 24704
}

// Message S2CChallengeLimitFight 24704
func (x *S2CChallengeLimitFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ChallengeLimitFight] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}
