package mhyc

import (
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// Everyday 每日奖励
func Everyday() {
	// 每日在线时长奖励
	go func() {
		t := time.NewTimer(time.Second)
		defer t.Stop()
		f := func() {
			Fight.Lock()
			defer Fight.Unlock()
			task := &S2CGetActTask{}
			go func() {
				_ = CLI.GetActTask(&C2SGetActTask{ActId: 11002})
			}()
			_ = Receive.Wait(task, s3)
			s := 0
			for i := range task.Task {
				ret := &S2CGetTaskPrize{}
				go func(tid int32) {
					_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 6, Multi: 1, TaskId: tid})
				}(task.Task[i].Id)
				_ = Receive.Wait(ret, s3)
				if ret.Tag == 0 || ret.Tag == 5032 {
					s++
				}
			}
			if s == len(task.Task) {
				t.Reset(time.Hour)
			} else {
				t.Reset(RandMillisecond(60, 120))
			}
		}
		for range t.C {
			f()
		}
	}()
	// 等级大礼
	go func() {
		t := time.NewTimer(time.Second)
		defer t.Stop()
		f := func() bool {
			Fight.Lock()
			defer Fight.Unlock()
			task := &S2CGetActTask{}
			go func() {
				_ = CLI.GetActTask(&C2SGetActTask{ActId: 11011})
			}()
			_ = Receive.Wait(task, s3)
			s := 0
			for i := range task.Task {
				ret := &S2CGetTaskPrize{}
				go func(tid int32) {
					_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 6, Multi: 1, TaskId: tid})
				}(task.Task[i].Id)
				_ = Receive.Wait(ret, s3)
				if ret.Tag == 0 || ret.Tag == 5032 {
					s++
				}
			}
			if s == len(task.Task) {
				return true
			}
			t.Reset(time.Hour)
			return false
		}
		for range t.C {
			if f() {
				return
			}
		}
	}()
	// 战力大礼
	go func() {
		t := time.NewTimer(time.Second)
		defer t.Stop()
		f := func() bool {
			Fight.Lock()
			defer Fight.Unlock()
			task := &S2CGetActTask{}
			go func() {
				_ = CLI.GetActTask(&C2SGetActTask{ActId: 11012})
			}()
			_ = Receive.Wait(task, s3)
			s := 0
			for i := range task.Task {
				ret := &S2CGetTaskPrize{}
				go func(tid int32) {
					_ = CLI.GetTaskPrize(&C2SGetTaskPrize{TaskType: 6, Multi: 1, TaskId: tid})
				}(task.Task[i].Id)
				_ = Receive.Wait(ret, s3)
				if ret.Tag == 0 || ret.Tag == 5032 {
					s++
				}
			}
			if s == len(task.Task) {
				return true
			}
			t.Reset(time.Hour)
			return false
		}
		for range t.C {
			if f() {
				return
			}
		}
	}()
	t := time.NewTimer(ms100)
	f := func() {
		Fight.Lock()
		defer Fight.Unlock()
		// 充值->1元秒杀->每日礼
		go func() {
			_ = CLI.ActGiftNewReceive(DefineGiftRechargeEveryDay)
		}()
		_ = Receive.Wait(&S2CActGiftNewReceive{}, s3)
		// 排名—>本区榜->膜拜
		if RoleInfo.Get("Respect").Int64() == 0 {
			go func() {
				_ = CLI.Respect(DefineRespectL)
			}()
			_ = Receive.Wait(&S2CRespect{}, s3)
		}
		// 排名—>跨服榜->膜拜
		if RoleInfo.Get("RespectUnion").Int64() == 0 {
			go func() {
				_ = CLI.Respect(DefineRespectG)
			}()
			_ = Receive.Wait(&S2CRespect{}, s3)
		}
		// SVIP 每日礼包
		if RoleInfo.Get("VipDayGift").Int64() == 0 {
			Receive.Action(CLI.GetVipDayGift)
			_ = Receive.Wait(&S2CGetVipDayGift{}, s3)
		}
		// 寻宝
		s := 0
		for i := 1; i <= 8; i++ {
			data := &S2CGetActXunBaoData{}
			id := int32(i) + 500
			go func(id int32) {
				_ = CLI.GetActXunBaoData(&C2SGetActXunBaoData{ActId: id})
			}(id)
			_ = Receive.Wait(data, s3)
			if data.HaveFreeTime == 1 {
				go func(id int32) {
					_ = CLI.XunBaoDraw(&C2SActXunBaoDraw{ActId: id, Type: 1, AutoBuy: 0})
				}(id)
				_ = Receive.Wait(&S2CActXunBaoDraw{}, s3)
				s++
			}
		}
		if s > 0 {
			go func() {
				_ = CLI.WareHouseReceiveItem(5)
			}()
			_ = Receive.Wait(&S2CWareHouseReceiveItem{}, s3)
		}
		// 特权卡 -> 至尊卡
		if RoleInfo.Get("VipDayGift").Int64() == 0 {
			Receive.Action(CLI.LifeCardDayPrize)
			_ = Receive.Wait(&S2CLifeCardDayPrize{}, s3)
		}
		// 每日签到
		if RoleInfo.Get("HaveSign").Int64() == 0 {
			Receive.Action(CLI.EverydaySign)
			_ = Receive.Wait(&S2CSign{}, s3)
		}
		// 签到 领取签到奖励 @TODO: 无法确定是否领取
		for i := 1; i <= 4; i++ {
			go func(i int) {
				_ = CLI.TotalSignPrize(int32(i))
			}(i)
			_ = Receive.Wait(&S2CTotalSignPrize{}, s3)
		}
		// 商城购物 免费 @TODO: 无法确定是否领取
		go func() {
			_ = CLI.ShopBuy(DefineShopBuyFree)
		}()
		_ = Receive.Wait(&S2CShopBuy{}, s3)
		// 膜拜 宗主
		Receive.Action(CLI.Worship)
		_ = Receive.Wait(&S2CWorship{}, s3)
	}
	for range t.C {
		f()
		t.Reset(RandMillisecond(60, 180)) // 1 ~ 3 分钟
	}
}

// ActGiftNewReceive 充值->1元秒杀->每日礼
func (c *Connect) ActGiftNewReceive(act *C2SActGiftNewReceive) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(12011, body)
}

// Worship 膜拜 宗主
func (c *Connect) Worship() error {
	body, err := proto.Marshal(&C2SWorship{})
	if err != nil {
		return err
	}
	return c.send(19007, body)
}

// Respect 排名->膜拜
func (c *Connect) Respect(act *C2SRespect) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(13, body)
}

// GetVipDayGift 每日礼包
func (c *Connect) GetVipDayGift() error {
	body, err := proto.Marshal(DefineVipDayGift)
	if err != nil {
		return err
	}
	return c.send(136, body)
}

// XunBaoDraw 寻宝
func (c *Connect) XunBaoDraw(act *C2SActXunBaoDraw) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(11035, body)
}

// GetActXunBaoData 寻宝数据
func (c *Connect) GetActXunBaoData(act *C2SGetActXunBaoData) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(11033, body)
}

// LifeCardDayPrize 特权卡 -> 至尊卡
func (c *Connect) LifeCardDayPrize() error {
	body, err := proto.Marshal(DefineLifeCardDayPrize)
	if err != nil {
		return err
	}
	return c.send(22405, body)
}

// EverydaySign 每日签到
func (c *Connect) EverydaySign() error {
	body, err := proto.Marshal(DefineSign)
	if err != nil {
		return err
	}
	return c.send(22302, body)
}

// TotalSignPrize 累记签到奖励
func (c *Connect) TotalSignPrize(id int32) error {
	body, err := proto.Marshal(&C2STotalSignPrize{Id: id})
	if err != nil {
		return err
	}
	return c.send(22306, body)
}

// GetActTask 任务奖励
func (c *Connect) GetActTask(act *C2SGetActTask) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(12151, body)
}

// GetTaskPrize 领取奖励
func (c *Connect) GetTaskPrize(act *C2SGetTaskPrize) error {
	body, err := proto.Marshal(act)
	if err != nil {
		return err
	}
	return c.send(703, body)
}

////////////////////////////////////////////////////////////

func (x *S2CActGiftNewReceive) ID() uint16 {
	return 12012
}

// Message S2CActGiftNewReceive 12012
func (x *S2CActGiftNewReceive) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ActGiftNewReceive] tag=%v aid=%v gid=%v", x.Tag, x.Aid, x.Gid)
}

////////////////////////////////////////////////////////////

func (x *S2CRespect) ID() uint16 {
	return 14
}

// Message S2CRespect 14
func (x *S2CRespect) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][Respect] tag=%d type=%d prize=%v", x.Tag, x.Type, x.Prize)
}

////////////////////////////////////////////////////////////

func (x *S2CGetVipDayGift) ID() uint16 {
	return 137
}

// Message S2CGetVipDayGift 137
func (x *S2CGetVipDayGift) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetVipDayGift] tag=%d", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CActXunBaoDraw) ID() uint16 {
	return 11036
}

// Message S2CActXunBaoDraw 11036
func (x *S2CActXunBaoDraw) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ActXunBaoDraw] tag=%d %v", x.Tag, x)
}

////////////////////////////////////////////////////////////

func (x *S2CGetActXunBaoData) ID() uint16 {
	return 11034
}

// Message S2CGetActXunBaoData 11034
func (x *S2CGetActXunBaoData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetActXunBaoData] tag=%d %v", x.Tag, x)
}

////////////////////////////////////////////////////////////

func (x *S2CLifeCardDayPrize) ID() uint16 {
	return 22406
}

func (x *S2CLifeCardDayPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][LifeCardDayPrize] tag=%d", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CSign) ID() uint16 {
	return 22303
}

func (x *S2CSign) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][Sign] tag=%d", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CShopBuy) ID() uint16 {
	return 433
}

func (x *S2CShopBuy) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ShopBuy] tag=%d %v", x.Tag, x)
}

////////////////////////////////////////////////////////////

func (x *S2CTotalSignPrize) ID() uint16 {
	return 22307
}

func (x *S2CTotalSignPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][TotalSignPrize] tag=%d", x.Tag)
}

////////////////////////////////////////////////////////////

func (x *S2CGetActTask) ID() uint16 {
	return 12152
}

func (x *S2CGetActTask) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetActTask] tag=%d %v", x.Tag, x)
}

////////////////////////////////////////////////////////////

func (x *S2CWorship) ID() uint16 {
	return 19008
}

func (x *S2CWorship) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][CWorship] tag=%v %v", x.Tag, x)
}

////////////////////////////////////////////////////////////

func (x *S2CGetTaskPrize) ID() uint16 {
	return 704
}

func (x *S2CGetTaskPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetTaskPrize] tag=%d %v", x.Tag, x)
}
