package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// Everyday 每日奖励
func Everyday(ctx context.Context) {
	// 每日在线时长奖励
	go func() {
		t := time.NewTimer(time.Second)
		defer t.Stop()
		f := func() time.Duration {
			Fight.Lock()
			am := SetAction(ctx, "每日-在线时长奖励")
			defer func() {
				am.End()
				Fight.Unlock()
			}()
			task := &S2CGetActTask{}
			go func() {
				_ = CLI.GetActTask(&C2SGetActTask{ActId: 11002})
			}()
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, task, s3)
			count := len(task.Task)
			s := 0
			i := 0
			return am.RunAction(ctx, func() (time.Duration, time.Duration) {
				go func(tid int32) {
					_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 6, Multi: 1, TaskId: tid})
				}(task.Task[i].Id)
				ret := &S2CGetTaskPrize{}
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, ret, s3); err != nil {
					return 0, RandMillisecond(30, 60)
				}
				// 任务完成或已完成
				if ret.Tag == 0 || ret.Tag == 5032 {
					s++
				}
				// 任务未完成
				if ret.Tag == 5033 {
					return 0, RandMillisecond(300, 600)
				}
				i++
				if i >= count {
					if s == count {
						return 0, TomorrowDuration(RandMillisecond(1800, 3600))
					}
					return 0, RandMillisecond(300, 600)
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
	}()
	// 修仙 - 境界 任务
	go func() {
		t := time.NewTimer(time.Second)
		defer t.Stop()
		f := func() time.Duration {
			Fight.Lock()
			am := SetAction(ctx, "修仙-境界-任务")
			defer func() {
				am.End()
				Fight.Unlock()
			}()
			tc := time.NewTimer(ms10)
			defer tc.Stop()
			i := 0
			r := false
			Receive.Action(CLI.RealmTask)
			task := &S2CRealmTask{}
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, task, s3); err != nil {
				return RandMillisecond(30, 60)
			}
			ReTime := am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
				if i >= len(task.Tasks) {
					loop = 0
					next = RandMillisecond(600, 1800)
					return
				}
				r = false
				tk := task.Tasks[i]
				if tk.S != 1 {
					i++
					loop = ms100
					next = 0
					return
				}
				r = true
				go func(tid int32) {
					_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 24, Multi: 1, TaskId: tid})
				}(task.Tasks[i].Id)
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CGetTaskPrize{}, s3); err != nil {
					loop = 0
					next = RandMillisecond(30, 60)
					return
				}
				i++
				loop = ms100
				next = 0
				return
			})
			if r {
				return ms100
			}
			// 突破
			Receive.Action(CLI.RealmOverFulfil)
			_ = Receive.Wait(&S2CRealmOverfulfil{}, s3)
			return ReTime
		}
		for {
			select {
			case <-t.C:
				t.Reset(f())
			case <-ctx.Done():
				return
			}
		}
	}()
	// 等级大礼
	go func() {
		t := time.NewTimer(time.Second)
		defer t.Stop()
		f := func() time.Duration {
			Fight.Lock()
			am := SetAction(ctx, "等级大礼")
			defer func() {
				am.End()
				Fight.Unlock()
			}()
			go func() {
				_ = CLI.GetActTask(&C2SGetActTask{ActId: 11011})
			}()
			task := &S2CGetActTask{}
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, task, s3); err != nil {
				return RandMillisecond(30, 60)
			}
			count := len(task.Task)
			s := 0
			i := 0
			return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
				go func(tid int32) {
					_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 6, Multi: 1, TaskId: tid})
				}(task.Task[i].Id)
				ret := &S2CGetTaskPrize{}
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, ret, s3); err != nil {
					loop = 0
					next = RandMillisecond(30, 60)
					return
				}
				if ret.Tag == 0 || ret.Tag == 5032 {
					s++
				}
				if ret.Tag == 5033 {
					loop = 0
					next = RandMillisecond(300, 600)
					return
				}
				i++
				if i >= count {
					if s == count {
						return 0, TomorrowDuration(RandMillisecond(1800, 3600))
					}
					return 0, RandMillisecond(300, 600)
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
	}()
	// 战力大礼
	go func() {
		t := time.NewTimer(time.Second)
		defer t.Stop()
		f := func() time.Duration {
			Fight.Lock()
			am := SetAction(ctx, "战力大礼")
			defer func() {
				am.End()
				Fight.Unlock()
			}()
			task := &S2CGetActTask{}
			go func() {
				_ = CLI.GetActTask(&C2SGetActTask{ActId: 11012})
			}()
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, task, s3); err != nil {
				return RandMillisecond(30, 60)
			}
			count := len(task.Task)
			s := 0
			i := 0
			return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
				ret := &S2CGetTaskPrize{}
				go func(tid int32) {
					_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 6, Multi: 1, TaskId: tid})
				}(task.Task[i].Id)
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, ret, s3); err != nil {
					loop = 0
					next = RandMillisecond(30, 60)
					return
				}
				if ret.Tag == 0 || ret.Tag == 5032 {
					s++
				}
				if ret.Tag == 5033 {
					loop = 0
					next = RandMillisecond(300, 600)
					return
				}
				i++
				if i >= count {
					if s == count {
						return 0, TomorrowDuration(RandMillisecond(1800, 3600))
					}
					return 0, RandMillisecond(300, 600)
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
	}()
	// 我要变强
	go func() {
		t := time.NewTimer(time.Second)
		defer t.Stop()
		f := func() time.Duration {
			Fight.Lock()
			am := SetAction(ctx, "我要变强")
			defer func() {
				am.End()
				Fight.Unlock()
			}()
			s := 0
			i := 1
			return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
				ret := &S2CGetTaskPrize{}
				go func(i int32) {
					_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 2, Multi: 1, TaskId: i})
				}(int32(i))
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, ret, s3); err != nil {
					loop = 0
					next = RandMillisecond(30, 60)
					return
				}
				if ret.Tag == 0 || ret.Tag == 5032 {
					s++
				}
				if ret.Tag == 5033 {
					loop = 0
					next = RandMillisecond(300, 600)
					return
				}
				i++
				if i >= 5 {
					loop = 0
					next = RandMillisecond(1800, 3600)
					return
				}
				loop = ms100
				next = 0
				return
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
	}()
	// 定时领取主线任务奖励
	go func() {
		t := time.NewTimer(ms100)
		defer t.Stop()
		f := func() time.Duration {
			Fight.Lock()
			am := SetAction(ctx, "主线任务奖励")
			defer func() {
				am.End()
				Fight.Unlock()
			}()
			return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
				h := &S2CGetHistoryTaskPrize{}
				Receive.Action(CLI.GetHistoryTaskPrize)
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, h, s3); err != nil {
					loop = 0
					next = RandMillisecond(30, 60)
					return
				}
				if h.Tag != 0 {
					loop = 0
					next = RandMillisecond(1800, 3600)
					return
				}
				loop = ms100
				next = 0
				return
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
	}()
	// 仙缘副本 待仙缘上线 即该扫荡
	go ListenMessageCall(ctx, &S2CPartnerOnline{}, func(data []byte) {
		go func() {
			Fight.Lock()
			am := SetAction(ctx, "仙缘副本")
			defer func() {
				am.End()
				Fight.Unlock()
			}()
			am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
				Receive.Action(CLI.OneKeyWeddingIns)
				_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2COneKeyWeddingIns{}, s3)
				return 0, 0
			})
		}()
	})
	t := time.NewTimer(ms100)
	defer t.Stop()
	f := func() time.Duration {
		Fight.Lock()
		am := SetAction(ctx, "每日一次性任务")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			// 充值->1元秒杀->每日礼
			go func() {
				_ = CLI.ActGiftNewReceive(DefineGiftRechargeEveryDay)
			}()
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CActGiftNewReceive{}, s3)
			// 排名—>本区榜->膜拜
			if RoleInfo.Get("Respect").Int64() == 0 {
				go func() {
					_ = CLI.Respect(DefineRespectL)
				}()
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CRespect{}, s3); err != nil {
					loop = 0
					next = RandMillisecond(6, 10)
					return
				}
			}
			// 排名—>跨服榜->膜拜
			if RoleInfo.Get("RespectUnion").Int64() == 0 {
				go func() {
					_ = CLI.Respect(DefineRespectG)
				}()
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CRespect{}, s3); err != nil {
					loop = 0
					next = RandMillisecond(6, 10)
					return
				}
			}
			// SVIP 每日礼包
			if RoleInfo.Get("VipDayGift").Int64() == 0 {
				Receive.Action(CLI.GetVipDayGift)
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CGetVipDayGift{}, s3); err != nil {
					loop = 0
					next = RandMillisecond(6, 10)
					return
				}
			}
			// 寻宝
			s := 0
			for i := 1; i <= 8; i++ {
				data := &S2CGetActXunBaoData{}
				id := int32(i) + 500
				go func(id int32) {
					_ = CLI.GetActXunBaoData(&C2SGetActXunBaoData{ActId: id})
				}(id)
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, data, s3); err != nil {
					loop = 0
					next = RandMillisecond(6, 10)
					return
				}
				if data.HaveFreeTime == 1 {
					go func(id int32) {
						_ = CLI.XunBaoDraw(&C2SActXunBaoDraw{ActId: id, Type: 1, AutoBuy: 0})
					}(id)
					if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CActXunBaoDraw{}, s3); err != nil {
						loop = 0
						next = RandMillisecond(6, 10)
						return
					}
					s++
				}
			}
			if s > 0 {
				go func() {
					_ = CLI.WareHouseReceiveItem(5)
				}()
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CWareHouseReceiveItem{}, s3); err != nil {
					loop = 0
					next = RandMillisecond(6, 10)
					return
				}
			}
			// 特权卡 -> 至尊卡
			if RoleInfo.Get("LifeCardDayPrize").Int64() == 0 {
				Receive.Action(CLI.LifeCardDayPrize)
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CLifeCardDayPrize{}, s3); err != nil {
					loop = 0
					next = RandMillisecond(6, 10)
					return
				}
			}
			// 每日签到
			if RoleInfo.Get("HaveSign").Int64() == 0 {
				Receive.Action(CLI.EverydaySign)
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CSign{}, s3); err != nil {
					loop = 0
					next = RandMillisecond(6, 10)
					return
				}
			}
			// 签到 领取签到奖励 @TODO: 无法确定是否领取
			for i := 1; i <= 4; i++ {
				go func(i int) {
					_ = CLI.TotalSignPrize(int32(i))
				}(i)
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CTotalSignPrize{}, s3); err != nil {
					loop = 0
					next = RandMillisecond(6, 10)
					return
				}
			}
			// 商城购物 免费 @TODO: 无法确定是否领取
			go func() {
				_ = CLI.ShopBuy(DefineShopBuyFree)
			}()
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CShopBuy{}, s3); err != nil {
				loop = 0
				next = RandMillisecond(6, 10)
				return
			}
			// 膜拜 宗主
			if RoleInfo.Get("SectWorship").Int64() == 0 {
				Receive.Action(CLI.Worship)
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CWorship{}, s3); err != nil {
					loop = 0
					next = RandMillisecond(6, 10)
					return
				}
			}
			// 仙宗 - 仙殿 - 仙宗声望 -> 每日奉碌
			if RoleInfo.Get("SectPrestigeRecv").Int64() == 0 {
				Receive.Action(CLI.SectPrestigeRecv)
				if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CSectPrestigeRecv{}, s3); err != nil {
					loop = 0
					next = RandMillisecond(6, 10)
					return
				}
			}
			loop = 0
			next = TomorrowDuration(RandMillisecond(1800, 3600))
			return
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

////////////////////////////////////////////////////////////

// OneKeyWeddingIns 仙缘 扫荡
func (c *Connect) OneKeyWeddingIns() error {
	body, err := proto.Marshal(&C2SOneKeyWeddingIns{})
	if err != nil {
		return err
	}
	return c.send(22360, body)
}

func (x *S2COneKeyWeddingIns) ID() uint16 {
	return 22361
}

func (x *S2COneKeyWeddingIns) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][OneKeyWeddingIns] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CPartnerOnline) ID() uint16 {
	return 22345
}

func (x *S2CPartnerOnline) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][PartnerOnline] user_id=%v nick=%v", x.UserId, x.Nick)
}

////////////////////////////////////////////////////////////

// RealmOverFulfil 突破
func (c *Connect) RealmOverFulfil() error {
	body, err := proto.Marshal(&C2SRealmOverfulfil{})
	if err != nil {
		return err
	}
	return c.send(22004, body)
}

func (x *S2CRealmOverfulfil) ID() uint16 {
	return 22005
}

func (x *S2CRealmOverfulfil) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][RealmOverfulfil] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

// GetHistoryTaskPrize 主线任务奖励
func (c *Connect) GetHistoryTaskPrize() error {
	body, err := proto.Marshal(&C2SGetHistoryTaskPrize{TaskId: 0})
	if err != nil {
		return err
	}
	return c.send(713, body)
}

func (x *S2CGetHistoryTaskPrize) ID() uint16 {
	return 714
}

func (x *S2CGetHistoryTaskPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetHistoryTaskPrize] tag=%v tag_msg=%s raw=%v", x.Tag, GetTagMsg(x.Tag), x)
}

////////////////////////////////////////////////////////////

// RealmTask 修仙 - 境界 任务
func (c *Connect) RealmTask() error {
	body, err := proto.Marshal(&C2SRealmTask{})
	if err != nil {
		return err
	}
	return c.send(22012, body)
}

func (x *S2CRealmTask) ID() uint16 {
	return 22013
}

// Message S2CRealmTask 22013
func (x *S2CRealmTask) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][RealmTask] task=%v", x.Tasks)
}

////////////////////////////////////////////////////////////

// ActGiftNewReceive 充值->1元秒杀->每日礼
func (c *Connect) ActGiftNewReceive(act *C2SActGiftNewReceive) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	log.Printf("[C][ActGiftNewReceive] aid=%v gid=%v", act.Aid, act.Gid)
	return c.send(12011, body)
}

// Worship 膜拜 宗主
func (c *Connect) Worship() error {
	body, err := proto.Marshal(&C2SWorship{})
	if err != nil {
		return err
	}
	log.Println("[C][Worship]")
	return c.send(19007, body)
}

// Respect 排名->膜拜
func (c *Connect) Respect(act *C2SRespect) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	log.Printf("[C][Respect] type=%v", act.Type)
	return c.send(13, body)
}

// GetVipDayGift 每日礼包
func (c *Connect) GetVipDayGift() error {
	body, err := proto.Marshal(DefineVipDayGift)
	if err != nil {
		return err
	}
	log.Println("[C][GetVipDayGift]")
	return c.send(136, body)
}

// XunBaoDraw 寻宝
func (c *Connect) XunBaoDraw(act *C2SActXunBaoDraw) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	log.Printf("[C][XunBaoDraw] act_id=%v type=%v auto_buy=%v", act.ActId, act.Type, act.AutoBuy)
	return c.send(11035, body)
}

// GetActXunBaoData 寻宝数据
func (c *Connect) GetActXunBaoData(act *C2SGetActXunBaoData) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	log.Printf("[C][GetActXunBaoData] act_id=%v", act.ActId)
	return c.send(11033, body)
}

// LifeCardDayPrize 特权卡 -> 至尊卡
func (c *Connect) LifeCardDayPrize() error {
	body, err := proto.Marshal(DefineLifeCardDayPrize)
	if err != nil {
		return err
	}
	log.Println("[C][LifeCardDayPrize]")
	return c.send(22405, body)
}

// EverydaySign 每日签到
func (c *Connect) EverydaySign() error {
	body, err := proto.Marshal(DefineSign)
	if err != nil {
		return err
	}
	log.Println("[C][EverydaySign]")
	return c.send(22302, body)
}

// TotalSignPrize 累记签到奖励
func (c *Connect) TotalSignPrize(id int32) error {
	body, err := proto.Marshal(&C2STotalSignPrize{Id: id})
	if err != nil {
		return err
	}
	log.Printf("[C][TotalSignPrize] id=%v", id)
	return c.send(22306, body)
}

// GetActTask 任务奖励
func (c *Connect) GetActTask(act *C2SGetActTask) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	log.Printf("[C][GetActTask] act_id=%v", act.ActId)
	return c.send(12151, body)
}

// GetTaskPrize 领取奖励
func (c *Connect) GetTaskPrize(act *C2SGetTaskPrize) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	log.Printf("[C][GetTaskPrize] task_id=%v task_type=%v multi=%v", act.TaskId, act.TaskType, act.Multi)
	return c.send(703, body)
}

////////////////////////////////////////////////////////////

func (x *S2CActGiftNewReceive) ID() uint16 {
	return 12012
}

// Message S2CActGiftNewReceive 12012
func (x *S2CActGiftNewReceive) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ActGiftNewReceive] tag=%v tag_msg=%s aid=%v gid=%v", x.Tag, GetTagMsg(x.Tag), x.Aid, x.Gid)
}

////////////////////////////////////////////////////////////

func (x *S2CRespect) ID() uint16 {
	return 14
}

// Message S2CRespect 14
func (x *S2CRespect) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][Respect] tag=%d tag_msg=%s type=%d prize=%v", x.Tag, GetTagMsg(x.Tag), x.Type, x.Prize)
}

////////////////////////////////////////////////////////////

func (x *S2CGetVipDayGift) ID() uint16 {
	return 137
}

// Message S2CGetVipDayGift 137
func (x *S2CGetVipDayGift) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetVipDayGift] tag=%d tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CActXunBaoDraw) ID() uint16 {
	return 11036
}

// Message S2CActXunBaoDraw 11036
func (x *S2CActXunBaoDraw) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ActXunBaoDraw] tag=%d tag_msg=%s %v", x.Tag, GetTagMsg(x.Tag), x)
}

////////////////////////////////////////////////////////////

func (x *S2CGetActXunBaoData) ID() uint16 {
	return 11034
}

// Message S2CGetActXunBaoData 11034
func (x *S2CGetActXunBaoData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetActXunBaoData] tag=%d tag_msg=%s %v", x.Tag, GetTagMsg(x.Tag), x)
}

////////////////////////////////////////////////////////////

func (x *S2CLifeCardDayPrize) ID() uint16 {
	return 22406
}

func (x *S2CLifeCardDayPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][LifeCardDayPrize] tag=%d tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CSign) ID() uint16 {
	return 22303
}

func (x *S2CSign) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][Sign] tag=%d tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CTotalSignPrize) ID() uint16 {
	return 22307
}

func (x *S2CTotalSignPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][TotalSignPrize] tag=%d tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CGetActTask) ID() uint16 {
	return 12152
}

func (x *S2CGetActTask) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetActTask] tag=%d tag_msg=%s %v", x.Tag, GetTagMsg(x.Tag), x)
}

////////////////////////////////////////////////////////////

func (x *S2CWorship) ID() uint16 {
	return 19008
}

func (x *S2CWorship) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][CWorship] tag=%v tag_msg=%s %v", x.Tag, GetTagMsg(x.Tag), x)
}

////////////////////////////////////////////////////////////

func (x *S2CGetTaskPrize) ID() uint16 {
	return 704
}

func (x *S2CGetTaskPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetTaskPrize] tag=%d tag_msg=%s %v", GetTagMsg(x.Tag), x.Tag, x)
}
