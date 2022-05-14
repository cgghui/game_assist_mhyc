package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"sort"
	"time"
)

func illusionSweep(ctx context.Context) time.Duration {
	Fight.Lock()
	am := SetAction(ctx, "跨服-幻境一键扫荡")
	defer func() {
		am.End()
		Fight.Unlock()
	}()
	irList := []int32{1, 2, 11}
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
	if RoleInfo.Get("YiJiDayBeEntrustedTimes").Int64() <= 0 {
		return TomorrowDuration(RandMillisecond(1800, 3600))
	}
	if RoleInfo.Get("YiJiBeEntrustedNum").Int64() < 3 {
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
		if RoleInfo.Get("YiJiBeEntrustedNum").Int64() <= 0 {
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

// KuaFu 跨服
func KuaFu(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	t2 := time.NewTimer(ms50)
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
