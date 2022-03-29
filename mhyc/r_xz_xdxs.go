package mhyc

// 仙宗 - 仙殿 - 仙宗悬赏

import (
	"google.golang.org/protobuf/proto"
	"log"
	"sort"
	"time"
)

// XianDianXDXS - 仙殿 - 仙宗悬赏
func XianDianXDXS() {
	t := time.NewTimer(ms100)
	f := func() time.Duration {
		task := &S2CPlayerXZXS{}
		Receive.Action(CLI.PlayerXZXS)
		if err := Receive.Wait(task, s3); err != nil {
			return time.Second
		}
		if len(task.Tasks) == 0 {
			return time.Unix(task.ResetTimestamp, 0).Local().Add(30 * time.Minute).Sub(time.Now())
		}
		PQ := false            // 需要一键派遣时为true
		LS := make([]int64, 0) // 等领取或待接收任务
		for _, tk := range task.Tasks {
			if tk.TaskState == 0 {
				PQ = true
				continue
			}
			LS = append(LS, tk.TaskTimestamp)
		}
		if len(LS) > 0 {
			sort.Slice(LS, func(i, j int) bool {
				return LS[i] > LS[j]
			})
			cur := time.Now()
			ttm := time.Unix(LS[0], 0).Local()
			// 一键领取时间到了
			if cur.After(ttm) {
				Receive.Action(CLI.XZXSGetTaskPrize)
			} else {
				time.AfterFunc(ttm.Sub(cur), func() { _ = CLI.XZXSGetTaskPrize() })
			}
			_ = Receive.Wait(&S2CXZXSGetTaskPrize{}, 24*time.Hour)
		}
		if PQ {
			info := &S2CXZXSGetAllCanStartTask{}
			Receive.Action(CLI.XZXSGetAllCanStartTask)
			_ = Receive.Wait(info, s3)
			if info.Tag == 0 && len(info.Data) > 0 {
				go func() {
					_ = CLI.XZXSOneKeyStartTask(info)
				}()
				_ = Receive.Wait(&S2CXZXSOneKeyStartTask{}, s30)
			}
		}
		return RandMillisecond(600, 900)
	}
	for range t.C {
		t.Reset(f())
	}
}

// PlayerXZXS 仙宗 - 仙殿 - 仙宗悬赏
func (c *Connect) PlayerXZXS() error {
	body, err := proto.Marshal(&C2SPlayerXZXS{})
	if err != nil {
		return err
	}
	log.Println("[C][C2SPlayerXZXS]")
	return c.send(25401, body)
}

// XZXSGetAllCanStartTask 仙宗 - 仙殿 - 仙宗悬赏 -> 一键派遣任务
func (c *Connect) XZXSGetAllCanStartTask() error {
	body, err := proto.Marshal(&C2SXZXSGetAllCanStartTask{})
	if err != nil {
		return err
	}
	log.Println("[C][XZXSGetAllCanStartTask]")
	return c.send(25407, body)
}

// XZXSOneKeyStartTask 仙宗 - 仙殿 - 仙宗悬赏 -> 开始任务
func (c *Connect) XZXSOneKeyStartTask(start *S2CXZXSGetAllCanStartTask) error {
	body, err := proto.Marshal(&C2SXZXSOneKeyStartTask{Datas: start.Data})
	if err != nil {
		return err
	}
	log.Println("[C][XZXSOneKeyStartTask]")
	return c.send(25409, body)
}

// XZXSGetTaskPrize 仙宗 - 仙殿 - 仙宗悬赏 -> 一键领取奖励
func (c *Connect) XZXSGetTaskPrize() error {
	body, err := proto.Marshal(&C2SXZXSGetTaskPrize{TaskId: 0})
	if err != nil {
		return err
	}
	log.Println("[C][XZXSGetTaskPrize]")
	return c.send(25411, body)
}

////////////////////////////////////////////////////////////

func (x *S2CPlayerXZXS) ID() uint16 {
	return 25402
}

// Message PlayerXZXS 25402
func (x *S2CPlayerXZXS) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][PlayerXZXS] level=%v max_level=%v cur_exp=%v", x.Level, x.MaxLevel, x.CurExp)
}

////////////////////////////////////////////////////////////

func (x *S2CXZXSGetAllCanStartTask) ID() uint16 {
	return 25408
}

// Message S2CXZXSGetTaskPrize 25408
func (x *S2CXZXSGetAllCanStartTask) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XZXSGetAllCanStartTask] tag=%d %v", x.Tag, x)
}

////////////////////////////////////////////////////////////

func (x *S2CXZXSGetTaskPrize) ID() uint16 {
	return 25412
}

// Message S2CXZXSGetTaskPrize 25412
func (x *S2CXZXSGetTaskPrize) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XZXSGetTaskPrize] tag=%d %v", x.Tag, x)
}

////////////////////////////////////////////////////////////

func (x *S2CXZXSOneKeyStartTask) ID() uint16 {
	return 25410
}

// Message S2CXZXSGetTaskPrize 25412
func (x *S2CXZXSOneKeyStartTask) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][XZXSOneKeyStartTask] %v", x)
}
