package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"sync"
	"time"
)

// actJXSCTime 极限生存时间
func actJXSCTime() time.Duration {
	cur := time.Now()
	y := cur.Year()
	m := cur.Month()
	d := cur.Day()
	actStartTime := []time.Time{
		time.Date(y, m, d, 11, 00, 0, 0, time.Local).Add(time.Second),
		time.Date(y, m, d, 15, 00, 0, 0, time.Local).Add(time.Second),
		time.Date(y, m, d, 19, 00, 0, 0, time.Local).Add(time.Second),
	}
	for _, ast := range actStartTime {
		if cur.Before(ast) {
			return ast.Sub(cur)
		}
		if cur.Before(ast.Add(3 * time.Minute)) {
			return 0
		}
	}
	return TomorrowDuration(3 * time.Hour)
}

var monsterIDS = make([]int64, 0)
var monsterMUX = &sync.Mutex{}

func monsterLeaveIn(id int64) bool {
	for _, _id := range monsterIDS {
		if _id == id {
			return true
		}
	}
	return false
}

func monsterLeaveAdd(id int64) {
	monsterMUX.Lock()
	defer monsterMUX.Unlock()
	monsterIDS = append(monsterIDS, id)
}

func monsterLeaveDel(id int64) bool {
	monsterMUX.Lock()
	defer monsterMUX.Unlock()
	add := false
	for {
		add = false
		for i, _id := range monsterIDS {
			if _id == id {
				add = true
				monsterIDS = append(monsterIDS[:i], monsterIDS[i+1:]...)
			}
		}
		if add == false {
			return true
		}
	}
}

func jxsc(ctx context.Context) time.Duration {
	if td := actJXSCTime(); td != 0 {
		return td
	}
	if RoleInfo.Get("JXSC_Join_Times").Int64() >= 4 {
		return TomorrowDuration(RandMillisecond(1800, 3600))
	}
	Fight.Lock()
	am := SetAction(ctx, "HuoDongJXSC", 20*time.Minute)
	defer func() {
		am.End()
		Fight.Unlock()
	}()
	//
	monster := make(chan *S2CMonsterEnterMap, 100)
	go func() {
		defer close(monster)
		ListenMessageCall(am.Ctx, &S2CMonsterEnterMap{}, func(data []byte) {
			enter := &S2CMonsterEnterMap{}
			enter.Message(data)
			monsterLeaveDel(enter.Id)
			monster <- enter
		})
	}()
	//
	monsterIDS = make([]int64, 0)
	go func() {
		ListenMessageCall(am.Ctx, &S2CMonsterLeaveMap{}, func(data []byte) {
			leave := &S2CMonsterLeaveMap{}
			leave.Message(data)
			monsterLeaveAdd(leave.Id)
		})
	}()
	go ListenMessage(am.Ctx, &S2CJXSCMyScore{})
	stage := int32(0)
	go ListenMessageCall(am.Ctx, &S2CJXSCStageChange{}, func(data []byte) {
		r := &S2CJXSCStageChange{}
		r.Message(data)
		stage = r.Stage
	})
	go ListenMessageCall(am.Ctx, &S2CJXSCLeaveScene{}, func(_ []byte) {
		am.End()
	})
	// 进入活动
	go func() {
		_ = CLI.JoinActive(&C2SJoinActive{AId: 5})
	}()
	join := &S2CJoinActive{}
	_ = Receive.WaitWithContextOrTimeout(am.Ctx, join, s3)
	if join.Tag == 50502 {
		return TomorrowDuration(RandMillisecond(3600, 7200))
	}
	defer func() {
		// 离开
		go func() {
			_ = CLI.LeaveActive(&C2SLeaveActive{AId: 5})
		}()
		_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CLeaveActive{}, s3)
	}()
	Receive.Action(CLI.JXSCKeyNum)
	var info S2CJXSCKeyNum
	if err := Receive.WaitWithContextOrTimeout(am.Ctx, &info, s3); err != nil {
		return time.Second
	}
	Receive.Action(CLI.JXSCMyScore)
	var my S2CJXSCMyScore
	if err := Receive.WaitWithContextOrTimeout(am.Ctx, &my, s3); err != nil {
		return time.Second
	}
	Receive.Action(CLI.JXSCSkinChange)
	if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CJXSCSkinChange{}, s3); err != nil {
		return time.Second
	}
	// TODO: 此地观查
	tm := time.NewTimer(time.Hour)
	defer tm.Stop()
	for m := range monster {
		if m == nil {
			break
		}
		if monsterLeaveIn(m.Id) {
			log.Println("[S][monster] 怪兽离开 id=", m.Id, " x=", m.X, " y=", m.Y)
			continue
		}
		am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			go func() {
				_ = CLI.StartMove(&C2SStartMove{P: []int32{int32(m.X), int32(m.Y)}})
			}()
			//_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CStartMove{}, s3)
			tm.Reset(ms50)
			<-tm.C
			if stage == 6 {
				go func() {
					_ = CLI.JXSCOpenBox(int32(m.Id))
				}()
				_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CJXSCOpenBox{}, s3)
			} else {
				_, _ = FightAction(am.Ctx, m.Id, 8)
			}
			return 0, 0
		})
		tm.Reset(RandMillisecond(0, 3))
		<-tm.C
	}
	return RandMillisecond(1, 3)
}

// JXSC 跨服 极限生存
func JXSC(ctx context.Context) {
	t1 := time.NewTimer(ms100)
	defer t1.Stop()
	for {
		select {
		case <-t1.C:
			t1.Reset(jxsc(ctx))
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
	log.Printf("[S][JXSCSkinChange] tag=%v tag_msg=%s id=%v", x.Tag, GetTagMsg(x.Tag), x.Id)
}

////////////////////////////////////////////////////////////

func (c *Connect) JXSCLeaveScene() error {
	body, err := proto.Marshal(&C2SJXSCLeaveScene{})
	if err != nil {
		return err
	}
	log.Println("[C][JXSCLeaveScene]")
	return c.send(23209, body)
}

func (x *S2CJXSCLeaveScene) ID() uint16 {
	return 23210
}

// Message S2CJXSCLeaveScene Code:23210
func (x *S2CJXSCLeaveScene) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JXSCLeaveScene] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}

////////////////////////////////////////////////////////////

func (x *S2CJXSCStageChange) ID() uint16 {
	return 23202
}

// Message S2CJXSCStageChange Code:23202
func (x *S2CJXSCStageChange) Message(data []byte) {
	if err := proto.Unmarshal(data, x); err != nil {
		return
	}
	t := time.Unix(x.EndTimestamp, 0).Local()
	log.Printf("[S][JXSCStageChange] stage=%v end_timestamp=%s", x.Stage, t.Format("2006-01-02 15:04:05"))
}

////////////////////////////////////////////////////////////

// JXSCOpenBox 开箱子
func (c *Connect) JXSCOpenBox(id int32) error {
	body, err := proto.Marshal(&C2SJXSCOpenBox{Id: id})
	if err != nil {
		return err
	}
	log.Println("[C][JXSCOpenBox]")
	return c.send(23213, body)
}

func (x *S2CJXSCOpenBox) ID() uint16 {
	return 23214
}

// Message S2CJXSCOpenBox Code:23214
func (x *S2CJXSCOpenBox) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][JXSCOpenBox] tag=%v tag_msg=%s", x.Tag, GetTagMsg(x.Tag))
}
