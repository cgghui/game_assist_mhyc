package mhyc

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"sort"
	"time"
)

const BossMultiID = 8 // 多人BOSS
const BossHomeID = 7  // 跨服 - BOSS之家 - 7层

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
		Fight.Lock()
		defer Fight.Unlock()
		// 检测是否有挑战次数
		if RoleInfo.Get("MultiBoss_Times").Int64() == 0 {
			// 无
			if RoleInfo.Get("MultiBoss_Add_Times").Int64() == 10 {
				return TomorrowDuration(RandMillisecond(30000, 30600))
			}
			// 有
			mn := time.Unix(RoleInfo.Get("MultiBoss_NextTime").Int64(), 0).Local().Add(time.Minute)
			cur := time.Now()
			if cur.Before(mn) {
				return mn.Add(time.Minute).Sub(cur)
			}
		}
		// 监听相关消息
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		listens := []HandleMessage{
			&S2CMultiBossGetDamageLog{},
			&S2CMultiBossLeaveScene{},
			&S2CMultiBossInfo{},
		}
		for i := range listens {
			go ListenMessage(ctx, listens[i])
		}
		// 获取BOSS信息
		info := &S2CMultiBossInfo{}
		Receive.Action(CLI.MultiBossInfo)
		if err := Receive.Wait(info, s3); err != nil {
			return ms100
		}
		for _, item := range info.Items {
			// 检测BOSS是否在冷却
			if item.Id == int32(BossMultiID) {
				rt := time.Unix(item.ReliveTimestamp, 0).Local().Add(time.Minute)
				cur := time.Now()
				if cur.Before(rt) {
					return rt.Add(time.Minute).Sub(cur)
				}
			}
		}
		go func() {
			_ = CLI.MultiBossJoinScene(&C2SMultiBossJoinScene{Id: int32(BossMultiID)})
		}()
		go func() {
			_ = Receive.Wait(&S2CMultiBossJoinScene{}, s30)
		}()
		enter := &S2CMonsterEnterMap{}
		_ = Receive.Wait(enter, s30)
		// loop 战斗
		tc := time.NewTimer(0)
		defer tc.Stop()
		for range tc.C {
			ret := &S2CStartFight{}
			go func() {
				_ = CLI.StartFight(&C2SStartFight{Id: enter.Id, Type: int64(BossMultiID)})
			}()
			if err := Receive.Wait(ret, s3); err != nil {
				return ms100
			}
			if ret.Tag == 4022 { // 逃跑
				break
			}
			tc.Reset(ms500)
		}
		//
		return ms100
	}
	for range t.C {
		t.Reset(f())
	}
}

// XuanShangBoss 悬赏BOSS
func XuanShangBoss() {
	t := time.NewTimer(ms100)
	f := func() time.Duration {
		Fight.Lock()
		defer Fight.Unlock()
		info := &S2CXuanShangBossInfo{}
		Receive.Action(CLI.XuanShangBossInfo)
		if err := Receive.Wait(info, s3); err != nil {
			return ms100
		}
		// 没有挑战次数
		if info.LeftKillTimes <= 0 {
			Receive.Action(CLI.XuanShangBossScoreReward)
			_ = Receive.Wait(&S2CXuanShangBossScoreReward{}, s3)
			return TomorrowDuration(30600 * time.Second)
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
			refresh := &S2CXuanShangBossRefresh{}
			Receive.Action(CLI.XuanShangBossRefresh)
			_ = Receive.Wait(refresh, s3)
			if refresh.XuanShangID <= 3 { // 刷新后乃然低于3，退出重来
				return ms100
			}
		}
		// 接受悬赏
		accept := &S2CXuanShangBossAccept{}
		Receive.Action(CLI.XuanShangBossAccept)
		if err := Receive.Wait(accept, s3); err != nil {
			return ms100
		}
		// 进入战场
		go func() {
			_ = CLI.XuanShangBossJoinScene(accept.BossID)
		}()
		_ = Receive.Wait(&S2CXuanShangBossJoinScene{}, s3)
		// 开始战斗
		go func() {
			_ = CLI.StartFight(&C2SStartFight{Id: accept.BossID, Type: 8})
		}()
		_ = Receive.Wait(&S2CBattlefieldReport{}, s3)
		//
		return ms500
	}
	for range t.C {
		t.Reset(f())
	}
}

func BossGlobal() {
	t := time.NewTimer(ms100)
	f := func() time.Duration {
		Fight.Lock()
		defer Fight.Unlock()
		Receive.Action(CLI.BossGlobalJoinActive)
		info := &S2CJoinActive{}
		if err := Receive.Wait(info, s3); err != nil {
			return ms100
		}
		if info.Tag != 0 {
			return time.Second
		}
		go func() {
			_ = CLI.StartFight(&C2SStartFight{Id: 385, Type: 8})
		}()
		_ = Receive.Wait(&S2CBattlefieldReport{}, s3)
		return ms500
	}
	for range t.C {
		t.Reset(f())
	}
}

func BossHome() {
	t := time.NewTimer(ms100)
	f := func() time.Duration {
		if RoleInfo.Get("BossHome_BodyPower").Int64() < 10 {
			// 领取奖励，明天再战
			Receive.Action(CLI.HomeBossReceiveTempBag)
			_ = Receive.Wait(&S2CHomeBossReceiveTempBag{}, s3)
			return TomorrowDuration(RandMillisecond(30000, 30600))
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// 地图怪
		monster := make([]*S2CMonsterEnterMap, 0)
		go ListenMessageCall(ctx, &S2CMonsterEnterMap{}, func(data []byte) {
			var r S2CMonsterEnterMap
			_ = proto.Unmarshal(data, &r)
			monster = append(monster, &r)
		})
		// 怪信息
		bossInfoChan := make(chan *S2CHomeBossInfo)
		defer close(bossInfoChan)
		go func() {
			info := &S2CHomeBossInfo{}
			if err := Receive.Wait(info, s3); err != nil {
				bossInfoChan <- nil
			} else {
				bossInfoChan <- info
			}
		}()
		join := &S2CBossHomeJoinScene{}
		Receive.Action(CLI.BossHomeJoinScene)
		_ = Receive.Wait(join, s3)
		bossInfo := <-bossInfoChan // 等待BOSS信息返回
		// 挑战体力不足
		if join.Tag == 4049 {
			return ms100
		}
		// 地图内无怪时
		// 尝试等待所有怪冷却后再战
		if len(monster) == 0 {
			if bossInfo == nil {
				return ms500
			}
			var timeList = make([]int64, 0)
			for _, info := range bossInfo.Items {
				timeList = append(timeList, info.ReliveTimestamp)
			}
			sort.Slice(timeList, func(i, j int) bool {
				return timeList[i] > timeList[j]
			})
			ttm := time.Unix(timeList[0], 0).Local()
			cur := time.Now()
			if cur.Before(ttm) {
				return ttm.Add(s10).Sub(cur)
			}
			return s30
		}
		// 打怪
		for i := range monster {
			for {
				go func(i int) {
					_ = CLI.StartFight(&C2SStartFight{Id: monster[i].Id, Type: 8})
				}(i)
				r := &S2CBattlefieldReport{}
				if err := Receive.Wait(r, s6); err != nil { // 无战斗报告反馈
					return ms100
				}
				if r.Win == 1 {
					break
				}
			}
		}
		return ms500
	}
	for range t.C {
		t.Reset(f())
	}
}

func BossXLD() {
	t := time.NewTimer(ms100)
	f := func() time.Duration {
		info := &S2CXLDBossInfo{}
		Receive.Action(CLI.XLDBossInfo)
		if err := Receive.Wait(info, s3); err != nil {
			return ms100
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
			return ttm.Add(s10).Sub(cur)
		}
		bs := &S2CXLDBossSweep{}
		Receive.Action(CLI.XLDBossSweep)
		if err := Receive.Wait(bs, s3); err != nil {
			return ms100
		}
		if bs.Tag == 57015 {
			return TomorrowDuration(RandMillisecond(30000, 30600))
		}
		return time.Second
	}
	for range t.C {
		t.Reset(f())
	}
}

func BossXSD() {
	t := time.NewTimer(ms100)
	f := func() time.Duration {
		if RoleInfo.Get("XsdXsdDayFightTimes").Int64() <= 0 {
			return TomorrowDuration(RandMillisecond(30000, 30600))
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// 地图怪
		monster := make([]*S2CMonsterEnterMap, 0)
		go ListenMessageCall(ctx, &S2CMonsterEnterMap{}, func(data []byte) {
			var r S2CMonsterEnterMap
			_ = proto.Unmarshal(data, &r)
			monster = append(monster, &r)
		})
		// 怪信息
		bossInfoChan := make(chan *S2CXsdBossInfo)
		defer close(bossInfoChan)
		go func() {
			info := &S2CXsdBossInfo{}
			if err := Receive.Wait(info, s3); err != nil {
				bossInfoChan <- nil
			} else {
				bossInfoChan <- info
			}
		}()
		join := &S2CXsdBossJoinScene{}
		go func() {
			_ = CLI.XsdBossJoinScene(&C2SXsdBossJoinScene{XsdId: 1, BossId: 1})
		}()
		_ = Receive.Wait(join, s3)
		boss := <-bossInfoChan // 等待BOSS信息返回
		fmt.Println(boss)
		// 打怪
		if len(monster) > 0 {
			// 按怪的血量排序，优先攻击血量多的怪（奖励多些）
			var HP = make([]int64, 0)
			for i := range monster {
				HP = append(HP, monster[i].Hp)
			}
			sort.Slice(HP, func(i, j int) bool {
				return HP[i] > HP[j]
			})
			for _, hp := range HP {
				// 找到同等血量的怪
				idx := -1
				for i, m := range monster {
					if m.Hp == hp {
						idx = i
						break
					}
				}
				if idx == -1 {
					continue
				}
				// 开打
				for {
					go func(i int) {
						_ = CLI.StartFight(&C2SStartFight{Id: monster[i].Id, Type: 8})
					}(idx)
					r := &S2CBattlefieldReport{}
					if err := Receive.Wait(r, s3); err != nil {
						return ms100
					}
					if r.Win == 1 { // 斗报胜利
						break
					}
				}
			}
		}
		return s3
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
	log.Printf("[S][XuanShangBossRefresh] tag=%v xuan_shang_id=%v", x.Tag, x.XuanShangID)
}

////////////////////////////////////////////////////////////

func (x *S2CXuanShangBossAccept) ID() uint16 {
	return 12458
}

// Message S2CXuanShangBossAccept 12458
func (x *S2CXuanShangBossAccept) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XuanShangBossRefresh] tag=%v boss_id=%v", x.Tag, x.BossID)
}

////////////////////////////////////////////////////////////

func (x *S2CXuanShangBossJoinScene) ID() uint16 {
	return 12460
}

// Message S2CXuanShangBossJoinScene 12460
func (x *S2CXuanShangBossJoinScene) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XuanShangBossJoinScene] tag=%v boss_id=%v", x.Tag, x.BossID)
}

////////////////////////////////////////////////////////////

func (x *S2CXuanShangBossScoreReward) ID() uint16 {
	return 12467
}

// Message S2CXuanShangBossScoreReward 12467
func (x *S2CXuanShangBossScoreReward) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XuanShangBossScoreReward] tag=%v score_reward=%v", x.Tag, x.ScoreRewardGet)
}

////////////////////////////////////////////////////////////

func (x *S2CBossHomeJoinScene) ID() uint16 {
	return 15032
}

// Message S2CBossHomeJoinScene 15032
func (x *S2CBossHomeJoinScene) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][BossHomeJoinScene] tag=%v home_id=%v", x.Tag, x.HomeId)
}

////////////////////////////////////////////////////////////

func (x *S2CJoinActive) ID() uint16 {
	return 1508
}

// Message S2CJoinActive 1508
func (x *S2CJoinActive) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JoinActive] tag=%v aid=%v", x.Tag, x.AId)
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
	log.Printf("[S][XLDBossSweep] tag=%v", x.Tag)
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
	log.Printf("[S][XsdBossJoinScene] tag=%v xsd_id=%v boss_id=%v", x.Tag, x.XsdId, x.BossId)
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
	log.Printf("[S][XsdCollect] tag=%v %v", x.Tag, x)
}
