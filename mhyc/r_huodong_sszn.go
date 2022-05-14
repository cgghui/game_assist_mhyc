package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// actSsznTime 神兽之怒时间
func actSsznTime() time.Duration {
	cur := time.Now()
	y := cur.Year()
	m := cur.Month()
	d := cur.Day()
	actStartTime := []time.Time{
		time.Date(y, m, d, 11, 50, 0, 0, time.Local).Add(ms500),
		time.Date(y, m, d, 13, 50, 0, 0, time.Local).Add(ms500),
		time.Date(y, m, d, 17, 50, 0, 0, time.Local).Add(ms500),
	}
	for _, ast := range actStartTime {
		if cur.Before(ast) {
			return ast.Sub(cur)
		}
		if cur.Before(ast.Add(30 * time.Minute)) {
			return 0
		}
	}
	return TomorrowDuration(time.Hour)
}

// HuoDongSSZN 活动<神兽之怒>
func HuoDongSSZN(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	defer t1.Stop()
	f1 := func() time.Duration {
		if td := actSsznTime(); td != 0 {
			return td
		}
		Fight.Lock()
		am := SetAction(ctx, "本服灵兽园-神兽之怒")
		defer func() {
			am.End()
			Fight.Unlock()
		}()
		// 进入
		Receive.Action(CLI.EnterAnimalPark)
		ret := &S2CEnterAnimalPark{}
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, ret, s10); err != nil {
			return time.Second
		}
		defer func() {
			Receive.Action(CLI.LeaveAnimalPark)
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CLeaveAnimalPark{}, s3)
		}()
		go func() {
			_ = CLI.StartMove(&C2SStartMove{P: []int32{44, 41}})
		}()
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CStartMove{}, s3); err != nil {
			return RandMillisecond(0, 1)
		}
		Receive.Action(CLI.FightBoss)
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CFightBoss{}, s3); err != nil {
			return RandMillisecond(0, 1)
		}
		return time.Hour
	}
	for {
		select {
		case <-t1.C:
			t1.Reset(f1())
		case <-ctx.Done():
			return
		}
	}
}

////////////////////////////////////////////////////////////

// FightBoss 战斗BOSS
func (c *Connect) FightBoss() error {
	body, err := proto.Marshal(&C2SFightBoss{})
	if err != nil {
		return err
	}
	log.Printf("[C][FightBoss]")
	return c.send(19091, body)
}

func (x *S2CFightBoss) ID() uint16 {
	return 19092
}

// Message S2CFightBoss Code:19092
func (x *S2CFightBoss) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][FightBoss] tag=%v", x.Tag)
}
