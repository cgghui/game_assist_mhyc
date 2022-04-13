package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func actXSTime() time.Duration {
	cur := time.Now()
	ast := time.Date(cur.Year(), cur.Month(), cur.Day(), 12, 00, 3, 0, time.Local)
	if cur.Before(ast) {
		return ast.Sub(cur)
	}
	if cur.Before(ast.Add(240 * time.Minute)) {
		return 0
	}
	return TomorrowDuration(43203 * time.Second)
}

// HuoDongXS 活动<仙山>
func HuoDongXS(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	defer t1.Stop()
	f1 := func() time.Duration {
		if td := actXSTime(); td != 0 {
			return td
		}
		// 仙宗信息
		Receive.Action(CLI.SectInfo)
		_ = Receive.Wait(&S2CSectInfo{}, s3)
		// 进入活动
		go func() {
			_ = CLI.JoinActive(&C2SJoinActive{AId: 4})
		}()
		_ = Receive.Wait(&S2CJoinActive{}, s3)
		defer func() {
			// 离开
			go func() {
				_ = CLI.LeaveActive(&C2SLeaveActive{AId: 4})
			}()
			_ = Receive.Wait(&S2CLeaveActive{}, s3)
		}()
		//
		ui := RoleInfo.Get("UserId").Int64()
		fv := RoleInfo.Get("FightValue").Int64()
		//
		currentI := int32(0)
		currentP := int32(0)
		newI := int32(0)
		newP := int32(0)
		newN := ""
		for i := int32(1); i <= 5; i++ {
			info := &S2CGetAllIMInfo{}
			for {
				go func() {
					_ = CLI.GetAllIMInfo(i)
				}()
				if err := Receive.Wait(info, s3); err != nil {
					continue
				}
				break
			}
			isEnd := false
			for _, player := range info.Players {
				if player.UserId == ui {
					currentI = i
					currentP = player.Pos
					isEnd = true
					break
				}
				if fv > player.Fv {
					newI = i
					newP = player.Pos
					newN = player.SectName
					isEnd = true
					break
				}
			}
			if isEnd {
				break
			}
		}
		if currentI == 0 && currentP == 0 {
			go func() {
				_ = CLI.SectIMSeize(newI, newP, newN)
			}()
			s := &S2CSectIMSeize{}
			_ = Receive.Wait(s, s10)
		}
		return RandMillisecond(60, 120)
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

// SectInfo 仙山信息
func (c *Connect) SectInfo() error {
	body, err := proto.Marshal(&C2SSectInfo{})
	if err != nil {
		return err
	}
	log.Println("[C][SectInfo]")
	return c.send(19015, body)
}

func (x *S2CSectInfo) ID() uint16 {
	return 19016
}

// Message S2CSectInfo Code:19016
func (x *S2CSectInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SectInfo] sect=%v", x.Sect)
}

////////////////////////////////////////////////////////////

// GetAllIMInfo 仙山信息
func (c *Connect) GetAllIMInfo(id int32) error {
	body, err := proto.Marshal(&C2SGetAllIMInfo{Id: id})
	if err != nil {
		return err
	}
	log.Printf("[C][GetAllIMInfo] id=%v", id)
	return c.send(19063, body)
}

func (x *S2CGetAllIMInfo) ID() uint16 {
	return 19064
}

// Message S2CGetAllIMInfo Code:19064
func (x *S2CGetAllIMInfo) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GetAllIMInfo] tag=%v", x.Tag)
}

////////////////////////////////////////////////////////////

// SectIMSeize 仙山抢夺
func (c *Connect) SectIMSeize(id, pos int32, sectName string) error {
	body, err := proto.Marshal(&C2SSectIMSeize{Id: id, Pos: pos, SectName: sectName})
	if err != nil {
		return err
	}
	log.Printf("[C][SectIMSeize] sect_name=%s id=%d pos=%d", sectName, id, pos)
	return c.send(19057, body)
}

func (x *S2CSectIMSeize) ID() uint16 {
	return 19058
}

// Message S2CSectIMSeize Code:19058
func (x *S2CSectIMSeize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SectIMSeize] tag=%v", x.Tag)
}
