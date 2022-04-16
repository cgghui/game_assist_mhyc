package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// actJXSCTime 极限生存时间
func actJXSCTime() time.Duration {
	cur := time.Now()
	actStartTime := []time.Time{
		time.Date(cur.Year(), cur.Month(), cur.Day(), 11, 00, 0, 0, time.Local).Add(s3),
		time.Date(cur.Year(), cur.Month(), cur.Day(), 13, 00, 0, 0, time.Local).Add(s3),
		time.Date(cur.Year(), cur.Month(), cur.Day(), 19, 00, 0, 0, time.Local).Add(s3),
	}
	for _, ast := range actStartTime {
		if cur.Before(ast) {
			return ast.Sub(cur)
		}
		if cur.Before(ast.Add(10 * time.Minute)) {
			return 0
		}
	}
	return TomorrowDuration(3 * time.Hour)
}

func jxsc() time.Duration {
	Fight.Lock()
	defer Fight.Unlock()
	// 进入活动
	go func() {
		_ = CLI.JoinActive(&C2SJoinActive{AId: 5})
	}()
	_ = Receive.Wait(&S2CJoinActive{}, s3)
	defer func() {
		// 离开
		go func() {
			_ = CLI.LeaveActive(&C2SLeaveActive{AId: 5})
		}()
		_ = Receive.Wait(&S2CLeaveActive{}, s3)
	}()
	Receive.Action(CLI.JXSCKeyNum)
	var info S2CJXSCKeyNum
	if err := Receive.Wait(&info, s3); err != nil {
		return time.Second
	}
	Receive.Action(CLI.JXSCMyScore)
	var my S2CJXSCMyScore
	if err := Receive.Wait(&my, s3); err != nil {
		return time.Second
	}
	Receive.Action(CLI.JXSCSkinChange)
	if err := Receive.Wait(&S2CJXSCSkinChange{}, s3); err != nil {
		return time.Second
	}
	monster := make(chan *S2CMonsterEnterMap, 50)
	defer close(monster)
	go ListenMessageCallEx(&S2CMonsterEnterMap{}, func(data []byte) bool {
		var enter S2CMonsterEnterMap
		if err := proto.Unmarshal(data, &enter); err == nil {
			monster <- &enter
			return false
		}
		return true
	})
	// TODO: 此地观查
	for m := range monster {
		FightAction(m.Id, 8)
	}
	return ms500
}

// JXSC 跨服 极限生存
func JXSC(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	defer t1.Stop()
	for {
		select {
		case <-t1.C:
			t1.Reset(jxsc())
		case <-ctx.Done():
			return
		}
	}
}

////////////////////////////////////////////////////////////

// JXSCKeyNum 列表
func (c *Connect) JXSCKeyNum() error {
	body, err := proto.Marshal(&C2SJXSCKeyNum{})
	if err != nil {
		return err
	}
	log.Println("[C][JXSCKeyNum]")
	return c.send(23215, body)
}

func (x *S2CJXSCKeyNum) ID() uint16 {
	return 23216
}

// Message S2CJXSCKeyNum Code:23216
func (x *S2CJXSCKeyNum) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JXSCKeyNum] key_num=%v", x.KeyNum)
}

////////////////////////////////////////////////////////////

// JXSCMyScore 列表
func (c *Connect) JXSCMyScore() error {
	body, err := proto.Marshal(&C2SJXSCMyScore{})
	if err != nil {
		return err
	}
	log.Println("[C][JXSCMyScore]")
	return c.send(23217, body)
}

func (x *S2CJXSCMyScore) ID() uint16 {
	return 23218
}

// Message S2CJXSCMyScore Code:23218
func (x *S2CJXSCMyScore) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JXSCMyScore] my_rank=%v my_score=%v", x.MyRank, x.MyScore)
}

////////////////////////////////////////////////////////////

// JXSCSkinChange 列表
func (c *Connect) JXSCSkinChange() error {
	body, err := proto.Marshal(&C2SJXSCSkinChange{})
	if err != nil {
		return err
	}
	log.Println("[C][JXSCSkinChange]")
	return c.send(23205, body)
}

func (x *S2CJXSCSkinChange) ID() uint16 {
	return 23206
}

// Message S2CJXSCSkinChange Code:23206
func (x *S2CJXSCSkinChange) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JXSCSkinChange] tag=%v id=%v", x.Tag, x.Id)
}
