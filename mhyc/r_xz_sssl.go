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
			return TomorrowDuration(RandMillisecond(1800, 3600))
		}
		Fight.Lock()
		am := SetAction(ctx, "仙宗-仙殿-神兽试炼")
		defer func() {
			Receive.Action(CLI.LeaveActive100)
			_ = Receive.WaitWithContextOrTimeout(am.Ctx, &S2CLeaveActive{}, s3)
			am.End()
			Fight.Unlock()
		}()
		Receive.Action(CLI.GodAnimalPassData)
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CGodAnimalPassData{}, s3); err != nil {
			return RandMillisecond(3, 6)
		}
		Receive.Action(CLI.JoinActive100)
		if err := Receive.WaitWithContextOrTimeout(am.Ctx, &S2CJoinActive{}, s3); err != nil {
			return RandMillisecond(3, 6)
		}
		return am.RunAction(ctx, func() (loop time.Duration, next time.Duration) {
			Receive.Action(CLI.ChallengeGodAnimal)
			info := &S2CChallengeGodAnimal{}
			if err := Receive.WaitWithContextOrTimeout(am.Ctx, info, s3); err != nil {
				return 0, RandMillisecond(3, 6)
			}
			if info.Tag == 48025 {
				return 0, TomorrowDuration(RandMillisecond(1800, 3600))
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
	log.Printf("[S][ChallengeGodAnimal] tag=%v tag_msg=%s is_win=%v", x.Tag, GetTagMsg(x.Tag), x.IsWin)
}

////////////////////////////////////////////////////////////

func (x *S2CGodAnimalPassData) ID() uint16 {
	return 19050
}

// Message S2CGodAnimalPassData 19050
func (x *S2CGodAnimalPassData) Message(data []byte) {
	_ = proto.Unmarshal(data, x)
	log.Printf("[S][GodAnimalPassData] tag=%v tag_msg=%s %v", x.Tag, GetTagMsg(x.Tag), x)
}
