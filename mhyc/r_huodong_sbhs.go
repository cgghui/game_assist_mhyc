package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// actSbhsTime 双倍护送时间
func actSbhsTime() time.Duration {
	cur := time.Now()
	actStartTime := []time.Time{
		time.Date(cur.Year(), cur.Month(), cur.Day(), 11, 30, 10, 0, time.Local),
		time.Date(cur.Year(), cur.Month(), cur.Day(), 18, 30, 10, 0, time.Local),
	}
	for _, ast := range actStartTime {
		if cur.Before(ast) {
			return ast.Sub(cur)
		}
		if cur.Before(ast.Add(150 * time.Minute)) {
			return 0
		}
	}
	return TomorrowDuration(3 * time.Hour)
}

// actSbhs 双倍护送
func actSbhs() time.Duration {
	if td := actSbhsTime(); td != 0 {
		return td
	}
	Fight.Lock()
	defer Fight.Unlock()
	Receive.Action(CLI.GetWestExp)
	var west S2CGetWestExp
	if err := Receive.Wait(&west, s3); err != nil {
		return ms500
	}
	if west.Tag == 819 {
		return time.Minute
	}
	n := 0
	for {
		if west.I == 4 || west.I == 5 || n >= 10 {
			break
		}
		Receive.Action(CLI.GetWestExpRef)
		_ = Receive.Wait(&S2CGetWestExp{}, s3)
		n++
	}
	Receive.Action(CLI.StartWestExp)
	_ = Receive.Wait(&S2CStartWestExp{}, s3)
	return time.Hour
}

// actHsPlayer 护送玩家 拦截
func actHsPlayer() time.Duration {
	if td := actSbhsTime(); td != 0 {
		return td
	}
	Fight.Lock()
	defer Fight.Unlock()
	Receive.Action(CLI.GetProtectPlayer)
	var pp S2CGetProtectPlayer
	if err := Receive.Wait(&pp, s3); err != nil {
		return ms500
	}
	self := RoleInfo.Get("FightValue").Int64()
	for _, player := range pp.List {
		if player.Fv >= self {
			continue
		}
		s, _ := FightActionRob(player.Uid)
		if s == nil {
			continue
		}
		if s.Tag == 809 {
			break
		}
	}
	return time.Hour
}

// HuoDongSBHS 活动<双倍护送>
func HuoDongSBHS(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	defer t1.Stop()
	t2 := time.NewTimer(ms100)
	defer t2.Stop()
	for {
		select {
		case <-t1.C:
			t1.Reset(actSbhs())
		case <-t2.C:
			t2.Reset(actHsPlayer())
		case <-ctx.Done():
			return
		}
	}
}

////////////////////////////////////////////////////////////

func FightActionRob(uid int32) (*S2CSendRob, *S2CBattlefieldReport) {
	c := make(chan *S2CSendRob)
	defer close(c)
	go func() {
		sf := &S2CSendRob{}
		if err := Receive.Wait(sf, s3); err != nil {
			c <- nil
		} else {
			c <- sf
		}
	}()
	go func() {
		_ = CLI.SendRob(uid)
	}()
	r := &S2CBattlefieldReport{}
	if err := Receive.Wait(r, s3); err != nil {
		r = nil
	}
	return <-c, r
}

////////////////////////////////////////////////////////////

// SendRob 拦截
func (c *Connect) SendRob(uid int32) error {
	body, err := proto.Marshal(&C2SSendRob{U: uid})
	if err != nil {
		return err
	}
	log.Printf("[C][SendRob] user_id=%v", uid)
	return c.send(476, body)
}

func (x *S2CSendRob) ID() uint16 {
	return 477
}

// Message S2CSendRob Code:477
func (x *S2CSendRob) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SendRob] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

// GetProtectPlayer 护送列表
func (c *Connect) GetProtectPlayer() error {
	body, err := proto.Marshal(&C2SGetProtectPlayer{})
	if err != nil {
		return err
	}
	log.Printf("[C][C2SGetProtectPlayer] get_type=%v", 1)
	return c.send(461, body)
}

func (x *S2CGetProtectPlayer) ID() uint16 {
	return 462
}

// Message S2CGetProtectPlayer Code:462
func (x *S2CGetProtectPlayer) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetProtectPlayer] list=%v", x.List)
}

////////////////////////////////////////////////////////////

// GetWestExp 护送 免费
func (c *Connect) GetWestExp() error {
	body, err := proto.Marshal(&C2SGetWestExp{GetType: 0})
	if err != nil {
		return err
	}
	log.Printf("[C][GetWestExp] get_type=%v", 0)
	return c.send(463, body)
}

// GetWestExpRef 护送 用50元宝刷新
func (c *Connect) GetWestExpRef() error {
	body, err := proto.Marshal(&C2SGetWestExp{GetType: 1})
	if err != nil {
		return err
	}
	log.Printf("[C][GetWestExp] get_type=%v", 1)
	return c.send(463, body)
}

func (x *S2CGetWestExp) ID() uint16 {
	return 464
}

// Message S2CGetWestExp Code:464
func (x *S2CGetWestExp) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetWestExp] %v", x)
}

////////////////////////////////////////////////////////////

// StartWestExp 护送动作
func (c *Connect) StartWestExp() error {
	body, err := proto.Marshal(&C2SStartWestExp{})
	if err != nil {
		return err
	}
	log.Printf("[C][StartWestExp]")
	return c.send(474, body)
}

func (x *S2CStartWestExp) ID() uint16 {
	return 475
}

// Message S2CStartWestExp Code:475
func (x *S2CStartWestExp) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][StartWestExp] tag=%v", x.Tag)
}
