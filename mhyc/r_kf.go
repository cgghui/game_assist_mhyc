package mhyc

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"sort"
	"time"
)

func illusionSweep() time.Duration {
	Fight.Lock()
	defer Fight.Unlock()
	for _, ir := range []int32{1, 2, 11} {
		go func(ir int32) {
			_ = CLI.IllusionSweep(&C2SIllusionSweep{IllusionType: ir})
		}(ir)
		r := &S2CIllusionSweep{}
		_ = Receive.Wait(r)
		fmt.Println(r)
	}
	return TomorrowDuration(RandMillisecond(30000, 30600))
}

func yiJi() time.Duration {
	Fight.Lock()
	defer Fight.Unlock()
	if RoleInfo.Get("YiJiAdscTimes").Int64() <= 0 {
		return TomorrowDuration(RandMillisecond(30000, 30600))
	}
	Receive.Action(CLI.YiJiInfo)
	var info S2CYiJiInfo
	_ = Receive.Wait(&info, s3)
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
		_ = Receive.Wait(&S2CYiJiLeaveScene{}, s3)
	}()
	mc := make(chan *S2CMonsterEnterMap)
	go func() {
		go func() {
			var monster S2CMonsterEnterMap
			_ = Receive.Wait(&monster, s3)
			mc <- &monster
			close(mc)
		}()
		_ = CLI.YiJiJoinScene(id)
	}()
	var join S2CYiJiJoinScene
	_ = Receive.Wait(&join, s3)
	monster := <-mc
	tc := time.NewTimer(ts0)
	defer tc.Stop()
	for range tc.C {
		s, r := FightAction(monster.Id, 8)
		if s == nil {
			break
		}
		if s.Tag == 17002 {
			break
		}
		if r != nil && r.Win == 1 {
			break
		}
		tc.Reset(time.Second)
	}
	//
	return ms500
}

// KuaFu 跨服
func KuaFu(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	t4 := time.NewTimer(ms10)
	defer t1.Stop()
	defer t4.Stop()
	for {
		select {
		case <-t1.C:
			t1.Reset(yiJi())
		case <-t4.C:
			t4.Reset(illusionSweep())
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
	log.Printf("[S][离开遗迹场景] tag=%v id=%v", x.Tag, x.Id)
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
	log.Printf("[S][进入遗迹场景] tag=%v id=%v", x.Tag, x.Id)
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
