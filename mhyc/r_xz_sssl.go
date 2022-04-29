package mhyc

import (
	"context"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

// XianDianSSSL 仙宗 - 仙殿 - 神兽试炼
func XianDianSSSL(ctx context.Context) {
	t := time.NewTimer(ms100)
	defer t.Stop()
	f := func() time.Duration {
		if RoleInfo.Get("SectGodAnimalChallenge").Int64() >= 35 {
			curr := time.Now()
			next := SelfWeekMonday(curr).Add(169 * time.Hour) // 168小时 = 7天
			return next.Sub(curr)
		}
		Fight.Lock()
		defer func() {
			Receive.Action(CLI.LeaveActive100)
			_ = Receive.Wait(&S2CLeaveActive{}, s3)
			Fight.Unlock()
		}()
		Receive.Action(CLI.GodAnimalPassData)
		_ = Receive.Wait(&S2CGodAnimalPassData{}, s3)
		Receive.Action(CLI.JoinActive100)
		_ = Receive.Wait(&S2CJoinActive{}, s3)
		k := time.NewTimer(ms100)
		defer k.Stop()
		for {
			select {
			case <-k.C:
				info := &S2CChallengeGodAnimal{}
				Receive.Action(CLI.ChallengeGodAnimal)
				_ = Receive.Wait(info, s3)
				if info.Tag == 48025 {
					return RandMillisecond(86400, 115200)
				}
				k.Reset(ms500)
			case <-ctx.Done():
				return s3
			}
		}
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

// SectGodAnimalData 仙宗 - 仙殿 - 神兽试炼 - 信息
func (c *Connect) SectGodAnimalData() error {
	body, err := proto.Marshal(&C2SChallengeGodAnimal{})
	if err != nil {
		return err
	}
	log.Println("[C][SectGodAnimalData]")
	return c.send(19043, body)
}

func (c *Connect) ChallengeGodAnimal() error {
	body, err := proto.Marshal(&C2SChallengeGodAnimal{})
	if err != nil {
		return err
	}
	log.Println("[C][ChallengeGodAnimal]")
	return c.send(19045, body)
}

func (c *Connect) GodAnimalPassData() error {
	body, err := proto.Marshal(&C2SGodAnimalPassData{})
	if err != nil {
		return err
	}
	log.Println("[C][GodAnimalPassData]")
	return c.send(19049, body)
}

func (c *Connect) JoinActive100() error {
	return c.JoinActive(&C2SJoinActive{AId: 100})
}

func (c *Connect) LeaveActive100() error {
	return c.LeaveActive(&C2SLeaveActive{AId: 100})
}

////////////////////////////////////////////////////////////

func (x *S2CSectGodAnimalData) ID() uint16 {
	return 19044
}

// Message S2CSectGodAnimalData 19044
func (x *S2CSectGodAnimalData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][SectGodAnimalData] %v", x)
}

////////////////////////////////////////////////////////////

func (x *S2CChallengeGodAnimal) ID() uint16 {
	return 19046
}

// Message S2CChallengeGodAnimal 19046
func (x *S2CChallengeGodAnimal) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][ChallengeGodAnimal] tag=%v is_win=%v", x.Tag, x.IsWin)
}

////////////////////////////////////////////////////////////

func (x *S2CGodAnimalPassData) ID() uint16 {
	return 19050
}

// Message S2CGodAnimalPassData 19050
func (x *S2CGodAnimalPassData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GodAnimalPassData] tag=%v %v", x.Tag, x)
}
