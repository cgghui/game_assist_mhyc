package mhyc

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"sort"
	"time"
)

type CfgPrefectWarBoss struct {
	BossId int32
	CpId   int32
	Name   string
}

var cfgPrefectWarBossList []CfgPrefectWarBoss

func init() {
	if err := json.Unmarshal(cfg1PrefectWarBoss, &cfgPrefectWarBossList); err != nil {
		panic(err)
	}
}

func illusionSweep(ctx context.Context) time.Duration {
	Fight.Lock()
	am := SetAction(ctx, "跨服-幻境一键扫荡")
	defer func() {
		am.End()
		Fight.Unlock()
	}()
	irList := []int32{1, 2, 11, 21}
	i := 0
	return am.RunAction(am.Ctx, func() (loop time.Duration, next time.Duration) {
		if i >= len(irList) {
			return 0, TomorrowDuration(RandMillisecond(1800, 3600))
		}
		go func() {
			_ = CLI.IllusionSweep(&C2SIllusionSweep{IllusionType: irList[i]})
		}()
		r := &S2CIllusionSweep{}
		_ = Receive.WaitWithContextOrTimeout(am.Ctx, r, s10)
		i++
		return ms100, 0
	})
}

func yiJi(ctx context.Context) time.Duration {
	Fight.Lock()
	am := SetAction(ctx, "跨服-遗迹之战")
	defer func() {
		am.End()
		Fight.Unlock()
	}()
	if RoleInfo.Get("YiJiAdscTimes").Int64() <= 0 {
		return TomorrowDuration(RandMillisecond(1800, 3600))
	}
	Receive.Action(CLI.YiJiInfo)
	var info S2CYiJiInfo
	if err := Receive.WaitWithContextOrTimeout(am.Ctx, &info, s3); err != nil {
		return RandMillisecond(1, 3)
	}
	if len(info.Items) == 0 {
		return time.Second
	}
	timeList := make([]int64, 0)
	id := int32(0)
	for _, item := range info.Items {
		if item.ReliveTimestamp > 0 {
			timeList = append(timeList, item.ReliveTimestamp)
		}
		if item.BossState == 1 {
			id = item.Id
			break
		}
	}
	if id == 0 {
		if len(timeList) == 0 {
			return time.Second
		}
		sort.Slice(timeList, func(i, j int) bool {
			return timeList[i] < timeList[j]
		})
		ttm := time.Unix(timeList[0], 0).Local()
		cur := time.Now()
		if cur.Before(ttm) {
			return ttm.Add(time.Second).Sub(cur)
		}
		return ms500
	}
	// 进入
	defer func() {
		go func() {
			_ = CLI.YiJiLeaveScene(id)
		}()
		_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CYiJiLeaveScene{}, s3)
	}()
	mc := make(chan *S2CMonsterEnterMap)
	go func() {
		go func() {
			var monster S2CMonsterEnterMap
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, &monster, s3); err != nil {
				mc <- nil
			} else {
				mc <- &monster
			}
			close(mc)
		}()
		_ = CLI.YiJiJoinScene(id)
	}()
	var join S2CYiJiJoinScene
	if err := Receive.WaitWithContextOrTimeout(am.Ctx, &join, s3); err != nil {
		return RandMillisecond(0, 3)
	}
	monster := <-mc
	if monster == nil {
		return RandMillisecond(0, 3)
	}
	return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
		s, r := FightAction(am.Ctx, monster.Id, 8)
		if s == nil {
			return 0, RandMillisecond(0, 3)
		}
		if s.Tag == 17002 {
			return 0, RandMillisecond(0, 3)
		}
		if r != nil && r.Win == 1 {
			return 0, RandMillisecond(0, 3)
		}
		return time.Second, 0
	})
}

func yiJiWT(ctx context.Context) time.Duration {
	Fight.Lock()
	am := SetAction(ctx, "跨服-遗迹委托之战")
	defer func() {
		am.End()
		Fight.Unlock()
	}()
	a1 := RoleInfo.Get("YiJiDayBeEntrustedTimes").Int64()
	a2 := RoleInfo.Get("YiJiBeEntrustedNum").Int64()
	a3 := RoleInfo.Get("YiJiDayEntrustTimes").Int64()
	if a1 == 0 && a2 == 0 && a3 == 0 {
		return TomorrowDuration(RandMillisecond(1800, 3600))
	}
	if a2 < 3 {
		Receive.Action(CLI.GetEntrustWallList)
		var list S2CGetEntrustWallList
		_ = Receive.WaitWithContextOrTimeout(am.Ctx, &list, s3)
		if len(list.List) > 0 {
			for i := range list.List {
				go func(i int) {
					_ = CLI.ReceiveEntrust(list.List[i].User.UserId)
				}(i)
				var rec S2CReceiveEntrust
				_ = Receive.WaitWithContextOrTimeout(am.Ctx, &rec, s3)
			}
		}
		if a2 <= 0 {
			return RandMillisecond(1, 3)
		}
	}

	Receive.Action(CLI.YiJiInfo)
	var info S2CYiJiInfo
	if err := Receive.WaitWithContextOrTimeout(am.Ctx, &info, s3); err != nil {
		return RandMillisecond(1, 3)
	}
	if len(info.Items) == 0 {
		return time.Second
	}
	timeList := make([]int64, 0)
	id := int32(0)
	for _, item := range info.Items {
		if item.ReliveTimestamp > 0 {
			timeList = append(timeList, item.ReliveTimestamp)
		}
		if item.BossState == 1 {
			id = item.Id
			break
		}
	}
	if id == 0 {
		if len(timeList) == 0 {
			return time.Second
		}
		sort.Slice(timeList, func(i, j int) bool {
			return timeList[i] < timeList[j]
		})
		ttm := time.Unix(timeList[0], 0).Local()
		cur := time.Now()
		if cur.Before(ttm) {
			return ttm.Add(time.Second).Sub(cur)
		}
		return ms500
	}
	// 进入
	defer func() {
		go func() {
			_ = CLI.YiJiLeaveScene(id)
		}()
		_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CYiJiLeaveScene{}, s3)
	}()
	mc := make(chan *S2CMonsterEnterMap)
	go func() {
		go func() {
			var monster S2CMonsterEnterMap
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, &monster, s3); err != nil {
				mc <- nil
			} else {
				mc <- &monster
			}
			close(mc)
		}()
		_ = CLI.YiJiJoinScene(id)
	}()
	var join S2CYiJiJoinScene
	if err := Receive.WaitWithContextOrTimeout(am.Ctx, &join, s3); err != nil {
		return RandMillisecond(0, 3)
	}
	monster := <-mc
	if monster == nil {
		return RandMillisecond(0, 3)
	}
	return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
		s, r := FightAction(am.Ctx, monster.Id, 8)
		if s == nil {
			return 0, RandMillisecond(0, 3)
		}
		if s.Tag == 17002 {
			return 0, RandMillisecond(0, 3)
		}
		if r != nil && r.Win == 1 {
			return 0, RandMillisecond(0, 3)
		}
		return time.Second, 0
	})
}

func KFZZJZ(ctx context.Context) time.Duration {
	Fight.Lock()
	am := SetAction(ctx, "跨服-征战九州")
	defer func() {
		am.End()
		Fight.Unlock()
	}()
	Receive.Action(CLI.GetPrefectWarData)
	data := S2CGetPrefectWarData{}
	_ = Receive.WaitWithContextOrTimeout(am.Ctx, &data, s3)
	bossId := int32(3)
	cp := int32(0)
	pt := int32(0)
	for _, item := range data.TabItems {
		if item.TabId == 0 {
			cp = item.CurCpId
			pt = item.PassTimes
			if item.CurCpId == 9 {
				ttm := time.Unix(item.StateResetTimestamp, 0)
				cur := time.Now()
				if cur.Before(ttm) {
					//
					go func() {
						_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 44, Multi: 1, TaskId: 1})
					}()
					_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CGetTaskPrize{}, s3)
					//
					go func() {
						_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 44, Multi: 1, TaskId: 2})
					}()
					_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CGetTaskPrize{}, s3)
					//
					go func() {
						_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 44, Multi: 1, TaskId: 3})
					}()
					_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CGetTaskPrize{}, s3)
					//
					go func() {
						_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 44, Multi: 1, TaskId: 4})
					}()
					_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CGetTaskPrize{}, s3)
					//
					go func() {
						_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 44, Multi: 1, TaskId: 5})
					}()
					_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CGetTaskPrize{}, s3)
					//
					return ttm.Add(ms100).Sub(cur)
				}
			}
		}
	}
	cp++
	go func() {
		_ = CLI.CreateTeam(&C2SCreateTeam{Key1: int64(cp), Key2: int64(bossId), Key3: int64(pt), Key4: 0, IsCross: 1, FightLimit: 0, FuncId: 829})
	}()
	var ct S2CCreateTeam
	if err := Receive.WaitWithContextOrTimeout(am.Ctx, &ct, s3); err != nil {
		return RandMillisecond(1, 3)
	}
	tNext := time.NewTimer(s15)
	defer tNext.Stop()
	go ListenMessageCall(am.Ctx, &S2CTeamInfo{}, func(data []byte) {
		info := &S2CTeamInfo{}
		info.Message(data)
		if len(info.Players) >= 3 {
			tNext.Reset(ms10)
		}
	})
	go ListenMessageCall(am.Ctx, &S2CDisbandTeam{}, func(_ []byte) {
		am.End()
	})
	name := ""
	for _, boss := range cfgPrefectWarBossList {
		if boss.BossId == bossId && boss.CpId == cp {
			name = boss.Name
			break
		}
	}
	param1 := "<color=#46ff69>豪杰</c>·征战九州·" + name
	param2 := fmt.Sprintf("%d|%d|%d|%d|%d", ct.Team.FuncId, ct.Team.Key1, ct.Team.Key2, 0, ct.Team.TeamId)
	params := []C2SCommonShout{
		{
			Param1:    param1,
			Param2:    param2,
			NoticeId:  288,
			ChannelId: 10003,
		},
		{
			Param1:    param1,
			Param2:    param2,
			NoticeId:  288,
			ChannelId: 10004,
		},
		{
			Param1:    param1,
			Param2:    param2,
			NoticeId:  288,
			ChannelId: 10005,
		},
		{
			Param1:    param1,
			Param2:    param2,
			NoticeId:  288,
			ChannelId: 10006,
		},
	}
	i := 0
	count := len(params)
	am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
		if i >= count {
			return 0, 0
		}
		go func(i int) {
			_ = CLI.CommonShout(&params[i])
		}(i)
		_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CCommonShout{}, s3)
		i++
		return ms100, 0
	})
	am.TimeWait(ctx, tNext)
	return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
		if cp >= 10 {
			return 0, time.Minute
		}
		go func() {
			_ = CLI.PrefectWarFight(&C2SPrefectWarFight{TabId: 0, BossId: 3, CpId: cp})
		}()
		ret := &S2CPrefectWarFight{}
		_ = Receive.WaitWithContextOrTimeout(am.Ctx, ret, s3)
		if ret.Tag != 0 {
			return time.Second, 0
		}
		tNext.Reset(9 * time.Second)
		am.TimeWait(ctx, tNext)
		Receive.Action(CLI.GetPrefectWarData)
		_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CGetPrefectWarData{}, s3)
		cp++
		return s3, 0
	})
}

// KuaFu 跨服
func KuaFu(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	t2 := time.NewTimer(ms50)
	t3 := time.NewTimer(ms100)
	t4 := time.NewTimer(ms10)
	defer t1.Stop()
	defer t2.Stop()
	defer t4.Stop()
	for {
		select {
		case <-t1.C:
			t1.Reset(yiJi(ctx))
		case <-t2.C:
			t2.Reset(yiJiWT(ctx))
		case <-t3.C:
			t3.Reset(KFZZJZ(ctx))
		case <-t4.C:
			t4.Reset(illusionSweep(ctx))
		case <-ctx.Done():
			return
		}
	}
}

////////////////////////////////////////////////////////////

// YiJiLeaveScene 离开遗迹场景
func (c *Connect) YiJiLeaveScene(id int32) error {
	body, err := proto.Marshal(&C2SYiJiLeaveScene{Id: id})
	if err != nil {
		return err
	}
	log.Printf("[C][离开遗迹场景] id=%v", id)
	return c.send(25317, body)
}

func (x *S2CYiJiLeaveScene) ID() uint16 {
	return 25318
}

// Message S2CYiJiLeaveScene Code:25318
func (x *S2CYiJiLeaveScene) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][离开遗迹场景] tag=%v tag_msg=%s id=%v", x.Tag, GetTagMsg(x.Tag), x.Id)
}

////////////////////////////////////////////////////////////

// YiJiJoinScene 进入遗迹场景
func (c *Connect) YiJiJoinScene(id int32) error {
	body, err := proto.Marshal(&C2SYiJiJoinScene{Id: id})
	if err != nil {
		return err
	}
	log.Printf("[C][进入遗迹场景] id=%v", id)
	return c.send(25315, body)
}

func (x *S2CYiJiJoinScene) ID() uint16 {
	return 25316
}

// Message S2CYiJiJoinScene Code:25316
func (x *S2CYiJiJoinScene) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][进入遗迹场景] tag=%v tag_msg=%s id=%v", x.Tag, GetTagMsg(x.Tag), x.Id)
}

////////////////////////////////////////////////////////////

// YiJiInfo 遗迹信息
func (c *Connect) YiJiInfo() error {
	body, err := proto.Marshal(&C2SYiJiInfo{})
	if err != nil {
		return err
	}
	log.Printf("[C][遗迹信息]")
	return c.send(25301, body)
}

func (x *S2CYiJiInfo) ID() uint16 {
	return 25302
}

// Message S2CYiJiInfo Code:25724
func (x *S2CYiJiInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][遗迹信息] items=%v", x.Items)
}

////////////////////////////////////////////////////////////

// EntrustInfo 遗迹信息
func (c *Connect) EntrustInfo() error {
	body, err := proto.Marshal(&C2SEntrustInfo{EntrustId: 8})
	if err != nil {
		return err
	}
	log.Printf("[C][遗迹委托信息]")
	return c.send(25361, body)
}

func (x *S2CEntrustInfo) ID() uint16 {
	return 25302
}

// Message S2CEntrustInfo Code:25362
func (x *S2CEntrustInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][遗迹委托信息] entrust_id=%v", x.EntrustId)
}

////////////////////////////////////////////////////////////

// GetEntrustWallList 遗迹信息
func (c *Connect) GetEntrustWallList() error {
	body, err := proto.Marshal(&C2SGetEntrustWallList{EntrustId: 8})
	if err != nil {
		return err
	}
	log.Printf("[C][遗迹委托列表]")
	return c.send(25351, body)
}

func (x *S2CGetEntrustWallList) ID() uint16 {
	return 25352
}

// Message S2CGetEntrustWallList Code:25352
func (x *S2CGetEntrustWallList) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][遗迹委托列表] entrust_id=%v list=%v", x.EntrustId, x.List)
}

////////////////////////////////////////////////////////////

// ReceiveEntrust 接取委托
func (c *Connect) ReceiveEntrust(uid int64) error {
	body, err := proto.Marshal(&C2SReceiveEntrust{EntrustId: 8, UserId: uid})
	if err != nil {
		return err
	}
	log.Printf("[C][接取委托] user_id=%d", uid)
	return c.send(25357, body)
}

func (x *S2CReceiveEntrust) ID() uint16 {
	return 25358
}

// Message S2CReceiveEntrust Code:25358
func (x *S2CReceiveEntrust) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][接取委托] tag=%v tag_msg=%s  user_id=%d", x.Tag, GetTagMsg(x.Tag), x.UserId)
}

////////////////////////////////////////////////////////////

// IllusionSweep 扫荡
func (c *Connect) IllusionSweep(t *C2SIllusionSweep) error {
	body, err := proto.Marshal(t)
	if err != nil {
		return err
	}
	log.Printf("[C][IllusionSweep] illusion_type=%v", t.IllusionType)
	return c.send(25723, body)
}

func (x *S2CIllusionSweep) ID() uint16 {
	return 25724
}

// Message S2CIllusionSweep Code:25724
func (x *S2CIllusionSweep) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][IllusionSweep] %v", x)
}

////////////////////////////////////////////////////////////

// PrefectWarFight 扫荡
func (c *Connect) PrefectWarFight(f *C2SPrefectWarFight) error {
	body, err := proto.Marshal(f)
	if err != nil {
		return err
	}
	log.Printf("[C][PrefectWarFight] boss_id=%v cp_id=%v tab_id=%v", f.BossId, f.CpId, f.TabId)
	return c.send(28407, body)
}

func (x *S2CPrefectWarFight) ID() uint16 {
	return 28408
}

// Message S2CPrefectWarFight Code:28408
func (x *S2CPrefectWarFight) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][PrefectWarFight] tag=%v tag_msg=%s boss_id=%v cp_id=%v tab_id=%v", x.Tag, GetTagMsg(x.Tag), x.BossId, x.CpId, x.TabId)
}

////////////////////////////////////////////////////////////

// GetPrefectWarData 扫荡
func (c *Connect) GetPrefectWarData() error {
	body, err := proto.Marshal(&C2SGetPrefectWarData{})
	if err != nil {
		return err
	}
	log.Println("[C][GetPrefectWarData]")
	return c.send(28401, body)
}

func (x *S2CGetPrefectWarData) ID() uint16 {
	return 28402
}

// Message S2CGetPrefectWarData Code:28402
func (x *S2CGetPrefectWarData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetPrefectWarData] tab_items=%v", x.TabItems)
}
